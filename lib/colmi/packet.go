package rcColmi

import (
	rcLog "ringcli/lib/log"
)

func MakePacket(command byte, data []byte) []byte {

	packet := make([]byte, 16, 16)
	packet[0] = command

	if len(data) > 0 {
		if len(data) > 14 {
			rcLog.ReportErrorAndExit(3, "Colmi packet payload must be 14 bytes or less")
		}

		for i := range len(data) {
			packet[i+1] = data[i]
		}
	}

	packet[15] = checksum(packet)
	return packet
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
