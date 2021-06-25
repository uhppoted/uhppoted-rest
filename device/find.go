package device

import (
	"context"
	"fmt"
	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppoted-lib/uhppoted"
	"github.com/uhppoted/uhppoted-rest/errors"
	"net"
	"net/http"
)

type device struct {
	DeviceID   uint32 `json:"device-id"`
	DeviceType string `json:"device-type"`
}

func GetDevices(impl uhppoted.IUHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	rq := uhppoted.GetDevicesRequest{}

	response, err := impl.GetDevices(rq)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("get-devices", "Error searching for active devices"),
			err
	} else if response == nil {
		return http.StatusOK, nil, nil
	}

	devices := make([]device, 0)

	for k, v := range response.Devices {
		devices = append(devices, device{
			DeviceID:   k,
			DeviceType: v.DeviceType,
		})
	}

	return http.StatusOK, struct {
		Devices []device `json:"devices"`
	}{
		Devices: devices,
	}, nil
}

func GetDevice(impl uhppoted.IUHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	deviceID := ctx.Value("device-id").(uint32)

	rq := uhppoted.GetDeviceRequest{
		DeviceID: uhppoted.DeviceID(deviceID),
	}

	response, err := impl.GetDevice(rq)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("get-device", fmt.Sprintf("Could not retrieve device information for %d", deviceID)),
			err
	} else if response == nil {
		return http.StatusOK, nil, nil
	}

	device := struct {
		DeviceType string           `json:"device-type"`
		IPAddress  net.IP           `json:"ip-address"`
		SubnetMask net.IP           `json:"subnet-mask"`
		Gateway    net.IP           `json:"gateway-address"`
		MacAddress types.MacAddress `json:"mac-address"`
		Version    types.Version    `json:"version"`
		Date       types.Date       `json:"date"`
	}{
		DeviceType: response.DeviceType,
		IPAddress:  response.IpAddress,
		SubnetMask: response.SubnetMask,
		Gateway:    response.Gateway,
		MacAddress: response.MacAddress,
		Version:    response.Version,
		Date:       response.Date,
	}

	return http.StatusOK, struct {
		Device interface{} `json:"device"`
	}{
		Device: device,
	}, nil
}
