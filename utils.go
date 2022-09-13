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

// TODO: move this to someplace normal
var sugar *zap.SugaredLogger

func InitializeLogger() {
	logger, _ := zap.NewProduction()
	sugar = logger.Sugar()
}
