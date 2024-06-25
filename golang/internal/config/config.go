package config

import (
	"time"

	_ "github.com/spf13/viper"
)

/*  ==============================  */
/*  === FROM MY templates REPO ===  */
/*  ==============================  */

// Constaraints on ENV , Flags , Config.file and Default values
//
// #ENV
// - Must be defined, otherwise program shouldn't start
// - Lifetime constants, shouldnt be overridden in runtime
// - Can Be *defaulted only in build\makefile, where they explicitly defined
//
// #Config.file
// - Must exist, have same structure as config.Config, otherwise program shouldn't start
// - May be overridden in runtime or exist in multiple variants across sessions
// - Cannot Be Defaulted
//
// #--Flags
// - May not be defined, program should start,
// 		then SHOULD be defaulted by:
// 		[false , 0 , -1 , "NO" , "stop"]
// 		and any other kind of negative value
// - Program-Lifetime constants, cannot be changed in any way in runtime
// - Can Be Defaulted in build\makefile and in code \
// 		Code variant defaulting is preferred

type Config struct {
	Logger struct {
		Cores []struct {
			Name         string `mapstructure:"Name"`         // name of the core, for idientification reasons
			EncoderLevel string `mapstructure:"EncoderLevel"` // production or development
			Path         string `mapstructure:"Path"`         // everything that getLogFile can handle
			Level        int    `mapstructure:"Level"`        // level of the core , might be negative
		} `mapstructure:"Cores"`
	} `mapstructure:"Logger"`
	Service struct {
		Timeout time.Duration `mapstructure:"Timeout"`
	} `mapstructure:"Service"`
}

var C Config

// Initialaise config process
func InitConfig() {
	setEnv()
	setFlags()
	handleConfigFile()
	bindFlags()
	fillGlobalConfig()
	setElse()
	override()
}

// Two main functions you should change in config code are:
// setEnv() and setFlags()
// See ./example/example.go for additional hints

// Set ENV
// Immediately validate thorough utilitary register*()
func setEnv() {
	registerENV("CONFIG_FILE")
}

// Set flags and explicitly define defaults
// Defaults, as stated in constraints, should be *negative
func setFlags() {}

// Callback on config change , aliases etc.
func setElse() {}

// Do not use, this violates constraints
// If there any way to not override - do not override (C) Me
func override() {}

// *defaults
// "","false","no","stop" for string
// 0 for int
// 0.0 for float
// false for bool
