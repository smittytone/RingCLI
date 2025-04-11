package ringcliColmi

import (
	"fmt"
	utils "ringcli/lib/utils"
	"time"
)

type BloodOxygenDatapoint struct {
	Maximum   int
	Minimum   int
	Timestamp string
	Time      time.Time
}

type BloodOxygenDataSet struct {
	Rates []BloodOxygenDatapoint
	Time  time.Time
	Index int
}

type BloodOxygenData struct {
	Data []BloodOxygenDataSet
	Time time.Time
	Raw  []byte
}

func MakeBloodOxygenGetRequest() []byte {

	return MakeDataPacket(DATA_REQUEST_ID_OXYGEN)
}

func ParseBloodOxygenDataResponse(packet []byte) *BloodOxygenData {

	/* SAMPLE
	   Header: 188 (Magic) 42 (ID) 98 (data length LSB) 00 (data length MSB) 183 (CRC LSB) 63 (CRC MSB)
	   Data:   02 - 2 days ago
	           00 00 00 00 00 00 00 00 00 00 00 00 00 00 93 93 96 96 96 96 96 96 99 99 96 96 96 96 97 97 99 99 96 96 98 98 98 98 00 00 00 00 00 00 99 99 00 00
	           01  - 1 day ago
	           00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 99 99 00 00 99 99 98 98 97 97 97 97 97 97 97 97 99 99 99 99 99 99 97 97 96 96 98 98 00 00 00 00

	   Later:
	           02
	           00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 99 99 00 00 99 99 98 98 97 97 97 97 97 97 97 97 99 99 99 99 99 99 97 97 96 96 98 98 00 00 00 00
	           01
	           00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 96 96 00 00 00 00 96 96 96 96 00 00 98 98 96 96 99 99 96 96 97 97 92 92 92 92 96 96
	           00
	           96 96 90 90 96 96 98 98 97 97 98 98 98 98 96 96 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	*/

	if packet[1] == DATA_REQUEST_ID_OXYGEN {
		// Max data length, including 'days previous' two-byte markers
		dataLength := int(packet[3])<<8 | int(packet[2])
		dataCount := dataLength / 49

		// Get data time range
		maxHoursPrevious := time.Duration(packet[6]) * time.Hour * -24
		startTime := utils.StartToday(time.Now()).Add(maxHoursPrevious).UTC()
		hour := startTime.Hour()

		// Ignore CRC for now
		// crc :=int(packet[4]) << 8 | int(packet[5])
		var crc uint16 = 0
		for i := range dataLength {
			crc += uint16(packet[6+i])
		}

		// Instantiate the return struct
		data := BloodOxygenData{
			Data: make([]BloodOxygenDataSet, 0, dataCount),
			Raw:  packet[7:],
			Time: startTime,
			//DayCount: int(packet[7]),
		}

		index := 6
		for index < len(packet) {
			// Get over count byte
			count := int(packet[index])
			index += 1

			// Get data time range
			hoursPrevious := time.Duration(count) * time.Hour * -24
			dataTime := utils.StartToday(time.Now()).Add(hoursPrevious).UTC()
			hour = startTime.Hour()

			set := BloodOxygenDataSet{
				Index: count,
				Rates: make([]BloodOxygenDatapoint, 0, 24),
				Time:  dataTime,
			}

			for j := range 24 {
				point := BloodOxygenDatapoint{
					Maximum:   int(packet[index+(j*2)]),
					Minimum:   int(packet[index+1+(j*2)]),
					Timestamp: fmt.Sprintf("%02d:00", hour),
					Time:      dataTime,
				}

				set.Rates = append(set.Rates, point)
				hour = (hour + 1) % 24
				dataTime = dataTime.Add(time.Hour)
			}

			data.Data = append(data.Data, set)
			index += 48
		}

		return &data
	}

	return nil
}
