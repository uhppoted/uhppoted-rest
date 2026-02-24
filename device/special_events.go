package device

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/uhppoted/uhppoted-lib/uhppoted"

	"github.com/uhppoted/uhppoted-rest/errors"
	"github.com/uhppoted/uhppoted-rest/lib"
)

// Implements the record-special-events REST API. Extracts the 'enabled' value from the
// request body and invokes the uhppoted-lib.RecordSpecialEvents API function to update
// the controller 'record special events' flag.
func SpecialEvents(impl uhppoted.IUHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, any, error) {
	deviceID := ctx.Value(lib.DeviceID).(uint32)

	blob, err := io.ReadAll(r.Body)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("special-events", "Error reading request"),
			err
	}

	body := struct {
		Enabled *bool `json:"enabled"`
	}{}

	err = json.Unmarshal(blob, &body)
	if err != nil {
		return http.StatusBadRequest,
			errors.NewRESTError("special-events", "Error parsing request"),
			err
	}

	if body.Enabled == nil {
		return http.StatusBadRequest,
			errors.NewRESTError("special-events", "Missing/invalid 'enabled'"),
			fmt.Errorf("missing/invalid 'enabled' in request body (%s)", string(blob))
	}

	enabled := *body.Enabled

	updated, err := impl.RecordSpecialEvents(deviceID, enabled)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("special-events", "Error setting 'record special events'"),
			err
	}

	if !updated {
		return http.StatusInternalServerError,
			errors.NewRESTError("special-events", "Failed to update 'record special events'"),
			fmt.Errorf("failed to update 'record special events'")
	}

	return http.StatusOK, &struct {
		Enabled bool `json:"enabled"`
	}{
		Enabled: enabled,
	}, nil
}
