package acl

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/uhppoted/uhppote-core/types"
	api "github.com/uhppoted/uhppoted-lib/acl"

	"github.com/uhppoted/uhppoted-rest/log"
)

type permission struct {
	CardNumber uint32     `json:"card-number"`
	From       types.Date `json:"start-date"`
	To         types.Date `json:"end-date"`
	Doors      []door     `json:"doors,omitempty"`
}

type door struct {
	Door    string `json:"door,omitempty"`
	Profile int    `json:"profile,omitempty"`
}

func PermissionsToTable(p []permission) (*api.Table, error) {
	header := []string{"Card Number", "From", "To"}
	records := [][]string{}
	index := map[string]int{}

	for _, r := range p {
		for _, door := range r.Doors {
			d := clean(door.Door)
			if _, ok := index[d]; !ok {
				index[d] = 3 + len(index)
				header = append(header, door.Door)
			}
		}
	}

	for _, r := range p {
		// ... discard cards with invalid start/end dates
		if r.From.IsZero() {
			log.Warnf("ACL", "card %v discarded (missing start date)", r.CardNumber)
			continue
		}

		if r.To.IsZero() {
			log.Warnf("ACL", "card %v discarded (missing end date)", r.CardNumber)
			continue
		}

		// ... create ACL record
		record := make([]string, len(header))
		record[0] = fmt.Sprintf("%v", r.CardNumber)
		record[1] = fmt.Sprintf("%v", r.From)
		record[2] = fmt.Sprintf("%v", r.To)
		for i := 3; i < len(record); i++ {
			record[i] = "N"
		}

		for _, door := range r.Doors {
			d := clean(door.Door)
			if ix, ok := index[d]; ok {
				if door.Profile >= 2 && door.Profile <= 254 {
					record[ix] = strconv.Itoa(door.Profile)
				} else {
					record[ix] = "Y"
				}

				continue
			}

			return nil, fmt.Errorf("card %v: unindexed door '%v'", r.CardNumber, door.Door)
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
		return nil, fmt.Errorf("invalid ACL table - missing 'card number'")
	}

	if index.from == 0 {
		return nil, fmt.Errorf("invalid ACL table - missing 'from' date")
	}

	if index.to == 0 {
		return nil, fmt.Errorf("invalid ACL table - missing 'to' date")
	}

	permissions := []permission{}
	for _, row := range table.Records {
		cardID, err := strconv.ParseUint(row[index.cardnumber-1], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid ACL table - invalid 'card number':%s (%w)", row[index.cardnumber-1], err)
		}

		from, err := types.DateFromString(row[index.from-1])
		if err != nil {
			log.Warnf("ACL", "card %v  invalid from date '%s' (%v)", cardID, row[index.from-1], err)
		}

		to, err := types.DateFromString(row[index.to-1])
		if err != nil {
			log.Warnf("ACL", "card %v  invalid to date   '%s' (%v)", cardID, row[index.to-1], err)
		}

		doors := []door{}
		for k, v := range index.doors {
			switch {
			case row[k-1] == "Y":
				doors = append(doors, door{
					Door: v,
				})

			case regexp.MustCompile("[0-9]+").MatchString(row[k-1]):
				profile, _ := strconv.Atoi(row[k-1])
				doors = append(doors, door{
					Door:    fmt.Sprintf("%v:%v", v, profile),
					Profile: profile,
				})
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
