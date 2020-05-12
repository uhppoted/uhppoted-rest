package device

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppoted-api/uhppoted"
	"github.com/uhppoted/uhppoted-rest/errors"
	"io/ioutil"
	"net/http"
)

func GetTime(impl *uhppoted.UHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	deviceID := ctx.Value("device-id").(uint32)

	rq := uhppoted.GetTimeRequest{
		DeviceID: uhppoted.DeviceID(deviceID),
	}

	response, err := impl.GetTime(rq)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("get-time", fmt.Sprintf("Error retrieving device time for %v", deviceID)),
			err
	} else if response == nil {
		return http.StatusOK, nil, nil
	}

	return http.StatusOK, struct {
		DateTime types.DateTime `json:"datetime"`
	}{
		DateTime: response.DateTime,
	}, nil
}

func SetTime(impl *uhppoted.UHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	deviceID := ctx.Value("device-id").(uint32)

	blob, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("set-time", "Error reading request"),
			err
	}

	body := struct {
		DateTime *types.DateTime `json:"datetime"`
	}{}

	err = json.Unmarshal(blob, &body)
	if err != nil {
		return http.StatusBadRequest,
			errors.NewRESTError("set-time", "Error parsing request"),
			err
	}

	if body.DateTime == nil {
		return http.StatusBadRequest,
			errors.NewRESTError("set-time", "Missing/invalid date/time"),
			fmt.Errorf("Missing/invalid date-time in  request body (%s)", string(blob))
	}

	rq := uhppoted.SetTimeRequest{
		DeviceID: uhppoted.DeviceID(deviceID),
		DateTime: *body.DateTime,
	}

	response, err := impl.SetTime(rq)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("set-time", "Error setting device date/time"),
			err
	} else if response == nil {
		return http.StatusOK, nil, nil
	}

	return http.StatusOK, &struct {
		DateTime types.DateTime `json:"datetime"`
	}{
		DateTime: response.DateTime,
	}, nil
}
