package rest

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

func (d *dispatcher) authorized(r *http.Request) ([]string, error) {
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
			return d.basic(resource, action, credentials)

		case "Bearer":
			return d.bearer(resource, action, credentials)

		default:
			return []string{}, fmt.Errorf("Unsupported authorization scheme: '%s'", scheme)
		}
	}

	return []string{".*"}, nil
}

func (d *dispatcher) basic(resource, action, credentials string) ([]string, error) {
	var uid string
	var pwd string

	plaintext, err := base64.StdEncoding.DecodeString(credentials)
	if err != nil {
		return []string{}, err
	}

	tokens := strings.Split(string(plaintext), ":")
	if len(tokens) > 0 {
		uid = tokens[0]
	}

	if len(tokens) > 1 {
		pwd = tokens[1]
	}

	if err := d.auth.VerifyPassword(uid, pwd); err != nil {
		return []string{}, err
	}

	if err := d.auth.Authorize(resource, action, uid); err != nil {
		return []string{}, err
	}

	return d.auth.Cards(uid), nil
}

func (d *dispatcher) bearer(resource, action, credentials string) ([]string, error) {
	plaintext, err := base64.StdEncoding.DecodeString(credentials)
	if err != nil {
		return []string{}, err
	}

	if strings.HasPrefix(string(plaintext), "HOTP:") {
		return d.hotp(resource, action, string(plaintext)[5:])

	}

	return []string{}, fmt.Errorf("Unsupported bearer token")
}

func (d *dispatcher) hotp(resource, action, credentials string) ([]string, error) {
	var uid string
	var otp string

	tokens := strings.Split(credentials, ":")
	if len(tokens) > 0 {
		uid = tokens[0]
	}

	if len(tokens) > 1 {
		otp = tokens[1]
	}

	if err := d.auth.VerifyHOTP(uid, otp); err != nil {
		return []string{}, err
	}

	if err := d.auth.Authorize(resource, action, uid); err != nil {
		return []string{}, err
	}

	return d.auth.Cards(uid), nil
}
