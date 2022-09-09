package nelly

import "go.uber.org/zap"

func ConvertChanToPointerChan[T any](regularChan chan T) chan *T {
	output := make(chan *T)
	go func() {
		for input := range regularChan {
			output <- &input
		}
	}()
	return output
}

func InitializeLogger() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	return logger.Sugar()
}
