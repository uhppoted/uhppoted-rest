package device

import (
	"context"
	"fmt"
	"github.com/uhppoted/uhppoted-lib/uhppoted"
	"github.com/uhppoted/uhppoted-rest/errors"
	"net/http"
)

func GetStatus(impl *uhppoted.UHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	deviceID := ctx.Value("device-id").(uint32)

	rq := uhppoted.GetStatusRequest{
		DeviceID: uhppoted.DeviceID(deviceID),
	}

	response, err := impl.GetStatus(rq)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("get-status", fmt.Sprintf("Error retrieving device status for %v", deviceID)),
			err
	} else if response == nil {
		return http.StatusOK, nil, nil
	}

	return http.StatusOK, struct {
		Status uhppoted.Status `json:"status"`
	}{
		Status: response.Status,
	}, nil
}
