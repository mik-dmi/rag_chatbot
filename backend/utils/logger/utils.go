package lg

import "go.uber.org/zap"

func NewProductionLogger() (*zap.SugaredLogger, error) {
	prodConfig := zap.NewProductionConfig()
	// Write production logs to both stdout and a dedicated production log file.
	prodConfig.OutputPaths = []string{"stdout", "productionlogfile.log"}
	logger, err := prodConfig.Build()
	if err != nil {
		return nil, err
	}
	return logger.Sugar(), nil
}

func NewDevelopmentLogger() (*zap.SugaredLogger, error) {
	devConfig := zap.NewDevelopmentConfig()
	// Write development logs to both stdout and a dedicated development log file.
	devConfig.OutputPaths = []string{"stdout", "developmentlogfile.log"}
	logger, err := devConfig.Build()
	if err != nil {
		return nil, err
	}
	return logger.Sugar(), nil
}
