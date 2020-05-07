package device

import (
	"context"
	"fmt"
	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
	"github.com/uhppoted/uhppoted-api/uhppoted"
	"github.com/uhppoted/uhppoted-rest/errors"
	"net/http"
)

type event struct {
	Index      uint32         `json:"event-id"`
	Type       uint8          `json:"event-type"`
	Granted    bool           `json:"access-granted"`
	Door       uint8          `json:"door-id"`
	DoorOpened bool           `json:"door-opened"`
	UserID     uint32         `json:"user-id"`
	Timestamp  types.DateTime `json:"timestamp"`
	Result     uint8          `json:"event-result"`
}

func GetEvents(impl *uhppoted.UHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (interface{}, *errors.IError) {
	deviceID := ctx.Value("device-id").(uint32)

	rq := uhppoted.GetEventRangeRequest{
		DeviceID: uhppoted.DeviceID(deviceID),
		Start:    nil,
		End:      nil,
	}

	response, err := impl.GetEventRange(rq)
	if err != nil {
		return nil, errors.Errorf(err, deviceID, "get-events", fmt.Sprintf("Error retrieving event indices from %v", deviceID))
	} else if response == nil {
		return nil, errors.Errorf(errors.RequestFailed, deviceID, "get-events", fmt.Sprintf("Error retrieving event indices from %v", deviceID))
	}

	return &struct {
		Events struct {
			First uint32 `json:"first"`
			Last  uint32 `json:"last"`
		} `json:"events"`
	}{
		Events: struct {
			First uint32 `json:"first"`
			Last  uint32 `json:"last"`
		}{
			First: response.Events.First,
			Last:  response.Events.Last,
		},
	}, nil
}

func GetEvent(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	deviceID := ctx.Value("device-id").(uint32)
	eventID := ctx.Value("event-id").(uint32)

	record, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetEvent(deviceID, eventID)
	if err != nil {
		warn(ctx, deviceID, "get-event", err)
		http.Error(w, "Error retrieving event", http.StatusInternalServerError)
		return
	}

	if record == nil {
		http.Error(w, "Event record does not exist", http.StatusNotFound)
		return
	}

	if record.Index != eventID {
		http.Error(w, "Event record does not exist", http.StatusNotFound)
		return
	}

	response := struct {
		Event event `json:"event"`
	}{
		Event: event{
			Index:      record.Index,
			Type:       record.Type,
			Granted:    record.Granted,
			Door:       record.Door,
			DoorOpened: record.DoorOpened,
			UserID:     record.UserID,
			Timestamp:  record.Timestamp,
			Result:     record.Result,
		},
	}

	reply(ctx, w, response)
}
