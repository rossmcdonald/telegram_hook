# Telegam Logrus Hook

[![Go Report Card](https://goreportcard.com/badge/github.com/rossmcdonald/telegram_hook)](https://goreportcard.com/report/github.com/rossmcdonald/telegram_hook)

This hook emits log messages (and corresponding fields) to the
Telegram API
for [logrus](https://github.com/Sirupsen/logrus). Currently this hook
will only emit messages for the following levels:

* `ERROR`
* `FATAL`
* `PANIC`

## Installation

Install the package with:

```
go get github.com/rossmcdonald/telegram_hook
```

## Usage

See the tests for working examples. Also:

```go
import (
	log "github.com/Sirupsen/logrus"
	"github.com/rossmcdonald/telegram_hook"
)

func main() {
	hook, err := telegram_hook.NewTelegramHook(
		"MyCoolApp",
		"MYTELEGRAMTOKEN",
		"@mycoolusername",
	)
	if err != nil {
		log.Fatalf("Encountered error when creating Telegram hook: %s", err)
	}
	log.AddHook(hook)
	
	// Receive messages on failures
	log.Errorf("Uh oh...")
	...
	
}
```

## License

MIT
