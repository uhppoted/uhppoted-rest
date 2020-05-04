package device

import (
	"context"
	"fmt"
	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
	"github.com/uhppoted/uhppoted-api/uhppoted"
	"github.com/uhppoted/uhppoted-rest/errors"
	"net"
	"net/http"
	"strconv"
	"strings"
)

type device struct {
	SerialNumber uint32 `json:"device-id"`
	DeviceType   string `json:"device-type"`
}

func GetDevices(impl *uhppoted.UHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (interface{}, *errors.IError) {
	debug(ctx, 0, "get-devices", r)

	rq := uhppoted.GetDevicesRequest{}

	response, err := impl.GetDevices(rq)
	if err != nil {
		return nil, errors.Errorf(err, 0, "get-device", "Error searching for active devices")
	} else if response == nil {
		return nil, nil
	}

	devices := make([]device, 0)

	for k, _ := range response.Devices {
		devices = append(devices, device{
			SerialNumber: k,
			DeviceType:   identify(k),
		})
	}

	return struct {
		Devices []device `json:"devices"`
	}{
		Devices: devices,
	}, nil
}

func GetDevice(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	deviceID := ctx.Value("device-id").(uint32)

	device, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).FindDevice(deviceID)

	if err != nil {
		warn(ctx, deviceID, "get-device", err)
		http.Error(w, "Error retrieving device list", http.StatusInternalServerError)
		return
	}

	if device == nil {
		http.Error(w, fmt.Sprintf("No device with ID '%v'", deviceID), http.StatusNotFound)
		return
	}

	response := struct {
		SerialNumber types.SerialNumber `json:"serial-number"`
		DeviceType   string             `json:"device-type"`
		IPAddress    net.IP             `json:"ip-address"`
		SubnetMask   net.IP             `json:"subnet-mask"`
		Gateway      net.IP             `json:"gateway-address"`
		MacAddress   types.MacAddress   `json:"mac-address"`
		Version      types.Version      `json:"version"`
		Date         types.Date         `json:"date"`
	}{
		SerialNumber: device.SerialNumber,
		DeviceType:   identify(uint32(device.SerialNumber)),
		IPAddress:    device.IpAddress,
		SubnetMask:   device.SubnetMask,
		Gateway:      device.Gateway,
		MacAddress:   device.MacAddress,
		Version:      device.Version,
		Date:         device.Date,
	}

	reply(ctx, w, response)
}

func identify(deviceID uint32) string {
	id := strconv.FormatUint(uint64(deviceID), 10)

	if strings.HasPrefix(id, "4") {
		return "UTO311-L04"
	}

	if strings.HasPrefix(id, "3") {
		return "UTO311-L03"
	}

	if strings.HasPrefix(id, "2") {
		return "UTO311-L02"
	}

	if strings.HasPrefix(id, "1") {
		return "UTO311-L01"
	}

	return "UTO311-L0x"
}
