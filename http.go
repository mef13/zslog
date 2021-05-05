package zslog

import (
	"bytes"
	"net/http"
	"net/url"
	"time"

	"github.com/rs/zerolog"
)

type HttpWriter struct {
	client *http.Client
	url string
	method string
	levels levels
}

func (hw *HttpWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	if hw.levels.Contains(level) {
		return hw.Write(p)
	}
	return len(p), nil
}

func (hw *HttpWriter) Write(p []byte) (n int, err error) {
	req, err := http.NewRequest(hw.method, hw.url, bytes.NewBuffer(p))
	if err != nil {
		return 0, err
	}
	response, err := hw.client.Do(req)
	if err != nil {
		return 0, err
	}

	_ = response.Body.Close()

	return len(p), nil
}

func GetHttpWriter(method string, endpoint string, l levels) zerolog.LevelWriter {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil
	}
	if u.Host == "" {
		return nil
	}
	h := &HttpWriter{
		client: &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 50,
			},
			Timeout: 6 * time.Second,
		},
		url:    u.String(),
		//TODO: validate method
		method: method,
		levels: l,
	}
	return h
}
