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

func Show(impl *uhppoted.UHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (interface{}, *errors.IError) {
	url := r.URL.Path

	matches := regexp.MustCompile("^/uhppote/acl/card/([0-9]+)$").FindStringSubmatch(url)
	if matches == nil || len(matches) < 2 {
		return nil, errors.Errorf(fmt.Errorf("%w: Missing card number/door", uhppoted.BadRequest), 0, "show", "Missing card number/door")
	}

	cardID, err := strconv.ParseUint(matches[1], 10, 32)
	if err != nil {
		return nil, errors.Errorf(fmt.Errorf("%w: Invalid card number (%s)", uhppoted.BadRequest, matches[1]), 0, "show", "Invalid card number")
	}

	u := ctx.Value("uhppote").(*uhppote.UHPPOTE)
	devices := ctx.Value("devices").([]*uhppote.Device)

	acl, err := api.GetCard(u, devices, uint32(cardID))
	if err != nil {
		return nil, errors.Errorf(fmt.Errorf("%w: Error retrieving card access permissions", err), 0, "show", "Error retrieving card access permissions")
	}

	response := []struct {
		Door      string     `json:"door"`
		StartDate types.Date `json:"start-date"`
		EndDate   types.Date `json:"end-date"`
	}{}

	for k, v := range acl {
		response = append(response, struct {
			Door      string     `json:"door"`
			StartDate types.Date `json:"start-date"`
			EndDate   types.Date `json:"end-date"`
		}{
			Door:      k,
			StartDate: v.From,
			EndDate:   v.To,
		})
	}

	return &response, nil
}
