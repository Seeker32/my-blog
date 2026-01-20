package initialization

import (
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LoggerConfig struct {
	LogDir     string `mapstructure:"log_dir"`
	MinLevel   string `mapstructure:"min_level"`
	MaxSizeMB  int    `mapstructure:"max_size_mb"`
	MaxBackUps int    `mapstructure:"max_backups"`
	MaxAgeDays int    `mapstructure:"max_age_days"`
	Compress   bool   `mapstructure:"compress"`
}

func InitLogger(cfg *LoggerConfig) *zap.Logger {
	// 解析最低日志等级
	minLevel := parseLogLevel(cfg.MinLevel)

	// 判断并转换为绝对路径
	logDir := cfg.LogDir
	if !filepath.IsAbs(logDir) {
		absPath, err := filepath.Abs(logDir)
		if err != nil {
			// 如果转换失败，使用当前工作目录
			logDir = filepath.Join(".", logDir)
		} else {
			logDir = absPath
		}
	}

	// 确保日志目录存在
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		_ = os.MkdirAll(logDir, 0755)
	}

	// 配置 zapcore 编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var cores []zapcore.Core

	// 控制台日志
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		minLevel,
	)
	cores = append(cores, consoleCore)

	allLevels := []zapcore.Level{
		zapcore.DebugLevel,
		zapcore.InfoLevel,
		zapcore.WarnLevel,
		zapcore.ErrorLevel,
		zapcore.PanicLevel,
		zapcore.FatalLevel,
	}

	// 为每个日志等级配置单独的日志文件
	for _, level := range allLevels {
		if level < minLevel {
			continue
		}

		filename := filepath.Join(logDir, strings.ToLower(level.String())+".log")

		fileWriter := &lumberjack.Logger{
			Filename:   filename,
			MaxSize:    cfg.MaxSizeMB,
			MaxBackups: cfg.MaxBackUps,
			MaxAge:     cfg.MaxAgeDays,
			Compress:   cfg.Compress,
		}

		levelEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl == level
		})

		// 创建 Core
		fileCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(fileWriter),
			levelEnabler,
		)

		cores = append(cores, fileCore)
	}

	core := zapcore.NewTee(cores...)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))

	zap.ReplaceGlobals(logger)
	return logger
}

// 解析日志等级
func parseLogLevel(levelStr string) zapcore.Level {
	switch strings.ToLower(levelStr) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}
