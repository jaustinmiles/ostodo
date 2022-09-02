package common

import "go.uber.org/zap"

func GetLogger() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	return logger.Sugar()
}
