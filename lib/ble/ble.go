package ringcliBLE

import (
	//"fmt"
	ring "ringcli/lib/colmi"
	errors "ringcli/lib/errors"
	log "ringcli/lib/log"
	"time"
	"tinygo.org/x/bluetooth"
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
	BLE_GAP_SERVICE_ID                               uint16 = 0x1800
	BLE_GAP_SERVICE_DEVICE_NAME_CHAR_ID              uint16 = 0x2A00
)

var (
	connectTimer     *time.Timer
	syncTimer        *time.Timer
	UARTInfoReceived bool = false
	isConnected      bool = false
	currentDevice    *bluetooth.Device
	PollCount        int = 0
)

func Open() *bluetooth.Adapter {

	ble := bluetooth.DefaultAdapter
	if ble.Enable() != nil {
		log.ReportErrorAndExit(errors.ERROR_CODE_BLE, "Failed to enable BLE")
	}

	return ble
}

func Connect(adapter *bluetooth.Adapter, ringAddress bluetooth.Address) bluetooth.Device {

	// Establish a timer so we're not trying to connect forever
	connectTimer = time.NewTimer(CONNECT_TIMEOUT_S * time.Second)
	defer connectTimer.Stop()
	go func() {
		<-connectTimer.C
		log.ReportErrorAndExit(errors.ERROR_CODE_BLE, "Could not connect to %s, or connection timed out", ringAddress.String())
	}()

	// Attempt to connect
	device, err := adapter.Connect(ringAddress, bluetooth.ConnectionParams{})
	if err != nil {
		connectTimer.Stop()
		log.ReportErrorAndExit(errors.ERROR_CODE_BLE, "Could not connect to %s", ringAddress)
	}

	isConnected = true
	currentDevice = &device
	return device
}

func Disconnect(device bluetooth.Device) {

	err := device.Disconnect()
	if err != nil {
		log.ReportErrorAndExit(errors.ERROR_CODE_BLE, "Could not cleanly disconnect from %s", device.Address.String())
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
		log.ClearPrompt()
		log.ReportErrorAndExit(errors.ERROR_CODE_BLE, "Could not get ring BLE service list: %s", err.Error())
	}

	return services
}

func Characteristics(service bluetooth.DeviceService, uuids []bluetooth.UUID) []bluetooth.DeviceCharacteristic {

	characteristics, err := service.DiscoverCharacteristics(uuids)
	if err != nil {
		// FROM 0.1.5
		// Do not exit on error, instead return an empty array to inform
		// the caller that one of the characteristics could not be found.
		// Callers should issue characteristic UUIDs one at a time to learn
		// which specifics UUIDs are not supported by the ring firmware
		return []bluetooth.DeviceCharacteristic{}
	}

	return characteristics
}

