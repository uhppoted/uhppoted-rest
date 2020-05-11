package acl

import (
	"context"
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
		return nil, errors.ErrorX(err, "get-acl", http.StatusInternalServerError, "Error retrieving access control list")
	}

	table, err := api.MakeTable(acl, devices)
	if err != nil {
		return nil, errors.ErrorX(err, "get-acl", http.StatusInternalServerError, "Error processing access control list")
	}

	permissions, err := PermissionsFromTable(table)
	if err != nil {
		return nil, errors.ErrorX(err, "get-acl", http.StatusInternalServerError, "Error processing access control table")
	}

	return &struct {
		ACL []permission `json:"acl"`
	}{
		ACL: permissions,
	}, nil
}
