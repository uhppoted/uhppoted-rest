package device

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppoted-api/uhppoted"
	rerrors "github.com/uhppoted/uhppoted-rest/errors"
)

func GetTimeProfile(impl *uhppoted.UHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	deviceID, err := getDeviceID(r)
	if err != nil {
		return http.StatusBadRequest,
			rerrors.NewRESTError("set-time-profile", fmt.Sprintf("Error:  %v", err)),
			err
	}

	profileID, err := getTimeProfileID(r)
	if err != nil {
		return http.StatusBadRequest,
			rerrors.NewRESTError("get-time-profile", fmt.Sprintf("Error:  %v", err)),
			err
	}

	rq := uhppoted.GetTimeProfileRequest{
		DeviceID:  deviceID,
		ProfileID: profileID,
	}

	response, err := impl.GetTimeProfile(rq)
	if err != nil {
		return http.StatusInternalServerError,
			rerrors.NewRESTError("get-time-profile", fmt.Sprintf("Error retrieving time profile %v from controller %v", profileID, deviceID)),
			err
	} else if response == nil {
		return http.StatusInternalServerError,
			rerrors.NewRESTError("get-card", fmt.Sprintf("Error retrieving time profile %v from controller %v", profileID, deviceID)),
			fmt.Errorf("No response returned to request for time profile %v from controller %v", profileID, deviceID)
	}

	return http.StatusOK, struct {
		TimeProfile interface{} `json:"time-profile"`
	}{
		TimeProfile: response.TimeProfile,
	}, nil
}

func PutTimeProfile(impl *uhppoted.UHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	deviceID, err := getDeviceID(r)
	if err != nil {
		return http.StatusBadRequest,
			rerrors.NewRESTError("set-time-profile", fmt.Sprintf("Error:  %v", err)),
			err
	}

	profileID, err := getTimeProfileID(r)
	if err != nil {
		return http.StatusBadRequest,
			rerrors.NewRESTError("set-time-profile", fmt.Sprintf("Error:  %v", err)),
			err
	}

	blob, err := io.ReadAll(r.Body)
	if err != nil {
		return http.StatusInternalServerError,
			rerrors.NewRESTError("set-time-profile", "Error reading request"),
			err
	}

	profile := types.TimeProfile{}
	if err = json.Unmarshal(blob, &profile); err != nil {
		return http.StatusBadRequest,
			rerrors.NewRESTError("set-time-profile", fmt.Sprintf("Error parsing request (%v)", err)),
			err
	}

	profile.ID = profileID

	rq := uhppoted.SetTimeProfileRequest{
		DeviceID:    deviceID,
		TimeProfile: profile,
	}

	if _, err = impl.SetTimeProfile(rq); err != nil {
		if errors.Unwrap(err) == nil {
			return http.StatusBadRequest,
				rerrors.NewRESTError("set-time-profile", fmt.Sprintf("%v", err)),
				err
		} else {
			return http.StatusInternalServerError,
				rerrors.NewRESTError("set-time-profile", fmt.Sprintf("Error creating/updating time profile %v on controller %v", profileID, deviceID)),
				err
		}
	}

	return http.StatusOK, nil, nil
}

func ClearTimeProfiles(impl *uhppoted.UHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	deviceID, err := getDeviceID(r)
	if err != nil {
		return http.StatusBadRequest,
			rerrors.NewRESTError("clear-time-profiles", fmt.Sprintf("Error:  %v", err)),
			err
	}

	rq := uhppoted.ClearTimeProfilesRequest{
		DeviceID: deviceID,
	}

	response, err := impl.ClearTimeProfiles(rq)
	if err != nil {
		return http.StatusInternalServerError,
			rerrors.NewRESTError("clear-time-profiles", fmt.Sprintf("Error clearing all time profiles from controller %v", deviceID)),
			err
	} else if response == nil {
		return http.StatusInternalServerError,
			rerrors.NewRESTError("clear-time-profiles", fmt.Sprintf("Error clearing all time profiles from %v", deviceID)),
			fmt.Errorf("No response for clear-time-profiles request from controller %v", deviceID)
	}

	if !response.Cleared {
		return http.StatusInternalServerError,
			rerrors.NewRESTError("clear-time-profiles", fmt.Sprintf("Failed to clear all time profiles from controller %v", deviceID)),
			fmt.Errorf("clear-time-profiles from controler %v returned %v", deviceID, response.Cleared)
	}

	return http.StatusOK, nil, nil
}
