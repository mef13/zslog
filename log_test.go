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
