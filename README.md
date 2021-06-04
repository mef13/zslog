# zslog

[![Actions Status](https://github.com/mef13/zslog/workflows/build/badge.svg)](https://github.com/mef13/zslog/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/mef13/zslog)](https://goreportcard.com/report/github.com/mef13/zslog)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/mef13/zslog)](https://github.com/mef13/zslog/tags)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/mef13/zslog)
[![Codecov](https://img.shields.io/codecov/c/github/mef13/zslog)](https://app.codecov.io/gh/mef13/zslog)
![GitHub](https://img.shields.io/github/license/mef13/zslog)

Simply configuring [zerolog](https://github.com/rs/zerolog) with output to console, [lumberjack.v2](https://github.com/natefinch/lumberjack), [sentry-go](https://github.com/getsentry/sentry-go), http.

## Features

* [Easy use](#example)
* [Error Logging with stacktrace](#error-logging-with-stacktrace)
* Contextual Logging to `Payload` sentry context
* [Set log level for each outputs](#set-log-level-for-output)

## Example

```go
package main

import (
    "github.com/mef13/zslog"
    "time"
)

func main() {
    fileConf := zslog.FileConfig{
        MaxSize:    1,
        MaxBackups: 10,
    }
    zslog.InitLogger(
        zslog.StdOut(zslog.NewLevels().SetMaxLevel(zslog.WarnLevel)),
        zslog.StdErr(zslog.NewLevels().SetMinLevel(zslog.ErrorLevel)),
        zslog.File("/var/log/myapp/err.log", fileConf, zslog.NewLevels().SetMinLevel(zslog.ErrorLevel)),
        zslog.Sentry(zslog.SentryConfig{
            Dsn:              "https://public@sentry.example.com/1",
            Release:          "0.0",
        }, 3*time.Second, zslog.NewLevels().SetMinLevel(zslog.ErrorLevel)))

    //Flush sentry buffer before exit
    defer zslog.Close()

    zslog.Info().Msg("hello world")
}
```

### Set log level for output

Set specific levels: 

```go
zslog.StdOut(zslog.NewLevels(zslog.ErrorLevel, zslog.WarnLevel))
```

or set maximum(minimum) level:

```go
zslog.StdOut(zslog.NewLevels().SetMaxLevel(zslog.WarnLevel))
```

### Error Logging with Stacktrace

For add stacktrace use `Stack()` function:

```go
package main

import (
	"github.com/mef13/zslog"
	"github.com/pkg/errors"
	"time"
)

func main() {
	zslog.InitLogger(
		zslog.Sentry(zslog.SentryConfig{
			Dsn:     "https://public@sentry.example.com/1",
			Release: "0.0",
		}, 3*time.Second, zslog.NewLevels().SetMinLevel(zslog.ErrorLevel)))

	//Flush sentry buffer before exit
	defer zslog.Close()

	err = outer()
	zslog.Error().Stack().Err(err).Msg("outer error")
}

func inner() error {
	return errors.New("Ah s*, here we go again")
}

func middle() error {
	err := inner()
	if err != nil {
		return err
	}
	return nil
}

func outer() error {
	err := middle()
	if err != nil {
		return err
	}
	return nil
}
```
