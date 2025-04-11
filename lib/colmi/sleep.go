package ringcliColmi

import (
	//"fmt"
	"time"
	utils "ringcli/lib/utils"
)

type SleepDatapoint struct {
	Type     int
	Duration int
}

type SleepPeriod struct {
	Sleep     []SleepDatapoint
	StartTime time.Time
	EndTime   time.Time
}

type SleepData struct {
	Data []SleepPeriod
	Time time.Time
	Raw  []byte
}

func MakeSleepGetRequest() []byte {

	return MakeDataPacket(DATA_REQUEST_ID_SLEEP)
}

func ParseSleepDataResponse(packet []byte) *SleepData {

	/* SAMPLE
        Header: 188 (Magic) 39 (ID) 39 (data length LSB) 00 (data length MSB) 216 (CRC LSB) 128 (CRC MSB) 1 (days' data available)
        Day Data:	000 (days previous) 036 (day record bytes)
					119 005 (sleep start mins past midnight)
					201 001 (sleep end mins past midnight)

					Sleep Data:
					002 067 005 006 002 064 003 032 002 016 004 032
					002 096 003 016 002 032 004 032 002 032 003 016
					002 016 004 016 005 012 002 013
	*/



	if packet[1] == DATA_REQUEST_ID_SLEEP {
		// Max data length, including 'days previous' two-byte markers
		dataLength := int(packet[3]) << 8 | int(packet[2])

		// Ignore CRC for now
		// crc :=int(packet[4]) << 8 | int(packet[5])
		var crc uint16 = 0
		for i := range dataLength {
			crc += uint16(packet[6 + i])
		}

		// Instantiate the return struct
		days := int(packet[6])
		data := SleepData{
			Data: make([]SleepPeriod, 0, days),
			Raw: packet[6:],
			Time: time.Now(),
		}

		index := 7
		for index < len(packet) {
			// Get data time range
			hoursPrevious := time.Duration(packet[index]) * time.Hour * -24
			midnightTime := utils.StartToday(time.Now()).Add(hoursPrevious).UTC()

			// Get the number of data points
			dataCount := (int(packet[index + 1]) - 4) >> 1
			//fmt.Println(index,"/",len(packet), days, dataCount)

			start := int(packet[index + 3]) << 8 | int(packet[index + 2])
			startMins := time.Duration(start) * time.Minute
			end := int(packet[index + 5]) << 8 | int(packet[index + 4])
			endMins := time.Duration(end) * time.Minute

			period := SleepPeriod{
				StartTime: midnightTime.Add(startMins),
				EndTime:midnightTime.Add(endMins),
				Sleep: make([]SleepDatapoint, 0, dataCount),
			}

			index += 6
			//fmt.Println(index,"/",len(packet))
			for _ = range dataCount {
				point := SleepDatapoint{
					Type: int(packet[index]),
					Duration: int(packet[index + 1]),
				}

				period.Sleep = append(period.Sleep, point)
				index += 2
			}

			data.Data = append(data.Data, period)
		}

		return &data
	}

	return nil
}

func GetSleepType(sleepType int) string {

	switch sleepType {
	case SLEEP_TYPE_NO_DATA:
		return SLEEP_STRING_NO_DATA
	case SLEEP_TYPE_LIGHT:
		return SLEEP_STRING_LIGHT
	case SLEEP_TYPE_DEEP:
		return SLEEP_STRING_DEEP
	case SLEEP_TYPE_REM:
		return SLEEP_STRING_REM
	case SLEEP_TYPE_AWAKE:
		return SLEEP_STRING_AWAKE
	default:
		return SLEEP_STRING_ERROR
	}
}