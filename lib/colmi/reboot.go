package rcColmi

func MakeRebootReq() []byte {

	return MakePacket(COMMAND_REBOOT, make([]byte, 0, 0))
}
