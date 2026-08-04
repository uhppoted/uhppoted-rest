package acl

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"codeberg.org/uhppoted/uhppoted-core/types"
	"codeberg.org/uhppoted/uhppoted-core/uhppote"
	api "codeberg.org/uhppoted/uhppoted-lib/acl"
	"codeberg.org/uhppoted/uhppoted-lib/uhppoted"

	"codeberg.org/uhppoted/uhppoted-rest/errors"
	"codeberg.org/uhppoted/uhppoted-rest/lib"
)

func Show(impl uhppoted.IUHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, any, error) {
	url := r.URL.Path

	matches := regexp.MustCompile("^/uhppote/acl/card/([0-9]+)$").FindStringSubmatch(url)
	if len(matches) < 2 {
		return http.StatusBadRequest,
			errors.NewRESTError("show", "Missing card number/door"),
			fmt.Errorf("missing card number/door in request URL %s)", url)
	}

	cardID, err := strconv.ParseUint(matches[1], 10, 32)
	if err != nil {
		return http.StatusBadRequest,
			errors.NewRESTError("show", "Invalid card number"),
			err
	}

	u := ctx.Value(lib.Uhppote).(uhppote.IUHPPOTE)
	devices := ctx.Value(lib.Devices).([]uhppote.Device)

	acl, err := api.GetCard(u, devices, uint32(cardID))
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("show", "Error retrieving card access permissions"),
			err
	}

	permissions := []struct {
		Door      string     `json:"door"`
		Profile   int        `json:"profile,omitempty"`
		StartDate types.Date `json:"start-date"`
		EndDate   types.Date `json:"end-date"`
		FirstCard bool       `json:"first-card,omitempty"`
	}{}

	for k, v := range acl {
		permissions = append(permissions, struct {
			Door      string     `json:"door"`
			Profile   int        `json:"profile,omitempty"`
			StartDate types.Date `json:"start-date"`
			EndDate   types.Date `json:"end-date"`
			FirstCard bool       `json:"first-card,omitempty"`
		}{
			Door:      k,
			StartDate: v.From,
			EndDate:   v.To,
			Profile:   v.Profile,
			FirstCard: v.FirstCard,
		})
	}

	return http.StatusOK, &struct {
		Permissions []struct {
			Door      string     `json:"door"`
			Profile   int        `json:"profile,omitempty"`
			StartDate types.Date `json:"start-date"`
			EndDate   types.Date `json:"end-date"`
			FirstCard bool       `json:"first-card,omitempty"`
		} `json:"permissions"`
	}{
		Permissions: permissions,
	}, nil
}
