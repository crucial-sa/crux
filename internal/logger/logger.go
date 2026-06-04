package logger

import "go.uber.org/zap"

var Zap *zap.Logger

func Init(debug bool) {
	var err error
	var logger *zap.Logger

	if debug {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}

	Zap = logger
}

func Sync() {
	if Zap != nil {
		Zap.Sync()
	}
}
