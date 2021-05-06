package zslog

import (
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/getsentry/sentry-go"
)

//TODO: test sentry stacktrace
func TestSentry(t *testing.T) {
	transport := &TransportMock{}
	log, err := New(Sentry(SentryConfig{Transport: transport}, time.Second, NewLevels()))
	if err != nil {
		t.Fatal(err)
	}
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
