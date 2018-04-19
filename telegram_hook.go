package telegram_hook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// TelegramHook to send logs via the Telegram API.
type TelegramHook struct {
	AppName     string
	c           *http.Client
	authToken   string
	targetID    string
	apiEndpoint string
	async bool
}

// apiRequest encapsulates the request structure we are sending to the
// Telegram API.
type apiRequest struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode,omitempty"`
}

// apiResponse encapsulates the response structure received from the
// Telegram API.
type apiResponse struct {
	Ok        bool         `json:"ok"`
	ErrorCode *int         `json:"error_code,omitempty"`
	Desc      *string      `json:"description,omitempty"`
	Result    *interface{} `json:"result,omitempty"`
}

// Config defines a method for additional configuration when instantiating TelegramHook
type Config func(*TelegramHook)

// NewTelegramHook creates a new instance of a hook targeting the
// Telegram API.
func NewTelegramHook(appName, authToken, targetID string, config ...Config) (*TelegramHook, error) {
	client := &http.Client{}
	return NewTelegramHookWithClient(appName, authToken, targetID, client, config...)
}

// NewTelegramHook creates a new instance of a hook targeting the
// Telegram API with custom http.Client.
func NewTelegramHookWithClient(appName, authToken, targetID string, client *http.Client, config ...Config) (*TelegramHook, error) {
	apiEndpoint := fmt.Sprintf(
		"https://api.telegram.org/bot%s",
		authToken,
	)
	h := TelegramHook{
		AppName:     appName,
		c:           client,
		authToken:   authToken,
		targetID:    targetID,
		apiEndpoint: apiEndpoint,
		async: false,
	}

	for _, c := range config {
		c(&h)
	}

	// Verify the API token is valid and correct before continuing
	err := h.verifyToken()
	if err != nil {
		return nil, err
	}

	return &h, nil
}

// Async sets logging to telegram as asynchronous
func WithAsync(b bool) Config {
	return func(hook *TelegramHook) {
		hook.async = b
	}
}

// Timeout sets http call timeout for telegram client
func WithTimeout(t time.Duration) Config {
	return func(hook *TelegramHook) {
		if t > 0 {
			hook.c.Timeout = t
		}
	}
}

// verifyToken issues a test request to the Telegram API to ensure the
// provided token is correct and valid.
func (hook *TelegramHook) verifyToken() error {
	endpoint := strings.Join([]string{hook.apiEndpoint, "getme"}, "/")
	res, err := hook.c.Get(endpoint)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	apiRes := apiResponse{}
	err = json.NewDecoder(res.Body).Decode(&apiRes)
	if err != nil {
		return err
	}

	if !apiRes.Ok {
		// Received an error from the Telegram API
		msg := "Received error response from Telegram API"
		if apiRes.ErrorCode != nil {
			msg = fmt.Sprintf("%s (error code %d)", msg, *apiRes.ErrorCode)
		}
		if apiRes.Desc != nil {
			msg = fmt.Sprintf("%s: %s", msg, *apiRes.Desc)
		}
		j, _ := json.MarshalIndent(apiRes, "", "\t")
		msg = fmt.Sprintf("%s\n%s", msg, j)
		return fmt.Errorf(msg)
	}

	return nil
}

// sendMessage issues the provided message to the Telegram API.
func (hook *TelegramHook) sendMessage(msg string) error {
	apiReq := apiRequest{
		ChatID:    hook.targetID,
		Text:      msg,
		ParseMode: "HTML",
	}
	b, err := json.Marshal(apiReq)
	if err != nil {
		return err
	}

	res, err := hook.c.Post(
		strings.Join([]string{hook.apiEndpoint, "sendmessage"}, "/"),
		"application/json",
		bytes.NewReader(b),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Encountered error when issuing request to Telegram API, %v", err)
		return err
	}
	defer res.Body.Close()

	apiRes := apiResponse{}
	err = json.NewDecoder(res.Body).Decode(&apiRes)
	if err != nil {
		return err
	}

	if !apiRes.Ok {
		// Received an error from the Telegram API
		msg := "Received error response from Telegram API"
		if apiRes.ErrorCode != nil {
			msg = fmt.Sprintf("%s (error code %d)", msg, *apiRes.ErrorCode)
		}
		if apiRes.Desc != nil {
			msg = fmt.Sprintf("%s: %s", msg, *apiRes.Desc)
		}
		return fmt.Errorf(msg)
	}

	return nil
}

// createMessage crafts an HTML-formatted message to send to the
// Telegram API.
func (hook *TelegramHook) createMessage(entry *logrus.Entry) string {
	var msg string

	switch entry.Level {
	case logrus.PanicLevel:
		msg = "<b>PANIC</b>"
	case logrus.FatalLevel:
		msg = "<b>FATAL</b>"
	case logrus.ErrorLevel:
		msg = "<b>ERROR</b>"
	}

	msg = strings.Join([]string{msg, hook.AppName}, "@")
	msg = strings.Join([]string{msg, entry.Message}, " - ")
	fields, err := json.MarshalIndent(entry.Data, "", "\t")
	if err == nil {
		msg = strings.Join([]string{msg, "<pre>"}, "\n")
		msg = strings.Join([]string{msg, string(fields)}, "\n")
		msg = strings.Join([]string{msg, "</pre>"}, "\n")
	}
	return msg
}

// Fire emits a log message to the Telegram API.
func (hook *TelegramHook) Fire(entry *logrus.Entry) error {
	msg := hook.createMessage(entry)

	if hook.async {
		go hook.sendMessage(msg)
		return nil
	}

	err := hook.sendMessage(msg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to send message, %v", err)
		return err
	}

	return nil
}

// Levels returns the log levels that the hook should be enabled for.
func (hook *TelegramHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}
