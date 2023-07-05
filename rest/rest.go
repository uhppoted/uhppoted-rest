package rest

import (
	"compress/gzip"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/uhppoted/uhppote-core/uhppote"
	"github.com/uhppoted/uhppoted-lib/uhppoted"

	"github.com/uhppoted/uhppoted-rest/acl"
	"github.com/uhppoted/uhppoted-rest/auth"
	"github.com/uhppoted/uhppoted-rest/device"
	"github.com/uhppoted/uhppoted-rest/lib"
	"github.com/uhppoted/uhppoted-rest/log"
)

// OpenAPI is a container for the runtime flags for the Open API user interface
// implementation. Intended for development use only.
type OpenAPI struct {
	// Enabled enables the Open API user interface if true. Should be false in production.
	Enabled bool
	// Directory sets the directory for the Open API user interface HTTP resources.
	Directory string
}

// RESTD is a container for the runtime information for the REST daemon. Isn't really exported
// but (temporarily) capitalized here pending a better name.
type RESTD struct {
	// HTTPEnabled enables HTTP connections to the REST daemon.
	HTTPEnabled bool

	//HTTPPort is the HTTP port assigned to the REST daemon.
	HTTPPort uint16

	// HTTPSEnabled enables HTTPS connections to the REST daemon.
	HTTPSEnabled bool

	//HTTPSPort is the HTTPS port assigned to the REST daemon.
	HTTPSPort uint16

	//TLSKeyFile is the path the the HTTPS server key PEM file.
	TLSKeyFile string

	//TLSKeyFile is the path the the HTTPS server certificate PEM file.
	TLSCertificateFile string

	//CACertificateFile is the path the the HTTPS CA certificate PEM file used to verify client certificates.
	CACertificateFile string

	// RequireClientCertificates makes client TLS certificates mandatory. Otherwise the client certificate is validated
	// if presented.
	RequireClientCertificates bool

	//CORSEnabled allows CORS requests if true. Should be false in production.
	CORSEnabled bool

	//AuthEnabled enables validation of the Authorization header.
	AuthEnabled bool

	//AuthUsers is the filepath for the 'users' file that stores the authorized users
	AuthUsers string

	//AuthGroups is the filepath for the 'groups' file that stores the group permissions
	AuthGroups string

	//HOTPWindow is the allowable 'window' for HOTP counter values
	HOTPWindow uint64

	//HOTPCounters is the filepath for the 'counters' file that stores the HOTP counters
	HOTPCounters string

	//Protocol version sets the version used for encoding JSON structs in responses
	Protocol string

	//OpenAPI runtime flags.
	OpenAPI
}

type handlerfn func(uhppoted.IUHPPOTED, context.Context, http.ResponseWriter, *http.Request) (int, interface{}, error)

type handler struct {
	re     *regexp.Regexp
	method string
	fn     handlerfn
}

type dispatcher struct {
	corsEnabled bool
	uhppote     uhppote.IUHPPOTE
	uhppoted    *uhppoted.UHPPOTED
	devices     []uhppote.Device
	handlers    []handler
	auth        auth.IAuth
	openapi     http.Handler
}

