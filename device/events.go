package device

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

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
	var deviceID uint32

	if matches := regexp.MustCompile("^/uhppote/device/([0-9]+)(?:$|/.*$)").FindStringSubmatch(r.URL.Path); matches == nil {
		return http.StatusBadRequest, nil, errors.NewRESTError("get-events", "Missing device ID")
	} else if v, err := strconv.ParseUint(matches[1], 10, 32); err != nil {
		return http.StatusBadRequest, nil, errors.NewRESTError("get-events", fmt.Sprintf("Invalid device ID (%v)", matches[1]))
	} else {
		deviceID = uint32(v)
	}

	first, last, current, err := impl.GetEventIndices(deviceID)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("get-events", fmt.Sprintf("Error retrieving event indices from %v", deviceID)),
			err
	}

	response := struct {
		Events struct {
			First   uint32 `json:"first,omitempty"`
			Last    uint32 `json:"last,omitempty"`
			Current uint32 `json:"current,omitempty"`
		} `json:"events"`
	}{}

	response.Events.First = first
	response.Events.Last = last
	response.Events.Current = current

	return http.StatusOK, &response, nil
}

func GetEvent(impl uhppoted.IUHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	var url = r.URL.Path
	var deviceID uint32
	var index string

	if matches := regexp.MustCompile("^/uhppote/device/([0-9]+)(?:$|/.*$)").FindStringSubmatch(url); matches == nil {
		return http.StatusBadRequest, nil, errors.NewRESTError("get-events", "Missing device ID")
	} else if v, err := strconv.ParseUint(matches[1], 10, 32); err != nil {
		return http.StatusBadRequest, nil, errors.NewRESTError("get-events", fmt.Sprintf("Invalid device ID(%v)", matches[1]))
	} else {
		deviceID = uint32(v)
	}

	if matches := regexp.MustCompile("^/uhppote/device/[0-9]+/events/([0-9]+|first|last|current|next)$").FindStringSubmatch(url); matches == nil {
		return http.StatusBadRequest, nil, errors.NewRESTError("get-events", "Missing event index")
	} else {
		index = matches[1]
	}

	// .. get event indices
	first, last, current, err := impl.GetEventIndices(deviceID)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("get-events", fmt.Sprintf("Error retrieving event indices from %v", deviceID)),
			err
	}

	// ... get event
	switch index {
	case "first":
		return getEvent(impl, deviceID, first)

	case "last":
		return getEvent(impl, deviceID, last)

	case "current":
		return getEvent(impl, deviceID, current)

	case "next":
		return getNextEvent(impl, deviceID)

	default:
		if v, err := strconv.ParseUint(index, 10, 32); err != nil {
			return http.StatusBadRequest,
				nil,
				errors.NewRESTError("get-events", fmt.Sprintf("Invalid event index (%v)", index))
		} else {
			return getEvent(impl, deviceID, uint32(v))
		}
	}
}

func getEvent(impl uhppoted.IUHPPOTED, deviceID uint32, index uint32) (int, interface{}, error) {
	event, err := impl.GetEvent(deviceID, index)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("get-event", fmt.Sprintf("%v", err)),
			err
	} else if event == nil {
		return http.StatusNotFound,
			errors.NewRESTError("get-event", fmt.Sprintf("Error retrieving event %v from %v", index, deviceID)),
			fmt.Errorf("No response returned to request for event %v from device %v", index, deviceID)
	}

	return http.StatusOK, struct {
		Event interface{} `json:"event"`
	}{
		Event: event,
	}, nil
}

func getNextEvent(impl uhppoted.IUHPPOTED, deviceID uint32) (int, interface{}, error) {
	event, err := impl.GetNextEvent(deviceID)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("get-event", fmt.Sprintf("%v", err)),
			err
	} else if event == nil {
		return http.StatusNotFound,
			errors.NewRESTError("get-event", fmt.Sprintf("Error retrieving 'next' event from %v", deviceID)),
			fmt.Errorf("No response returned to request for 'next' event from device %v", deviceID)
	}

	return http.StatusOK, struct {
		Event interface{} `json:"event"`
	}{
		Event: event,
	}, nil
}
