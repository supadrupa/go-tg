package tg

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoginURL_New(t *testing.T) {

	url := NewLoginURL("https://test.com").
		WithForwardText("Try It").
		WithBotUsername("channely_bot").
		WithRequestWriteAccess(true)

	assert.Equal(t,
		&LoginURL{
			URL:                "https://test.com",
			ForwardText:        "Try It",
			BotUsername:        "channely_bot",
			RequestWriteAccess: true,
		},
		url,
	)
}

func TestInlineKeyboardButton_New(t *testing.T) {
	for _, tt := range []struct {
		Actual   InlineKeyboardButton
		Excepted InlineKeyboardButton
	}{
		{
			NewInlineKeyboardButtonURL("test", "https://google.com"),
			InlineKeyboardButton{
				Text: "test",
				URL:  "https://google.com",
			},
		},
		{
			NewInlineKeyboardButtonLogin("test", NewLoginURL("https://google.com")),
			InlineKeyboardButton{
				Text:     "test",
				LoginURL: &LoginURL{URL: "https://google.com"},
			},
		},
		{
			NewInlineKeyboardButtonCallback("test", "test"),
			InlineKeyboardButton{
				Text:         "test",
				CallbackData: "test",
			},
		},
		{
			NewInlineKeyboardButtonSwitchInline("test", "test"),
			InlineKeyboardButton{
				Text:              "test",
				SwitchInlineQuery: "test",
			},
		},
		{
			NewInlineKeyboardButtonSwitchInlineCurrent("test", "test"),
			InlineKeyboardButton{
				Text:                         "test",
				SwitchInlineQueryCurrentChat: "test",
			},
		},
		{
			NewInlineKeyboardButtonPay("test"),
			InlineKeyboardButton{
				Text: "test",
				Pay:  true,
			},
		},
	} {
		assert.Equal(t, tt.Excepted, tt.Actual)
	}
}

func TestInlineKeyboardMarkup_New(t *testing.T) {
	assert.Equal(t,
		InlineKeyboardMarkup{
			[]InlineKeyboardRow{
				{
					{Text: "test", URL: "https://google.com"},
				},
			},
		},

		NewInlineKeyboardMarkup(
			NewInlineKeyboardRow(
				NewInlineKeyboardButtonURL("test", "https://google.com"),
			),
		),
	)
}

func TestKeyboardButton_New(t *testing.T) {
	for _, tt := range []struct {
		Actual   KeyboardButton
		Excepted KeyboardButton
	}{
		{
			NewKeyboardButton("test"),
			KeyboardButton{Text: "test"},
		},
		{
			NewKeyboardButtonContact("test"),
			KeyboardButton{Text: "test", RequestConact: true},
		},
		{
			NewKeyboardButtonLocation("test"),
			KeyboardButton{Text: "test", RequestLocation: true},
		},
	} {
		assert.Equal(t, tt.Excepted, tt.Actual)
	}
}

func TestReplyKeyboardMarkup(t *testing.T) {
	assert.Equal(t,
		ReplyKeyboardMarkup{
			Keyboard: []KeyboardRow{
				{
					KeyboardButton{Text: "Simple"},
				},
				{
					KeyboardButton{Text: "Location", RequestLocation: true},
					KeyboardButton{Text: "Contact", RequestConact: true},
				},
			},
			Resize:    true,
			Selective: true,
			OneTime:   true,
		},
		NewReplyKeyboardMarkup(
			NewKeyboardRow(
				NewKeyboardButton("Simple"),
			),
			NewKeyboardRow(
				NewKeyboardButtonLocation("Location"),
				NewKeyboardButtonContact("Contact"),
			),
		).WithResize(true).WithOneTime(true).WithSelective(true),
	)
}

func TestForceReply(t *testing.T) {
	forceReply := NewForceReply().WithSelective(true)

	raw, err := json.Marshal(forceReply)
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"force_reply": true,
		"selective": true
	}`, string(raw))
}

func TestReplyKeyboardRemove(t *testing.T) {
	forceReply := NewReplyKeyboardRemove().WithSelective(true)

	raw, err := json.Marshal(forceReply)
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"remove_keyboard": true,
		"selective": true
	}`, string(raw))
}
