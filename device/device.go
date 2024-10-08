package device

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

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
