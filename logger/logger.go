package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// var Logger *zap.Logger

func extendedColorLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	var color string
	switch level {
	case zapcore.DebugLevel:
		color = "\033[32m" // 亮青色
	case zapcore.InfoLevel:
		color = "\033[32m" // 亮绿色
	case zapcore.WarnLevel:
		color = "\033[33m" // 亮黄色
	case zapcore.ErrorLevel:
		color = "\033[31m"                // 亮红色
		color = "\033[48;5;88m\033[1;37m" // 深红底亮白字（加粗）
	default:
		color = "\033[0m"
	}

	enc.AppendString(color + level.CapitalString() + "\033[0m")
}

func InitLogger() *zap.Logger {

	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeTime = func(t time.Time, pae zapcore.PrimitiveArrayEncoder) {
		pae.AppendString(t.Format("2006/01/02 - 15:04:05.000"))
	}
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncoderConfig.TimeKey = "timestamp"

	// 设置输出目标
	// consoleCore := zapcore.NewCore(
	// 	zapcore.NewJSONEncoder(config.EncoderConfig),
	// 	zapcore.AddSync(os.Stdout),
	// 	zap.InfoLevel,
	// )

	// 构建 Logger
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	// Logger = zap.New(consoleCore, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	defer logger.Sync() // / 确保程序退出时刷新日志

	return logger

}

// func ProductionLogger() *zap.Logger {
// 	// 自定义编码器（级别前置+紧凑格式）
// 	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
// 		LevelKey:        "level",
// 		TimeKey:         "timestamp",
// 		CallerKey:       "caller",
// 		MessageKey:      "msg",
// 		StacktraceKey:   "stack",
// 		LineEnding:      zapcore.DefaultLineEnding,
// 		EncodeLevel:     zapcore.CapitalColorLevelEncoder,
// 		EncodeTime:      zapcore.ISO8601TimeEncoder,
// 		EncodeDuration:  zapcore.MillisDurationEncoder,
// 		EncodeCaller:    zapcore.ShortCallerEncoder,
// 	})

// 	// 日志级别动态配置
// 	atomicLevel := zap.NewAtomicLevelAt(zap.InfoLevel)

// 	// 控制台核心（带颜色）
// 	consoleCore := zapcore.NewCore(
// 		encoder,
// 		zapcore.Lock(os.Stdout),
// 		atomicLevel,
// 	)

// 	// 文件核心（JSON格式+轮转）
// 	fileCore := zapcore.NewCore(
// 		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
// 		zapcore.AddSync(&lumberjack.Logger{
// 			Filename:   "logs/app.log",
// 			MaxSize:    100, // MB
// 			MaxBackups: 7,
// 			MaxAge:     30, // days
// 		}),
// 		atomicLevel,
// 	)

// 	core := zapcore.NewTee(consoleCore, fileCore)

// 	return zap.New(core,
// 		zap.AddCaller(),
// 		zap.AddStacktrace(zap.ErrorLevel),
// 		zap.Fields(zap.String("app", "my-im")),
// 	)
// }
