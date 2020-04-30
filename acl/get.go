package acl

import (
	"context"
	"github.com/uhppoted/uhppote-core/uhppote"
	api "github.com/uhppoted/uhppoted-api/acl"
	"net/http"
)

func GetACL(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	u := ctx.Value("uhppote").(*uhppote.UHPPOTE)
	devices := ctx.Value("devices").([]*uhppote.Device)

	acl, err := api.GetACL(u, devices)
	if err != nil {
		warn(ctx, "get-acl", err)
		http.Error(w, "Error retrieving access control list", http.StatusInternalServerError)
		return
	}

	table, err := api.MakeTable(acl, devices)
	if err != nil {
		warn(ctx, "get-acl", err)
		http.Error(w, "Error processing access control list", http.StatusInternalServerError)
		return
	}

	response := permissions{}

	if err = response.fromTable(table); err != nil {
		warn(ctx, "get-acl", err)
		http.Error(w, "Error processing access control table", http.StatusInternalServerError)
		return
	}

	reply(ctx, w, response)
}
