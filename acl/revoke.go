package acl

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/uhppoted/uhppote-core/uhppote"
	api "github.com/uhppoted/uhppoted-lib/acl"
	"github.com/uhppoted/uhppoted-lib/uhppoted"

	"github.com/uhppoted/uhppoted-rest/errors"
	"github.com/uhppoted/uhppoted-rest/lib"
)

func Revoke(impl uhppoted.IUHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, any, error) {
	url := r.URL.Path

	matches := regexp.MustCompile(`^/uhppote/acl/card/([0-9]+)/door/(\S.*)$`).FindStringSubmatch(url)
	if len(matches) < 3 {
		return http.StatusBadRequest,
			errors.NewRESTError("revoke", "Missing card number/door"),
			fmt.Errorf("missing card number/door in request URL (%s)", url)
	}

	cardID, err := strconv.ParseUint(matches[1], 10, 32)
	if err != nil {
		return http.StatusBadRequest,
			errors.NewRESTError("revoke", "Invalid card number"),
			err
	}

	door := matches[2]
	if strings.TrimSpace(door) == "" {
		return http.StatusBadRequest,
			errors.NewRESTError("revoke", "Invalid door"),
			fmt.Errorf("invalid door (%s) in request URL", matches[1])
	}

	u := ctx.Value(lib.Uhppote).(uhppote.IUHPPOTE)
	devices := ctx.Value(lib.Devices).([]uhppote.Device)

	err = api.Revoke(u, devices, uint32(cardID), []string{door})
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("revoke", "Error revoking card access permissions"),
			err
	}

	return http.StatusOK, nil, nil
}
