package acl

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/uhppoted/uhppote-core/uhppote"
	api "github.com/uhppoted/uhppoted-api/acl"
	"github.com/uhppoted/uhppoted-api/uhppoted"
	"github.com/uhppoted/uhppoted-rest/errors"
)

func PutACL(impl *uhppoted.UHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	blob, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("put-acl", "Error reading request"),
			err
	}

	body := struct {
		ACL []permission `json:"acl"`
	}{}

	if err = json.Unmarshal(blob, &body); err != nil {
		return http.StatusBadRequest,
			errors.NewRESTError("put-acl", "Invalid request format"),
			err
	}

	table, err := PermissionsToTable(body.ACL)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("put-acl", "Error parsing request"),
			err
	}

	u := ctx.Value("uhppote").(*uhppote.UHPPOTE)
	devices := ctx.Value("devices").([]*uhppote.Device)

	acl, _, err := api.ParseTable(table, devices, true)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("put-acl", "Error processing access control list"),
			err
	}

	rpt, errs := api.PutACL(u, *acl, false)
	if len(errs) > 0 {
		return http.StatusInternalServerError,
			errors.NewRESTError("put-table", "Error storing access control list"),
			fmt.Errorf("Error(s) storing access control list (%v)", errs)
	}

	report := []struct {
		DeviceID  uint32 `json:"device-id"`
		Unchanged int    `json:"unchanged"`
		Updated   int    `json:"updated"`
		Added     int    `json:"added"`
		Deleted   int    `json:"deleted"`
		Failed    int    `json:"failed"`
		Errors    int    `json:"errors"`
	}{}

	for k, v := range rpt {
		report = append(report, struct {
			DeviceID  uint32 `json:"device-id"`
			Unchanged int    `json:"unchanged"`
			Updated   int    `json:"updated"`
			Added     int    `json:"added"`
			Deleted   int    `json:"deleted"`
			Failed    int    `json:"failed"`
			Errors    int    `json:"errors"`
		}{
			DeviceID:  k,
			Unchanged: len(v.Unchanged),
			Updated:   len(v.Updated),
			Added:     len(v.Added),
			Deleted:   len(v.Deleted),
			Failed:    len(v.Failed),
			Errors:    len(v.Errors),
		})
	}

	return http.StatusOK, &struct {
		Report []struct {
			DeviceID  uint32 `json:"device-id"`
			Unchanged int    `json:"unchanged"`
			Updated   int    `json:"updated"`
			Added     int    `json:"added"`
			Deleted   int    `json:"deleted"`
			Failed    int    `json:"failed"`
			Errors    int    `json:"errors"`
		} `json:"report"`
	}{
		Report: report,
	}, nil
}
