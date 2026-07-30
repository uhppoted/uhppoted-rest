package acl

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"codeberg.org/uhppoted/uhppoted-core/types"
	api "codeberg.org/uhppoted/uhppoted-lib/acl"

	"codeberg.org/uhppoted/uhppoted-rest/log"
)

type permission struct {
	CardNumber uint32     `json:"card-number"`
	PIN        uint32     `json:"PIN,omitempty"`
	From       types.Date `json:"start-date"`
	To         types.Date `json:"end-date"`
	Doors      []door     `json:"doors,omitempty"`
	FirstCard  bool       `json:"first-card,omitempty"`
}

type door struct {
	Door    string `json:"door,omitempty"`
	Profile int    `json:"profile,omitempty"`
}

func PermissionsToTable(p []permission, withPINs bool, withFirstCard bool) (*api.Table, error) {
	header := []string{"Card Number", "From", "To"}
	records := [][]string{}
	index := map[string]int{}

	if withPINs {
		header = []string{"Card Number", "PIN", "From", "To"}
	}

	for _, r := range p {
		for _, door := range r.Doors {
			d := clean(door.Door)
			if _, ok := index[d]; !ok {
				if withPINs {
					index[d] = 4 + len(index)
				} else {
					index[d] = 3 + len(index)
				}

				header = append(header, door.Door)
			}
		}
	}

	if withFirstCard {
		header = append(header, "FirstCard")
		if withPINs {
			index["firstcard"] = 4 + len(index)
		} else {
			index["firstcard"] = 3 + len(index)
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
		ix := 0

		record := make([]string, len(header))

		record[ix] = fmt.Sprintf("%v", r.CardNumber)
		ix++

		if withPINs {
			if r.PIN > 0 {
				record[ix] = fmt.Sprintf("%v", r.PIN)
			} else {
				record[ix] = ""
			}
			ix++
		}

		record[ix] = fmt.Sprintf("%v", r.From)
		ix++

		record[ix] = fmt.Sprintf("%v", r.To)
		ix++

		for i := ix; i < len(record); i++ {
			record[ix] = "N"
			ix++
		}

		for _, door := range r.Doors {
			d := clean(door.Door)
			if k, ok := index[d]; ok {
				if door.Profile >= 2 && door.Profile <= 254 {
					record[k] = strconv.Itoa(door.Profile)
				} else {
					record[k] = "Y"
				}

				ix++
				continue
			}

			return nil, fmt.Errorf("card %v: unindexed door '%v'", r.CardNumber, door.Door)
		}

		if withFirstCard {
			if k, ok := index["firstcard"]; ok {
				if r.FirstCard {
					record[k] = "Y"
				} else {
					record[k] = "N"
				}
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

func PermissionsFromTable(table api.Table, withPINs bool, withFirstCard bool) ([]permission, error) {
	index := struct {
		cardnumber int
		PIN        int
		from       int
		to         int
		doors      map[int]string
		firstcard  int
	}{
		cardnumber: 0,
		PIN:        0,
		from:       0,
		to:         0,
		doors:      map[int]string{},
		firstcard:  0,
	}

	for i, h := range table.Header {
		switch clean(h) {
		case "cardnumber":
			index.cardnumber = i + 1

		case "pin":
			if withPINs {
				index.PIN = i + 1
			}

		case "from":
			index.from = i + 1

		case "to":
			index.to = i + 1

		case "firstcard":
			if withFirstCard {
				index.firstcard = i + 1
			} else {
				index.doors[i+1] = h
			}

		default:
			index.doors[i+1] = h
		}
	}

	if index.cardnumber == 0 {
		return nil, fmt.Errorf("invalid ACL table - missing 'card number'")
	}

	if withPINs && index.PIN == 0 {
		return nil, fmt.Errorf("invalid ACL table - missing 'PIN'")
	}

	if index.from == 0 {
		return nil, fmt.Errorf("invalid ACL table - missing 'from' date")
	}

	if index.to == 0 {
		return nil, fmt.Errorf("invalid ACL table - missing 'to' date")
	}

	if withFirstCard && index.firstcard == 0 {
		return nil, fmt.Errorf("invalid ACL table - missing 'firstcard'")
	}

	permissions := []permission{}
	for _, row := range table.Records {
		cardID, err := strconv.ParseUint(row[index.cardnumber-1], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid ACL table - invalid 'card number':%s (%w)", row[index.cardnumber-1], err)
		}

		PIN := uint64(0)
		if withPINs {
			if v := clean(row[index.PIN-1]); v != "" {
				if pin, err := strconv.ParseUint(v, 10, 32); err != nil {
					return nil, fmt.Errorf("invalid ACL table - invalid 'PIN':%s (%w)", row[index.PIN-1], err)
				} else {
					PIN = pin
				}
			}
		}

		from, err := types.ParseDate(row[index.from-1])
		if err != nil {
			log.Warnf("ACL", "card %v  invalid from date '%s' (%v)", cardID, row[index.from-1], err)
		}

		to, err := types.ParseDate(row[index.to-1])
		if err != nil {
			log.Warnf("ACL", "card %v  invalid to date   '%s' (%v)", cardID, row[index.to-1], err)
		}

		firstcard := false
		if withFirstCard {
			if v := clean(row[index.firstcard-1]); v == "y" {
				firstcard = true
			}
		}

		doors := []door{}

		keys := []int{}
		for k := range index.doors {
			keys = append(keys, k)
		}

		sort.Ints(keys)

		for _, k := range keys {
			v := index.doors[k]
			switch {
			case clean(row[k-1]) == "y":
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
			PIN:        uint32(PIN),
			From:       from,
			To:         to,
			Doors:      doors,
			FirstCard:  firstcard,
		})
	}

	return permissions, nil
}

func clean(s string) string {
	return strings.ReplaceAll(strings.ToLower(s), " ", "")
}
