package acl

import (
	"context"
	"fmt"
	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
	api "github.com/uhppoted/uhppoted-api/acl"
	"net/http"
	"regexp"
	"strconv"
)

func Show(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path

	matches := regexp.MustCompile("^/uhppote/acl/card/([0-9]+)$").FindStringSubmatch(url)
	if matches == nil || len(matches) < 2 {
		warn(ctx, "show", fmt.Errorf("Missing card number"))
		http.Error(w, "Invalid request: missing card number", http.StatusBadRequest)
		return
	}

	cardID, err := strconv.ParseUint(matches[1], 10, 32)
	if err != nil {
		warn(ctx, "show", fmt.Errorf("Invalid card number '%s' (%w)", matches[1], err))
		http.Error(w, "Invalid card number", http.StatusBadRequest)
		return
	}

	u := ctx.Value("uhppote").(*uhppote.UHPPOTE)
	devices := ctx.Value("devices").([]*uhppote.Device)

	acl, err := api.GetCard(u, devices, uint32(cardID))
	if err != nil {
		warn(ctx, "ACL::show", err)
		http.Error(w, "Error retrieving card access permissions", http.StatusInternalServerError)
		return
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

	reply(ctx, w, response)
}
