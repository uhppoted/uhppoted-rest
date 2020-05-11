package acl

import (
	"context"
	"encoding/json"
	"github.com/uhppoted/uhppote-core/uhppote"
	api "github.com/uhppoted/uhppoted-api/acl"
	"github.com/uhppoted/uhppoted-api/uhppoted"
	"github.com/uhppoted/uhppoted-rest/errors"
	"io/ioutil"
	"net/http"
)

func PutACL(impl *uhppoted.UHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (interface{}, *errors.IError) {
	blob, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.ErrorX(err, "put-acl", http.StatusInternalServerError, "Error reading request")
	}

	body := struct {
		ACL []permission `json:"acl"`
	}{}

	if err = json.Unmarshal(blob, &body); err != nil {
		return nil, errors.ErrorX(err, "put-acl", http.StatusBadRequest, "Invalid request format")
	}

	table, err := PermissionsToTable(body.ACL)
	if err != nil {
		return nil, errors.ErrorX(err, "put-acl", http.StatusInternalServerError, "Error parsing request")
	}

	u := ctx.Value("uhppote").(*uhppote.UHPPOTE)
	devices := ctx.Value("devices").([]*uhppote.Device)

	acl, err := api.ParseTable(*table, devices)
	if err != nil {
		return nil, errors.ErrorX(err, "put-acl", http.StatusInternalServerError, "Error processing access control list")
	}

	rpt, err := api.PutACL(u, *acl)
	if err != nil {
		return nil, errors.ErrorX(err, "put-table", http.StatusInternalServerError, "Error storing access control list")
	}

	report := []struct {
		DeviceID  uint32 `json:"device-id"`
		Unchanged int    `json:"unchanged"`
		Updated   int    `json:"updated"`
		Added     int    `json:"added"`
		Deleted   int    `json:"deleted"`
		Failed    int    `json:"failed"`
	}{}

	for k, v := range rpt {
		report = append(report, struct {
			DeviceID  uint32 `json:"device-id"`
			Unchanged int    `json:"unchanged"`
			Updated   int    `json:"updated"`
			Added     int    `json:"added"`
			Deleted   int    `json:"deleted"`
			Failed    int    `json:"failed"`
		}{
			DeviceID:  k,
			Unchanged: v.Unchanged,
			Updated:   v.Updated,
			Added:     v.Added,
			Deleted:   v.Deleted,
			Failed:    v.Failed,
		})
	}

	return &struct {
		Report []struct {
			DeviceID  uint32 `json:"device-id"`
			Unchanged int    `json:"unchanged"`
			Updated   int    `json:"updated"`
			Added     int    `json:"added"`
			Deleted   int    `json:"deleted"`
			Failed    int    `json:"failed"`
		} `json:"report"`
	}{
		Report: report,
	}, nil
}
