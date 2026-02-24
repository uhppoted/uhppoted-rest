package device

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppoted-lib/uhppoted"
	rerrors "github.com/uhppoted/uhppoted-rest/errors"
)

func PutTaskList(impl uhppoted.IUHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, any, error) {
	deviceID, err := getDeviceID(r)
	if err != nil {
		return http.StatusBadRequest,
			rerrors.NewRESTError("set-task-list", fmt.Sprintf("Error:  %v", err)),
			err
	}

	blob, err := io.ReadAll(r.Body)
	if err != nil {
		return http.StatusInternalServerError,
			rerrors.NewRESTError("set-task-list", "Error reading request"),
			err
	}

	body := struct {
		Tasks []types.Task `json:"tasks"`
	}{}

	if err = json.Unmarshal(blob, &body); err != nil {
		return http.StatusBadRequest,
			rerrors.NewRESTError("set-task-list", "Invalid request format"),
			err
	}

	rq := uhppoted.PutTaskListRequest{
		DeviceID: deviceID,
		Tasks:    body.Tasks,
	}

	response, code, err := impl.PutTaskList(rq)
	if err != nil {
		if code == http.StatusBadRequest {
			return http.StatusBadRequest,
				rerrors.NewRESTError("set-task-list", fmt.Sprintf("Error: %v", err)),
				err
		} else {
			return http.StatusInternalServerError,
				rerrors.NewRESTError("set-task-list", fmt.Sprintf("Error updating task list controller %v", deviceID)),
				err
		}
	} else if response == nil {
		return http.StatusInternalServerError,
			rerrors.NewRESTError("set-task-list", fmt.Sprintf("Error updating task list on controller %v", deviceID)),
			fmt.Errorf("no response returned to task list update on controller %v", deviceID)
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
