package rcColmi

import (
	"encoding/binary"
	"fmt"
	"time"
	// App
	rcLog "ringcli/lib/log"
)

type HeartRateDatapoint struct {
	Timestamp string
	Rate int
}

type HeartRateData struct {
	Rates []HeartRateDatapoint
	DataRange int
	Raw []byte
	Timestamp time.Time
}

var (
	lastPacket int = 0
	packetRange int = 0
	initialTime time.Time
	hrd HeartRateData
	data []byte
)

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

func MakeHeartRateReadReq(target time.Time) []byte {

	payload := make([]byte, 4, 4)
	timestamp := uint32(target.UnixMilli() / 1000)
	binary.LittleEndian.PutUint32(payload[0:], timestamp)
	return MakePacket(COMMAND_HEART_RATE_READ, payload)
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

func ParseHeartRateDataResp(packet []byte) *HeartRateData {

	/* SAMPLE
	[21 0 24 5 0 0 0 0 0 0 0 0 0 0 0 50]
	[21 1 212 148 234 103 64 0 0 0 0 0 0 0 0 15]
	[21 2 0 0 0 63 0 0 0 0 0 0 0 0 0 86]
	[21 3 0 0 63 0 0 0 0 0 0 0 0 0 0 87]
	[21 4 0 0 0 0 0 0 0 0 0 0 0 0 0 25]
	[21 5 0 0 0 0 0 0 0 0 0 0 0 0 0 26]
	[21 6 0 0 0 0 0 0 0 0 0 0 0 0 0 27]
	[21 7 0 0 0 0 0 0 0 0 0 0 0 0 0 28]
	[21 8 0 0 0 0 0 0 0 0 0 0 0 0 0 29]
	[21 9 0 0 0 0 0 0 0 0 0 0 0 0 0 30]
	[21 10 0 0 0 0 0 0 0 80 0 0 0 0 0 111]
	[21 11 0 0 0 0 0 0 86 0 0 0 0 0 0 118]
	[21 12 0 0 0 0 0 87 0 0 0 0 0 0 0 120]
	[21 13 0 0 0 0 81 0 0 0 0 0 89 0 0 204]
	[21 14 0 0 0 82 0 0 0 0 0 0 0 0 0 117]
	[21 15 0 0 0 0 0 0 0 0 0 0 0 0 0 36]
	[21 16 0 0 0 0 0 0 0 0 0 0 0 0 0 37]
	[21 17 0 0 0 0 0 0 0 0 0 0 0 0 0 38]
	[21 18 0 0 0 0 0 0 0 0 0 0 0 0 0 39]
	[21 19 0 0 0 0 0 0 0 0 0 0 0 0 0 40]
	[21 20 0 0 0 0 0 0 0 0 0 0 0 0 0 41]
	[21 21 0 0 0 0 0 0 0 0 0 0 0 0 0 42]
	[21 22 0 0 0 0 0 0 0 0 0 0 0 0 0 43]
	[21 23 0 0 0 0 0 0 0 0 0 0 0 0 0 44]
	*/

	if packet[0] == 0x15 {
		if packet[1] == 0xFF {
			// ERROR
			rcLog.ReportError("Input heartrate data packet malformed")
			return nil
		}

		if !VerifyChecksum(packet) {
			rcLog.ReportError("Checkum fail")
		}

		// Header packet
		if packet[1] == 0 && packetIndex == 0 {
			lastPacket = int(packet[2]) - 1
			packetRange = int(packet[3]) // What is this for???
			data = make([]byte, 0, 255)
			//rates = make([]HeartRateCount, 0, packet[2] - 2)
			packetIndex += 1
			return nil
		}

		// Timestamp + Data packet
		if packetIndex == 1 {
			// First four bytes form a timestamp, rest is data TODO
			ts := (int(packet[5]) << 24) | (int(packet[4]) << 16) | (int(packet[3]) << 8) | int(packet[2])
			initialTime = time.UnixMilli(int64(ts) * 1000)
			data = append(data, packet[6:15]...)
			packetIndex +=1
			return nil
		}

		// Data packet
		data = append(data, packet[2:15]...)
		packetIndex += 1

		if int(packet[1]) == lastPacket {
			// Return data
			hrd = HeartRateData{
				Rates: packageData(initialTime),
				DataRange: packetRange,
				Raw: data,
				Timestamp: initialTime,
			}

			return &hrd
		}
	}

	return nil
}

func packageData(startTime time.Time) []HeartRateDatapoint {

	count := (len(data) / 6) / 5
	minuteDelta := 60 / count
	fmt.Println(count, minuteDelta)
	results := make([]HeartRateDatapoint, 0, 24)
	done := false
	index := 0
	hour := 0
	min := 0
	for !done {
		var hrdp  HeartRateDatapoint = HeartRateDatapoint{
			Rate: int(data[index]),
			Timestamp: fmt.Sprintf("%02d:%02d", hour, min),
		}

		results = append(results, hrdp)

		min += minuteDelta
		if min >= 60 {
			min = 60 - min
			hour += 1
		}

		index += 6
		done = index >= len(data)
	}

	return results
}