package acl

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
	api "github.com/uhppoted/uhppoted-api/acl"
	"github.com/uhppoted/uhppoted-api/uhppoted"
	"github.com/uhppoted/uhppoted-rest/errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func Grant(impl *uhppoted.UHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	url := r.URL.Path

	matches := regexp.MustCompile("^/uhppote/acl/card/([0-9]+)/door/(\\S.*)$").FindStringSubmatch(url)
	if matches == nil || len(matches) < 3 {
		return http.StatusBadRequest,
			errors.NewRESTError("grant", "Missing card number/door"),
			fmt.Errorf("Missing card number/door (%s)", url)
	}

	cardID, err := strconv.ParseUint(matches[1], 10, 32)
	if err != nil {
		return http.StatusBadRequest,
			errors.NewRESTError("grant", fmt.Sprintf("Invalid card number (%s)", matches[1])),
			fmt.Errorf("Invalid card number (%s)", matches[1])
	}

	door := matches[2]
	if strings.TrimSpace(door) == "" {
		return http.StatusBadRequest,
			errors.NewRESTError("grant", fmt.Sprintf("Invalid door (%s)", matches[1])),
			fmt.Errorf("Invalid door (%s)", matches[1])
	}

	blob, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("grant", "Error reading request"),
			err
	}

	body := struct {
		From    *types.Date `json:"start-date"`
		To      *types.Date `json:"end-date"`
		Profile int         `json:"profile"`
	}{}

	if err = json.Unmarshal(blob, &body); err != nil {
		return http.StatusBadRequest,
			errors.NewRESTError("grant", "Invalid request format"),
			err
	}

	if body.From == nil {
		return http.StatusBadRequest,
			errors.NewRESTError("grant", "Missing/invalid start date"),
			fmt.Errorf("Missing/invalid start date in request body")
	}

	if body.To == nil {
		return http.StatusBadRequest,
			errors.NewRESTError("grant", "Missing/invalid end date"),
			fmt.Errorf("Missing/invalid end date in request body")
	}

	u := ctx.Value("uhppote").(uhppote.IUHPPOTE)
	devices := ctx.Value("devices").([]uhppote.Device)

	err = api.Grant(u, devices, uint32(cardID), *body.From, *body.To, body.Profile, []string{door})
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("grant", fmt.Sprintf("%v", err)),
			err
	}

	return http.StatusOK, nil, nil
}
