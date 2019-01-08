go-logger
=========

[![Build Status](https://travis-ci.org/savsgio/go-logger.svg?branch=master)](https://travis-ci.org/savsgio/go-logger)
[![Coverage Status](https://coveralls.io/repos/github/savsgio/go-logger/badge.svg?branch=master)](https://coveralls.io/github/savsgio/go-logger?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/savsgio/go-logger)](https://goreportcard.com/report/github.com/savsgio/go-logger)
[![GitHub release](https://img.shields.io/github/release/savsgio/go-logger.svg)](https://github.com/savsgio/go-logger/releases)

Lighweight wrapper for oficial Golang Log to adds support of levels for the log and reduce extra-allocations to zero.

## Benchmarks
```
Benchmark_Printf-8       3000000               399 ns/op              67 B/op          0 allocs/op
Benchmark_Errorf-8       3000000               479 ns/op             158 B/op          0 allocs/op
Benchmark_Warningf-8     3000000               458 ns/op             160 B/op          0 allocs/op
Benchmark_Infof-8        3000000               467 ns/op             156 B/op          0 allocs/op
Benchmark_Debugf-8       3000000               459 ns/op             158 B/op          0 allocs/op
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
