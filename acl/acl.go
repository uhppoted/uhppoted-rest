package acl

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"github.com/uhppoted/uhppote-core/types"
	api "github.com/uhppoted/uhppoted-api/acl"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type permission struct {
	CardNumber uint32      `json:"card-number"`
	From       *types.Date `json:"start-date"`
	To         *types.Date `json:"end-date"`
	Doors      []string    `json:"doors"`
}

type permissions []permission

func reply(ctx context.Context, w http.ResponseWriter, response interface{}) {
	b, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error generating response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if len(b) > 1024 && ctx.Value("compression") == "gzip" {
		w.Header().Set("Content-Encoding", "gzip")
		encoder := gzip.NewWriter(w)
		encoder.Write(b)
		encoder.Flush()
	} else {
		w.Write(b)
	}
}

func (p *permissions) toTable() (*api.Table, error) {
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

func (p *permissions) fromTable(table *api.Table) error {
	index := struct {
		cardnumber int
		from       int
		to         int
		doors      map[int]string
	}{
		cardnumber: 0,
		from:       0,
		to:         0,
		doors:      map[int]string{},
	}

	for i, h := range table.Header {
		switch clean(h) {
		case "cardnumber":
			index.cardnumber = i + 1
		case "from":
			index.from = i + 1
		case "to":
			index.to = i + 1
		default:
			index.doors[i+1] = h
		}
	}

	if index.cardnumber == 0 {
		return fmt.Errorf("Invalid ACL table - missing 'card number'")
	}

	if index.from == 0 {
		return fmt.Errorf("Invalid ACL table - missing 'from' date")
	}

	if index.to == 0 {
		return fmt.Errorf("Invalid ACL table - missing 'to' date")
	}

	for _, row := range table.Records {
		cardID, err := strconv.ParseUint(row[index.cardnumber-1], 10, 32)
		if err != nil {
			return fmt.Errorf("Invalid ACL table - invalid 'card number':%s (%w)", row[index.cardnumber-1], err)
		}

		from, err := types.DateFromString(row[index.from-1])
		if err != nil {
			return fmt.Errorf("Invalid ACL table - invalid 'from' date:%s (%w)", row[index.from-1], err)
		}

		to, err := types.DateFromString(row[index.to-1])
		if err != nil {
			return fmt.Errorf("Invalid ACL table - invalid 'to' date:%s (%w)", row[index.to-1], err)
		}

		doors := []string{}
		for k, v := range index.doors {
			if row[k-1] == "Y" {
				doors = append(doors, v)
			}
		}

		*p = append(*p, permission{
			CardNumber: uint32(cardID),
			From:       from,
			To:         to,
			Doors:      doors,
		})
	}

	return nil
}

func clean(s string) string {
	return strings.ReplaceAll(strings.ToLower(s), " ", "")
}

func warn(ctx context.Context, operation string, err error) {
	ctx.Value("log").(*log.Logger).Printf("WARN  %-20s %v\n", operation, err)
}
