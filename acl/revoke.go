package acl

import (
	"context"
	"fmt"
	"github.com/uhppoted/uhppote-core/uhppote"
	api "github.com/uhppoted/uhppoted-api/acl"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func Revoke(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path

	matches := regexp.MustCompile("^/uhppote/acl/card/([0-9]+)/doors/(\\S.*)$").FindStringSubmatch(url)
	if matches == nil || len(matches) < 3 {
		warn(ctx, "revoke", fmt.Errorf("Missing card number/door"))
		http.Error(w, "Invalid request: missing card number/door", http.StatusBadRequest)
		return
	}

	cardID, err := strconv.ParseUint(matches[1], 10, 32)
	if err != nil {
		warn(ctx, "revoke", fmt.Errorf("Invalid card number '%s' (%w)", matches[1], err))
		http.Error(w, "Invalid card number", http.StatusBadRequest)
		return
	}

	doors := []string{}
	tokens := strings.Split(matches[2], ",")
	for _, s := range tokens {
		if d := strings.TrimSpace(s); d != "" {
			doors = append(doors, d)
		}
	}

	if len(doors) == 0 {
		warn(ctx, "revoke", fmt.Errorf("Invalid list of doors '%s' (%w)", matches[2], err))
		http.Error(w, "Invalid list of doors", http.StatusBadRequest)
		return
	}

	u := ctx.Value("uhppote").(*uhppote.UHPPOTE)
	devices := ctx.Value("devices").([]*uhppote.Device)

	err = api.Revoke(u, devices, uint32(cardID), doors)
	if err != nil {
		warn(ctx, "revoke", err)
		http.Error(w, "Error revoking card access permissions", http.StatusInternalServerError)
		return
	}
}
