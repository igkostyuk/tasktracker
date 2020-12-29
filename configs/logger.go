package configs

import (
	"go.uber.org/zap"
)

func GetLoggerConfig(build string) zap.Config {
	switch build {
	case "production":
		return productionConfig()
	default:
		return developmentConfig()
	}
}

// nolint:exhaustivestruct
func developmentConfig() zap.Config {
	return zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: true,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:          "json",
		EncoderConfig:     zap.NewDevelopmentEncoderConfig(),
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
		DisableStacktrace: false,
	}
}

// nolint:exhaustivestruct
func productionConfig() zap.Config {
	return zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:          "json",
		EncoderConfig:     zap.NewProductionEncoderConfig(),
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
		DisableStacktrace: true,
	}
}
