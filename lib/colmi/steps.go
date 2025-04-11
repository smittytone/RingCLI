package ringcliColmi

import (
	"time"
)

type DateStamp struct {
	Year      int
	Month     int
	Day       int
	Hour      int
	Minutes   int
	TimePhase int // A 15 minute segment within a day
}

type ActivityInfo struct {
	Timestamp   DateStamp
	Steps       int
	Distance    int
	Calories    int
	NoData      bool
	IsDone      bool
	NewCalories bool
}

var (
	packetIndex     int  = 0
	doScaleCalories bool = false
)

func MakeStepsRequest(offset int) []byte {

	// Make sure the offset is within byte range
	if offset < -128 || offset > 127 {
		return nil
	}

	return MakePacket(COMMAND_GET_ACTIVITY_DATA, []byte{byte(offset), 0x0F, 0x00, 0x5F, 0x01})
}

func ParseStepsResponse(packet []byte) ActivityInfo {

	/* SAMPLE
	[67  240 12 1 0 0 0 0 0 0 0 0 0 0 0 64]
	[67  37 3 40 0 0 12 142 0 32 0 20 0 0 0 97]
	[67  37 3 40 24 1 12 87 0 15 0 12 0 0 0 42]
	[67  37 3 40 28 2 12 95 2 109 0 87 0 0 0 226]
	[67  37 3 40 32 3 12 87 2 119 0 86 0 0 0 232]
	[67  37 3 40 36 4 12 153 0 31 0 21 0 0 0 148]
	[67  37 3 40 40 5 12 100 2 152 0 88 0 0 0 34]
	[67  37 3 40 44 6 12 129 4 20 1 165 0 0 0 16]
	[67  37 3 40 48 7 12 237 7 198 1 34 1 0 0 180]
	[67  37 3 40 52 8 12 184 10 130 2 137 1 0 0 171]
	[67  37 3 40 56 9 12 144 3 210 0 131 0 0 0 200]
	[67  37 3 40 60 10 12 168 3 200 0 134 0 0 0 222]
	[67  37 3 40 64 11 12 69 0 15 0 9 0 0 0 71]
	*/

	// Nothing will be coming - no heart rate taken yet as the ring has not been worn yet
	if packet[1] == 0xFF && packet[5] == 0 {
		return ActivityInfo{NoData: true}
	}

	// This seems to be a header signal packet included to set the calorie scale factor
	if packet[1] == 0xF0 && packet[5] == 0 {
		if packet[3] == 0x01 {
			doScaleCalories = true
		}

		// No valid data yet so return an empty struct
		return ActivityInfo{}
	}

	// Assemble data from the packet
	timeIndex := int(packet[4])
	packetTime := DateStamp{
		Year:      2000 + bcdToDecimal(int(packet[1])),
		Month:     bcdToDecimal(int(packet[2])),
		Day:       bcdToDecimal(int(packet[3])),
		Hour:      timeIndex / 4,
		Minutes:   timeIndex % 4 * 15,
		TimePhase: timeIndex,
	}

	info := ActivityInfo{
		Timestamp: packetTime,
		Steps:     int(packet[10])<<8 + int(packet[9]),
		Distance:  int(packet[12])<<8 + int(packet[11]),
		Calories:  int(packet[8])<<8 + int(packet[7]),
		NoData:    false,
		IsDone:    false,
	}

	if doScaleCalories {
		info.Calories *= 10
		info.NewCalories = true
	}

	// Packet management
	currentPacket := int(packet[5])
	maxPacket := int(packet[6]) - 1

	// Check for end of block
	info.IsDone = (currentPacket == maxPacket)

	// Return the data for accrual
	return info
}

func bcdToDecimal(bcd int) int {

	return (((bcd >> 4) & 15) * 10) + (bcd & 15)
}

func TimestampFromNow() DateStamp {

	now := time.Now()
	stamp := DateStamp{
		Year:      now.Year(),
		Month:     int(now.Month()),
		Day:       now.Day(),
		Hour:      now.Hour(),
		Minutes:   now.Minute(),
		TimePhase: 0,
	}

	return stamp
}
