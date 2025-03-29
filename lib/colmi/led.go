package rcColmi

func MakeLedFlashReq() []byte {

	return MakePacket(0x10, make([]byte, 0, 0))
}
