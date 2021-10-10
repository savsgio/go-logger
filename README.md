# go-logger

[![Test status](https://github.com/savsgio/go-logger/actions/workflows/test.yml/badge.svg?branch=master)](https://github.com/savsgio/go-logger/actions?workflow=test)
[![Coverage Status](https://coveralls.io/repos/github/savsgio/go-logger/badge.svg?branch=master)](https://coveralls.io/github/savsgio/go-logger?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/savsgio/go-logger)](https://goreportcard.com/report/github.com/savsgio/go-logger)
[![GoDev](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/savsgio/go-logger/v2)
[![GitHub release](https://img.shields.io/github/release/savsgio/go-logger.svg)](https://github.com/savsgio/go-logger/releases)

Fast, lightweight and structured logger for Go.

## Install

```bash
go get github.com/savsgio/go-logger/v3
```

## Supported Go versions:

- 1.17.x
- 1.16.x
- 1.15.x
- 1.14.x
- 1.13.x

## Levels:

| Level   | Code (constant) | Value (str)         |
| ------- | --------------- | ------------------- |
| Fatal   | logger.FATAL    | fatal / FATAL       |
| Error   | logger.ERROR    | error / ERROR       |
| Warning | logger.WARNING  | warning / WARNING   |
| Info    | logger.INFO     | info / INFO         |
| Debug   | logger.DEBUG    | debug / DEBUG       |

**The default level for std logger is _logger.INFO_**

## Output (_Important_)

By default, output of log is `os.Stderr`, but you can customize it with other `io.Writer`.

### Format

Example of format ouput:

```text
# Standar instance of logger
2018/03/16 12:26:48 - DEBUG - Listening on http://0.0.0.0:8000

# Own instance of logger
2018/03/16 12:26:48 - <name of instance> - DEBUG - Hello gopher
```

## How to use:

Call logger ever you want with:

```go
logger.Debugf("Listening on %s", "http://0.0.0.0:8000")
```

If you want use your own instance of logger:

```go
myLog := logger.New("myInstance", logger.DEBUG, &bytes.Buffer{})  // Change level

myLog.Warning("Hello gopher")
```

### Example

```go
import bytes
import "github.com/savsgio/go-logger"


func myFunc(){
    logger.SetLevel(logger.DEBUG) // Optional (default: logger.INFO)

    ....
    logger.Infof("Hi, you are using %s/%s", "savsgio", "go-logger")
    ....

    customOutput = &bytes.Buffer{}
    myLog := logger.New("myInstance", logger.INFO, customOutput)  // Change level

    myLog.Warning("Hello gopher")
}
```

## Contributing

**Feel free to contribute it or fork me...** :wink:
