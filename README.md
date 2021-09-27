go-logger
=========

[![Test status](https://github.com/savsgio/go-logger/actions/workflows/test.yml/badge.svg?branch=master)](https://github.com/savsgio/go-logger/actions?workflow=test)
[![Coverage Status](https://coveralls.io/repos/github/savsgio/go-logger/badge.svg?branch=master)](https://coveralls.io/github/savsgio/go-logger?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/savsgio/go-logger)](https://goreportcard.com/report/github.com/savsgio/go-logger)
[![GoDev](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/savsgio/go-logger/v2)
[![GitHub release](https://img.shields.io/github/release/savsgio/go-logger.svg)](https://github.com/savsgio/go-logger/releases)

Lighweight wrapper for oficial Golang Log to adds support of levels for the log and reduce extra-allocations to zero.

## Install

- **WITH** Go modules:

```bash
go get github.com/savsgio/go-logger/v2
```

- **WITHOUT** Go modules:

```bash
go get github.com/savsgio/go-logger
```

## Supported Go versions:

- 1.16.x
- 1.15.x
- 1.14.x
- 1.13.x
- 1.12.x
- 1.11.x

## Benchmarks
```
Benchmark_Printf-12              3526083               335 ns/op             124 B/op          0 allocs/op
Benchmark_Errorf-12              3443049               375 ns/op             137 B/op          0 allocs/op
Benchmark_Warningf-12            3712971               317 ns/op             129 B/op          0 allocs/op
Benchmark_Infof-12               3668157               316 ns/op             128 B/op          0 allocs/op
Benchmark_Debugf-12              3719518               317 ns/op             127 B/op          0 allocs/op
```

## Levels:

|Literal |Code (constant) |Value (str) |
|--------|----------------|------------|
|Fatal   |logger.FATAL    |fatal       |
|Error   |logger.ERROR    |error       |
|Warning |logger.WARNING  |warning     |
|Info    |logger.INFO     |info        |
|Debug   |logger.DEBUG    |debug       |

**The default level for std logger is *logger.INFO***

## Output (*Important*)

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

Contributing
------------

**Feel free to contribute it or fork me...** :wink:
