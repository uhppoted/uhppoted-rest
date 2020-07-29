package rest

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

func (d *dispatcher) authorized(r *http.Request) error {
	if d.auth.Enabled() {
		var scheme string
		var credentials string

		for key, headers := range r.Header {
			if http.CanonicalHeaderKey(key) == "Authorization" {
				for _, header := range headers {
					tokens := strings.Split(header, " ")
					if len(tokens) > 0 {
						scheme = tokens[0]
					}

					if len(tokens) > 1 {
						credentials = tokens[1]
					}
				}
			}
		}

		resource := strings.TrimSpace(r.URL.Path)
		action := strings.ToLower(r.Method)

		switch scheme {
		case "Basic":
			if err := d.basic(resource, action, credentials); err != nil {
				return err
			}

		default:
			return fmt.Errorf("Unsupported authorization scheme: '%s'", scheme)
		}
	}

	return nil
}

func (d *dispatcher) basic(resource, action, credentials string) error {
	var uid string
	var pwd string

	plaintext, err := base64.StdEncoding.DecodeString(credentials)
	if err != nil {
		return err
	}

	tokens := strings.Split(string(plaintext), ":")
	if len(tokens) > 0 {
		uid = tokens[0]
	}

	if len(tokens) > 1 {
		pwd = tokens[1]
	}

	return d.auth.Authorize(resource, action, uid, pwd)
}
