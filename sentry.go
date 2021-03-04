package zslog

import (
	"encoding/json"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
	"reflect"
	"time"
	"unsafe"
)

var levelsMapping = map[zerolog.Level]sentry.Level{
	zerolog.DebugLevel: sentry.LevelDebug,
	zerolog.InfoLevel:  sentry.LevelInfo,
	zerolog.WarnLevel:  sentry.LevelWarning,
	zerolog.ErrorLevel: sentry.LevelError,
	zerolog.FatalLevel: sentry.LevelFatal,
	zerolog.PanicLevel: sentry.LevelFatal,
}

type SentryConfig sentry.ClientOptions

type SentryWriter struct {
	client *sentry.Client

	levels       map[zerolog.Level]sentry.Level
	flushTimeout time.Duration
}

func Sentry(sentryConf SentryConfig, flushTimeout time.Duration, lvls levels) zerolog.LevelWriter {
	client, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:              sentryConf.Dsn,
		Debug:            sentryConf.Debug,
		AttachStacktrace: sentryConf.AttachStacktrace,
		SampleRate:       sentryConf.SampleRate,
		TracesSampleRate: sentryConf.TracesSampleRate,
		TracesSampler:    sentryConf.TracesSampler,
		IgnoreErrors:     sentryConf.IgnoreErrors,
		BeforeSend:       sentryConf.BeforeSend,
		BeforeBreadcrumb: sentryConf.BeforeBreadcrumb,
		Integrations:     sentryConf.Integrations,
		DebugWriter:      sentryConf.DebugWriter,
		Transport:        sentryConf.Transport,
		ServerName:       sentryConf.ServerName,
		Release:          sentryConf.Release,
		Dist:             sentryConf.Dist,
		Environment:      sentryConf.Environment,
		MaxBreadcrumbs:   sentryConf.MaxBreadcrumbs,
		HTTPClient:       sentryConf.HTTPClient,
		HTTPTransport:    sentryConf.HTTPTransport,
		HTTPProxy:        sentryConf.HTTPProxy,
		HTTPSProxy:       sentryConf.HTTPSProxy,
		CaCerts:          sentryConf.CaCerts,
	})
	if err != nil {
		return nil
	}

	levels := make(map[zerolog.Level]sentry.Level, len(lvls))
	for _, lvl := range lvls {
		l, ok := levelsMapping[lvl]
		if ok {
			levels[lvl] = l
		}
	}

	timeout := 3 * time.Second
	if flushTimeout != 0 {
		timeout = flushTimeout
	}

	return &SentryWriter{
		client:       client,
		levels:       levels,
		flushTimeout: timeout,
	}
}

func (sw *SentryWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (sw *SentryWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	slevel, ok := sw.levels[level]
	if !ok {
		return len(p), nil
	}
	event, err := sw.parseEvent(p, slevel)
	if err != nil {
		return len(p), nil
	}
	sw.client.CaptureEvent(event, nil, nil)
	// should flush before os.Exit
	if event.Level == sentry.LevelFatal {
		sw.client.Flush(sw.flushTimeout)
	}
	return len(p), nil
}

func (sw *SentryWriter) Close() error {
	sw.client.Flush(sw.flushTimeout)
	return nil
}

func (sw *SentryWriter) parseEvent(data []byte, level sentry.Level) (*sentry.Event, error) {
	event := sentry.Event{
		Timestamp: time.Now().UTC(),
		Level:     level,
		Logger:    "github.com/rs/zerolog",
		Contexts:  make(map[string]interface{}),
	}

	isStack := false
	var errExept sentry.Exception
	payload := make(map[string]string)

	err := jsonparser.ObjectEach(data, func(key, value []byte, vt jsonparser.ValueType, offset int) error {
		switch string(key) {
		case zerolog.LevelFieldName, zerolog.TimestampFieldName:
			break
		case zerolog.MessageFieldName:
			event.Message = bytesToStrUnsafe(value)
		case zerolog.ErrorFieldName:
			errExept = sentry.Exception{
				Value:      bytesToStrUnsafe(value),
				Stacktrace: sentry.NewStacktrace(),
			}
		case zerolog.ErrorStackFieldName:
			var e errWithStackTrace
			err := json.Unmarshal(value, &e)
			if err != nil {
				fmt.Println("Error unmarshal: ", err)
				break
			}
			event.Exception = append(event.Exception, sentry.Exception{
				Value:      e.Err,
				Stacktrace: e.Stacktrace,
			})
			isStack = true
		default:
			payload[string(key)] = bytesToStrUnsafe(value)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(payload) != 0{
		event.Contexts["payload"] = payload
	}
	if !isStack && errExept.Value != "" {
		event.Exception = append(event.Exception, errExept)
	}

	return &event, nil
}

func bytesToStrUnsafe(data []byte) string {
	h := (*reflect.SliceHeader)(unsafe.Pointer(&data))
	return *(*string)(unsafe.Pointer(&reflect.StringHeader{Data: h.Data, Len: h.Len}))
}
