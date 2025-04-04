package ringCLI_Colmi

import (
	rcErrors "ringcli/lib/errors"
	rcLog "ringcli/lib/log"
	rcUtils "ringcli/lib/utils"
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
		rcLog.ReportErrorAndExit(rcErrors.ERROR_CODE_BAD_SET_TIME, "Supplied time out of range (must be 2001 or greater)")
	}

	payload := make([]byte, 7, 7)
	payload[0] = rcUtils.ToBCD(targetDate.Year() % 2000)
	payload[1] = rcUtils.ToBCD(int(targetDate.Month()))
	payload[2] = rcUtils.ToBCD(targetDate.Day())
	payload[3] = rcUtils.ToBCD(targetDate.Hour())
	payload[4] = rcUtils.ToBCD(targetDate.Minute())
	payload[5] = rcUtils.ToBCD(targetDate.Second())
	payload[6] = LANGUAGE_ENGLISH
	return payload
}
