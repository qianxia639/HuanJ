package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLogger() {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncoderConfig.TimeKey = "timestamp"
	config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)

	// 设置输出目标
	// consoleCore := zapcore.NewCore(
	// 	zapcore.NewJSONEncoder(config.EncoderConfig),
	// 	zapcore.AddSync(os.Stdout),
	// 	zap.InfoLevel,
	// )

	// 构建 Logger
	var err error
	Logger, err = config.Build()
	if err != nil {
		panic(err)
	}

	// Logger = zap.New(consoleCore, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	defer Logger.Sync() // / 确保程序退出时刷新日志

	// 替换全局 logger
	zap.ReplaceGlobals(Logger)

}
