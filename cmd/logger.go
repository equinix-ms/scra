package cmd

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func getLogger() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	if viper.GetBool("debug") {
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	return config.Build()
}
