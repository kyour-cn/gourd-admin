package initialize

import (
	"app/internal/config"
	"app/internal/global"
	"app/internal/modules/common/dblog"
	"context"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"
)

// InitLog 初始化日志
func InitLog() error {
	conf, err := config.GetLogConfig()
	if err != nil {
		return err
	}

	// 日志文件输出
	output := &lumberjack.Logger{
		Filename:  conf.LogFile,
		MaxSize:   conf.MaxSize, // megabytes
		LocalTime: true,
	}

	writers := []io.Writer{output}
	if conf.Console {
		// 日志是否输出到控制台
		writers = append(writers, os.Stdout)
	}
	logWriter := io.MultiWriter(writers...)

	// 设置日志级别
	level := slog.LevelVar{}
	if err = level.UnmarshalText([]byte(conf.Level)); err != nil {
		level.Set(slog.LevelInfo) // 默认日志级别为info
	}

	options := &slog.HandlerOptions{
		Level: level.Level(), // 设置日志级别
	}

	var handler slog.Handler
	if conf.Encoding == "json" {
		handler = slog.NewJSONHandler(logWriter, options)
	} else if conf.Encoding == "text" {
		handler = slog.NewTextHandler(logWriter, options)
	} else { // default
		handler = DefaultHandler{
			Level:  level.Level(),
			Writer: logWriter,
		}
	}

	slog.SetDefault(slog.New(handler))

	return nil
}

// DefaultHandler 自定义日志处理器
type DefaultHandler struct {
	Level  slog.Level
	Writer io.Writer
}

func (DefaultHandler) WithAttrs([]slog.Attr) slog.Handler { return nil }

func (DefaultHandler) WithGroup(string) slog.Handler { return nil }

func (h DefaultHandler) Enabled(_ context.Context, l slog.Level) bool {
	return l >= h.Level
}

func (h DefaultHandler) Handle(_ context.Context, r slog.Record) error {
	dt := time.Now().Format("2006-01-02 15:04:05")
	msg := r.Message

	// 输出日志属性
	r.Attrs(func(a slog.Attr) bool {
		msg += " " + a.Key + "=" + a.Value.String()
		return true
	})

	// TODO: 记录日志到数据库，暂未确定这种做法
	if global.GetDb("default") != nil {
		typeLabel := strings.ToLower(r.Level.String())
		_ = dblog.New(typeLabel).
			Write(r.Level.String()+" "+msg[:min(30, len(msg))], msg)
	}

	_, err := h.Writer.Write([]byte(dt + " " + r.Level.String() + " " + msg + "\n"))
	return err
}
