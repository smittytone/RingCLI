package ringcliErrors

import (
	"fmt"
)

type RingcliError struct {
	Message string
	Code    int
}

func (e *RingcliError) Error() string {

	return fmt.Sprintf("%s (%d)", e.Message, e.Code)
}

const (
	ERROR_CODE_NONE                     int = 0
	ERROR_CODE_BAD_PARAMS               int = 1
	ERROR_CODE_BLE                      int = 2
	ERROR_CODE_SCAN_TIMEOUT             int = 3
	ERROR_CODE_BAD_ACTIVITY_TIME_OFFSET int = 4
	ERROR_CODE_BAD_HEART_DATA_REQUEST   int = 5
	ERROR_CODE_BAD_SET_TIME             int = 6
	ERROR_CODE_BAD_BCD_INPUT_VALUE      int = 7
	ERROR_CODE_BINDING_FILE_ERROR       int = 10
)
