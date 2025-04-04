package ringcliColmi

import (
	"errors"
	rcErrors "ringcli/lib/errors"
	log "ringcli/lib/log"
	utils "ringcli/lib/utils"
	"time"
)

func MakeTimeSetRequest(targetDate time.Time) []byte {

	payload := makeTimeReqPayload(targetDate)
	return MakePacket(COMMAND_SET_TIME, payload)
}

func makeTimeReqPayload(targetDate time.Time) []byte {

	timezone, _ := targetDate.Zone()
	if timezone != "UTC" {
		targetDate = targetDate.UTC()
	}

	if targetDate.Year() <= 2000 {
		log.ReportErrorAndExit(rcErrors.ERROR_CODE_BAD_SET_TIME, "Supplied time out of range (must be 2001 or greater)")
	}

	payload := make([]byte, 7, 7)
	payload[0] = makePayloadByte(targetDate.Year() % 2000)
	payload[1] = makePayloadByte(int(targetDate.Month()))
	payload[2] = makePayloadByte(targetDate.Day())
	payload[3] = makePayloadByte(targetDate.Hour())
	payload[4] = makePayloadByte(targetDate.Minute())
	payload[5] = makePayloadByte(targetDate.Second())
	payload[6] = LANGUAGE_ENGLISH
	return payload
}

func makePayloadByte(value int) byte {

	var ringError *rcErrors.RingcliError
	result, err := utils.ToBCD(value)
	if err != nil && errors.As(err, &ringError) {
		log.ReportErrorAndExit(ringError.Code, ringError.Message)
	}

	return result
}
