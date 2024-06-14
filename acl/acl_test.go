package acl

import (
	"reflect"
	"testing"

	"github.com/uhppoted/uhppote-core/types"
	lib "github.com/uhppoted/uhppoted-lib/acl"
)

func TestPermissionsToTable(t *testing.T) {
	p := []permission{
		permission{
			CardNumber: 10058400,
			From:       types.MustParseDate("2023-01-01"),
			To:         types.MustParseDate("2023-12-31"),
			Doors: []door{
				door{Door: "Entrance"},
				door{Door: "Upstairs"},
				door{Door: "Downstairs"},
			},
		},
		permission{
			CardNumber: 10058401,
			From:       types.MustParseDate("2023-02-01"),
			To:         types.MustParseDate("2023-11-30"),
			Doors: []door{
				door{Door: "D1"},
				door{Door: "Upstairs"},
				door{Door: "D4"},
			},
		},
		permission{
			CardNumber: 10058402,
			From:       types.Date{},
			To:         types.MustParseDate("2023-12-31"),
			Doors: []door{
				door{Door: "Entrance"},
				door{Door: "Upstairs"},
				door{Door: "Downstairs"},
			},
		},
		permission{
			CardNumber: 10058403,
			From:       types.MustParseDate("2023-01-01"),
			To:         types.Date{},
			Doors: []door{
				door{Door: "Entrance"},
				door{Door: "Upstairs"},
				door{Door: "Downstairs"},
			},
		},
	}

	expected := lib.Table{
		Header: []string{"Card Number", "From", "To", "Entrance", "Upstairs", "Downstairs", "D1", "D4"},
		Records: [][]string{
			[]string{"10058400", "2023-01-01", "2023-12-31", "Y", "Y", "Y", "N", "N"},
			[]string{"10058401", "2023-02-01", "2023-11-30", "N", "Y", "N", "Y", "Y"},
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
