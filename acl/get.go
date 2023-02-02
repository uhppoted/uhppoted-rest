package acl

import (
	"context"
	"fmt"
	"net/http"

	"github.com/uhppoted/uhppote-core/uhppote"
	api "github.com/uhppoted/uhppoted-lib/acl"
	"github.com/uhppoted/uhppoted-lib/uhppoted"

	"github.com/uhppoted/uhppoted-rest/errors"
	"github.com/uhppoted/uhppoted-rest/lib"
)

func GetACL(impl uhppoted.IUHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	u := ctx.Value(lib.Uhppote).(uhppote.IUHPPOTE)
	devices := ctx.Value(lib.Devices).([]uhppote.Device)

	acl, errs := api.GetACL(u, devices)
	if len(errs) > 0 {
		return http.StatusInternalServerError,
			errors.NewRESTError("get-acl", "Error retrieving access control list"),
			fmt.Errorf("Error(s) retrieving access control list (%v)", errs)
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
