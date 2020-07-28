package rest

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

var permissions = map[*regexp.Regexp][]string{
	regexp.MustCompile("/uhppote/device/[0-9]+/door/[0-9]+/swipes::POST"): []string{"admin"},
}

var users = map[string]string{
	"him": "ghjkl",
	"her": "uiop",
	"it":  "asdf",
	"me":  "qwerty",
}

var groups = map[string][]string{
	"admin": []string{"me"},
}

func (d *dispatcher) authorized(r *http.Request) error {
	if d.authEnabled {
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

		resource := strings.TrimSpace(r.URL.Path) + "::" + strings.ToUpper(r.Method)

		for re, groups := range permissions {
			if re.Match([]byte(resource)) {
				switch scheme {
				case "Basic":
					if err := basic(resource, credentials, groups); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func basic(resource, credentials string, allowed []string) error {
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

	for _, g := range allowed {
		l, _ := groups[g]
		for _, u := range l {
			if uid == u {
				for k, v := range users {
					if k == uid && v == pwd {
						return nil
					}
				}
			}
		}

	}

	return fmt.Errorf("Invalid credentials %v, %v", uid, pwd)
}
