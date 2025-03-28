package rcUtilsCommands

import (
	"fmt"
	//"time"

	"github.com/spf13/cobra"
	"tinygo.org/x/bluetooth"

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
	Long: "Get ring info",
	Run:    getInfo,
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
	infoService := rcBLE.DeviceInfoService(device)

	// Get the device data
	requestDeviceInfo(infoService)

	// Get the battery data
	requestBatteryInfo(device)

	// Output received ring data
	outputInfo()
}

func requestBatteryInfo(ble bluetooth.Device) {

	requestPacket := rcColmi.MakeBatteryReq()
	rcBLE.RequestDataViaUART(ble, requestPacket, receiveBatteryInfo)
}

func receiveBatteryInfo(receivedData []byte) {

	if receivedData[0] == 0x03 {
		deviceInfo.battery = rcColmi.ParseBatteryResp(receivedData)
		rcBLE.UARTInfoReceived = true
	}
}

func requestDeviceInfo(service bluetooth.DeviceService) {

	// Get a list of characteristics within the service
	characteristics := rcBLE.Characteristics(service, []bluetooth.UUID{
		bluetooth.New16BitUUID(DEVICE_INFO_SERVICE_MANUFACTURER_CHAR_ID),
		bluetooth.New16BitUUID(DEVICE_INFO_SERVICE_FIRMWARE_VERSION_CHAR_ID),
		bluetooth.New16BitUUID(DEVICE_INFO_SERVICE_HARDWARE_VERSION_CHAR_ID),
		bluetooth.New16BitUUID(DEVICE_INFO_SERVICE_SYSTEM_ID_CHAR_ID),
		bluetooth.New16BitUUID(DEVICE_INFO_SERVICE_PNP_ID_CHAR_ID),
	})

	// Iterate over characteristics
	for _, characteristic := range characteristics {
		var data = make([]byte, 64, 64)
		_, err := characteristic.Read(data)
		if err == nil {
			c := characteristic.UUID()

			if c == bluetooth.New16BitUUID(DEVICE_INFO_SERVICE_MANUFACTURER_CHAR_ID) {
				deviceInfo.maker = string(data)
			}

			if c == bluetooth.New16BitUUID(DEVICE_INFO_SERVICE_FIRMWARE_VERSION_CHAR_ID) {
				deviceInfo.firmware = string(data)
			}

			if c == bluetooth.New16BitUUID(DEVICE_INFO_SERVICE_HARDWARE_VERSION_CHAR_ID) {
				deviceInfo.hardware = string(data)
			}

			if c == bluetooth.New16BitUUID(DEVICE_INFO_SERVICE_SYSTEM_ID_CHAR_ID) {
				deviceInfo.system = decodeSysId(data)
			}

			if c == bluetooth.New16BitUUID(DEVICE_INFO_SERVICE_PNP_ID_CHAR_ID) {
				deviceInfo.pnp = decodePnP(data)
			}
		}
	}
}

func outputInfo() {

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