// Run configures and starts the REST daemon HTTP and HTTPS request listeners. It returns once the listen
// connections have been closed.
func (r *RESTD) Run(u uhppote.IUHPPOTE, devices []uhppote.Device) {
	auth, err := auth.NewAuthProvider(r.AuthEnabled, r.AuthUsers, r.AuthGroups, r.HOTPCounters, r.HOTPWindow)
	if err != nil {
		log.Fatalf("RESTD", "error initialising AuthProvider (%v)", err)
	}

	device.SetProtocol(r.Protocol)

	d := dispatcher{
		uhppote: u,
		uhppoted: &uhppoted.UHPPOTED{
			UHPPOTE:         u,
			ListenBatchSize: 32,
		},
		devices: devices,

		handlers: []handler{
			handler{regexp.MustCompile("^/uhppote/device$"), http.MethodGet, device.GetDevices},
			handler{regexp.MustCompile("^/uhppote/device/[0-9]+$"), http.MethodGet, device.GetDevice},
			handler{regexp.MustCompile("^/uhppote/device/[0-9]+/status$"), http.MethodGet, device.GetStatus},
			handler{regexp.MustCompile("^/uhppote/device/[0-9]+/time$"), http.MethodGet, device.GetTime},
			handler{regexp.MustCompile("^/uhppote/device/[0-9]+/time$"), http.MethodPut, device.SetTime},
			handler{regexp.MustCompile("^/uhppote/device/[0-9]+/door/[1-4]$"), http.MethodGet, device.GetDoor},
			handler{regexp.MustCompile("^/uhppote/device/[0-9]+/door/[1-4]/swipes$"), http.MethodPost, device.OpenDoor},
			handler{regexp.MustCompile("^/uhppote/device/[0-9]+/door/[1-4]/delay$"), http.MethodPut, device.SetDoorDelay},
			handler{regexp.MustCompile("^/uhppote/device/[0-9]+/door/[1-4]/control$"), http.MethodPut, device.SetDoorControl},
			handler{regexp.MustCompile("^/uhppote/device/[0-9]+/door/interlock$"), http.MethodPut, device.SetInterlock},
			handler{regexp.MustCompile("^/uhppote/device/[0-9]+/door/keypads$"), http.MethodPut, device.ActivateKeypads},

			handler{regexp.MustCompile("^/uhppote/device/[0-9]+/cards$"), http.MethodGet, device.GetCards},
			handler{regexp.MustCompile("^/uhppote/device/[0-9]+/cards$"), http.MethodDelete, device.DeleteCards},
			handler{regexp.MustCompile("^/uhppote/device/[0-9]+/card/[0-9]+$"), http.MethodGet, device.GetCard},
			handler{regexp.MustCompile("^/uhppote/device/[0-9]+/card/[0-9]+$"), http.MethodPut, device.PutCard},
			handler{regexp.MustCompile("^/uhppote/device/[0-9]+/card/[0-9]+$"), http.MethodDelete, device.DeleteCard},

			handler{regexp.MustCompile("^/uhppote/device/[0-9]+/time-profile/[0-9]+$"), http.MethodGet, device.GetTimeProfile},
			handler{regexp.MustCompile("^/uhppote/device/[0-9]+/time-profile/[0-9]+$"), http.MethodPut, device.PutTimeProfile},
			handler{regexp.MustCompile("^/uhppote/device/[0-9]+/time-profiles"), http.MethodGet, device.GetTimeProfiles},
			handler{regexp.MustCompile("^/uhppote/device/[0-9]+/time-profiles"), http.MethodPut, device.PutTimeProfiles},
			handler{regexp.MustCompile("^/uhppote/device/[0-9]+/time-profiles"), http.MethodDelete, device.ClearTimeProfiles},

			handler{regexp.MustCompile("^/uhppote/device/[0-9]+/tasklist"), http.MethodPut, device.PutTaskList},

			handler{regexp.MustCompile("^/uhppote/device/[0-9]+/events$"), http.MethodGet, device.GetEvents},
			handler{regexp.MustCompile("^/uhppote/device/[0-9]+/events/([0-9]+)$"), http.MethodGet, device.GetEvents},
			handler{regexp.MustCompile("^/uhppote/device/[0-9]+/event/([0-9]+|first|last|current|next)$"), http.MethodGet, device.GetEvent},
			handler{regexp.MustCompile("^/uhppote/device/[0-9]+/special-events$"), http.MethodPut, device.SpecialEvents},

			handler{regexp.MustCompile("^/uhppote/acl$"), http.MethodGet, acl.GetACL},
			handler{regexp.MustCompile("^/uhppote/acl$"), http.MethodPut, acl.PutACL},
			handler{regexp.MustCompile("^/uhppote/acl/card/[0-9]+$"), http.MethodGet, acl.Show},
			handler{regexp.MustCompile(`^/uhppote/acl/card/[0-9]+/door/\S.*$`), http.MethodPut, acl.Grant},
			handler{regexp.MustCompile(`^/uhppote/acl/card/[0-9]+/door/\S.*$`), http.MethodDelete, acl.Revoke},
		},

		corsEnabled: r.CORSEnabled,
		auth:        auth,
		openapi:     http.NotFoundHandler(),
	}

	if r.OpenAPI.Enabled {
		d.openapi = http.StripPrefix("/openapi", http.FileServer(http.Dir(r.OpenAPI.Directory)))
	}

	var wg sync.WaitGroup

	if r.HTTPEnabled {
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Infof("RESTD", "... listening on port %d\n", r.HTTPPort)
			if err := http.ListenAndServe(fmt.Sprintf(":%d", r.HTTPPort), &d); err != nil {
				log.Fatalf("RESTD", "%v", err)
			}
		}()
	}

	if r.HTTPSEnabled {
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Infof("RESTD", "... listening on port %d\n", r.HTTPSPort)

			ca, err := os.ReadFile(r.CACertificateFile)
			if err != nil {
				log.Fatalf("RESTD", "%v", err)
			}

			certificates := x509.NewCertPool()
			if !certificates.AppendCertsFromPEM(ca) {
				log.Fatalf("RESTD", "Unable failed to parse CA certificate")
			}

			tlsConfig := tls.Config{
				ClientAuth: tls.RequireAndVerifyClientCert,
				ClientCAs:  certificates,
				CipherSuites: []uint16{
					tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				},
				PreferServerCipherSuites: true,
				MinVersion:               tls.VersionTLS12,
			}

			if r.RequireClientCertificates {
				tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
			} else {
				tlsConfig.ClientAuth = tls.VerifyClientCertIfGiven
			}

			tlsConfig.BuildNameToCertificate()

			httpsd := &http.Server{
				Addr:      fmt.Sprintf(":%d", r.HTTPSPort),
				Handler:   &d,
				TLSConfig: &tlsConfig,
			}

			if err := httpsd.ListenAndServeTLS(r.TLSCertificateFile, r.TLSKeyFile); err != nil {
				log.Fatalf("RESTD", "%v", err)
			}
		}()
	}

	wg.Wait()
}

