package rcUtilsCommands

import (
	"fmt"
	// External code
	"github.com/spf13/cobra"
	"tinygo.org/x/bluetooth"
	// Library code
	rcBLE "ringcli/lib/ble"
	rcErrors "ringcli/lib/errors"
	rcLog "ringcli/lib/log"
	rcColmi "ringcli/lib/colmi"
)

type DeviceInfo struct {
	maker string
	firmware string
	hardware string
	name string
	system string
	pnp string
	battery rcColmi.BatteryInfo
}

var (
	batteryInfoReceived bool = false
	deviceInfo DeviceInfo = DeviceInfo{}
)

// Define the `scan` subcommand.
var InfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get ring info",
	Long:  "Get ring info",
	Run:   getInfo,
}

func getInfo(cmd *cobra.Command, args []string) {

	// Bail when no ID data is provided
	if ringName == "" && ringAddress == "" {
		rcLog.ReportErrorAndExit(rcErrors.ERROR_CODE_BAD_PARAMS, "No name or address supplied")
	}

	// Enable BLE
	ble := rcBLE.Open()
	deviceInfo.battery.Level = 0

	// Generate the ring BLE address and connect to it
	bleAddress := rcBLE.AddressFromString(ringAddress)
	device := rcBLE.Connect(ble, bleAddress)
	defer rcBLE.Disconnect(device)
	infoService := rcBLE.DeviceInfoService(device)

	// Get the device data
	requestDeviceInfo(infoService)

	// Get the battery data
	requestBatteryInfo(device)

	// Output received ring data
	outputRingInfo()
}

func requestBatteryInfo(ble bluetooth.Device) {

	requestPacket := rcColmi.MakeBatteryReq()
	rcBLE.RequestDataViaCommandUART(ble, requestPacket, receiveBatteryInfo, 1)
}

func receiveBatteryInfo(receivedData []byte) {

	if receivedData[0] == rcColmi.COMMAND_BATTERY_INFO {
		deviceInfo.battery = rcColmi.ParseBatteryResp(receivedData)
		rcBLE.UARTInfoReceived = true
	}
}

func requestDeviceInfo(service bluetooth.DeviceService) {

	uuidvendor := rcBLE.UUIDFromUInt16(rcBLE.BLE_DEVICE_INFO_SERVICE_MANUFACTURER_CHAR_ID)
	uuidfirmware := rcBLE.UUIDFromUInt16(rcBLE.BLE_DEVICE_INFO_SERVICE_FIRMWARE_VERSION_CHAR_ID)
	uuidhardware := rcBLE.UUIDFromUInt16(rcBLE.BLE_DEVICE_INFO_SERVICE_HARDWARE_VERSION_CHAR_ID)
	uuidsystemid := rcBLE.UUIDFromUInt16(rcBLE.BLE_DEVICE_INFO_SERVICE_SYSTEM_ID_CHAR_ID)
	uuidpnpid := rcBLE.UUIDFromUInt16(rcBLE.BLE_DEVICE_INFO_SERVICE_PNP_ID_CHAR_ID)

	// Get a list of characteristics within the service
	characteristics := rcBLE.Characteristics(service, []bluetooth.UUID{
		uuidvendor,
		uuidfirmware,
		uuidhardware,
		uuidsystemid,
		uuidpnpid,
	})

	// Iterate over characteristics
	for _, characteristic := range characteristics {
		var data = make([]byte, 64, 64)
		_, err := characteristic.Read(data)
		if err == nil {
			c := characteristic.UUID()

			if c == uuidvendor {
				deviceInfo.maker = string(data)
			}

			if c == uuidfirmware {
				deviceInfo.firmware = string(data)
			}

			if c == uuidhardware {
				deviceInfo.hardware = string(data)
			}

			if c == uuidsystemid {
				deviceInfo.system = decodeSysId(data)
			}

			if c == uuidpnpid {
				deviceInfo.pnp = decodePnP(data)
			}
		}
	}
}

func outputRingInfo() {

	chargeState := "not charging"
	if deviceInfo.battery.IsCharging {
		chargeState = "charging"
	}

	rcLog.Report("Ring Info:")
	rcLog.Report("         Battery: %d%% (%s)", deviceInfo.battery.Level, chargeState)
	rcLog.Report("Firmware Version: %s", deviceInfo.firmware)
	rcLog.Report("Hardware Version: %s", deviceInfo.hardware)
	rcLog.Report("    Manufacturer: %s", deviceInfo.maker)
	rcLog.Report("       System ID: %s", deviceInfo.system)
	rcLog.Report("          PnP ID: %s", deviceInfo.pnp)
}

func decodePnP(data []byte) string {

	vendorId := int(data[2]) << 8 + int(data[1])
	productID := int(data[4]) << 8 + int(data[3])
	productVersion := int(data[6]) << 8 + int(data[5])

	return fmt.Sprintf("Vendor ID 0x%04X Product ID 0x%04X Product Version 0x%04X", vendorId, productID, productVersion)
}

func decodeSysId(data []byte) string {

	total := 0
	for i := range(8) {
		total += int(data[i]) << (56 - (i * 8))
	}

	return fmt.Sprintf("0x%08X", total)
}