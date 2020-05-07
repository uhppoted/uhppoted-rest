package device

import (
	"context"
	"fmt"
	"github.com/uhppoted/uhppote-core/types"
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

func GetEvent(impl *uhppoted.UHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (interface{}, *errors.IError) {
	deviceID := ctx.Value("device-id").(uint32)
	eventID := ctx.Value("event-id").(uint32)

	rq := uhppoted.GetEventRequest{
		DeviceID: uhppoted.DeviceID(deviceID),
		EventID:  eventID,
	}

	response, err := impl.GetEvent(rq)
	if err != nil {
		return nil, errors.Errorf(err, deviceID, "get-event", fmt.Sprintf("Error retrieving event %v from %v", eventID, deviceID))
	}

	if response == nil {
		return nil, errors.Errorf(fmt.Errorf("%w: No record for event %v", uhppoted.NotFound, eventID), deviceID, "get-event", fmt.Sprintf("Error retrieving event %v from %v", eventID, deviceID))
	}

	return response.Event, nil
}