// Close gracefully releases any long-held resources on terminating the REST daemon. The current
// implementation is a placeholder.
func Close() {
}

func (d *dispatcher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	compression := "none"

	for key, headers := range r.Header {
		if http.CanonicalHeaderKey(key) == "Accept-Encoding" {
			for _, header := range headers {
				encodings := strings.Split(header, ",")
				for _, encoding := range encodings {
					if strings.TrimSpace(encoding) == "gzip" {
						compression = "gzip"
					}
				}
			}
		}
	}

	if d.corsEnabled {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// CORS pre-flight request ?
		if r.Method == http.MethodOptions {
			return
		}
	}

	// OpenAPI ?

	if strings.HasPrefix(r.URL.Path, "/openapi") {
		d.openapi.ServeHTTP(w, r)
		return
	}

	// Dispatch to handler
	url := r.URL.Path
	for _, h := range d.handlers {
		if h.re.MatchString(url) && r.Method == h.method {
			cards, err := d.authorized(r)
			if err != nil {
				log.Warnf("RESTD", "%v", err)
				http.Error(w, "Access denied", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(context.Background(), lib.Uhppote, d.uhppote)
			ctx = context.WithValue(ctx, lib.Devices, d.devices)
			ctx = context.WithValue(ctx, lib.Compression, compression)
			ctx = context.WithValue(ctx, lib.AuthorizedCards, cards)
			ctx = parse(ctx, r)

			status, response, err := h.fn(d.uhppoted, ctx, w, r)
			if err != nil {
				log.Warnf("RESTD", "%v", err)
			}

			reply(ctx, w, status, response)
			return
		}
	}

	// Fall-through handler
	http.Error(w, "Unsupported API", http.StatusNotImplemented)
}

func parse(ctx context.Context, r *http.Request) context.Context {
	url := r.URL.Path

	matches := regexp.MustCompile("^/uhppote/device/([0-9]+)(?:$|/.*$)").FindStringSubmatch(url)
	if matches != nil {
		deviceID, err := strconv.ParseUint(matches[1], 10, 32)
		if err == nil {
			ctx = context.WithValue(ctx, lib.DeviceID, uint32(deviceID))
		}
	}

	matches = regexp.MustCompile("^/uhppote/device/[0-9]+/door/([1-4])(?:$|/.*$)").FindStringSubmatch(url)
	if matches != nil {
		door, err := strconv.ParseUint(matches[1], 10, 8)
		if err == nil {
			ctx = context.WithValue(ctx, lib.Door, uint8(door))
		}
	}

	matches = regexp.MustCompile("^/uhppote/device/[0-9]+/card/([0-9]+)$").FindStringSubmatch(url)
	if matches != nil {
		cardNumber, err := strconv.ParseUint(matches[1], 10, 32)
		if err == nil {
			ctx = context.WithValue(ctx, lib.CardNumber, uint32(cardNumber))
		}
	}

	return ctx
}

func reply(ctx context.Context, w http.ResponseWriter, status int, response interface{}) {
	var err error
	b := []byte{}
	if response != nil {
		b, err = json.Marshal(response)
		if err != nil {
			http.Error(w, "Error generating response", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")

	if len(b) > 1024 && ctx.Value(lib.Compression) == "gzip" {
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(status)
		encoder := gzip.NewWriter(w)
		encoder.Write(b)
		encoder.Flush()
	} else {
		w.WriteHeader(status)
		w.Write(b)
	}
}
