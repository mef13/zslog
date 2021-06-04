package zslog

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
)

func TestConsoleLogger(t *testing.T) {
	t.Run("Numbers", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log, err := New(&LevelWriter{writer: zerolog.ConsoleWriter{Out: buf, NoColor: true, TimeFormat: "2006"}, levels: NewLevels()})
		if err != nil {
			t.Fatal(err)
		}
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

func TestLevels(t *testing.T) {
	l := NewLevels(WarnLevel)
	if len(l) != 1 {
		t.Error("NewLevel return incorrect levels count")
	}
	if l[0] != WarnLevel {
		t.Errorf("got:%v, want:%v", l[0], WarnLevel)
	}
	l = l.SetMinLevel(ErrorLevel)
	if len(l) != 0 {
		t.Error("SetMinLevel return incorrect levels count")
	}
	if l.Contains(WarnLevel) {
		t.Error("SetMinLevel work incorrect")
	}
}
