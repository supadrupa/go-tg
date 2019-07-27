package tg

import "github.com/pkg/errors"

// UpdateID is the update‘s unique identifier.
// Update identifiers start from a certain positive number and increase sequentially.
// If there are no new updates for at least a week,
// then identifier of the next update will be chosen randomly instead of sequentially.
type UpdateID int

func (id UpdateID) Next() UpdateID {
	return id + 1
}

// This object represents an incoming update.
// At most one of the optional parameters can be present in any given update.
type Update struct {
	// Update ID
	ID UpdateID `json:"update_id"`

	// Optional. New incoming message of any kind — text, photo, sticker, etc.
	Message *Message `json:"message,omitempty"`

	// Optional. New version of a message that is known to the bot and was edited
	EditedMessage *Message `json:"edited_message,omitempty"`

	// Optional. New incoming channel post of any kind — text, photo, sticker, etc.
	ChannelPost *Message `json:"channel_post,omitempty"`

	// Optional. New version of a channel post that is known to the bot and was edited
	EditedChannelPost *Message `json:"edited_channel_post,omitempty"`

	// Optional. New incoming inline query
	InlineQuery *InlineQuery `json:"inline_query,omitempty"`

	// Optional. The result of an inline query that was chosen by a user and sent to their chat partner.
	ChosenInlineResult *ChosenInlineResult `json:"chosen_inline_result,omitempty"`

	// Optional. The result of an inline query that was chosen by a user and sent to their chat partner.
	CallbackQuery *CallbackQuery `json:"callback_query,omitempty"`

	// Optional. New incoming shipping query. Only for invoices with flexible price
	ShippingQuery *ShippingQuery `json:"shipping_query,omitempty"`

	// Optional. New incoming pre-checkout query. Contains full information about checkout
	PreCheckoutQuery *PreCheckoutQuery `json:"pre_checkout_query,omitempty"`

	// Optional. New poll state. Bots receive only updates about stopped polls and polls, which are sent by the bot.
	Poll *Poll `json:"poll,omitempty"`
}

// Type returns update type.
func (u Update) Type() UpdateType {
	switch {
	case u.Message != nil:
		return UpdateMessage
	case u.EditedMessage != nil:
		return UpdateEditedMessage
	case u.ChannelPost != nil:
		return UpdateChannelPost
	case u.EditedChannelPost != nil:
		return UpdateEditedChannelPost
	case u.InlineQuery != nil:
		return UpdateInlineQuery
	case u.ChosenInlineResult != nil:
		return UpdateChosenInlineResult
	case u.CallbackQuery != nil:
		return UpdateCallbackQuery
	case u.ShippingQuery != nil:
		return UpdateShippingQuery
	case u.PreCheckoutQuery != nil:
		return UpdatePreCheckoutQuery
	case u.Poll != nil:
		return UpdatePoll
	default:
		return UpdateType(0)
	}
}

// UpdateSlice it's just alias.
type UpdateSlice []Update

// UpdateType represents type of incoming Update.
type UpdateType int8

const (
	UpdateMessage UpdateType = iota + 1
	UpdateEditedMessage
	UpdateChannelPost
	UpdateEditedChannelPost
	UpdateInlineQuery
	UpdateChosenInlineResult
	UpdateCallbackQuery
	UpdateShippingQuery
	UpdatePreCheckoutQuery
	UpdatePoll
)

// UpdateTypes list of all possible values of UpdateType.
var UpdateTypes = []UpdateType{
	UpdateMessage,
	UpdateEditedMessage,
	UpdateChannelPost,
	UpdateEditedChannelPost,
	UpdateInlineQuery,
	UpdateChosenInlineResult,
	UpdateCallbackQuery,
	UpdateShippingQuery,
	UpdatePreCheckoutQuery,
	UpdatePoll,
}

// String returns name of update type.
func (ut UpdateType) String() string {
	switch ut {
	case UpdateMessage:
		return "message"
	case UpdateEditedMessage:
		return "edited_message"
	case UpdateChannelPost:
		return "channel_post"
	case UpdateEditedChannelPost:
		return "edited_channel_post"
	case UpdateInlineQuery:
		return "inline_query"
	case UpdateChosenInlineResult:
		return "chosen_inline_result"
	case UpdateCallbackQuery:
		return "callback_query"
	case UpdateShippingQuery:
		return "shipping_query"
	case UpdatePreCheckoutQuery:
		return "pre_checkout_query"
	case UpdatePoll:
		return "poll"
	default:
		return ""
	}
}

var errUpdateTypeUnknown = errors.New("UpdateType.MarshalText unknown value")

func (ut UpdateType) MarshalText() ([]byte, error) {
	val := ut.String()

	if val == "" {
		return nil, errUpdateTypeUnknown
	}

	return []byte(val), nil
}

func (ut *UpdateType) UnmarshalText(data []byte) error {
	updateType, err := ParseUpdateType(string(data))
	if err != nil {
		return err
	}

	*ut = updateType

	return nil
}

// ParseUpdateType from string
func ParseUpdateType(v string) (UpdateType, error) {
	switch v {
	case "message":
		return UpdateMessage, nil
	case "edited_message":
		return UpdateEditedMessage, nil
	case "channel_post":
		return UpdateChannelPost, nil
	case "edited_channel_post":
		return UpdateEditedChannelPost, nil
	case "inline_query":
		return UpdateInlineQuery, nil
	case "chosen_inline_result":
		return UpdateChosenInlineResult, nil
	case "callback_query":
		return UpdateCallbackQuery, nil
	case "shipping_query":
		return UpdateShippingQuery, nil
	case "pre_checkout_query":
		return UpdatePreCheckoutQuery, nil
	case "poll":
		return UpdatePoll, nil
	default:
		return UpdateType(0), errUpdateTypeUnknown
	}
}
