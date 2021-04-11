package util

import (
	"context"
	"strings"
	"time"

	"github.com/juju/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// TimeEncoder for logging time format.
func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.00"))
}

type key string

const TransIDKey key = "trid"

var mlog *MLogger

// MLogger ...
type MLogger struct {
	*zap.SugaredLogger
}

// With ...
func (m *MLogger) With(ctx context.Context) *zap.SugaredLogger {
	if ctx == nil {
		return m.SugaredLogger
	}
	id, ok := ctx.Value(TransIDKey).(string)
	if !ok || id == "" {
		return m.SugaredLogger
	}

	return m.SugaredLogger.With("tr", id)
}

// InitLog returns logger instance.
func InitLog(name string, env string) (log *MLogger, err error) {
	cfg := zap.NewDevelopmentConfig()
	cfg.Encoding = "console"
	cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	enccfg := zap.NewDevelopmentEncoderConfig()

	if LogOn(env, "log_info") {
		cfg.Encoding = "json"
		cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
		enccfg = zap.NewDevelopmentEncoderConfig()
	} else if LogOn(env, "log_error") {
		cfg = zap.NewProductionConfig()
		cfg.Encoding = "json"
		cfg.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
		enccfg = zap.NewProductionEncoderConfig()
	} else if LogOn(env, "log_fatal") {
		cfg = zap.NewProductionConfig()
		cfg.Encoding = "console"
		cfg.Level = zap.NewAtomicLevelAt(zapcore.FatalLevel)
		enccfg = zap.NewProductionEncoderConfig()
	}
	enccfg.EncodeTime = TimeEncoder
	enccfg.CallerKey = "caller"
	enccfg.LevelKey = ""
	cfg.EncoderConfig = enccfg

	logger, err := cfg.Build()
	if err != nil {
		return nil, errors.Annotatef(err, "InitLog")
	}
	defer logger.Sync()

	mlog = &MLogger{logger.Sugar()}
	return mlog, nil
}

// LogOn ...
func LogOn(level, target string) bool {
	l := strings.ToLower(level)
	return strings.Contains(l, target)
}
