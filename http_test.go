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
		w.Write([]byte(`OK`))
	}))
	log := New(&HttpWriter{
		client: server.Client(),
		url:    server.URL,
		method: "POST",
		levels: NewLevels(),
	})
	SetTimeFieldFormat("2006")
	defer server.Close()
	log.Info().Msg("msg")
}
