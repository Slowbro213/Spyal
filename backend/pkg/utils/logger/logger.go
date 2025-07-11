package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	logFilePath = "/var/log/alspy.log"
	logFileMode = 0o644
)

func NewLogger() (*zap.Logger, error) {
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, logFileMode)
	if err != nil {
		return nil, err
	}

	writeSyncer := zapcore.AddSync(file)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	core := zapcore.NewCore(encoder, writeSyncer, zap.InfoLevel)

	logger := zap.New(core)
	return logger, nil
}
