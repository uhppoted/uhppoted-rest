package device

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/uhppoted/uhppoted-lib/uhppoted"
	"github.com/uhppoted/uhppoted-rest/errors"
	"net/http"
)

func GetEvents(impl uhppoted.IUHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	url := r.URL.Path
	count := 0
	events := []interface{}{}

	deviceID, err := getDeviceID(r)
	if err != nil {
		return http.StatusBadRequest, nil, errors.NewRESTError("get-events", "Missing device ID")
	}

	if matches := regexp.MustCompile("^/uhppote/device/[0-9]+/events/([0-9]+)$").FindStringSubmatch(url); matches != nil {
		if N, err := strconv.ParseInt(matches[1], 10, 32); err == nil {
			count = int(N)
		}
	}

	if count > 0 {
		list, err := impl.GetEvents(deviceID, count)
		if err != nil {
			return http.StatusInternalServerError,
				errors.NewRESTError("get-event", fmt.Sprintf("%v", err)),
				err
		}

		for _, e := range list {
			events = append(events, e)
		}
	}

	first, last, current, err := impl.GetEventIndices(deviceID)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("get-events", fmt.Sprintf("Error retrieving event indices from %v", deviceID)),
			err
	}

	response := struct {
		Events struct {
			First   uint32        `json:"first,omitempty"`
			Last    uint32        `json:"last,omitempty"`
			Current uint32        `json:"current,omitempty"`
			Events  []interface{} `json:"events,omitempty"`
		} `json:"events"`
	}{}

	response.Events.First = first
	response.Events.Last = last
	response.Events.Current = current
	response.Events.Events = events

	return http.StatusOK, &response, nil
}

func GetEvent(impl uhppoted.IUHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	var url = r.URL.Path
	var index string

	deviceID, err := getDeviceID(r)
	if err != nil {
		return http.StatusBadRequest, nil, errors.NewRESTError("get-events", "Missing device ID")
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
	event, err := impl.GetEvents(deviceID, 1)
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
