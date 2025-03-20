package device

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppoted-lib/uhppoted"

	"github.com/uhppoted/uhppoted-rest/errors"
	"github.com/uhppoted/uhppoted-rest/lib"
)

func authorized(ctx context.Context, cardNumber uint32) bool {
	cards := ctx.Value(lib.AuthorizedCards).([]string)
	c := fmt.Sprintf("%v", cardNumber)
	for _, re := range cards {
		if ok, err := regexp.MatchString(re, c); ok && err == nil {
			return true
		}
	}

	return false
}

func getDeviceID(r *http.Request) (uint32, error) {
	matches := regexp.MustCompile("^/uhppote/device/([0-9]+)(?:$|/.*$)").FindStringSubmatch(r.URL.Path)
	if matches == nil {
		return 0, fmt.Errorf("missing device-id")
	}

	deviceID, err := strconv.ParseUint(matches[1], 10, 32)
	if err != nil {
		return 0, err
	}

	return uint32(deviceID), nil
}

func getTimeProfileID(r *http.Request) (uint8, error) {
	matches := regexp.MustCompile("^/uhppote/device/[0-9]+/time-profile/([0-9]+)$").FindStringSubmatch(r.URL.Path)
	if matches == nil {
		return 0, fmt.Errorf("missing time-profile-id")
	}

	profileID, err := strconv.ParseUint(matches[1], 10, 32)
	if err != nil {
		return 0, err
	}

	if profileID < 2 || profileID > 254 {
		return 0, fmt.Errorf("invalid time profile ID (%v) - valid range is [2..254]", profileID)
	}

	return uint8(profileID), nil
}

func GetAntiPassback(impl uhppoted.IUHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, any, error) {
	controller := ctx.Value(lib.DeviceID).(uint32)

	if antipassback, err := impl.GetAntiPassback(controller); err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("get-antipassback", "Error retrieving controller anti-passback"),
			err
	} else {
		return http.StatusOK, &struct {
			AntiPassback string `json:"anti-passback"`
		}{
			AntiPassback: fmt.Sprintf("%v", antipassback),
		}, nil
	}
}

func SetAntiPassback(impl uhppoted.IUHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, any, error) {
	controller := ctx.Value(lib.DeviceID).(uint32)

	blob, err := io.ReadAll(r.Body)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("set-antipassback", "Error reading request"),
			err
	}

	body := struct {
		AntiPassback string `json:"anti-passback"`
	}{}

	if err := json.Unmarshal(blob, &body); err != nil {
		return http.StatusBadRequest,
			errors.NewRESTError("set-antipassback", "Error parsing request"),
			err
	}

	if antipassback, err := func(s string) (types.AntiPassback, error) {
		v := regexp.MustCompile(`[ (),]+`).ReplaceAllString(s, "")

		switch v {
		case "disabled":
			return types.Disabled, nil

		case "1:2;3:4":
			return types.Readers12_34, nil

		case "13:24":
			return types.Readers13_24, nil

		case "1:23":
			return types.Readers1_23, nil

		case "1:234":
			return types.Readers1_234, nil

		default:
			return types.Disabled, fmt.Errorf("invalid anti-passback value (%v)", s)
		}

	}(body.AntiPassback); err != nil {
		return http.StatusBadRequest,
			errors.NewRESTError("set-antipassback", "Missing/invalid AntiPassback"),
			fmt.Errorf("missing/invalid anti-passback value in request body (%s)", string(blob))

	} else if ok, err := impl.SetAntiPassback(controller, antipassback); err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("set-antipassback", "Error setting controller anti-passback"),
			err

	} else if !ok {
		return http.StatusInternalServerError,
			errors.NewRESTError("set-antipassback", "Failed to set controller anti-passback"),
			err

	} else {
		return http.StatusOK, &struct {
			AntiPassback string `json:"anti-passback"`
		}{
			AntiPassback: fmt.Sprintf("%v", antipassback),
		}, nil
	}
}

func RestoreDefaultParameters(impl uhppoted.IUHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, any, error) {
	deviceID := ctx.Value(lib.DeviceID).(uint32)

	if err := impl.RestoreDefaultParameters(deviceID); err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("restore-default-parameters", "Error restoring controller default parameters"),
			err
	}

	return http.StatusOK, &struct {
	}{}, nil
}
