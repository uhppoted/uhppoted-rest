package device

import (
	"context"
	"fmt"
	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppoted-api/uhppoted"
	"github.com/uhppoted/uhppoted-rest/errors"
	"net"
	"net/http"
)

type device struct {
	DeviceID   uint32 `json:"device-id"`
	DeviceType string `json:"device-type"`
}

func GetDevices(impl *uhppoted.UHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (interface{}, *errors.IError) {
	rq := uhppoted.GetDevicesRequest{}

	response, err := impl.GetDevices(rq)
	if err != nil {
		return nil, errors.Errorf(err, 0, "get-devices", "Error searching for active devices")
	} else if response == nil {
		return nil, nil
	}

	devices := make([]device, 0)

	for k, v := range response.Devices {
		devices = append(devices, device{
			DeviceID:   k,
			DeviceType: v.DeviceType,
		})
	}

	return struct {
		Devices []device `json:"devices"`
	}{
		Devices: devices,
	}, nil
}

func GetDevice(impl *uhppoted.UHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (interface{}, *errors.IError) {
	deviceID := ctx.Value("device-id").(uint32)

	rq := uhppoted.GetDeviceRequest{
		DeviceID: uhppoted.DeviceID(deviceID),
	}

	response, err := impl.GetDevice(rq)
	if err != nil {
		return nil, errors.Errorf(err, deviceID, "get-device", fmt.Sprintf("Could not retrieve device information for %d", deviceID))
	}

	if response == nil {
		return nil, nil
	}

	reply := struct {
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

	return reply, nil
}
