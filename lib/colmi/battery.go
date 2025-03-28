package rcColmi

type BatteryInfo struct {
	Level int
	IsCharging bool
}

func MakeBatteryReq() []byte {

	return MakePacket(0x03, make([]byte, 0, 0))
}

func ParseBatteryResp(packet []byte) BatteryInfo {

	return BatteryInfo{Level: int(packet[1]), IsCharging: (packet[2] != 0)}
}
