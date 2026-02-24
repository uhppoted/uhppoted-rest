package device

import (
	"context"
	"fmt"
	"maps"
	"net/http"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppoted-lib/uhppoted"

	"github.com/uhppoted/uhppoted-rest/errors"
	"github.com/uhppoted/uhppoted-rest/lib"
)

type Status struct {
	DoorState      map[uint8]bool `json:"door-states"`
	DoorButton     map[uint8]bool `json:"door-buttons"`
	SystemError    uint8          `json:"system-error"`
	SystemDateTime types.DateTime `json:"system-datetime"`
	SequenceId     uint32         `json:"sequence-id"`
	SpecialInfo    uint8          `json:"special-info"`
	RelayState     uint8          `json:"relay-state"`
	InputState     uint8          `json:"input-state"`
	Event          any            `json:"event,omitempty"`
}

func GetStatus(impl uhppoted.IUHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, any, error) {
	deviceID := ctx.Value(lib.DeviceID).(uint32)

	reply, err := impl.GetStatus(deviceID)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("get-status", fmt.Sprintf("Error retrieving device status for %v", deviceID)),
			err
	} else if reply == nil {
		return http.StatusOK, nil, nil
	}

	response := struct {
		DeviceID uint32 `json:"device-id"`
		Status   Status `json:"status"`
	}{
		DeviceID: deviceID,
		Status: Status{
			DoorState:      map[uint8]bool{},
			DoorButton:     map[uint8]bool{},
			SystemError:    reply.SystemError,
			SystemDateTime: reply.SystemDateTime,
			SequenceId:     reply.SequenceId,
			SpecialInfo:    reply.SpecialInfo,
			RelayState:     reply.RelayState,
			InputState:     reply.InputState,
			Event:          nil,
		},
	}

	maps.Copy(response.Status.DoorState, reply.DoorState)

	maps.Copy(response.Status.DoorButton, reply.DoorButton)

	if !reply.Event.IsZero() {
		event := Transmogrify(reply.Event)
		response.Status.Event = event
	}

	return http.StatusOK, &response, nil
}
