package zslog

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
)

func TestConsoleLogger(t *testing.T) {
	t.Run("Numbers", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := New(&LevelWriter{writer: zerolog.ConsoleWriter{Out: buf, NoColor: true, TimeFormat: "2006"}, levels: NewLevels()})
		log.Info().
			Float64("float", 1.23).
			Uint64("small", 123).
			Uint64("big", 1152921504606846976).
			Msg("msg")
		if got, want := strings.TrimSpace(buf.String()), fmt.Sprintf("%s INF msg big=1152921504606846976 float=1.23 small=123", time.Now().Format("2006")); got != want {
			t.Errorf("\ngot:\n%s\nwant:\n%s", got, want)
		}
	})
}

//TODO: test sentry stacktrace
func TestSentry(t *testing.T) {
	transport := &TransportMock{}
	log := New(Sentry(SentryConfig{Transport: transport}, time.Second, NewLevels()))
	log.Debug().Str("s", "test").Int("i", 9).Msg("msg")
	if transport.lastEvent == nil {
		t.Fatal("Sentry not capture event")
	}
	t.Run("Message", func(t *testing.T) {
		if transport.lastEvent.Message != "msg" {
			t.Errorf("Sentry message want 'msg', but get '%s'", transport.lastEvent.Message)
		}
	})
	t.Run("Level", func(t *testing.T) {
		if transport.lastEvent.Level != "debug" {
			t.Errorf("Sentry want debug level, but get %s", transport.lastEvent.Level)
		}
	})
	t.Run("Payload", func(t *testing.T) {
		payload, ok := transport.lastEvent.Contexts["payload"]
		if !ok {
			t.Fatal("Payload not found")
		}
		p, ok:=payload.(map[string]string)
		if !ok {
			t.Fatalf("Invalid payload structure: %s", reflect.ValueOf(payload).Type().String())
		}
		if p["s"] != "test" {
			t.Errorf("Sentry want 'test' in 's' payload, but get '%s'", p["s"])
		}
		if p["i"] != "9" {
			t.Errorf("Sentry want '9' in 'i' payload, but get '%s'", p["i"])
		}
	})
}

type TransportMock struct {
	mu        sync.Mutex
	events    []*sentry.Event
	lastEvent *sentry.Event
}

func (t *TransportMock) Configure(options sentry.ClientOptions) {}
func (t *TransportMock) SendEvent(event *sentry.Event) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.events = append(t.events, event)
	t.lastEvent = event
}
func (t *TransportMock) Flush(timeout time.Duration) bool {
	return true
}
func (t *TransportMock) Events() []*sentry.Event {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.events
}
