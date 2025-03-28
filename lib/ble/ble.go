package rcBLE

import (
	"time"

	"tinygo.org/x/bluetooth"

	rcErrors "ringcli/lib/errors"
	rcLog "ringcli/lib/log"
)

const (
	CONNECT_TIMEOUT_S             = 60
	SYNC_TIMEOUT_S                = 10
	DEVICE_INFO_SERVICE_ID uint16 = 0x180A
)

var (
	connectTimer *time.Timer
	syncTimer *time.Timer
	UARTInfoReceived bool = false
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

	return device
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

func RequestDataViaUART(ble bluetooth.Device, requestPacket []byte, callback func([]byte)) {

	// Get the UART service
	uuid := UUIDFromString("6E40FFF0-B5A3-F393-E0A9-E50E24DCCA9E")
	service := Services(ble, []bluetooth.UUID{uuid})

	// Get the characteristics within the UART service
	tx := UUIDFromString("6E400002-B5A3-F393-E0A9-E50E24DCCA9E")
	rx := UUIDFromString("6E400003-B5A3-F393-E0A9-E50E24DCCA9E")
	characteristics := Characteristics(service[0], []bluetooth.UUID{tx, rx})

	// Enable notifications via RX
	UARTInfoReceived = false
	err := characteristics[1].EnableNotifications(callback)
	if err != nil {
		connectTimer.Stop()
		rcLog.ReportErrorAndExit(rcErrors.ERROR_CODE_BLE, "Could not enable UART notifications: %s", err.Error())
	}

	// Request data via the TX
	_, err = characteristics[0].WriteWithoutResponse(requestPacket)
	if err != nil {
		rcLog.ReportErrorAndExit(rcErrors.ERROR_CODE_BLE, "Could not write UART packet: %s", err.Error())
	}

	for !UARTInfoReceived {
		// Wait for the return packet(s)...
	}
}

func DeviceInfoService(bleDevice bluetooth.Device) bluetooth.DeviceService {

	uuid := UUIDFromUInt16(DEVICE_INFO_SERVICE_ID)
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