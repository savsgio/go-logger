# go-logger

[![Test status](https://github.com/savsgio/go-logger/actions/workflows/test.yml/badge.svg?branch=master)](https://github.com/savsgio/go-logger/actions?workflow=test)
[![Coverage Status](https://coveralls.io/repos/github/savsgio/go-logger/badge.svg?branch=master)](https://coveralls.io/github/savsgio/go-logger?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/savsgio/go-logger)](https://goreportcard.com/report/github.com/savsgio/go-logger)
[![GoDev](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/savsgio/go-logger/v3)
[![GitHub release](https://img.shields.io/github/release/savsgio/go-logger.svg)](https://github.com/savsgio/go-logger/releases)

Fast, lightweight, customizable and structured logger for Go.

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

| Level   | Code (constant) | Value (str)       |
| ------- | --------------- | ----------------- |
| Print   | logger.PRINT    |                   |
| Trace   | logger.TRACE    | trace / TRACE     |
| Fatal   | logger.FATAL    | fatal / FATAL     |
| Error   | logger.ERROR    | error / ERROR     |
| Warning | logger.WARNING  | warning / WARNING |
| Info    | logger.INFO     | info / INFO       |
| Debug   | logger.DEBUG    | debug / DEBUG     |

**NOTE:** _The default level of standard logger is **INFO**._

## Encoders:

- Text
- JSON
- Custom (your own encoder).

**NOTE:** _The default encoder of standard logger is **text**._

## Contributing

**Feel free to contribute it or fork me...** :wink:
