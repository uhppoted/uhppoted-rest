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

func Grant(impl *uhppoted.UHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (interface{}, *errors.IError) {
	url := r.URL.Path

	matches := regexp.MustCompile("^/uhppote/acl/card/([0-9]+)/door/(\\S.*)$").FindStringSubmatch(url)
	if matches == nil || len(matches) < 3 {
		return nil, errors.ErrorX(fmt.Errorf("Missing card number/door"), "grant", http.StatusBadRequest, "Missing card number/door")
	}

	cardID, err := strconv.ParseUint(matches[1], 10, 32)
	if err != nil {
		return nil, errors.ErrorX(fmt.Errorf("Invalid card number (%s)", matches[1]), "grant", http.StatusBadRequest, "Invalid card number")
	}

	door := matches[2]
	if strings.TrimSpace(door) == "" {
		return nil, errors.ErrorX(fmt.Errorf("Invalid door (%s)", matches[1]), "grant", http.StatusBadRequest, "Invalid door")
	}

	blob, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.ErrorX(err, "grant", http.StatusInternalServerError, "Error reading request")
	}

	body := struct {
		From *types.Date `json:"start-date"`
		To   *types.Date `json:"end-date"`
	}{}

	if err = json.Unmarshal(blob, &body); err != nil {
		return nil, errors.ErrorX(err, "grant", http.StatusBadRequest, "Invalid request format")
	}

	if body.From == nil {
		return nil, errors.ErrorX(fmt.Errorf("Missing/invalid start date in request body"), "grant", http.StatusBadRequest, "Missing/invalid start date")
	}

	if body.To == nil {
		return nil, errors.ErrorX(fmt.Errorf("Missing/invalid end date in request body"), "grant", http.StatusBadRequest, "Missing/invalid end date")
	}

	u := ctx.Value("uhppote").(*uhppote.UHPPOTE)
	devices := ctx.Value("devices").([]*uhppote.Device)

	err = api.Grant(u, devices, uint32(cardID), *body.From, *body.To, []string{door})
	if err != nil {
		return nil, errors.ErrorX(err, "grant", http.StatusInternalServerError, "Error granting card access permissions")
	}

	return nil, nil
}
