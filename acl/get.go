package acl

import (
	"context"
	"fmt"
	"github.com/uhppoted/uhppote-core/uhppote"
	api "github.com/uhppoted/uhppoted-api/acl"
	"github.com/uhppoted/uhppoted-api/uhppoted"
	"github.com/uhppoted/uhppoted-rest/errors"
	"net/http"
)

func GetACL(impl *uhppoted.UHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (interface{}, *errors.IError) {
	u := ctx.Value("uhppote").(*uhppote.UHPPOTE)
	devices := ctx.Value("devices").([]*uhppote.Device)

	acl, err := api.GetACL(u, devices)
	if err != nil {
		return nil, errors.Errorf(fmt.Errorf("%w: Error retrieving access control list", err), 0, "get-acl", "Error retrieving access control list")
	}

	table, err := api.MakeTable(acl, devices)
	if err != nil {
		return nil, errors.Errorf(fmt.Errorf("%w: Error processing access control list", err), 0, "get-acl", "Error processing access control list")
	}

	response := permissions{}

	if err = response.fromTable(table); err != nil {
		return nil, errors.Errorf(fmt.Errorf("%w: Error processing access control table", err), 0, "get-acl", "Error processing access control table")
	}

	return &response, nil
}
