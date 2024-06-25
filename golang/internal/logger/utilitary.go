package logger

import (
	"log"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/*  ==============================  */
/*  === FROM MY templates REPO ===  */
/*  ==============================  */

// TODO: There might be problems with /stderr and debugging go code via deluge
// TODO: Add network support
func getLogDest(path string) *os.File {

	// Trying to create log file
	logfile, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)

	if err != nil {

		// There is common case that directory doesn't exist
		// So we try to create it
		log.Println("Cannot create log file", err)
		log.Println("Trying to create directory")
		os.Mkdir(filepath.Dir(path), 0777)

		// Retry to create log file
		logfile, err = os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			log.Println("Unsuccessful logger initialization, cannot create log file ", err)
			return nil
		}
	}
	return logfile
}

// Be careful when changing config.logger.cores.encoderLevel in runtime
// Might Panic!
func mustSetEncoder(name string) zapcore.Encoder {

	if name == "production" {
		return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	}
	if name == "development" {
		return zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig())
	}
	// this panics may only happen if config is wrong
	// be very carefully when changing config.logger.cores.encoderLevel
	log.Fatal("Unknown encoder level: ", name)
	return nil
}
