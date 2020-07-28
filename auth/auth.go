package auth

import (
	"fmt"
	"regexp"
)

type IAuth interface {
	Authorize(resource, user, password string) error
}

type LocalAuth struct {
	permissions map[*regexp.Regexp][]string
	users       map[string]string
	groups      map[string][]string
}

func NewLocalAuth() *LocalAuth {
	return &LocalAuth{
		permissions: map[*regexp.Regexp][]string{
			regexp.MustCompile("/uhppote/device/[0-9]+/door/[0-9]+/swipes::POST"): []string{"admin"},
		},

		users: map[string]string{
			"him": "ghjkl",
			"her": "uiop",
			"it":  "asdf",
			"me":  "qwerty",
		},

		groups: map[string][]string{
			"admin": []string{"me"},
		},
	}
}

func (a *LocalAuth) Authorize(resource, user, password string) error {
	for re, groups := range a.permissions {
		if re.Match([]byte(resource)) {
			for _, g := range groups {
				l, _ := a.groups[g]
				for _, u := range l {
					if user == u {
						if pwd, ok := a.users[u]; ok && password == pwd {
							return nil
						}
					}
				}
			}
		}
	}

	return fmt.Errorf("Invalid credentials %v, %v", user, password)
}
