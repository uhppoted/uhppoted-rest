package acl

import (
	"reflect"
	"testing"

	"codeberg.org/uhppoted/uhppoted-core/types"
	lib "codeberg.org/uhppoted/uhppoted-lib/acl"
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

	table, err := PermissionsToTable(p, false, false)
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

func TestPermissionsToTableWithPIN(t *testing.T) {
	p := []permission{
		permission{
			CardNumber: 10058400,
			PIN:        7530,
			From:       types.MustParseDate("2026-01-01"),
			To:         types.MustParseDate("2026-12-31"),
			Doors: []door{
				door{Door: "Entrance"},
				door{Door: "Upstairs"},
				door{Door: "Downstairs"},
			},
		},
		permission{
			CardNumber: 10058401,
			From:       types.MustParseDate("2026-02-01"),
			To:         types.MustParseDate("2026-11-30"),
			Doors: []door{
				door{Door: "D1"},
				door{Door: "Upstairs"},
				door{Door: "D4"},
			},
		},
		permission{
			CardNumber: 10058402,
			From:       types.Date{},
			To:         types.MustParseDate("2026-12-31"),
			Doors: []door{
				door{Door: "Entrance"},
				door{Door: "Upstairs"},
				door{Door: "Downstairs"},
			},
		},
		permission{
			CardNumber: 10058403,
			From:       types.MustParseDate("2026-01-01"),
			To:         types.Date{},
			Doors: []door{
				door{Door: "Entrance"},
				door{Door: "Upstairs"},
				door{Door: "Downstairs"},
			},
		},
	}

	expected := lib.Table{
		Header: []string{"Card Number", "PIN", "From", "To", "Entrance", "Upstairs", "Downstairs", "D1", "D4"},
		Records: [][]string{
			[]string{"10058400", "7530", "2026-01-01", "2026-12-31", "Y", "Y", "Y", "N", "N"},
			[]string{"10058401", "", "2026-02-01", "2026-11-30", "N", "Y", "N", "Y", "Y"},
		},
	}

	table, err := PermissionsToTable(p, true, false)
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

func TestPermissionsToTableWithFirstCard(t *testing.T) {
	p := []permission{
		permission{
			CardNumber: 10058400,
			PIN:        7530,
			From:       types.MustParseDate("2026-01-01"),
			To:         types.MustParseDate("2026-12-31"),
			Doors: []door{
				door{Door: "Entrance"},
				door{Door: "Upstairs"},
				door{Door: "Downstairs"},
			},
			FirstCard: true,
		},
		permission{
			CardNumber: 10058401,
			From:       types.MustParseDate("2026-02-01"),
			To:         types.MustParseDate("2026-11-30"),
			Doors: []door{
				door{Door: "D1"},
				door{Door: "Upstairs"},
				door{Door: "D4"},
			},
			FirstCard: false,
		},
		permission{
			CardNumber: 10058402,
			From:       types.Date{},
			To:         types.MustParseDate("2026-12-31"),
			Doors: []door{
				door{Door: "Entrance"},
				door{Door: "Upstairs"},
				door{Door: "Downstairs"},
			},
		},
		permission{
			CardNumber: 10058403,
			From:       types.MustParseDate("2026-01-01"),
			To:         types.Date{},
			Doors: []door{
				door{Door: "Entrance"},
				door{Door: "Upstairs"},
				door{Door: "Downstairs"},
			},
		},
	}

	expected := lib.Table{
		Header: []string{"Card Number", "From", "To", "Entrance", "Upstairs", "Downstairs", "D1", "D4", "FirstCard"},
		Records: [][]string{
			[]string{"10058400", "2026-01-01", "2026-12-31", "Y", "Y", "Y", "N", "N", "Y"},
			[]string{"10058401", "2026-02-01", "2026-11-30", "N", "Y", "N", "Y", "Y", "N"},
		},
	}

	table, err := PermissionsToTable(p, false, true)
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

func TestPermissionsFromTable(t *testing.T) {
	table := lib.Table{
		Header: []string{"Card Number", "From", "To", "Entrance", "Upstairs", "Downstairs", "D1", "D4"},
		Records: [][]string{
			[]string{"10058400", "2023-01-01", "2023-12-31", "Y", "Y", "Y", "N", "N"},
			[]string{"10058401", "2023-02-01", "2023-11-30", "N", "Y", "N", "Y", "Y"},
			[]string{"10058402", "", "2023-11-30", "Y", "Y", "Y", "N", "N"},
			[]string{"10058403", "2023-02-01", "", "Y", "Y", "Y", "N", "N"},
		},
	}

	expected := []permission{
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
				door{Door: "Upstairs"},
				door{Door: "D1"},
				door{Door: "D4"},
			},
		},
		permission{
			CardNumber: 10058402,
			From:       types.Date{},
			To:         types.MustParseDate("2023-11-30"),
			Doors: []door{
				door{Door: "Entrance"},
				door{Door: "Upstairs"},
				door{Door: "Downstairs"},
			},
		},
		permission{
			CardNumber: 10058403,
			From:       types.MustParseDate("2023-02-01"),
			To:         types.Date{},
			Doors: []door{
				door{Door: "Entrance"},
				door{Door: "Upstairs"},
				door{Door: "Downstairs"},
			},
		},
	}

	permissions, err := PermissionsFromTable(table, false, false)
	if err != nil {
		t.Fatalf("Unexpected error converting permissions list to ACL table: %v", err)
	}

	if permissions == nil {
		t.Fatalf("PermissionsFromTable returned %v", permissions)
	}

	if !reflect.DeepEqual(permissions, expected) {
		t.Errorf("Incorrect ACL permissions\n   expected: %+v\n   got:      %+v", expected, permissions)
	}
}

func TestPermissionsFromTableWithPIN(t *testing.T) {
	table := lib.Table{
		Header: []string{"Card Number", "PIN", "From", "To", "Entrance", "Upstairs", "Downstairs", "D1", "D4"},
		Records: [][]string{
			[]string{"10058400", "7531", "2023-01-01", "2023-12-31", "Y", "Y", "Y", "N", "N"},
			[]string{"10058401", "", "2023-02-01", "2023-11-30", "N", "Y", "N", "Y", "Y"},
			[]string{"10058402", "", "", "2023-11-30", "Y", "Y", "Y", "N", "N"},
			[]string{"10058403", "", "2023-02-01", "", "Y", "Y", "Y", "N", "N"},
		},
	}

	expected := []permission{
		permission{
			CardNumber: 10058400,
			PIN:        7531,
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
			PIN:        0,
			From:       types.MustParseDate("2023-02-01"),
			To:         types.MustParseDate("2023-11-30"),
			Doors: []door{
				door{Door: "Upstairs"},
				door{Door: "D1"},
				door{Door: "D4"},
			},
		},
		permission{
			CardNumber: 10058402,
			PIN:        0,
			From:       types.Date{},
			To:         types.MustParseDate("2023-11-30"),
			Doors: []door{
				door{Door: "Entrance"},
				door{Door: "Upstairs"},
				door{Door: "Downstairs"},
			},
		},
		permission{
			CardNumber: 10058403,
			PIN:        0,
			From:       types.MustParseDate("2023-02-01"),
			To:         types.Date{},
			Doors: []door{
				door{Door: "Entrance"},
				door{Door: "Upstairs"},
				door{Door: "Downstairs"},
			},
		},
	}

	permissions, err := PermissionsFromTable(table, true, false)
	if err != nil {
		t.Fatalf("Unexpected error converting permissions list to ACL table: %v", err)
	}

	if permissions == nil {
		t.Fatalf("PermissionsFromTable returned %v", permissions)
	}

	if !reflect.DeepEqual(permissions, expected) {
		t.Errorf("Incorrect ACL permissions\n   expected: %+v\n   got:      %+v", expected, permissions)
	}
}

func TestPermissionsFromTableWithFirstCard(t *testing.T) {
	table := lib.Table{
		Header: []string{"Card Number", "From", "To", "Entrance", "Upstairs", "Downstairs", "D1", "D4", "FirstCard"},
		Records: [][]string{
			[]string{"10058400", "2023-01-01", "2023-12-31", "Y", "Y", "Y", "N", "N", "Y"},
			[]string{"10058401", "2023-02-01", "2023-11-30", "N", "Y", "N", "Y", "Y", "N"},
			[]string{"10058402", "", "2023-11-30", "Y", "Y", "Y", "N", "N", "N"},
			[]string{"10058403", "2023-02-01", "", "Y", "Y", "Y", "N", "N", "N"},
		},
	}

	expected := []permission{
		permission{
			CardNumber: 10058400,
			From:       types.MustParseDate("2023-01-01"),
			To:         types.MustParseDate("2023-12-31"),
			Doors: []door{
				door{Door: "Entrance"},
				door{Door: "Upstairs"},
				door{Door: "Downstairs"},
			},
			FirstCard: true,
		},
		permission{
			CardNumber: 10058401,
			PIN:        0,
			From:       types.MustParseDate("2023-02-01"),
			To:         types.MustParseDate("2023-11-30"),
			Doors: []door{
				door{Door: "Upstairs"},
				door{Door: "D1"},
				door{Door: "D4"},
			},
		},
		permission{
			CardNumber: 10058402,
			PIN:        0,
			From:       types.Date{},
			To:         types.MustParseDate("2023-11-30"),
			Doors: []door{
				door{Door: "Entrance"},
				door{Door: "Upstairs"},
				door{Door: "Downstairs"},
			},
		},
		permission{
			CardNumber: 10058403,
			PIN:        0,
			From:       types.MustParseDate("2023-02-01"),
			To:         types.Date{},
			Doors: []door{
				door{Door: "Entrance"},
				door{Door: "Upstairs"},
				door{Door: "Downstairs"},
			},
		},
	}

	permissions, err := PermissionsFromTable(table, false, true)
	if err != nil {
		t.Fatalf("Unexpected error converting permissions list to ACL table: %v", err)
	}

	if permissions == nil {
		t.Fatalf("PermissionsFromTable returned %v", permissions)
	}

	if !reflect.DeepEqual(permissions, expected) {
		t.Errorf("Incorrect ACL permissions\n   expected: %+v\n   got:      %+v", expected, permissions)
	}
}
