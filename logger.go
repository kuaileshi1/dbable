// @Title 请填写文件名称（需要改）
// @Description 请填写文件描述（需要改）
// @Author shigx 2022/3/17 4:44 下午
package dbable

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"path"
	"time"
)

type LogConfig struct {
	LogPath                   string        `yaml:"logPath"`                   // 日志存放位置
	MaxAge                    time.Duration `yaml:"maxAge"`                    // 日志文件最大保存时间 单位小时
	RotationTime              time.Duration `yaml:"rotationTime"`              // 日志切割时间
	LogLevel                  string        `yaml:"logLevel"`                  // 日志级别
	SlowThreshold             time.Duration `yaml:"slowThreshold"`             // 慢sql阈值
	IgnoreRecordNotFoundError bool          `yaml:"ignoreRecordNotFoundError"` // 忽略ErrRecordNotFound（记录未找到）错误
}

func getLogWriter(config *LogConfig) *rotatelogs.RotateLogs {
	// 检测目录是否存在，不存在创建
	if _, err := os.Stat(config.LogPath); os.IsNotExist(err) {
		os.MkdirAll(config.LogPath, 0755)
	}

	filename := path.Join(config.LogPath, "sql")

	logWriter, _ := rotatelogs.New(
		filename+".%Y-%m-%d.log",
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(config.MaxAge*time.Hour),
		rotatelogs.WithRotationTime(config.RotationTime*time.Hour),
	)

	return logWriter
}

func newLogger(config *LogConfig) logger.Interface {
	logWriter := getLogWriter(config)

	var logLevel logger.LogLevel
	switch config.LogLevel {
	case "info":
		logLevel = logger.Info
	case "warn":
		logLevel = logger.Warn
	case "error":
		logLevel = logger.Error
	default:
		logLevel = logger.Silent
	}

	newLogger := logger.New(
		log.New(logWriter, "", log.LstdFlags),
		logger.Config{
			SlowThreshold:             config.SlowThreshold * time.Millisecond,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: config.IgnoreRecordNotFoundError,
			Colorful:                  false,
		},
	)

	return newLogger
}
