package logster

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapAdapter struct {
	prefix string
	*zap.SugaredLogger
}

func (z *ZapAdapter) WithPrefix(prefix string) Logger {
	return &ZapAdapter{SugaredLogger: z.SugaredLogger, prefix: prefix}
}

func (z *ZapAdapter) WithField(key string, value interface{}) Logger {
	return z.with(z.prefix+key, value)
}

func (z *ZapAdapter) WithError(err error) Logger {
	if err != nil {
		return z.with(zap.String("error", err.Error()))
	}
	return z.with(zap.String("error", "<nil>"))
}

func (z *ZapAdapter) Write(p []byte) (n int, err error) {
	z.Warnw(string(p))
	return len(p), nil
}

func (z *ZapAdapter) with(args ...interface{}) Logger {
	return &ZapAdapter{SugaredLogger: z.With(args...), prefix: z.prefix}
}

func LogIfError(logger Logger, err error, msg string, args ...interface{}) error {
	if err != nil {
		logger.WithError(err).Errorf(msg, args...)
	}
	return err
}

var textToZapLevelMap = map[string]zapcore.Level{
	"panic": zapcore.PanicLevel,
	"fatal": zapcore.FatalLevel,
	"error": zapcore.ErrorLevel,
	"warn":  zapcore.WarnLevel,
	"info":  zapcore.InfoLevel,
	"debug": zapcore.DebugLevel,
}

func New(w zapcore.WriteSyncer, cfg Config) *ZapAdapter {
	if cfg.Project == "" {
		panic("logster: project field should be nonempty")
	}

	env := cfg.Env
	if env == "" {
		env = "local"
	}

	level, ok := textToZapLevelMap[cfg.Level]
	if !ok {
		level = zapcore.InfoLevel
	}

	fields := []zap.Field{
		zap.String("go_env", env),
		zap.String("go_project", cfg.Project),
	}

	var encoderCfg zapcore.EncoderConfig
	var enc zapcore.Encoder

	encoderCfg = zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = func(time time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(time.UTC().Format("2006-01-02T15:04:05.999Z07:00"))
	}
	enc = zapcore.NewJSONEncoder(encoderCfg)

	options := []zap.Option{
		zap.Fields(fields...),
		zap.AddCaller(),
	}

	core := zapcore.NewCore(enc, w, level)
	sugar := zap.New(core).WithOptions(options...).Sugar()

	return &ZapAdapter{SugaredLogger: sugar}
}
