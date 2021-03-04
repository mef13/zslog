# zslog
Simply configuring [zerolog](https://github.com/rs/zerolog) with output to console, [lumberjack.v2](https://github.com/natefinch/lumberjack), [sentry-go](https://github.com/getsentry/sentry-go)

## Example
```go
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
