package tg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTextMessage(t *testing.T) {
	t.Run("NewAndWith", func(t *testing.T) {
		assert.Equal(t,
			&TextMessage{
				Peer:                  UserID(1),
				Text:                  "test",
				ParseMode:             Markdown,
				DisableWebPagePreview: true,
				DisableNotification:   true,
				ReplyTo:               MessageID(1),
				ReplyMarkup:           NewForceReply(),
			},
			NewTextMessage(UserID(1), "test").
				WithParseMode(Markdown).
				WithNotification(false).
				WithWebPagePreview(false).
				WithReplyTo(MessageID(1)).
				WithReplyMarkup(NewForceReply()),
		)
	})

	t.Run("BuildSendRequest", func(t *testing.T) {
		msg := NewTextMessage(UserID(1), "test").
			WithParseMode(Markdown).
			WithNotification(false).
			WithWebPagePreview(false).
			WithReplyTo(MessageID(1)).
			WithReplyMarkup(NewForceReply())

		r, err := msg.BuildSendRequest()

		if assert.NoError(t, err) {
			args := extractArgs(r)

			assert.Equal(t, map[string]string{
				"chat_id":                  "1",
				"text":                     "test",
				"parse_mode":               "markdown",
				"disable_web_page_preview": "true",
				"disable_notification":     "true",
				"reply_markup":             `{"force_reply":true,"selective":false}`,
				"reply_to_message_id":      "1",
			}, args)
		}
	})
}
