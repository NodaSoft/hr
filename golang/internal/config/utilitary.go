package config

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// The difference between register() and bind() is
// that register() extends bind() logic
// in this case - validates that env is defined

// Adds must-validation to env binding
func registerENV(input ...string) {
	viper.BindEnv(input...)
	for _, env := range input {

		// Type-free validation
		// Not defined integer or bool would be "" as well
		envalue := viper.GetString(env)

		if envalue == "" {
			log.Fatalf("%s is not defined", env)
		}
	}
}

// Wraps viper.BindPFlags(pflag.CommandLine) into panic + os.Exit(1)
func bindFlags() {
	pflag.Parse()
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		log.Fatalf("Cannot bind flags: %v\n", err)
	}
}

// Wraps viper.BindPFlags(pflag.CommandLine) into panic + os.Exit(1)
func fillGlobalConfig() {
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Cannot read config file from CONFIG_FILE: %v\n", err)
	}
	err = viper.Unmarshal(&C)
	if err != nil {
		log.Fatalf("Cannot unmarshal config.file into config.C: %v\n", err)
	}
}

// Parse config file path for name and ext

func nameExtFromPath(path string) (name string, ext string) {
	nameAndExt := strings.Split(filepath.Base(path), ".")
	return nameAndExt[0], nameAndExt[1]
}

// Set config file name and extention
// Change only if something breaks
// For ./relative/path/to/config  and //full/path/to/config
// For config    .yaml .json .toml
// Works just fine
func handleConfigFile() {
	name, ext := nameExtFromPath(viper.GetString("CONFIG_FILE"))
	dir := filepath.Dir(viper.GetString("CONFIG_FILE"))
	viper.AddConfigPath(dir)
	viper.SetConfigName(name)
	viper.SetConfigType(ext)
}
