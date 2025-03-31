package rcColmi

func MakeShutdownReq() []byte {

	return MakePacket(COMMAND_SHUTDOWN, []byte{0x01})
}
