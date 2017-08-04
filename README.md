# Telegam Logrus Hook

### A Telegram Hook for Logrus

This hook emits log messages (and corresponding fields) to the
Telegram API. Currently this hook will only emit messages for the
following levels:

* ERROR
* FATAL
* PANIC

## Installation

Install the package with:

```
go get github.com/rossmcdonald/telegram_hook
```

## Usage

Example usage:

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
```
