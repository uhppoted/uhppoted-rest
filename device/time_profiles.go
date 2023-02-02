package device

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppoted-lib/uhppoted"
	rerrors "github.com/uhppoted/uhppoted-rest/errors"
)

func GetTimeProfile(impl uhppoted.IUHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
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
			rerrors.NewRESTError("get-time-profile", fmt.Sprintf("Error retrieving time profile %v from controller %v", profileID, deviceID)),
			fmt.Errorf("no response returned to request for time profile %v from controller %v", profileID, deviceID)
	}

	return http.StatusOK, struct {
		TimeProfile interface{} `json:"time-profile"`
	}{
		TimeProfile: response.TimeProfile,
	}, nil
}

func PutTimeProfile(impl uhppoted.IUHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
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

	rq := uhppoted.PutTimeProfileRequest{
		DeviceID:    deviceID,
		TimeProfile: profile,
	}

	if _, err = impl.PutTimeProfile(rq); err != nil {
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

func GetTimeProfiles(impl uhppoted.IUHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	deviceID, err := getDeviceID(r)
	if err != nil {
		return http.StatusBadRequest,
			rerrors.NewRESTError("get-time-profiles", fmt.Sprintf("Error:  %v", err)),
			err
	}

	rq := uhppoted.GetTimeProfilesRequest{
		DeviceID: deviceID,
	}

	response, err := impl.GetTimeProfiles(rq)
	if err != nil {
		return http.StatusInternalServerError,
			rerrors.NewRESTError("get-time-profiles", fmt.Sprintf("Error retrieving time profiles from controller %v", deviceID)),
			err
	} else if response == nil {
		return http.StatusInternalServerError,
			rerrors.NewRESTError("get-cards", fmt.Sprintf("Error retrieving time profiles from controller %v", deviceID)),
			fmt.Errorf("no response returned to request for time profiles from controller %v", deviceID)
	}

	return http.StatusOK, &struct {
		Profiles []types.TimeProfile `json:"profiles"`
	}{
		Profiles: response.Profiles,
	}, nil
}

func PutTimeProfiles(impl uhppoted.IUHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	deviceID, err := getDeviceID(r)
	if err != nil {
		return http.StatusBadRequest,
			rerrors.NewRESTError("set-time-profiles", fmt.Sprintf("Error:  %v", err)),
			err
	}

	blob, err := io.ReadAll(r.Body)
	if err != nil {
		return http.StatusInternalServerError,
			rerrors.NewRESTError("set-time-profiles", "Error reading request"),
			err
	}

	body := struct {
		Profiles []types.TimeProfile `json:"profiles"`
	}{}

	if err = json.Unmarshal(blob, &body); err != nil {
		return http.StatusBadRequest,
			rerrors.NewRESTError("set-time-profiles", "Invalid request format"),
			err
	}

	rq := uhppoted.PutTimeProfilesRequest{
		DeviceID: deviceID,
		Profiles: body.Profiles,
	}

	response, code, err := impl.PutTimeProfiles(rq)
	if err != nil {
		if code == http.StatusBadRequest {
			return http.StatusBadRequest,
				rerrors.NewRESTError("set-time-profiles", fmt.Sprintf("Error: %v", err)),
				err
		} else {
			return http.StatusInternalServerError,
				rerrors.NewRESTError("set-time-profiles", fmt.Sprintf("Error updating time profiles on controller %v", deviceID)),
				err
		}
	} else if response == nil {
		return http.StatusInternalServerError,
			rerrors.NewRESTError("set-time-profiles", fmt.Sprintf("Error updating time profiles on controller %v", deviceID)),
			fmt.Errorf("no response returned to time profile update on controller %v", deviceID)
	}

	warnings := []string{}
	for _, warning := range response.Warnings {
		warnings = append(warnings, fmt.Sprintf("%v", warning))
	}

	return http.StatusOK, struct {
		Warnings []string `json:"warnings"`
	}{
		Warnings: warnings,
	}, nil
}

func ClearTimeProfiles(impl uhppoted.IUHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
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
			fmt.Errorf("no response for clear-time-profiles request from controller %v", deviceID)
	}

	if !response.Cleared {
		return http.StatusInternalServerError,
			rerrors.NewRESTError("clear-time-profiles", fmt.Sprintf("Failed to clear all time profiles from controller %v", deviceID)),
			fmt.Errorf("clear-time-profiles from controler %v returned %v", deviceID, response.Cleared)
	}

	return http.StatusOK, nil, nil
}
