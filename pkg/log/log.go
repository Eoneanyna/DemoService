package log

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"runtime"
	"strings"
	"time"
)

var _ log.Logger = (*ZapLogger)(nil)

type ReportErr interface {
	ReportErr([]zap.Field) error
}
type ZapLogger struct {
	log       *zap.Logger
	Sync      func() error
	reportErr ReportErr
}

type LogLevelStr string

const (
	LogLevelStrDebug = "debug"
	LogLevelStrWarn  = "warn"
	LogLevelStrError = "error"
	LogLevelStrInfo  = "info"
)

// NewZapLogger return ZapLogger
func NewZapLogger(encoder zapcore.EncoderConfig, levelStr LogLevelStr, opts ...zap.Option) *ZapLogger {
	newEncoder := zapcore.NewJSONEncoder(encoder)
	writeSyncer := zapcore.NewMultiWriteSyncer(zapcore.AddSync(GetZapLoggergetWriter("./logs/project.log")))
	var level zap.AtomicLevel
	switch levelStr {
	case LogLevelStrDebug:
		level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case LogLevelStrWarn:
		level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case LogLevelStrError:
		level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case LogLevelStrInfo:
		level = zap.NewAtomicLevelAt(zapcore.InfoLevel)

	default:
		level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}
	if level.String() == "debug" {
		newEncoder = zapcore.NewConsoleEncoder(encoder)
		writeSyncer = zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(GetZapLoggergetWriter("./logs/project.log")),
			zapcore.AddSync(os.Stdout),
		)
	}
	core := zapcore.NewCore(newEncoder, writeSyncer, level)
	zapLogger := zap.New(core, opts...)
	return &ZapLogger{log: zapLogger, Sync: zapLogger.Sync}
}
func (l *ZapLogger) WithReporter(reporter ReportErr) {
	l.reportErr = reporter

}

// Log Implementation of logger interface
func (l *ZapLogger) Log(level log.Level, keyvals ...interface{}) error {
	if len(keyvals) == 0 || len(keyvals)%2 != 0 {
		l.log.Warn(fmt.Sprint("Keyvalues must appear in pairs: ", keyvals))
		return nil
	}

	// Zap.Field is used when keyvals pairs appear
	var data []zap.Field
	//是否需要上报错误
	var isReportErr bool
	for i := 0; i < len(keyvals); i += 2 {
		//获取key
		key := fmt.Sprint(keyvals[i])
		//获取value
		value := fmt.Sprint(keyvals[i+1])
		//判断key的是否是code,并且错误码是500(兼容老代码)，才上报错误
		if key == "code" && value == "500" {
			isReportErr = true
		}
		data = append(data, zap.Any(fmt.Sprint(keyvals[i]), fmt.Sprint(keyvals[i+1])))
	}
	switch level {
	case log.LevelDebug:
		l.log.Debug("", data...)
	case log.LevelInfo:
		l.log.Info("", data...)
	case log.LevelWarn:
		l.log.Warn("", data...)
	case log.LevelError:
		if l.reportErr != nil && isReportErr {
			_ = l.reportErr.ReportErr(data)
		}
		l.log.Error("", data...)
	}
	return nil
}

func GetZapLoggergetWriter(filename string) io.Writer {
	// 生成rotatelogs的Logger 实际生成的文件名 demo.log.YYmmddHH
	// demo.log是指向最新日志的链接
	// 保存7天内的日志，每1小时(整点)分割一次日志
	hook, err := rotatelogs.New(
		strings.Replace(filename, ".log", "", -1)+"-%Y%m%d.log", // 没有使用go风格反人类的format格式
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*7),
		rotatelogs.WithRotationTime(time.Hour*24),
	)

	if err != nil {
		panic(err)
	}
	return hook
}

func GetErrorLine(err error) string {
	_, file, line, ok := runtime.Caller(1)
	lineStr := err.Error()
	if ok {
		lineStr += "=>" + fmt.Sprintf("\n\t%s:%d\n\t", file, line)
	}
	return lineStr
}
