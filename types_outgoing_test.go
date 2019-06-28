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

func TestForwardMessage(t *testing.T) {
	t.Run("NewAndWith", func(t *testing.T) {
		assert.Equal(t,
			&ForwardMessage{
				Peer: UserID(2),
				Message: MessageLocation{
					Chat:    ChatID(1),
					Message: MessageID(1),
				},
				DisableNotification: true,
			},
			NewForwardMessage(UserID(2), MessageLocation{
				Chat:    ChatID(1),
				Message: MessageID(1),
			}).WithNotification(false),
		)
	})

	t.Run("BuildSendRequest", func(t *testing.T) {
		msg := NewForwardMessage(UserID(2), MessageLocation{
			Chat:    ChatID(1),
			Message: MessageID(1),
		}).WithNotification(false)

		r, err := msg.BuildSendRequest()

		if assert.NoError(t, err) {
			args := extractArgs(r)

			assert.Equal(t, map[string]string{
				"chat_id":              "2",
				"message_id":           "1",
				"from_chat_id":         "1",
				"disable_notification": "true",
			}, args)
		}
	})
}

func TestPhotoMessage(t *testing.T) {
	inputFile := NewInputFileBytes("test.png", []byte("no data"))

	t.Run("NewAndWith", func(t *testing.T) {
		assert.Equal(t,
			&PhotoMessage{
				Peer:                UserID(1),
				Photo:               inputFile,
				Caption:             "test",
				ParseMode:           Markdown,
				DisableNotification: true,
				ReplyTo:             MessageID(1),
				ReplyMarkup:         NewForceReply(),
			},
			NewPhotoMessage(UserID(1), inputFile).
				WithCaption("test").
				WithParseMode(Markdown).
				WithNotification(false).
				WithReplyTo(MessageID(1)).
				WithReplyMarkup(NewForceReply()),
		)
	})

	t.Run("BuildSendRequest", func(t *testing.T) {
		msg := NewPhotoMessage(UserID(1), inputFile).
			WithCaption("test").
			WithParseMode(Markdown).
			WithNotification(false).
			WithReplyTo(MessageID(1)).
			WithReplyMarkup(NewForceReply())

		r, err := msg.BuildSendRequest()

		if assert.NoError(t, err) {
			args := extractArgs(r)

			assert.Equal(t, map[string]string{
				"chat_id":              "1",
				"caption":              "test",
				"parse_mode":           "markdown",
				"disable_notification": "true",
				"reply_markup":         `{"force_reply":true,"selective":false}`,
				"reply_to_message_id":  "1",
			}, args)

			files := extractFiles(r)

			assert.Equal(t, map[string]InputFile{
				"photo": inputFile,
			}, files)
		}
	})
}
