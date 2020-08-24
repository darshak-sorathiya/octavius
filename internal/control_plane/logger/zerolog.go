package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"octavius/pkg/constant"
	"os"
	"time"
)

//Logger holds the pointer to zerolog.Logger object
type Logger struct {
	logger *zerolog.Logger
}

var log Logger

//Setup intializes the logger object
func Setup(configLogLevel string) {
	if (log != Logger{}) {
		return
	}
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
	logLevel, err := zerolog.ParseLevel(configLogLevel)
	if err != nil {
		fmt.Println("log level parsing problem")
	}
	logInit := zerolog.New(consoleWriter).With().Timestamp().CallerWithSkipFrameCount(constant.LoggerSkipFrameCount).Logger().Level(logLevel)
	zerolog.TimeFieldFormat = time.RFC822
	logInit.Level(logLevel)
	log = Logger{
		logger: &logInit,
	}
}

//Debug logs the message at debug level
func Debug(msg string) {
	log.logger.Debug().Msg(msg)
}

//Info logs the message at info level
func Info(msg string) {
	log.logger.Info().Msg(msg)
}

//Warn logs the message at warn level
func Warn(msg string) {
	log.logger.Warn().Msg(msg)
}

//Error logs the message and error at error level if err is not nil else it logs message at info level
func Error(err error, msg string) {
	log.logger.Err(err).Msg(msg)
}

//Fatal logs the message at fatal level followed by a os.Exit(1) call
func Fatal(msg string) {
	log.logger.Fatal().Msgf(msg)
}

//Panic logs the message and error at panic level and stops the ordinary flow of goroutine
func Panic(msg string, err error) {
	log.logger.Panic().Err(err).Msg(msg)
}