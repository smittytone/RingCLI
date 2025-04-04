package ringcliColmi

func MakeLedFlashRequest() []byte {

	return MakePacket(COMMAND_BATTERY_FLASH_LED, make([]byte, 0, 0))
}
