package log

import (
	"os"
	"time"

	"github.com/rs/zerolog"

	"github.com/tech-a-go-go/timeseries-aggregator/internal/config"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zerolog.Logger

func init() {
	// zerolog のtimeログのフォーマットを変更
	zerolog.TimeFieldFormat = time.RFC3339Nano

	// アプリケーションのログ設定
	var l zerolog.Logger
	if config.GetEnv().IsTest() {
		l = zerolog.New(os.Stdout)
	} else {
		applicationLogFilename := config.GetString("log.file")
		writer := &lumberjack.Logger{
			Filename: applicationLogFilename,
			MaxSize:  10, // megabytes
			// MaxBackups: 5,
			MaxAge:    3, //days
			LocalTime: true,
			Compress:  false, // disabled by default
		}
		l = zerolog.New(writer)
	}
	logger = &l
}

// GetLogger recorderアプリケーション用Loggerを返す
func GetLogger() *zerolog.Logger {
	return logger
}
