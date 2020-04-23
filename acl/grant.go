package acl

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
	api "github.com/uhppoted/uhppoted-api/acl"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

func Grant(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path

	matches := regexp.MustCompile("^/uhppote/acl/card/([0-9]+)$").FindStringSubmatch(url)
	if matches == nil {
		warn(ctx, "grant", fmt.Errorf("Missing card number"))
		http.Error(w, "Invalid request: missing card number", http.StatusBadRequest)
		return
	}

	cardID, err := strconv.ParseUint(matches[1], 10, 32)
	if err != nil {
		warn(ctx, "grant", fmt.Errorf("Invalid card number '%s' (%w)", matches[1], err))
		http.Error(w, "Invalid card number", http.StatusBadRequest)
		return
	}

	blob, err := ioutil.ReadAll(r.Body)
	if err != nil {
		warn(ctx, "grant", err)
		http.Error(w, "Error reading request", http.StatusInternalServerError)
		return
	}

	body := struct {
		From  *types.Date `json:"start-date"`
		To    *types.Date `json:"end-date"`
		Doors []string    `json:"doors"`
	}{}

	if err = json.Unmarshal(blob, &body); err != nil {
		warn(ctx, "grant", err)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if body.From == nil {
		warn(ctx, "grant", fmt.Errorf("Missing 'start-date'"))
		http.Error(w, "Invalid request: missing 'start-date'", http.StatusBadRequest)
		return
	}

	if body.To == nil {
		warn(ctx, "grant", fmt.Errorf("Missing 'end-date'"))
		http.Error(w, "Invalid request: missing 'end-date'", http.StatusBadRequest)
		return
	}

	u := ctx.Value("uhppote").(*uhppote.UHPPOTE)
	devices := ctx.Value("devices").([]*uhppote.Device)

	err = api.Grant(u, devices, uint32(cardID), *body.From, *body.To, body.Doors)
	if err != nil {
		warn(ctx, "grant", err)
		http.Error(w, "Error granting card access permissions", http.StatusInternalServerError)
		return
	}
}

func warn(ctx context.Context, operation string, err error) {
	ctx.Value("log").(*log.Logger).Printf("WARN  %-20s %v\n", operation, err)
}
