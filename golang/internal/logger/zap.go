package logger

import (
	"log"
	"taskhandler/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/*  ==============================  */
/*  === FROM MY templates REPO ===  */
/*  ==============================  */

// Should be initialised via InitGlobalLogger()
var globalLogger *zap.SugaredLogger

// Little wrapper for future ease of identification
type LevelWithName struct {
	zap.AtomicLevel
	Name string
}

func withName(name string, level zap.AtomicLevel) LevelWithName {
	return LevelWithName{level, name}
}

// []LevelWithName may be used to change specific output destination log levels
// Thread safe
func InitLogger() (*zap.SugaredLogger, []LevelWithName) {
	log.Println("Logger initialization started")

	// May dynamicly change log levels in runtime, will be returned from InitLogger()
	levels := make([]LevelWithName, 0, len(config.C.Logger.Cores))

	// Creating cores fully dynamic from config
	// stderr/stdout supported, network not supported
	// TODO: Add network support
	cores := make([]zapcore.Core, 0, len(config.C.Logger.Cores))
	for _, core := range config.C.Logger.Cores {
		levels = append(levels, withName(core.Name, zap.NewAtomicLevelAt(zapcore.Level(core.Level))))
		cores = append(cores,
			zapcore.NewCore(
				mustSetEncoder(core.EncoderLevel), // production or development
				getLogDest(core.Path),             // file or stderr/stdout
				levels[len(levels)-1],             // last level
			))
	}

	// Creating zap.Cores
	// And merging them
	core := zapcore.NewTee(cores...)
	// Creating Logger from cores
	// And sugaring
	logger := zap.New(core)
	sugarlogger := logger.Sugar()

	// First log message
	// That tells us that logger construction succeeded
	defer sugarlogger.Sync()
	sugarlogger.Debug("Logger construction succeeded")

	return sugarlogger, levels
}

// Useful for small apps where you want to log a bit
// Not sure about async
func InitGlobalLogger() {
	//Ignoring ability to change level in runtime for global usecase
	//TODO: it is not hard to add this feature, mb next time
	globalLogger, _ = InitLogger()
	SyncOnTimout()
}

func Debug(args ...any) {
	globalLogger.Debug(args...)
}

func Info(args ...any) {
	globalLogger.Info(args...)
}
