package acl

import (
	"context"
	"github.com/uhppoted/uhppote-core/uhppote"
	api "github.com/uhppoted/uhppoted-api/acl"
	"github.com/uhppoted/uhppoted-api/uhppoted"
	"github.com/uhppoted/uhppoted-rest/errors"
	"net/http"
)

func GetACL(impl *uhppoted.UHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	u := ctx.Value("uhppote").(*uhppote.UHPPOTE)
	devices := ctx.Value("devices").([]*uhppote.Device)

	acl, err := api.GetACL(u, devices)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("get-acl", "Error retrieving access control list"),
			err
	}

	table, err := api.MakeTable(acl, devices)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("get-acl", "Error processing access control list"),
			err
	}

	permissions, err := PermissionsFromTable(table)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("get-acl", "Error processing access control table"),
			err
	}

	return http.StatusOK, &struct {
		ACL []permission `json:"acl"`
	}{
		ACL: permissions,
	}, nil
}
