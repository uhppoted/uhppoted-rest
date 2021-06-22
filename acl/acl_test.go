package acl

import (
	"github.com/uhppoted/uhppote-core/types"
	api "github.com/uhppoted/uhppoted-lib/acl"
	"reflect"
	"testing"
	"time"
)

var date = func(s string) *types.Date {
	d, _ := time.ParseInLocation("2006-01-02", s, time.Local)
	p := types.Date(d)
	return &p
}

func TestPermissionsToTable(t *testing.T) {
	p := []permission{
		permission{
			CardNumber: 192837465,
			From:       date("2020-01-01"),
			To:         date("2020-12-31"),
			Doors: []door{
				door{Door: "Entrance"},
				door{Door: "Upstairs"},
				door{Door: "Downstairs"},
			},
		},
		permission{
			CardNumber: 729364646,
			From:       date("2020-02-01"),
			To:         date("2020-11-30"),
			Doors: []door{
				door{Door: "D1"},
				door{Door: "Upstairs"},
				door{Door: "D4"},
			},
		},
	}

	expected := api.Table{
		Header: []string{"Card Number", "From", "To", "Entrance", "Upstairs", "Downstairs", "D1", "D4"},
		Records: [][]string{
			[]string{"192837465", "2020-01-01", "2020-12-31", "Y", "Y", "Y", "N", "N"},
			[]string{"729364646", "2020-02-01", "2020-11-30", "N", "Y", "N", "Y", "Y"},
		},
	}

	table, err := PermissionsToTable(p)
	if err != nil {
		t.Fatalf("Unexpected error converting permissions list to ACL table: %v", err)
	}

	if table == nil {
		t.Fatalf("PermissionsToTable returned %v", table)
	}

	if !reflect.DeepEqual(*table, expected) {
		t.Errorf("Incorrect ACL table\n   expected: %+v\n   got:      %+v", expected, *table)
	}
}
