package log

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type GormLogger struct {
	ZapLogger                 *zap.Logger
	LogLevel                  gormlogger.LogLevel
	SlowThreshold             time.Duration
	SkipCallerLookup          bool
	IgnoreRecordNotFoundError bool
	prefix                    []interface{}
}

func NewGormLogger(file string, prefix []interface{}, loglevel gormlogger.LogLevel) *GormLogger {

	encoder := zapcore.EncoderConfig{
		TimeKey:       "t",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stack",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	level := zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	switch loglevel {
	case gormlogger.Error:
		level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case gormlogger.Info:
		level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case gormlogger.Warn:
		level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case gormlogger.Silent:
		level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	}
	zap.AddStacktrace(
		zap.NewAtomicLevelAt(zapcore.ErrorLevel))
	newEncoder := zapcore.NewJSONEncoder(encoder)
	writeSyncer := zapcore.NewMultiWriteSyncer(zapcore.AddSync(GetZapLoggergetWriter("./logs/project.log")))
	if level.String() == "debug" {
		newEncoder = zapcore.NewConsoleEncoder(encoder)
		writeSyncer = zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(GetZapLoggergetWriter("./logs/project.log")),
			zapcore.AddSync(os.Stdout),
		)
	}
	core := zapcore.NewCore(newEncoder, writeSyncer, level)
	zapLogger := zap.New(core, zap.AddCaller(),
		zap.AddCallerSkip(2))
	kvs := make([]interface{}, 0, len(prefix))
	kvs = append(kvs, prefix...)
	return &GormLogger{
		ZapLogger:                 zapLogger,
		LogLevel:                  loglevel,
		SlowThreshold:             1 * time.Second,
		SkipCallerLookup:          false,
		IgnoreRecordNotFoundError: false,
		prefix:                    kvs,
	}
}

func (l *GormLogger) SetAsDefault() {
	gormlogger.Default = l
}

func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

func (l *GormLogger) Info(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Info {
		return
	}
	l.logger().Sugar().Debugf(str, args...)
}

func (l *GormLogger) Warn(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Warn {
		return
	}
	l.logger().Sugar().Warnf(str, args...)
}

func (l *GormLogger) Error(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Error {
		return
	}
	l.logger().Sugar().Errorf(str, args...)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= 0 {
		return
	}
	elapsed := time.Since(begin)

	if len(l.prefix)%2 != 0 {
		l.logger().Warn(fmt.Sprint("Keyvalues must appear in pairs: ", l.prefix))
		return
	}

	// Zap.Field is used when keyvals pairs appear
	var data []zap.Field
	for i := 0; i < len(l.prefix); i += 2 {
		data = append(data, zap.Any(fmt.Sprint(l.prefix[i]), fmt.Sprint(l.prefix[i+1])))
	}
	switch {
	case err != nil && l.LogLevel >= gormlogger.Error && (!l.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		sql, rows := fc()
		data = append(data, zap.String("sql", sql), zap.Error(err), zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
		l.logger().Error("trace", data...)
	case l.SlowThreshold != 0 && elapsed > l.SlowThreshold && l.LogLevel >= gormlogger.Warn:
		sql, rows := fc()
		data = append(data, zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
		l.logger().Warn("trace", data...)
	case l.LogLevel >= gormlogger.Info:
		sql, rows := fc()
		data = append(data, zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
		l.logger().Debug("trace", data...)
	}
}

var (
	gormPackage    = filepath.Join("gorm.io", "gorm")
	zapgormPackage = filepath.Join("moul.io", "zapgorm2")
)

func (l *GormLogger) logger() *zap.Logger {
	for i := 2; i < 15; i++ {
		_, file, _, ok := runtime.Caller(i)
		switch {
		case !ok:
		case strings.HasSuffix(file, "_test.go"):
		case strings.Contains(file, gormPackage):
		case strings.Contains(file, zapgormPackage):
		default:
			return l.ZapLogger.WithOptions(zap.AddCallerSkip(i))
		}
	}
	return l.ZapLogger
}
