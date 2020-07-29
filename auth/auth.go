package auth

import (
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

func (p permission) String() string {
	return fmt.Sprintf("resource:`%s` action:`%s`", p.resource, p.action)
}

func NewAuthProvider(enabled bool, users, groups string, logger *log.Logger) (*AuthProvider, error) {
	separator := regexp.MustCompile(`\s*,\s*`)

	u := func(value string) (interface{}, error) {
		return separator.Split(value, -1), nil
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
		users:   kvs.NewKeyValueStore("permissions:users", u),
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

func (a *AuthProvider) Authorize(resource, action, user, password string) error {
	if !a.Enabled() {
		return nil
	}

	//	return fmt.Errorf("Invalid credentials %v, %v", user, password)

	groups, ok := a.users.Get(user)
	if !ok {
		return fmt.Errorf("%s: Not a member of any permissions groups", user)
	}

	for _, g := range groups.([]string) {
		if permissions, ok := a.groups.Get(g); ok {
			for _, q := range permissions.([]permission) {
				if q.resource.MatchString(resource) && q.action.MatchString(action) {
					return nil
				}
			}
		}
	}

	return fmt.Errorf("%s: Not authorised for %s:%s", user, resource, action)
}
