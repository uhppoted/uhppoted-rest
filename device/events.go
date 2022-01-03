package device

import (
	"context"
	"fmt"
	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppoted-lib/uhppoted"
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

func GetEvents(impl uhppoted.IUHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	deviceID := ctx.Value("device-id").(uint32)

	rq := uhppoted.GetEventsRequest{
		DeviceID: uhppoted.DeviceID(deviceID),
	}

	response, err := impl.GetEvents(rq)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("get-events", fmt.Sprintf("Error retrieving event indices from %v", deviceID)),
			err
	}

	if response == nil {
		return http.StatusNotFound,
			errors.NewRESTError("get-events", fmt.Sprintf("No events on device %v", deviceID)),
			fmt.Errorf("No events on device %v", deviceID)
	}

	return http.StatusOK, &struct {
		Events struct {
			First   uint32 `json:"first,omitempty"`
			Last    uint32 `json:"last,omitempty"`
			Current uint32 `json:"current,omitempty"`
		} `json:"events"`
	}{
		Events: response.Events,
	}, nil
}

func GetEvent(impl uhppoted.IUHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	var deviceID = ctx.Value("device-id").(uint32)
	var index uint32

	// .. get event indices
	events, err := impl.GetEvents(uhppoted.GetEventsRequest{DeviceID: uhppoted.DeviceID(deviceID)})
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("get-event", fmt.Sprintf("Error retrieving event indices from %v", deviceID)),
			err
	} else if events == nil {
		return http.StatusNotFound,
			errors.NewRESTError("get-events", fmt.Sprintf("No events on device %v", deviceID)),
			fmt.Errorf("No events on device %v", deviceID)
	}

	switch v := ctx.Value("event-index").(type) {
	case uint32:
		index = v

	case string:
		switch v {
		case "first":
			index = 0

		case "last":
			index = 0xffffffff

		case "current":
			index = events.Events.Current

		case "next":
			index = events.Events.Current + 1

		default:
			return http.StatusBadRequest,
				errors.NewRESTError("get-event", "Missing/invalid event index"),
				fmt.Errorf("Missing/invalid event index")
		}

	default:
		return http.StatusBadRequest,
			errors.NewRESTError("get-event", "Missing/invalid event index"),
			fmt.Errorf("Missing/invalid event index")
	}

	rq := uhppoted.GetEventRequest{
		DeviceID: uhppoted.DeviceID(deviceID),
		Index:    index,
	}

	response, err := impl.GetEvent(rq)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("get-event", fmt.Sprintf("%v", err)),
			err
	} else if response == nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("get-event", fmt.Sprintf("Error retrieving event %v from %v", index, deviceID)),
			fmt.Errorf("No response returned to request for event %v from device %v", index, deviceID)
	}

	return http.StatusOK, struct {
		Event interface{} `json:"event"`
	}{
		Event: response.Event,
	}, nil
}
