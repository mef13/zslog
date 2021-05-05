package zslog

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHttpWriter(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		if got, want := string(body), fmt.Sprintf("{\"level\":\"info\",\"time\":\"%s\",\"message\":\"msg\"}\n", time.Now().Format("2006")); got != want {
			t.Errorf("\ngot:\n%s\nwant:\n%s", got, want)
		}
		w.WriteHeader(200)
	}))
	defer server.Close()
	log, err := New(GetHttpWriter("POST", server.URL, NewLevels()))
	if err != nil {
		t.Fatal(err)
	}
	SetTimeFieldFormat("2006")
	log.Info().Msg("msg")

	t.Run("Nil writer", func(t *testing.T) {
		log, err = New(GetHttpWriter("POST", server.URL, NewLevels()), GetHttpWriter("POST", "1", NewLevels()))
		if err == nil {
			t.Error("an error was expected, but it did not appear")
		} else {
			if err != ErrInitWriter {
				t.Fatal(err)
			}
		}
		SetTimeFieldFormat("2006")
		log.Info().Msg("msg")
	})
}
