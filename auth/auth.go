package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/uhppoted/uhppoted-api/kvs"
)

type IAuth interface {
	Enabled() bool
	Authorize(resource, action, user, password string) error
}

type AuthProvider struct {
	enabled bool
	users   *kvs.KeyValueStore
	groups  *kvs.KeyValueStore
}

type permission struct {
	resource *regexp.Regexp
	action   *regexp.Regexp
}

type user struct {
	Password string   `json:"password"`
	Groups   []string `json:"groups"`
}

func (p permission) String() string {
	return fmt.Sprintf("resource:`%s` action:`%s`", p.resource, p.action)
}

func NewAuthProvider(enabled bool, users, groups string, logger *log.Logger) (*AuthProvider, error) {
	separator := regexp.MustCompile(`\s*,\s*`)

	f := func(value string) (interface{}, error) {
		u := user{}
		err := json.Unmarshal([]byte(value), &u)
		if err != nil {
			return nil, err
		}

		return &u, nil
	}

	g := func(value string) (interface{}, error) {
		permissions := []permission{}
		re := regexp.MustCompile(`(.*?):(.*)`)
		tokens := separator.Split(value, -1)
		for _, s := range tokens {
			if match := re.FindStringSubmatch(s); len(match) == 3 {
				resource, err := regexp.Compile("^" + strings.ReplaceAll(match[1], "*", ".*") + "$")
				if err != nil {
					return permissions, err
				}

				action, err := regexp.Compile("^" + strings.ReplaceAll(match[2], "*", ".*") + "$")
				if err != nil {
					return permissions, err
				}

				permissions = append(permissions, permission{
					resource: resource,
					action:   action,
				})
			}
		}

		return permissions, nil
	}

	provider := AuthProvider{
		enabled: enabled,
		users:   kvs.NewKeyValueStore("permissions:users", f),
		groups:  kvs.NewKeyValueStore("permissions:groups", g),
	}

	if enabled {
		err := provider.users.LoadFromFile(users)
		if err != nil {
			return nil, err
		}

		err = provider.groups.LoadFromFile(groups)
		if err != nil {
			return nil, err
		}

		provider.users.Watch(users, logger)
		provider.groups.Watch(groups, logger)
	}

	return &provider, nil
}

func (a *AuthProvider) Enabled() bool {
	if a == nil {
		return false
	}

	return a.enabled
}

func (a *AuthProvider) Authorize(resource, action, uid, pwd string) error {
	if !a.Enabled() {
		return nil
	}

	u, ok := a.users.Get(uid)
	if !ok {
		return fmt.Errorf("%s: Not a member of any permissions groups", uid)
	}

	if pwd != u.(*user).Password {
		return fmt.Errorf("Invalid credentials %v, %v", uid, pwd)
	}

	for _, g := range u.(*user).Groups {
		if permissions, ok := a.groups.Get(g); ok {
			for _, q := range permissions.([]permission) {
				if q.resource.MatchString(resource) && q.action.MatchString(action) {
					return nil
				}
			}
		}
	}

	return fmt.Errorf("%s: Not authorised for %s:%s", uid, resource, action)
}
