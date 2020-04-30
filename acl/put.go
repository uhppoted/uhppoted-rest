package acl

import (
	"context"
	"encoding/json"
	"github.com/uhppoted/uhppote-core/uhppote"
	api "github.com/uhppoted/uhppoted-api/acl"
	"io/ioutil"
	"net/http"
)

func PutACL(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	blob, err := ioutil.ReadAll(r.Body)
	if err != nil {
		warn(ctx, "put-acl", err)
		http.Error(w, "Error reading request", http.StatusInternalServerError)
		return
	}

	body := permissions{}

	if err = json.Unmarshal(blob, &body); err != nil {
		warn(ctx, "put-acl", err)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	table, err := body.toTable()
	if err != nil {
		warn(ctx, "put-acl", err)
		http.Error(w, "Error parsing request", http.StatusInternalServerError)
		return
	}

	u := ctx.Value("uhppote").(*uhppote.UHPPOTE)
	devices := ctx.Value("devices").([]*uhppote.Device)

	acl, err := api.ParseTable(*table, devices)
	if err != nil {
		warn(ctx, "put-acl", err)
		http.Error(w, "Error processing access control list", http.StatusInternalServerError)
		return
	}

	rpt, err := api.PutACL(u, *acl)
	if err != nil {
		warn(ctx, "put-acl", err)
		http.Error(w, "Error put-acling access control list", http.StatusInternalServerError)
		return
	}

	response := []struct {
		DeviceID  uint32 `json:"device-id"`
		Unchanged int    `json:"unchanged"`
		Updated   int    `json:"updated"`
		Added     int    `json:"added"`
		Deleted   int    `json:"deleted"`
		Failed    int    `json:"failed"`
	}{}

	for k, v := range rpt {
		response = append(response, struct {
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

	reply(ctx, w, response)
}
