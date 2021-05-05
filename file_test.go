package zslog

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"
)

func TestFileWriter(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "testFile")
	log, err := New(File(filePath, FileConfig{}, NewLevels()))
	if err != nil {
		t.Fatal(err)
	}
	SetTimeFieldFormat("2006")
	log.Trace().Msg("msg")
	log.Close()
	f, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(f), fmt.Sprintf("{\"level\":\"trace\",\"time\":\"%s\",\"message\":\"msg\"}\n", time.Now().Format("2006")); got != want {
		t.Errorf("\ngot:\n%s\nwant:\n%s", got, want)
	}
}
