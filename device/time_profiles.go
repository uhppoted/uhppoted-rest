package device

import (
	"context"
	"fmt"
	"net/http"

	"github.com/uhppoted/uhppoted-api/uhppoted"
	"github.com/uhppoted/uhppoted-rest/errors"
)

func GetTimeProfile(impl *uhppoted.UHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	deviceID, _ := getDeviceID(r)
	profileID, _ := getTimeProfileID(r)

	rq := uhppoted.GetTimeProfileRequest{
		DeviceID:  deviceID,
		ProfileID: profileID,
	}

	response, err := impl.GetTimeProfile(rq)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("get-time-profile", fmt.Sprintf("Error retrieving time profile %v from device %v", profileID, deviceID)),
			err
	} else if response == nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("get-card", fmt.Sprintf("Error retrieving time profile %v from device %v", profileID, deviceID)),
			fmt.Errorf("No response returned to request for time profile %v from device %v", profileID, deviceID)
	}

	return http.StatusOK, struct {
		TimeProfile interface{} `json:"time-profile"`
	}{
		TimeProfile: response.TimeProfile,
	}, nil
}
