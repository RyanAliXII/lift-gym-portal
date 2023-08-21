package logger

import (
	"sync"

	"go.uber.org/zap"
)

var logger * zap.Logger
var once sync.Once;

func Get()*zap.Logger{
	once.Do(func() {
		logger, _ = zap.NewProduction()
	})
	return logger
}