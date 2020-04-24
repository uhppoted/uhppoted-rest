package acl

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
	api "github.com/uhppoted/uhppoted-api/acl"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func Grant(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path

	matches := regexp.MustCompile("^/uhppote/acl/card/([0-9]+)/doors/(\\S.*)$").FindStringSubmatch(url)
	if matches == nil || len(matches) < 3 {
		warn(ctx, "grant", fmt.Errorf("Missing card number/door"))
		http.Error(w, "Invalid request: missing card number/door", http.StatusBadRequest)
		return
	}

	cardID, err := strconv.ParseUint(matches[1], 10, 32)
	if err != nil {
		warn(ctx, "grant", fmt.Errorf("Invalid card number '%s' (%w)", matches[1], err))
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
		warn(ctx, "grant", fmt.Errorf("Invalid list of doors '%s' (%w)", matches[2], err))
		http.Error(w, "Invalid list of doors", http.StatusBadRequest)
		return
	}

	fmt.Printf("%v\n", doors)

	blob, err := ioutil.ReadAll(r.Body)
	if err != nil {
		warn(ctx, "grant", err)
		http.Error(w, "Error reading request", http.StatusInternalServerError)
		return
	}

	body := struct {
		From *types.Date `json:"start-date"`
		To   *types.Date `json:"end-date"`
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

	err = api.Grant(u, devices, uint32(cardID), *body.From, *body.To, doors)
	if err != nil {
		warn(ctx, "ACL::grant", err)
		http.Error(w, "Error granting card access permissions", http.StatusInternalServerError)
		return
	}
}
