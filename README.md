# fxevent-zerolog

A structured logger for [Uber's Fx](https://github.com/uber-go/fx) that leverages the performance and flexibility of [zerolog](https://github.com/rs/zerolog). This package implements the `fxevent.Logger` interface, enabling seamless integration of zerolog into your Fx applications.

## Motivation

While Fx provides a powerful dependency injection framework, its default logging is not always ideal for high-performance or structured logging needs. This package bridges that gap, allowing you to use zerolog as the event logger for Fx, with full support for log levels, structured fields, and custom configuration.

## Features

- Implements `fxevent.Logger` for drop-in Fx compatibility
- Full support for zerolog's structured logging and log levels
- Customizable log and error levels
- Safe defaults (no panics on nil logger)
- Comprehensive test coverage

## Installation

```
go get github.com/amari/fxevent-zerolog
```

## Usage

```go
import (
	"github.com/rs/zerolog"
	"go.uber.org/fx"
	"github.com/amari/fxevent-zerolog"
)

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	app := fx.New(
        fx.Supply(&logger),
		fx.WithLogger(fxeventzerolog.New),
	)
	app.Run()
}
```

## API

See [zerolog.go](./zerolog.go) for full documentation and implementation details.

## Testing

This project includes comprehensive unit tests covering all major event types and behaviors. Run tests with:

```
go test -v
```

## License

This project is [MIT](./LICENSE) licensed.

