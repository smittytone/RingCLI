package ringcliColmi

import (
	"fmt"
	"time"
	utils "ringcli/lib/utils"
)

type BloodOxygenDatapoint struct {
	Maximum   int
	Minimum   int
	Timestamp string
}

type BloodOxygenData struct {
	Rates     []BloodOxygenDatapoint
	Raw       []byte
	Time      time.Time
	DaysPast  int
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
	*/

	if packet[1] == DATA_REQUEST_ID_OXYGEN {
		maxHoursPrevious := time.Duration(packet[6]) * time.Hour * -24
		startTime := utils.StartToday(time.Now()).Add(maxHoursPrevious).UTC()
		hour := startTime.Hour()

		// Max data length, including 'days previous' two-byte markers
		dataLength := int(packet[3]) << 8 | int(packet[2])

		// Ignore CRC for now
		// crc :=int(packet[4]) << 8 | int(packet[5])
		var crc uint16 = 0
		for i := range dataLength {
			crc += uint16(packet[6 + i])
		}

		// Instantiate the return struct
		data := BloodOxygenData{
			Rates: make([]BloodOxygenDatapoint, 0, dataLength >> 1),
			Raw: packet,
			Time: startTime,
			DaysPast: int(packet[7]),
		}

		index := 6
		for index < len(packet) {
			// Step over days byte
			index += 1
			for j := range 24 {
				point := BloodOxygenDatapoint{
					Maximum: int(packet[index + (j * 2)]),
					Minimum: int(packet[index + 1 + (j * 2)]),
					Timestamp: fmt.Sprintf("%02d:00", hour),
				}

				data.Rates = append(data.Rates, point)
				hour = (hour + 1) % 24
			}

			index += 48
		}

		return &data
	}

	return nil
}
