package device

import (
	"context"
	"fmt"
	"github.com/uhppoted/uhppoted-api/uhppoted"
	"github.com/uhppoted/uhppoted-rest/errors"
	"net/http"
)

func GetStatus(impl *uhppoted.UHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (interface{}, *errors.IError) {
	deviceID := ctx.Value("device-id").(uint32)

	rq := uhppoted.GetStatusRequest{
		DeviceID: uhppoted.DeviceID(deviceID),
	}

	response, err := impl.GetStatus(rq)
	if err != nil {
		return nil, errors.Errorf(err, deviceID, "get-status", fmt.Sprintf("Error retrieving status for %v", deviceID))
	}

	if response == nil {
		return nil, nil
	}

	return response.Status, nil
}
