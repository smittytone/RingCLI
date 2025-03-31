package rcColmi

func MakeHeartRatePeriodGetReq() []byte {

	return MakePacket(COMMAND_HEART_RATE_PERIOD, []byte{0x01})
}

func MakeHeartRatePeriodSetReq(isEnabled bool, period byte) []byte {

	// Set the byte value for the enabling flag
	var enabledValue byte = 1
	if !isEnabled {
		enabledValue = 2
	}

	payload := []byte{2, enabledValue, period}
	return MakePacket(COMMAND_HEART_RATE_PERIOD, payload)
}

func ParseHeartRatePeriodResp(packet []byte) (bool, byte) {

	/* SAMPLE SET
	[22 2 1 0 0 0 0 0 0 0 0 0 0 0 0 25]

	   SAMPLE GET
	[22 1 1 30 5 0 0 0 0 0 0 0 0 0 0 59]
	*/

	if packet[0] == COMMAND_HEART_RATE_PERIOD {
		enabled := false
		if packet[2] == 1 {
			enabled = true
		}

		return enabled, packet[3]
	}

	return false, 0
}
