package ringCLI_Colmi

func MakeShutdownRequest() []byte {

	return MakePacket(COMMAND_SHUTDOWN, []byte{0x01})
}
