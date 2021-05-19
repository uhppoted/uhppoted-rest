package device

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
)

func reply(ctx context.Context, w http.ResponseWriter, response interface{}) {
	b, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error generating response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if len(b) > 1024 && ctx.Value("compression") == "gzip" {
		w.Header().Set("Content-Encoding", "gzip")
		encoder := gzip.NewWriter(w)
		encoder.Write(b)
		encoder.Flush()
	} else {
		w.Write(b)
	}
}

func authorized(ctx context.Context, cardNumber uint32) bool {
	cards := ctx.Value("authorized-cards").([]string)
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

	profileID, err := strconv.ParseUint(matches[1], 10, 8)
	if err != nil {
		return 0, err
	}

	return uint8(profileID), nil
}
