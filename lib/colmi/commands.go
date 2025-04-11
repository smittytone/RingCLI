package ringcliColmi

const (
	COMMAND_SET_TIME              byte   = 0x01
	COMMAND_BATTERY_INFO          byte   = 0x03
	COMMAND_SHUTDOWN              byte   = 0x08
	COMMAND_BATTERY_FLASH_LED     byte   = 0x10
	COMMAND_HEART_RATE_READ       byte   = 0x15
	COMMAND_HEART_RATE_PERIOD     byte   = 0x16
	COMMAND_GET_ACTIVITY_DATA     byte   = 0x43
	COMMAND_GET_ACTIVITY_UNKNOWN  byte   = 0x73
	COMMAND_ERROR                 byte   = 0xFF

	UUID_BLE_COMMAND_UART_SERVICE string = "6E40FFF0-B5A3-F393-E0A9-E50E24DCCA9E"
	UUID_BLE_COMMAND_UART_TX_CHAR string = "6E400002-B5A3-F393-E0A9-E50E24DCCA9E"
	UUID_BLE_COMMAND_UART_RX_CHAR string = "6E400003-B5A3-F393-E0A9-E50E24DCCA9E"

	UUID_BLE_DATA_UART_SERVICE    string = "DE5BF728-D711-4E47-AF26-65E3012A5DC7"
	UUID_BLE_DATA_UART_TX_CHAR    string = "DE5BF72A-D711-4E47-AF26-65E3012A5DC7"
	UUID_BLE_DATA_UART_RX_CHAR    string = "DE5BF729-D711-4E47-AF26-65E3012A5DC7"

	DATA_REQUEST_ID_SLEEP         byte   = 0x27
	DATA_REQUEST_ID_OXYGEN        byte   = 0x2A
	DATA_REQUEST_MAGIC_VALUE      byte   = 0xBC

	LANGUAGE_CHINESE              byte   = 0x00
	LANGUAGE_ENGLISH              byte   = 0x01

	SLEEP_TYPE_NO_DATA            int    = 0x00
	SLEEP_TYPE_ERROR              int    = 0x01
	SLEEP_TYPE_LIGHT              int    = 0x02
	SLEEP_TYPE_DEEP               int    = 0x03
	SLEEP_TYPE_REM               int    = 0x04
	SLEEP_TYPE_AWAKE              int    = 0x05

	SLEEP_STRING_NO_DATA          string = "no data"
	SLEEP_STRING_ERROR            string = "error"
	SLEEP_STRING_LIGHT            string = "in light sleep"
	SLEEP_STRING_DEEP             string = "in deep sleep"
	SLEEP_STRING_REM              string = "in REM sleep (dreaming)"
	SLEEP_STRING_AWAKE            string = "awake"
)
