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

func GetTime(impl *uhppoted.UHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (interface{}, *errors.IError) {
	deviceID := ctx.Value("device-id").(uint32)

	rq := uhppoted.GetTimeRequest{
		DeviceID: uhppoted.DeviceID(deviceID),
	}

	response, err := impl.GetTime(rq)
	if err != nil {
		return nil, errors.Errorf(err, deviceID, "get-time", fmt.Sprintf("Error retrieving device time for %v", deviceID))
	}

	if response == nil {
		return nil, nil
	}

	return struct {
		DateTime types.DateTime `json:"datetime"`
	}{
		DateTime: response.DateTime,
	}, nil
}

func SetTime(impl *uhppoted.UHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (interface{}, *errors.IError) {
	deviceID := ctx.Value("device-id").(uint32)

	blob, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Errorf(fmt.Errorf("Error reading request (%w)", err), deviceID, "set-time", "Error reading request")
	}

	body := struct {
		DateTime *types.DateTime `json:"datetime"`
	}{}

	err = json.Unmarshal(blob, &body)
	if err != nil {
		return nil, errors.Errorf(fmt.Errorf("%w: %v", uhppoted.BadRequest, err), deviceID, "set-time", "Error parsing request")
	}

	if body.DateTime == nil {
		return nil, errors.Errorf(errors.InvalidDateTime, deviceID, "set-time", "Missing/invalid date/time")
	}

	rq := uhppoted.SetTimeRequest{
		DeviceID: uhppoted.DeviceID(deviceID),
		DateTime: *body.DateTime,
	}

	response, err := impl.SetTime(rq)
	if err != nil {
		return nil, errors.Errorf(err, deviceID, "set-time", "Error setting device date/time")
	}

	if response == nil {
		return nil, nil
	}

	return &struct {
		DateTime types.DateTime `json:"datetime"`
	}{
		DateTime: response.DateTime,
	}, nil
}
