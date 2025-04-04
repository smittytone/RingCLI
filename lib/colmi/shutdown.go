package ringcliColmi

func MakeShutdownRequest() []byte {

	return MakePacket(COMMAND_SHUTDOWN, []byte{0x01})
}
