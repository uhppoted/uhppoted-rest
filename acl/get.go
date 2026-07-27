package acl

import (
	"context"
	"fmt"
	"net/http"

	"codeberg.org/uhppoted/uhppoted-core/uhppote"
	api "codeberg.org/uhppoted/uhppoted-lib/acl"
	"codeberg.org/uhppoted/uhppoted-lib/uhppoted"

	"codeberg.org/uhppoted/uhppoted-rest/errors"
	"codeberg.org/uhppoted/uhppoted-rest/lib"
)

func GetACL(impl uhppoted.IUHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, any, error) {
	u := ctx.Value(lib.Uhppote).(uhppote.IUHPPOTE)
	controllers := ctx.Value(lib.Devices).([]uhppote.Device)
	withPINs := ctx.Value(lib.WithPINs).(bool)
	withFirstCard := ctx.Value(lib.WithFirstCard).(bool)

	acl, errs := api.GetACL(u, controllers)
	if len(errs) > 0 {
		return http.StatusInternalServerError,
			errors.NewRESTError("get-acl", "Error retrieving access control list"),
			fmt.Errorf("Error(s) retrieving access control list (%v)", errs)
	}

	table, err := api.MakeTable(acl, controllers, withPINs, withFirstCard)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("get-acl", "Error processing access control list"),
			err
	} else if table == nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("get-acl", "Error processing access control list"),
			fmt.Errorf("invalid table (%v)", table)
	}

	permissions, err := PermissionsFromTable(*table, withPINs, withFirstCard)
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