func RequestDataViaCommandUART(device bluetooth.Device, requestPacket []byte, callback func([]byte), maxWrites uint) {

	// Get the characteristics within the UART service
	characteristics := ReadyCommandUART(device)
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
			log.ReportErrorAndExit(errors.ERROR_CODE_BLE, "Could not enable UART notifications: %s", err.Error())
		}
	}

	// Send the command `maxWrites` times
	for writeCount < maxWrites {
		// Clear the 'data received' flag
		UARTInfoReceived = false

		// Request data via the TX
		_, err := characteristics[0].WriteWithoutResponse(requestPacket)
		if err != nil {
			log.ReportErrorAndExit(errors.ERROR_CODE_BLE, "Could not write UART packet: %s", err.Error())
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

func RequestPollViaCommandUART(device bluetooth.Device, requestPacket []byte, callback func([]byte)) []bluetooth.DeviceCharacteristic {

	// Get the characteristics within the UART service
	characteristics := ReadyCommandUART(device)

	// Enable notifications via RX
	err := characteristics[1].EnableNotifications(callback)
	if err != nil {
		connectTimer.Stop()
		log.ReportErrorAndExit(errors.ERROR_CODE_BLE, "Could not enable UART notifications: %s", err.Error())
	}

	// Request data via the TX
	_, err = characteristics[0].WriteWithoutResponse(requestPacket)
	if err != nil {
		log.ReportErrorAndExit(errors.ERROR_CODE_BLE, "Could not write UART packet: %s", err.Error())
	}

	return characteristics
}

func EndPollViaCommandUART(device bluetooth.Device, haltPacket []byte, characteristic bluetooth.DeviceCharacteristic) {

	// Halt polling via the TX
	_, err := characteristic.WriteWithoutResponse(haltPacket)
	if err != nil {
		log.ReportErrorAndExit(errors.ERROR_CODE_BLE, "Could not write UART packet: %s", err.Error())
	}
}

func ReadyCommandUART(device bluetooth.Device) []bluetooth.DeviceCharacteristic {

	// Get the UART service
	uuid := UUIDFromString(ring.UUID_BLE_COMMAND_UART_SERVICE)
	service := Services(device, []bluetooth.UUID{uuid})

	// Get the characteristics within the UART service
	tx := UUIDFromString(ring.UUID_BLE_COMMAND_UART_TX_CHAR)
	rx := UUIDFromString(ring.UUID_BLE_COMMAND_UART_RX_CHAR)
	return Characteristics(service[0], []bluetooth.UUID{tx, rx})
}

func RequestDataViaDataUART(device bluetooth.Device, requestPacket []byte, callback func([]byte)) {

	// Get the characteristics within the UART service
	characteristics := ReadyDataUART(device)

	// Enable notifications via RX
	err := characteristics[1].EnableNotifications(callback)
	if err != nil {
		connectTimer.Stop()
		log.ReportErrorAndExit(errors.ERROR_CODE_BLE, "Could not enable UART notifications: %s", err.Error())
	}

	// Request data via the TX
	_, err = characteristics[0].WriteWithoutResponse(requestPacket)
	if err != nil {
		log.ReportErrorAndExit(errors.ERROR_CODE_BLE, "Could not write UART packet: %s", err.Error())
	}

	// Wait (block) for the command response packet(s)
	for !UARTInfoReceived {
		// NOP
	}
}

func ReadyDataUART(device bluetooth.Device) []bluetooth.DeviceCharacteristic {

	// Get the UART service
	uuid := UUIDFromString(ring.UUID_BLE_DATA_UART_SERVICE)
	service := Services(device, []bluetooth.UUID{uuid})

	// Get the characteristics within the UART service
	tx := UUIDFromString(ring.UUID_BLE_DATA_UART_TX_CHAR)
	rx := UUIDFromString(ring.UUID_BLE_DATA_UART_RX_CHAR)
	return Characteristics(service[0], []bluetooth.UUID{tx, rx})
}

func DeviceInfoService(bleDevice bluetooth.Device) bluetooth.DeviceService {

	uuid := UUIDFromUInt16(BLE_DEVICE_INFO_SERVICE_ID)
	services := Services(bleDevice, []bluetooth.UUID{uuid})

	if len(services) == 0 {
		log.ReportErrorAndExit(1, "No GAP service")
	}

	return services[0]
}

func BeginScan(adapter *bluetooth.Adapter, callback func(*bluetooth.Adapter, bluetooth.ScanResult)) {

	if adapter.Scan(callback) != nil {
		log.ReportErrorAndExit(errors.ERROR_CODE_BLE, "Failed to initiate scan for rings")
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
		log.ReportErrorAndExit(errors.ERROR_CODE_BLE, "Could not convert UUID: %s", err.Error())
	}

	return convertedUUID
}

func PollRealtime(bleDevice bluetooth.Device, readingType byte, callback func([]byte), count int) {

	PollCount = 0

	// Make the initial poll data request
	characteristics := RequestPollViaCommandUART(bleDevice, ring.MakeStartPacket(readingType), callback)

	// Pause while we receive the data
	cont := false
	for PollCount < count {
		// Only if `readingType` is 0x01 and we want more than 10 readings, eg. to average over 30s
		// Otherwise (if `readingType` is 0x06) we just NOP to wait while `PollCount` updated elsewhere
		if readingType == 0x01 {
			if PollCount != 0 && !cont && PollCount%10 == 0 {
				_, err := characteristics[0].WriteWithoutResponse(ring.MakeContinuePacket(readingType))
				if err != nil {
					log.ReportErrorAndExit(errors.ERROR_CODE_BLE, "Could not write UART packet: %s", err.Error())
				}

				cont = true
			}

			if PollCount == 11 || PollCount == 21 {
				cont = false
			}
		}
	}

	// Issue the 'stop sending' request
	// NOTE This appears unsuccessful if `readingType` is 0x06
	//      When `readingType` is 0x01, this appears unnecessary - ring only sends 10 readings
	EndPollViaCommandUART(bleDevice, ring.MakeStopPacket(readingType), characteristics[0])
}
