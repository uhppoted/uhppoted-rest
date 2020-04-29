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
	"strings"
)

type permissions []struct {
	CardNumber uint32      `json:"card-number"`
	From       *types.Date `json:"start-date"`
	To         *types.Date `json:"end-date"`
	Doors      []string    `json:"doors"`
}

func Load(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	blob, err := ioutil.ReadAll(r.Body)
	if err != nil {
		warn(ctx, "load", err)
		http.Error(w, "Error reading request", http.StatusInternalServerError)
		return
	}

	body := permissions{}

	if err = json.Unmarshal(blob, &body); err != nil {
		warn(ctx, "load", err)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	table, err := body.parse()
	if err != nil {
		warn(ctx, "load", err)
		http.Error(w, "Error parsing request", http.StatusInternalServerError)
		return
	}

	u := ctx.Value("uhppote").(*uhppote.UHPPOTE)
	devices := ctx.Value("devices").([]*uhppote.Device)

	acl, err := api.ParseTable(*table, devices)
	if err != nil {
		warn(ctx, "ACL::load", err)
		http.Error(w, "Error processing access control list", http.StatusInternalServerError)
		return
	}

	_, err = api.PutACL(u, *acl)
	if err != nil {
		warn(ctx, "ACL::load", err)
		http.Error(w, "Error loading access control list", http.StatusInternalServerError)
		return
	}
}

func (p *permissions) parse() (*api.Table, error) {
	header := []string{"Card Number", "From", "To"}
	records := [][]string{}
	index := map[string]int{}

	for _, r := range *p {
		if r.From == nil {
			return nil, fmt.Errorf("Card %v: missing 'start-date'", r.CardNumber)
		}

		if r.To == nil {
			return nil, fmt.Errorf("Card %v: missing 'end-date'", r.CardNumber)
		}

		for _, door := range r.Doors {
			d := clean(door)
			if _, ok := index[d]; !ok {
				index[d] = len(index) + len(header)
			}
		}
	}

	for h, _ := range index {
		header = append(header, h)
	}

	for _, r := range *p {
		record := make([]string, len(header))
		record[0] = fmt.Sprintf("%v", r.CardNumber)
		record[1] = fmt.Sprintf("%s", r.From)
		record[2] = fmt.Sprintf("%s", r.To)
		for i := 3; i < len(record); i++ {
			record[i] = "N"
		}

		for _, door := range r.Doors {
			d := clean(door)
			if ix, ok := index[d]; !ok {
				return nil, fmt.Errorf("Card %v: unindexed door '%s'", r.CardNumber, door)
			} else {
				record[ix] = "Y"
			}
		}

		records = append(records, record)
	}

	table := api.Table{
		Header:  header,
		Records: records,
	}

	return &table, nil
}

func clean(s string) string {
	return strings.ReplaceAll(strings.ToLower(s), " ", "")
}
