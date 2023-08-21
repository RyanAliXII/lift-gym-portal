package applog

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger * zap.Logger
var once sync.Once;

func Get()*zap.Logger{
	once.Do(func() {
		cfg := zap.Config{
			Encoding:         "json",
			Level:            zap.NewAtomicLevel(),
			ErrorOutputPaths: []string{"stderr"},
			OutputPaths:      []string{"stdout", "app/logs/app.log"},
			EncoderConfig: zapcore.EncoderConfig{
				MessageKey:  "message",
				LevelKey:    "level",
				TimeKey:     "time",
				EncodeLevel: zapcore.CapitalLevelEncoder,
				EncodeTime:  zapcore.EpochTimeEncoder,
			},
		}
		var err error
		logger, err = cfg.Build()
		if err != nil {
			panic(err.Error())
		}
	})
	return logger
}