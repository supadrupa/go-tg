package tg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageIdentity(t *testing.T) {
	for _, tt := range []struct {
		Identity MessageIdentity
		Result   MessageID
	}{
		{MessageID(1), MessageID(1)},
		{Message{ID: MessageID(2)}, MessageID(2)},
	} {
		assert.Equal(t,
			tt.Result,
			tt.Identity.GetMessageID(),
		)
	}
}

func TestMessageLocation(t *testing.T) {

	for _, tt := range []struct {
		Identity MessageIdentityFull
		Peer     Peer
		Message  MessageIdentity
	}{
		{
			Identity: MessageLocation{
				Chat:    ChatID(1),
				Message: MessageID(1),
			},
			Peer:    ChatID(1),
			Message: MessageID(1),
		},
		{
			Identity: Message{
				ID: MessageID(2),
				Chat: Chat{
					ID: ChatID(2),
				},
			},
			Peer:    ChatID(2),
			Message: MessageID(2),
		},
	} {
		peer, msg := tt.Identity.GetMessageLocation()

		assert.Equal(t,
			tt.Peer,
			peer,
		)

		assert.Equal(t,
			tt.Message,
			msg,
		)
	}
}
