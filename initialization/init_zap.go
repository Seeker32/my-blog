package initialization

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Seeker32/my-blog/pkg"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger(cfg *pkg.LoggerConfig) *zap.Logger {
	minLevel := parseLogLevel(cfg.MinLevel)

	logDir := cfg.LogDir
	if !filepath.IsAbs(logDir) {
		absPath, err := filepath.Abs(logDir)
		if err != nil {
			logDir = filepath.Join(".", logDir)
		} else {
			logDir = absPath
		}
	}

	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		_ = os.MkdirAll(logDir, 0755)
	}

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

	for _, level := range allLevels {
		if level < minLevel {
			continue
		}

		filename := filepath.Join(logDir, strings.ToLower(level.String())+".log")

		fileWriter := &lumberjack.Logger{
			Filename:   filename,
			MaxSize:    cfg.MaxSizeMB,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAgeDays,
			Compress:   cfg.Compress,
		}

		levelEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl == level
		})

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