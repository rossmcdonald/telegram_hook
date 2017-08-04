package telegram_hook_test

import (
	"os"
	"testing"

	"github.com/rossmcdonald/telegram_hook"
)

func TestNewHook(t *testing.T) {
	_, err := telegram_hook.NewTelegramHook("", "", "")
	if err == nil {
		t.Errorf("No error on invalid Telegram API token.")
	}

	_, err = telegram_hook.NewTelegramHook("", os.Getenv("TELEGRAM_TOKEN"), "")
	if err != nil {
		t.Errorf("Error on valid Telegram API token: %s", err)
	}
}
