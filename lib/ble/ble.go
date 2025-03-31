package rcBLE

import (
	"time"

	"tinygo.org/x/bluetooth"

	rcColmi "ringcli/lib/colmi"
	rcErrors "ringcli/lib/errors"
	rcLog "ringcli/lib/log"
)

const (
	CONNECT_TIMEOUT_S                                       = 60
	SYNC_TIMEOUT_S                                          = 10
	BLE_DEVICE_INFO_SERVICE_ID                       uint16 = 0x180A
	BLE_DEVICE_INFO_SERVICE_MANUFACTURER_CHAR_ID     uint16 = 0x2A29
	BLE_DEVICE_INFO_SERVICE_FIRMWARE_VERSION_CHAR_ID uint16 = 0x2A26
	BLE_DEVICE_INFO_SERVICE_HARDWARE_VERSION_CHAR_ID uint16 = 0x2A27
	BLE_DEVICE_INFO_SERVICE_SYSTEM_ID_CHAR_ID        uint16 = 0x2A23
	BLE_DEVICE_INFO_SERVICE_PNP_ID_CHAR_ID           uint16 = 0x2A50
)

var (
	connectTimer *time.Timer
	syncTimer *time.Timer
	UARTInfoReceived bool = false
	isConnected bool = false
	currentDevice *bluetooth.Device
)

func Open() *bluetooth.Adapter {

	ble := bluetooth.DefaultAdapter
	if ble.Enable() != nil {
		rcLog.ReportErrorAndExit(rcErrors.ERROR_CODE_BLE, "Failed to enable BLE")
	}

	return ble
}

func Connect(adapter *bluetooth.Adapter, ringAddress bluetooth.Address) bluetooth.Device {

	// Establish a timer so we're not trying to connect forever
	connectTimer = time.NewTimer(CONNECT_TIMEOUT_S * time.Second)
	defer connectTimer.Stop()
	go func() {
		<-connectTimer.C
		rcLog.ReportErrorAndExit(rcErrors.ERROR_CODE_BLE, "Could not connect to %s, or connection timed out", ringAddress.String())
	}()

	// Attempt to connect
	device, err := adapter.Connect(ringAddress, bluetooth.ConnectionParams{})
	if err != nil {
		connectTimer.Stop()
		rcLog.ReportErrorAndExit(rcErrors.ERROR_CODE_BLE, "Could not connect to %s", ringAddress)
	}

	isConnected = true
	currentDevice = &device
	return device
}

func Disconnect(device bluetooth.Device) {

	err := device.Disconnect()
	if err != nil {
		rcLog.ReportErrorAndExit(rcErrors.ERROR_CODE_BLE, "Could not cleanly disconnect from %s", device.Address.String())
	}

	isConnected = false
}

func Clean() {

	if currentDevice != nil && isConnected {
		Disconnect(*currentDevice)
	}
}

func EnableAndConnect(ringAddress string) bluetooth.Device {

	return Connect(Open(), AddressFromString(ringAddress))
}

func Services(ble bluetooth.Device, uuids []bluetooth.UUID) []bluetooth.DeviceService {

	services, err := ble.DiscoverServices(uuids)
	if err != nil {
		connectTimer.Stop()
		rcLog.ReportErrorAndExit(rcErrors.ERROR_CODE_BLE, "Could not get ring BLE service list: %s", err.Error())
	}

	return services
}

func Characteristics(service bluetooth.DeviceService, uuids []bluetooth.UUID) []bluetooth.DeviceCharacteristic {

	characteristics, err := service.DiscoverCharacteristics(uuids)
	if err != nil {
		connectTimer.Stop()
		rcLog.ReportErrorAndExit(rcErrors.ERROR_CODE_BLE, "Could not get characteristic list for service %s: %s", service.UUID().String(), err.Error())
	}

	return characteristics
}

