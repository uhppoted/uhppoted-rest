package acl

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/uhppoted/uhppote-core/types"
	api "github.com/uhppoted/uhppoted-api/acl"
)

type permission struct {
	CardNumber uint32      `json:"card-number"`
	From       *types.Date `json:"start-date"`
	To         *types.Date `json:"end-date"`
	Doors      []string    `json:"doors"`
}

func PermissionsToTable(p []permission) (*api.Table, error) {
	header := []string{"Card Number", "From", "To"}
	records := [][]string{}
	index := map[string]int{}

	for _, r := range p {
		if r.From == nil {
			return nil, fmt.Errorf("Card %v: missing 'start-date'", r.CardNumber)
		}

		if r.To == nil {
			return nil, fmt.Errorf("Card %v: missing 'end-date'", r.CardNumber)
		}

		for _, door := range r.Doors {
			d := clean(door)
			if _, ok := index[d]; !ok {
				index[d] = 3 + len(index)
				header = append(header, door)
			}
		}
	}

	for _, r := range p {
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

func PermissionsFromTable(table *api.Table) ([]permission, error) {
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
		return nil, fmt.Errorf("Invalid ACL table - missing 'card number'")
	}

	if index.from == 0 {
		return nil, fmt.Errorf("Invalid ACL table - missing 'from' date")
	}

	if index.to == 0 {
		return nil, fmt.Errorf("Invalid ACL table - missing 'to' date")
	}

	permissions := []permission{}
	for _, row := range table.Records {
		cardID, err := strconv.ParseUint(row[index.cardnumber-1], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("Invalid ACL table - invalid 'card number':%s (%w)", row[index.cardnumber-1], err)
		}

		from, err := types.DateFromString(row[index.from-1])
		if err != nil {
			return nil, fmt.Errorf("Invalid ACL table - invalid 'from' date:%s (%w)", row[index.from-1], err)
		}

		to, err := types.DateFromString(row[index.to-1])
		if err != nil {
			return nil, fmt.Errorf("Invalid ACL table - invalid 'to' date:%s (%w)", row[index.to-1], err)
		}

		doors := []string{}
		for k, v := range index.doors {
			switch {
			case row[k-1] == "Y":
				doors = append(doors, v)

			case regexp.MustCompile("[0-9]+").MatchString(row[k-1]):
				profile, _ := strconv.Atoi(row[k-1])
				doors = append(doors, fmt.Sprintf("%v:%v", v, profile))
			}
		}

		permissions = append(permissions, permission{
			CardNumber: uint32(cardID),
			From:       from,
			To:         to,
			Doors:      doors,
		})
	}

	return permissions, nil
}

func clean(s string) string {
	return strings.ReplaceAll(strings.ToLower(s), " ", "")
}
