package ringcliColmi

type BatteryInfo struct {
	Level      int  // Battery level as a percentage
	IsCharging bool // Battery charging state
}

func MakeBatteryRequest() []byte {

	return MakePacket(COMMAND_BATTERY_INFO, make([]byte, 0, 0))
}

func ParseBatteryResponse(packet []byte) BatteryInfo {

	return BatteryInfo{Level: int(packet[1]), IsCharging: (packet[2] != 0)}
}
