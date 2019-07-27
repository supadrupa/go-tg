package tg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateID_Next(t *testing.T) {
	id := UpdateID(0)

	assert.Equal(t, UpdateID(1), id.Next())
}

func TestUpdate_Type(t *testing.T) {
	for _, tt := range []struct {
		Input  Update
		Output UpdateType
	}{
		// message
		{Update{Message: &Message{}}, UpdateMessage},

		// edited_message
		{Update{EditedMessage: &Message{}}, UpdateEditedMessage},

		// channel_post
		{Update{ChannelPost: &Message{}}, UpdateChannelPost},

		// edited_channel_post
		{Update{EditedChannelPost: &Message{}}, UpdateEditedChannelPost},

		// inline_query
		{Update{InlineQuery: &InlineQuery{}}, UpdateInlineQuery},

		// chosen_inline_result
		{Update{ChosenInlineResult: &ChosenInlineResult{}}, UpdateChosenInlineResult},

		// callback_query
		{Update{CallbackQuery: &CallbackQuery{}}, UpdateCallbackQuery},

		// shipping_query
		{Update{ShippingQuery: &ShippingQuery{}}, UpdateShippingQuery},

		// pre_checkout_query
		{Update{PreCheckoutQuery: &PreCheckoutQuery{}}, UpdatePreCheckoutQuery},

		// poll
		{Update{Poll: &Poll{}}, UpdatePoll},

		// unknown
		{Update{}, UpdateType(0)},
	} {
		assert.Equal(t,
			tt.Output,
			tt.Input.Type(),
		)
	}
}

func TestUpdateType_String(t *testing.T) {
	for _, tt := range []struct {
		Type       UpdateType
		Excepted   string
		MarshalErr bool
	}{
		{UpdateMessage, "message", false},
		{UpdateEditedMessage, "edited_message", false},
		{UpdateChannelPost, "channel_post", false},
		{UpdateEditedChannelPost, "edited_channel_post", false},
		{UpdateInlineQuery, "inline_query", false},
		{UpdateChosenInlineResult, "chosen_inline_result", false},
		{UpdateCallbackQuery, "callback_query", false},
		{UpdateShippingQuery, "shipping_query", false},
		{UpdatePreCheckoutQuery, "pre_checkout_query", false},
		{UpdatePoll, "poll", false},
		{UpdateType(0), "", true},
	} {
		assert.Equal(t,
			tt.Excepted,
			tt.Type.String(),
		)

		text, err := tt.Type.MarshalText()

		if tt.MarshalErr {
			assert.Error(t, err)
		} else {
			assert.Equal(t,
				tt.Excepted,
				string(text),
			)
		}

	}
}

func TestUpdateType_UnmarshalText(t *testing.T) {
	for _, tt := range []struct {
		Excepted     UpdateType
		Input        string
		UnmarshalErr bool
	}{
		{UpdateMessage, "message", false},
		{UpdateEditedMessage, "edited_message", false},
		{UpdateChannelPost, "channel_post", false},
		{UpdateEditedChannelPost, "edited_channel_post", false},
		{UpdateInlineQuery, "inline_query", false},
		{UpdateChosenInlineResult, "chosen_inline_result", false},
		{UpdateCallbackQuery, "callback_query", false},
		{UpdateShippingQuery, "shipping_query", false},
		{UpdatePreCheckoutQuery, "pre_checkout_query", false},
		{UpdatePoll, "poll", false},
		{UpdateType(0), "", true},
	} {
		var updateType UpdateType

		if tt.UnmarshalErr {
			assert.Error(t, updateType.UnmarshalText([]byte(tt.Input)))
		} else {
			err := updateType.UnmarshalText([]byte(tt.Input))
			if assert.NoError(t, err) {
				assert.Equal(t, tt.Excepted, updateType)
			}
		}
	}
}
