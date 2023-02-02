package device

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppoted-lib/locales"
	"github.com/uhppoted/uhppoted-lib/uhppoted"
	"github.com/uhppoted/uhppoted-rest/errors"
)

type Event struct {
	DeviceID      uint32         `json:"device-id"`
	Index         uint32         `json:"event-id"`
	Type          uint8          `json:"event-type"`
	TypeText      string         `json:"event-type-text"`
	Granted       bool           `json:"access-granted"`
	Door          uint8          `json:"door-id"`
	Direction     uint8          `json:"direction"`
	DirectionText string         `json:"direction-text"`
	CardNumber    uint32         `json:"card-number"`
	Timestamp     types.DateTime `json:"timestamp"`
	Reason        uint8          `json:"event-reason"`
	ReasonText    string         `json:"event-reason-text"`
}

func GetEvents(impl uhppoted.IUHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, any, error) {
	url := r.URL.Path
	count := 0
	events := []any{}

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
			First   uint32 `json:"first,omitempty"`
			Last    uint32 `json:"last,omitempty"`
			Current uint32 `json:"current,omitempty"`
			Events  []any  `json:"events,omitempty"`
		} `json:"events"`
	}{}

	response.Events.First = first
	response.Events.Last = last
	response.Events.Current = current
	response.Events.Events = events

	return http.StatusOK, &response, nil
}

func GetEvent(impl uhppoted.IUHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, any, error) {
	var url = r.URL.Path
	var index string

	deviceID, err := getDeviceID(r)
	if err != nil {
		return http.StatusBadRequest, nil, errors.NewRESTError("get-events", "Missing device ID")
	}

	if matches := regexp.MustCompile("^/uhppote/device/[0-9]+/event/([0-9]+|first|last|current|next)$").FindStringSubmatch(url); matches == nil {
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

func getEvent(impl uhppoted.IUHPPOTED, deviceID uint32, index uint32) (int, any, error) {
	event, err := impl.GetEvent(deviceID, index)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("get-event", fmt.Sprintf("%v", err)),
			err
	} else if event == nil {
		return http.StatusNotFound,
			errors.NewRESTError("get-event", fmt.Sprintf("Error retrieving event %v from %v", index, deviceID)),
			fmt.Errorf("no response returned to request for event %v from device %v", index, deviceID)
	}

	return http.StatusOK, struct {
		Event any `json:"event"`
	}{
		Event: Transmogrify(*event),
	}, nil
}

func getNextEvent(impl uhppoted.IUHPPOTED, deviceID uint32) (int, any, error) {
	response := struct {
		DeviceID uint32 `json:"device-id"`
		Event    any    `json:"event"`
	}{
		DeviceID: deviceID,
	}

	events, err := impl.GetEvents(deviceID, 1)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("get-event", fmt.Sprintf("%v", err)),
			err
	} else if events == nil {
		return http.StatusNotFound,
			errors.NewRESTError("get-event", fmt.Sprintf("Error retrieving 'next' event from %v", deviceID)),
			fmt.Errorf("no response returned to request for 'next' event from device %v", deviceID)
	} else if len(events) > 0 {
		response.Event = Transmogrify(events[0])
	}

	return http.StatusOK, &response, nil
}

func Transmogrify(e uhppoted.Event) any {
	lookup := func(key string) string {
		if v, ok := locales.Lookup(key); ok {
			return v
		}

		return ""
	}

	return Event{
		DeviceID:      e.DeviceID,
		Index:         e.Index,
		Type:          e.Type,
		TypeText:      lookup(fmt.Sprintf("event.type.%v", e.Type)),
		Granted:       e.Granted,
		Door:          e.Door,
		Direction:     e.Direction,
		DirectionText: lookup(fmt.Sprintf("event.direction.%v", e.Direction)),
		CardNumber:    e.CardNumber,
		Timestamp:     e.Timestamp,
		Reason:        e.Reason,
		ReasonText:    lookup(fmt.Sprintf("event.reason.%v", e.Reason)),
	}
}
