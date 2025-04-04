package rcUtilsCommands

import (
	"fmt"
	"github.com/spf13/cobra"
	ble "ringcli/lib/ble"
	ring "ringcli/lib/colmi"
	log "ringcli/lib/log"
	utils "ringcli/lib/utils"
	"tinygo.org/x/bluetooth"
)

// Standard record for device information
type DeviceInfo struct {
	maker    string
	firmware string
	hardware string
	name     string
	system   string
	pnp      string
	battery  ring.BatteryInfo
}

// Globals relevant only to this command
var (
	deviceInfo DeviceInfo = DeviceInfo{} // Device info record
)

// Define the `info` subcommand.
var InfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get ring info",
	Long:  "Get smart ring information, including battery state.",
	Run:   getInfo,
}

func getInfo(cmd *cobra.Command, args []string) {

	// Make sure we have a ring BLE address from the command line or store
	getRingAddress()

	bspCount = log.Raw("Retrieving ring information...  ")
	utils.AnimateCursor()

	// Enable BLE
	deviceInfo.battery.Level = 0
	device := ble.EnableAndConnect(ringAddress)
	defer ble.Disconnect(device)
	infoService := ble.DeviceInfoService(device)

	// Get the device data
	processDeviceInfo(infoService)

	// Get the battery data
	requestBatteryInfo(device)

	// Output received ring data
	outputRingInfo(false)
}

func requestBatteryInfo(device bluetooth.Device) {

	requestPacket := ring.MakeBatteryRequest()
	ble.RequestDataViaCommandUART(device, requestPacket, receiveBatteryInfo, 1)
}

func receiveBatteryInfo(receivedData []byte) {

	if receivedData[0] == ring.COMMAND_BATTERY_INFO {
		deviceInfo.battery = ring.ParseBatteryResponse(receivedData)
		ble.UARTInfoReceived = true
	}
}

func processDeviceInfo(service bluetooth.DeviceService) {

	// Set the BLE service, characteristic UUIDs
	uuidvendor := ble.UUIDFromUInt16(ble.BLE_DEVICE_INFO_SERVICE_MANUFACTURER_CHAR_ID)
	uuidfirmware := ble.UUIDFromUInt16(ble.BLE_DEVICE_INFO_SERVICE_FIRMWARE_VERSION_CHAR_ID)
	uuidhardware := ble.UUIDFromUInt16(ble.BLE_DEVICE_INFO_SERVICE_HARDWARE_VERSION_CHAR_ID)
	uuidsystemid := ble.UUIDFromUInt16(ble.BLE_DEVICE_INFO_SERVICE_SYSTEM_ID_CHAR_ID)
	uuidpnpid := ble.UUIDFromUInt16(ble.BLE_DEVICE_INFO_SERVICE_PNP_ID_CHAR_ID)

	// Get a list of characteristics within the service
	characteristics := ble.Characteristics(service, []bluetooth.UUID{
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
			switch c {
			case uuidvendor:
				deviceInfo.maker = string(data)
			case uuidfirmware:
				deviceInfo.firmware = string(data)
			case uuidhardware:
				deviceInfo.hardware = string(data)
			case uuidsystemid:
				deviceInfo.system = decodeSysId(data)
			case uuidpnpid:
				deviceInfo.pnp = decodePnP(data)
			}
		}
	}
}

func outputRingInfo(showBatteryOnly bool) {

	utils.StopAnimation()

	if bspCount > 0 {
		log.Backspaces(bspCount)
	}

	chargeState := "not charging"
	if deviceInfo.battery.IsCharging {
		chargeState = "charging"
	}

	if showBatteryOnly {
		log.Report("Battery state: %d%% (%s)", deviceInfo.battery.Level, chargeState)
		return
	}

	log.Report("Ring Info:                     ")
	log.Report("Firmware Version: %s", deviceInfo.firmware)
	log.Report("Hardware Version: %s", deviceInfo.hardware)
	log.Report("    Manufacturer: %s", deviceInfo.maker)
	log.Report("       System ID: %s", deviceInfo.system)
	log.Report("          PnP ID: %s", deviceInfo.pnp)
}

func decodePnP(data []byte) string {

	vendorId := int(data[2])<<8 + int(data[1])
	productID := int(data[4])<<8 + int(data[3])
	productVersion := int(data[6])<<8 + int(data[5])

	return fmt.Sprintf("Vendor ID 0x%04X Product ID 0x%04X Product Version 0x%04X", vendorId, productID, productVersion)
}

func decodeSysId(data []byte) string {

	total := 0
	for i := range 8 {
		total += int(data[i]) << (56 - (i * 8))
	}

	return fmt.Sprintf("0x%08X", total)
}
