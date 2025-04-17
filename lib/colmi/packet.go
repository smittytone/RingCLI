package ringcliColmi

import (
	log "ringcli/lib/log"
)

func MakePacket(command byte, data []byte) []byte {

	packet := make([]byte, 16, 16)
	packet[0] = command

	if len(data) > 0 {
		if len(data) > 14 {
			log.ReportErrorAndExit(3, "Colmi packet payload must be 14 bytes or less")
		}

		for i := range len(data) {
			packet[i+1] = data[i]
		}
	}

	packet[15] = checksum(packet)
	return packet
}

func MakeStartPacket(command byte) []byte {

	return MakePacket(COMMAND_START_REAL_TIME, []byte{command, REAL_TIME_ACTION_START})
}

func MakeContinuePacket(command byte) []byte {

	return MakePacket(COMMAND_START_REAL_TIME, []byte{command, REAL_TIME_ACTION_CONT})
}

func MakeStopPacket(command byte) []byte {

	return MakePacket(COMMAND_STOP_REAL_TIME, []byte{command})
}

func checksum(packet []byte) byte {

	// Add all the bytes together % 255
	var count byte = 0
	for _, aByte := range packet {
		count += aByte
	}

	return count
}

func VerifyChecksum(packet []byte) bool {

	chk := checksum(packet[0:15])
	return chk == packet[15]
}

func MakeDataPacket(command byte) []byte {

	packet := make([]byte, 6, 6)
	packet[0] = DATA_REQUEST_MAGIC_VALUE
	packet[1] = command
	packet[2] = 0x00
	packet[3] = 0x00
	packet[4] = 0xFF
	packet[5] = 0xFF
	return packet
}
