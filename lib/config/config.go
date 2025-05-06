package ringcliConfig

// This structure contains root-level settings.
type AppConfig struct {
	OutputToStdout bool
	OutputToJson   bool
	DoShowVersion  bool
	OutputToText   bool
}

var (
	Config = NewAppConfig()
)

func NewAppConfig() AppConfig {

	config := AppConfig{}
	config.OutputToStdout = false
	config.OutputToJson = false
	config.DoShowVersion = false
	config.OutputToText = false
	return config
}

func VerifyConfig() {

	if Config.OutputToStdout && Config.OutputToJson {
		Config.OutputToStdout = false
	}

	if Config.OutputToStdout && Config.OutputToText {
		Config.OutputToStdout = false
	}

	if Config.OutputToJson && Config.OutputToText {
		Config.OutputToText = false
	}
}
