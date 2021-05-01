package zslog

import (
	"io"
	"os"

	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const (
	// DebugLevel defines debug log level.
	DebugLevel = zerolog.DebugLevel
	// InfoLevel defines info log level.
	InfoLevel = zerolog.InfoLevel
	// WarnLevel defines warn log level.
	WarnLevel = zerolog.WarnLevel
	// ErrorLevel defines error log level.
	ErrorLevel = zerolog.ErrorLevel
	// FatalLevel defines fatal log level.
	FatalLevel = zerolog.FatalLevel
	// PanicLevel defines panic log level.
	PanicLevel = zerolog.PanicLevel
	// NoLevel defines an absent log level.
	NoLevel = zerolog.NoLevel
	// Disabled disables the logger.
	Disabled = zerolog.Disabled
	// TraceLevel defines trace log level.
	TraceLevel = zerolog.TraceLevel
)

type zlog struct {
	zerolog.Logger
	closers []io.Closer
	noSentryWriters zerolog.LevelWriter
}

var logger zlog

type levels []zerolog.Level

func NewLevels(lvls ...zerolog.Level) levels {
	if len(lvls) == 0 {
		return []zerolog.Level{
			zerolog.TraceLevel,
			zerolog.DebugLevel,
			zerolog.InfoLevel,
			zerolog.WarnLevel,
			zerolog.ErrorLevel,
			zerolog.FatalLevel,
			zerolog.PanicLevel}
	}
	var l []zerolog.Level
	for _, lvl := range lvls {
		l = append(l, lvl)
	}
	return l
}

func (l levels) SetMinLevel(level zerolog.Level) levels {
	var newL []zerolog.Level
	for _, n := range l {
		if n >= level {
			newL = append(newL, n)
		}
	}
	return newL
}

func (l levels) SetMaxLevel(level zerolog.Level) levels {
	var newL []zerolog.Level
	for _, n := range l {
		if n <= level {
			newL = append(newL, n)
		}
	}
	return newL
}

func (l levels) Contains(level zerolog.Level) bool {
	for _, n := range l {
		if level == n {
			return true
		}
	}
	return false
}

type LevelWriter struct {
	writer io.Writer
	levels levels
}

func (lw *LevelWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	if lw.levels.Contains(level) {
		return lw.writer.Write(p)
	}
	return len(p), nil
}

func (lw *LevelWriter) Write(p []byte) (n int, err error) {
	return lw.writer.Write(p)
}

func StdOut(l levels) zerolog.LevelWriter {
	lw := &LevelWriter{
		writer: zerolog.ConsoleWriter{
			Out: os.Stdout,
		},
		levels: l,
	}
	return lw
}

func StdErr(l levels) zerolog.LevelWriter {
	lw := &LevelWriter{
		writer: zerolog.ConsoleWriter{
			Out: os.Stderr,
		},
		levels: l,
	}
	return lw
}

func InitLogger(writers ...zerolog.LevelWriter) {
	var closers []io.Closer
	lwriters := zerolog.MultiLevelWriter()
	slwriters := zerolog.MultiLevelWriter()
	for _, w := range writers {
		if w != nil {
			lwriters = zerolog.MultiLevelWriter(lwriters, w)
			if c, ok := w.(io.Closer); ok {
				closers = append(closers, c)
			}
			if sw, ok := w.(*SentryWriter); !ok {
				slwriters = zerolog.MultiLevelWriter(slwriters, sw)
			}
		}
	}

	logger = zlog{
		Logger:          zerolog.New(lwriters).With().Timestamp().Logger(),
		closers:         closers,
		noSentryWriters: slwriters,
	}
	zerolog.ErrorStackMarshaler = func(err error) interface{} {
		es := errWithStackTrace{
			Err: err.Error(),
		}

		if _, ok := err.(stackTracer); !ok {
			err = errors.WithStack(err)
		}

		es.Stacktrace = sentry.ExtractStacktrace(err)
		return &es
	}
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

type errWithStackTrace struct {
	Err        string             `json:"error"`
	Stacktrace *sentry.Stacktrace `json:"stacktrace"`
}

func GetLogger() zerolog.Logger {
	return logger.Logger
}

func Trace() *zerolog.Event {
	return logger.Trace()
}
func Debug() *zerolog.Event {
	return logger.Debug()
}
func Info() *zerolog.Event {
	return logger.Info()
}
func Warn() *zerolog.Event {
	return logger.Warn()
}
func Error() *zerolog.Event {
	return logger.Error()
}
func Fatal() *zerolog.Event {
	return logger.Fatal()
}
func Panic() *zerolog.Event {
	return logger.Panic()
}

func Close() {
	for _, c := range logger.closers {
		c.Close()
	}
}

func SkipSentry() zerolog.Logger {
	return logger.Output(logger.noSentryWriters)
}
