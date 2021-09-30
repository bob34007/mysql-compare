package utils

import "go.uber.org/zap"

type fileName string

func (fn fileName) Logger() *zap.Logger {
	logger := zap.L().With(zap.String("filename",string(fn)))
	if len(fn) > 0 {
		logger = logger.Named(string(fn))
	}
	return logger
}
