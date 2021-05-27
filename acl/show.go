package acl

import (
	"context"
	"fmt"
	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
	api "github.com/uhppoted/uhppoted-api/acl"
	"github.com/uhppoted/uhppoted-api/uhppoted"
	"github.com/uhppoted/uhppoted-rest/errors"
	"net/http"
	"regexp"
	"strconv"
)

func Show(impl *uhppoted.UHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	url := r.URL.Path

	matches := regexp.MustCompile("^/uhppote/acl/card/([0-9]+)$").FindStringSubmatch(url)
	if matches == nil || len(matches) < 2 {
		return http.StatusBadRequest,
			errors.NewRESTError("show", "Missing card number/door"),
			fmt.Errorf("Missing card number/door in request URL %s)", url)
	}

	cardID, err := strconv.ParseUint(matches[1], 10, 32)
	if err != nil {
		return http.StatusBadRequest,
			errors.NewRESTError("show", "Invalid card number"),
			err
	}

	u := ctx.Value("uhppote").(uhppote.IUHPPOTE)
	devices := ctx.Value("devices").([]uhppote.Device)

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
	}{}

	for k, v := range acl {
		permissions = append(permissions, struct {
			Door      string     `json:"door"`
			Profile   int        `json:"profile,omitempty"`
			StartDate types.Date `json:"start-date"`
			EndDate   types.Date `json:"end-date"`
		}{
			Door:      k,
			StartDate: v.From,
			EndDate:   v.To,
		})
	}

	return http.StatusOK, &struct {
		Permissions []struct {
			Door      string     `json:"door"`
			Profile   int        `json:"profile,omitempty"`
			StartDate types.Date `json:"start-date"`
			EndDate   types.Date `json:"end-date"`
		} `json:"permissions"`
	}{
		Permissions: permissions,
	}, nil
}