func RequestDataViaCommandUART(ble bluetooth.Device, requestPacket []byte, callback func([]byte), maxWrites uint) {

	// Get the characteristics within the UART service
	characteristics := ReadyCommandUART(ble)
	noResponseExpected := false
	var writeCount uint = 0

	// Do we need to await responses?
	if maxWrites == 0 {
		// Increment so that the command is sent...
		maxWrites = 1

		// ...but flag that we can exit without waiting for a response
		noResponseExpected = true
	} else {
		// Enable notifications via RX
		err := characteristics[1].EnableNotifications(callback)
		if err != nil {
			connectTimer.Stop()
			rcLog.ReportErrorAndExit(rcErrors.ERROR_CODE_BLE, "Could not enable UART notifications: %s", err.Error())
		}
	}

	// Send the command `maxWrites` times
	for writeCount < maxWrites {
		// Clear the 'data received' flag
		UARTInfoReceived = false

		// Request data via the TX
		_, err := characteristics[0].WriteWithoutResponse(requestPacket)
		if err != nil {
			rcLog.ReportErrorAndExit(rcErrors.ERROR_CODE_BLE, "Could not write UART packet: %s", err.Error())
		}

		// No response coming? Bail early
		if noResponseExpected {
			return
		}

		// Wait (block) for the command response packet(s)
		for !UARTInfoReceived {
			// NOP
		}

		// And that's one write done...
		writeCount += 1
	}
}

func ReadyCommandUART(ble bluetooth.Device) []bluetooth.DeviceCharacteristic {

	// Get the UART service
	uuid := UUIDFromString(rcColmi.UUID_BLE_COMMAND_UART_SERVICE)
	service := Services(ble, []bluetooth.UUID{uuid})

	// Get the characteristics within the UART service
	tx := UUIDFromString(rcColmi.UUID_BLE_COMMAND_UART_TX_CHAR)
	rx := UUIDFromString(rcColmi.UUID_BLE_COMMAND_UART_RX_CHAR)
	return Characteristics(service[0], []bluetooth.UUID{tx, rx})
}

func ReadyDataUART(ble bluetooth.Device) []bluetooth.DeviceCharacteristic {

	// Get the UART service
	uuid := UUIDFromString(rcColmi.UUID_BLE_DATA_UART_SERVICE)
	service := Services(ble, []bluetooth.UUID{uuid})

	// Get the characteristics within the UART service
	tx := UUIDFromString(rcColmi.UUID_BLE_DATA_UART_TX_CHAR)
	rx := UUIDFromString(rcColmi.UUID_BLE_DATA_UART_RX_CHAR)
	return Characteristics(service[0], []bluetooth.UUID{tx, rx})
}

func DeviceInfoService(bleDevice bluetooth.Device) bluetooth.DeviceService {

	uuid := UUIDFromUInt16(BLE_DEVICE_INFO_SERVICE_ID)
	services := Services(bleDevice, []bluetooth.UUID{uuid})
	return services[0]
}

func BeginScan(adapter *bluetooth.Adapter, callback func(*bluetooth.Adapter, bluetooth.ScanResult)) {

	if adapter.Scan(callback) != nil {
		rcLog.ReportErrorAndExit(rcErrors.ERROR_CODE_BLE, "Failed to initiate scan for rings")
	}
}

func AddressFromString(ringAddress string) bluetooth.Address {

	var bleAddress bluetooth.Address
	bleAddress.Set(ringAddress)
	return bleAddress
}

func UUIDFromUInt16(uuid uint16) bluetooth.UUID {

	return bluetooth.New16BitUUID(uuid)
}

func UUIDFromString(uuid string) bluetooth.UUID {

	convertedUUID, err := bluetooth.ParseUUID(uuid)
	if err != nil {
		connectTimer.Stop()
		rcLog.ReportErrorAndExit(rcErrors.ERROR_CODE_BLE, "Could not convert UUID: %s", err.Error())
	}

	return convertedUUID
}