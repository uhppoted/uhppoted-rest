package acl

import (
	"context"
	"fmt"
	"github.com/uhppoted/uhppote-core/uhppote"
	api "github.com/uhppoted/uhppoted-api/acl"
	"github.com/uhppoted/uhppoted-api/uhppoted"
	"github.com/uhppoted/uhppoted-rest/errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func Revoke(impl *uhppoted.UHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (interface{}, *errors.IError) {
	url := r.URL.Path

	matches := regexp.MustCompile("^/uhppote/acl/card/([0-9]+)/door/(\\S.*)$").FindStringSubmatch(url)
	if matches == nil || len(matches) < 3 {
		return nil, errors.Errorf(fmt.Errorf("%w: Missing card number/door", uhppoted.BadRequest), 0, "revoke", "Missing card number/door")
	}

	cardID, err := strconv.ParseUint(matches[1], 10, 32)
	if err != nil {
		return nil, errors.Errorf(fmt.Errorf("%w: Invalid card number (%s)", uhppoted.BadRequest, matches[1]), 0, "revoke", "Invalid card number")
	}

	door := matches[2]
	if strings.TrimSpace(door) == "" {
		return nil, errors.Errorf(fmt.Errorf("%w: Invalid door (%s)", uhppoted.BadRequest, matches[1]), 0, "revoke", "Invalid door")
	}

	u := ctx.Value("uhppote").(*uhppote.UHPPOTE)
	devices := ctx.Value("devices").([]*uhppote.Device)

	err = api.Revoke(u, devices, uint32(cardID), []string{door})
	if err != nil {
		return nil, errors.Errorf(fmt.Errorf("%w: Error revoking card access permissions", err), 0, "revoke", "Error revoking card access permissions")
	}

	return nil, nil
}
