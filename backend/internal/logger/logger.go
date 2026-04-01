package logger

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"operation_admin/backend/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// New 负责创建同时输出到控制台和滚动文件的 Zap Logger。
func New(cfg config.LogConfig) (*zap.Logger, error) {
	// level 用于控制日志输出的最小级别。
	level, err := parseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	// encoderConfig 定义日志时间、级别和调用方的编码方式。
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	// consoleCore 负责把更易读的日志输出到标准输出。
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		level,
	)

	// cores 汇总所有需要同时启用的日志输出通道。
	cores := []zapcore.Core{consoleCore}

	if cfg.Filename != "" {
		if err := ensureLogDirectory(cfg.Filename); err != nil {
			return nil, err
		}

		// fileWriter 负责执行日志滚动切割与落盘。
		fileWriter := &lumberjack.Logger{
			Filename:   cfg.Filename,
			MaxSize:    cfg.MaxSizeMB,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAgeDays,
			Compress:   cfg.Compress,
		}

		// fileCore 负责把 JSON 日志写入滚动文件，便于后续检索与采集。
		fileCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(fileWriter),
			level,
		)

		cores = append(cores, fileCore)
	}

	return zap.New(
		zapcore.NewTee(cores...),
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	), nil
}

// Sync 负责在进程退出前刷新缓冲区中的日志内容。
func Sync(logger *zap.Logger) error {
	if err := logger.Sync(); err != nil && !isIgnorableSyncError(err) {
		return err
	}

	return nil
}

// parseLevel 负责把文本级别转换为 Zap 可识别的日志级别。
func parseLevel(rawLevel string) (zapcore.Level, error) {
	// normalizedLevel 用于兼容不同大小写的外部配置输入。
	normalizedLevel := strings.TrimSpace(strings.ToLower(rawLevel))

	switch normalizedLevel {
	case "", "info":
		return zapcore.InfoLevel, nil
	case "debug":
		return zapcore.DebugLevel, nil
	case "warn", "warning":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	default:
		return zapcore.InfoLevel, fmt.Errorf("不支持的日志级别: %s", rawLevel)
	}
}

// ensureLogDirectory 负责确保日志文件所在目录已存在。
func ensureLogDirectory(filename string) error {
	// logDirectory 是日志文件对应的父目录。
	logDirectory := filepath.Dir(filename)
	if logDirectory == "." || logDirectory == "" {
		return nil
	}

	if err := os.MkdirAll(logDirectory, 0o755); err != nil {
		return fmt.Errorf("创建日志目录失败: %w", err)
	}

	return nil
}

// isIgnorableSyncError 负责过滤 stdout/stderr 在某些平台上的无害 Sync 错误。
func isIgnorableSyncError(err error) bool {
	if err == nil {
		return false
	}

	return errors.Is(err, os.ErrInvalid) || strings.Contains(strings.ToLower(err.Error()), "invalid argument")
}
