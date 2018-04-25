# Telegam Logrus Hook

[![Go Report Card](https://goreportcard.com/badge/github.com/rossmcdonald/telegram_hook)](https://goreportcard.com/report/github.com/rossmcdonald/telegram_hook) [![GoDoc](https://godoc.org/github.com/rossmcdonald/telegram_hook?status.svg)](https://godoc.org/github.com/rossmcdonald/telegram_hook)

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
	"time"
	log "github.com/Sirupsen/logrus"
	"github.com/rossmcdonald/telegram_hook"
)

func main() {
	hook, err := telegram_hook.NewTelegramHook(
		"MyCoolApp",
		"MYTELEGRAMTOKEN",
		"@mycoolusername",
		telegram_hook.WithAsync(true),
		telegram_hook.WithTimeout(30 * time.Second),
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

Also you can set custom http.Client to use SOCKS5 proxy for example

```go
import (
	"time"
	"net/http"
	
	"golang.org/x/net/proxy"
	log "github.com/Sirupsen/logrus"
	"github.com/rossmcdonald/telegram_hook"
)

func main() {
	httpTransport := &http.Transport{}
    httpClient := &http.Client{Transport: httpTransport}
    dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:54321", nil, proxy.Direct)
    httpTransport.Dial = dialer.Dial
    
	hook, err := telegram_hook.NewTelegramHookWithClient(
		"MyCoolApp",
		"MYTELEGRAMTOKEN",
		"@mycoolusername",
		httpClient,
		telegram_hook.WithAsync(true),
		telegram_hook.WithTimeout(30 * time.Second),
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
