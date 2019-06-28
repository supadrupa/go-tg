package tg

import (
	"encoding/json"
)

// ReplyMarkup define generic interface for all reply markups.
//
// Types implementing this interface:
//  - InlineKeyboardMarkup
//  - KeyboardReplyMarkup
//  - ForceReply
//  - ReplyKeyboardRemove
type ReplyMarkup interface {
	EncodeReplyMarkup() (string, error)
}

// LoginURL represents a parameter of the inline keyboard button used to automatically authorize a user.
// Serves as a great replacement for the Telegram Login Widget when the user is coming from Telegram. All the user needs to do is tap/click a button and confirm that they want to log in.
//
// Telegram apps support these buttons as of version 5.7.
type LoginURL struct {
	// An HTTP URL to be opened with user authorization data added to the query string
	// when the button is pressed.
	// If the user refuses to provide authorization data,
	// the original URL without information about the user will be opened.
	URL string `json:"url"`

	// Optional. New text of the button in forwarded messages.
	ForwardText string `json:"forward_text,omitempty"`

	// Optional. Username of a bot, which will be used for user authorization.
	// If not specified, the current bot's username will be assumed.
	// The url's domain must be the same as the domain linked with the bot.
	BotUsername Username `json:"bot_username,omitempty"`

	// Optional. Pass true to request the permission for your bot
	// to send messages to the user.
	RequestWriteAccess bool `json:"request_write_access,omitempty"`
}

// NewLoginURL creates a new LoginURL object.
func NewLoginURL(url string) *LoginURL {
	return &LoginURL{
		URL: url,
	}
}

// WithForwardText sets new text of the button in forwarded messages.
func (url *LoginURL) WithForwardText(text string) *LoginURL {
	url.ForwardText = text
	return url
}

// WithBotUsername sets username of a bot, which will be used for user authorization.
func (url *LoginURL) WithBotUsername(username Username) *LoginURL {
	url.BotUsername = username
	return url
}

// WithRequestWriteAccess pass true to request the permission for your bot.
func (url *LoginURL) WithRequestWriteAccess(yes bool) *LoginURL {
	url.RequestWriteAccess = yes
	return url
}

// InlineKeyboardButton represents one button of an inline keyboard.
// You must use exactly one of the optional fields.
type InlineKeyboardButton struct {
	// Label text on the button
	Text string `json:"text"`

	// Optional. HTTP or tg:// url to be opened when button is pressed
	URL string `json:"url,omitempty"`

	// Optional. An HTTP URL used to automatically authorize the user. Can be used as a replacement for the Telegram Login Widget.
	LoginURL *LoginURL `json:"login_url,omitempty"`

	// Optional. Data to be sent in a callback query to the bot when button is pressed, 1-64 bytes
	CallbackData string `json:"callback_data,omitempty"`

	// Optional. If set, pressing the button will prompt the user
	// to select one of their chats, open that chat and
	// insert the bot‘s username and the specified inline query in the input field.
	// Can be empty, in which case just the bot’s username will be inserted.
	SwitchInlineQuery string `json:"switch_inline_query,omitempty"`

	// Optional. If set, pressing the button will insert the bot‘s username and
	// the specified inline query in the current chat's input field.
	// Can be empty, in which case only the bot’s username will be inserted.
	SwitchInlineQueryCurrentChat string `json:"switch_inline_query_current_chat,omitempty"`

	// Optional. Use only in invoice message.
	// Specify True, to send a Pay button.
	//
	// NOTE: This type of button must always be the first button in the first row.
	Pay bool `json:"pay,omitempty"`
}

// NewInlineKeyboardButtonURL creates inline keyboard button with a URL that will be opened when user click.
func NewInlineKeyboardButtonURL(text string, url string) InlineKeyboardButton {
	return InlineKeyboardButton{
		Text: text,
		URL:  url,
	}
}

// NewInlineKeyboardButtonLogin creates inline keyboard button used to automatically authorize the user.
func NewInlineKeyboardButtonLogin(text string, url *LoginURL) InlineKeyboardButton {
	return InlineKeyboardButton{
		Text:     text,
		LoginURL: url,
	}
}

// NewInlineKeyboardButtonCallback creates inline keyboard button data with specified callback data.
func NewInlineKeyboardButtonCallback(text string, cbd string) InlineKeyboardButton {
	return InlineKeyboardButton{
		Text:         text,
		CallbackData: cbd,
	}
}

// NewInlineKeyboardButtonSwitchInlineQuery creates inline keyboard for prompt the user
// to select one of their chats, open that chat and
// insert the bot‘s username and the specified inline query in the input field.
func NewInlineKeyboardButtonSwitchInline(text string, query string) InlineKeyboardButton {
	return InlineKeyboardButton{
		Text:              text,
		SwitchInlineQuery: query,
	}
}

// NewInlineKeyboardButtonSwitchInlineQueryCurrent creates inline keyboard
// for insert the bot‘s username and the specified inline query
// in the current chat's input field.
func NewInlineKeyboardButtonSwitchInlineCurrent(text string, query string) InlineKeyboardButton {
	return InlineKeyboardButton{
		Text:                         text,
		SwitchInlineQueryCurrentChat: query,
	}
}

// NewInlineKeyboardButtonPay creates Pay inline keyboard button.
// Use it only in invoice messages.
func NewInlineKeyboardButtonPay(text string) InlineKeyboardButton {
	return InlineKeyboardButton{
		Text: text,
		Pay:  true,
	}
}

// InlineKeyboardRow one row of InlineKeyboardButton's
type InlineKeyboardRow []InlineKeyboardButton

// NewInlineKeyboardRow creates a InlineKeyboardButton's row.
func NewInlineKeyboardRow(buttons ...InlineKeyboardButton) InlineKeyboardRow {
	return InlineKeyboardRow(buttons)
}

// InlineKeyboardMarkup represents an inline keyboard that appears right next to the message it belongs to.
type InlineKeyboardMarkup struct {
	InlineKeyboard []InlineKeyboardRow `json:"inline_keyboard"`
}

// TODO: the possibility of automatic placement of inline keyboard buttons

// NewInlineKeyboardMarkup creates a inline keyboard markup from provided rows.
func NewInlineKeyboardMarkup(rows ...InlineKeyboardRow) InlineKeyboardMarkup {
	return InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}

// EncodeReplyMarkup returns encoded (JSON) representation of InlineKeyboardMarkup.
func (ikb InlineKeyboardMarkup) EncodeReplyMarkup() (string, error) {
	obj, err := json.Marshal(ikb)
	return string(obj), err
}

// KeyboardButton represents one button of the reply keyboard.
// Optional fields are mutually exclusive.
type KeyboardButton struct {
	// Text of the button.
	// If none of the optional fields are used,
	// it will be sent as a message when the button is pressed
	Text string `json:"text"`

	// Optional. If True, the user's phone number
	// will be sent as a contact when the button is pressed.
	// Available in private chats only
	RequestConact bool `json:"request_conact,omitempty"`

	// Optional. If True, the user's current location
	// will be sent when the button is pressed.
	// Available in private chats only
	RequestLocation bool `json:"request_location,omitempty"`
}

// NewKeyboardButton creates a simple text keyboard button.
func NewKeyboardButton(text string) KeyboardButton {
	return KeyboardButton{
		Text: text,
	}
}

// NewKeyboardButtonContact creates a keyboard button for request user contact.
func NewKeyboardButtonContact(text string) KeyboardButton {
	return KeyboardButton{
		Text:          text,
		RequestConact: true,
	}
}

// NewKeyboardButtonContact creates a keyboard button for request user location.
func NewKeyboardButtonLocation(text string) KeyboardButton {
	return KeyboardButton{
		Text:            text,
		RequestLocation: true,
	}
}

// KeyboardRow one row of KeyboardButton's.
type KeyboardRow []KeyboardButton

// NewKeyboardRow creates KeyboardButton's row.
func NewKeyboardRow(buttons ...KeyboardButton) KeyboardRow {
	return KeyboardRow(buttons)
}

// ReplyKeyboardMarkup represents a custom keyboard with reply options.
type ReplyKeyboardMarkup struct {
	Keyboard []KeyboardRow `json:"keyboard"`

	// Optional. Requests clients to resize the keyboard vertically for optimal fit
	// (e.g., make the keyboard smaller if there are just two rows of buttons).
	// Defaults to false, in which case the custom keyboard is always
	// of the same height as the app's standard keyboard.
	Resize bool `json:"resize_keyboard,omitempty"`

	// Optional. Requests clients to hide the keyboard as soon as it's been used.
	// The keyboard will still be available, but clients will automatically display
	// the usual letter-keyboard in the chat – the user can press a special button
	// in the input field to see the custom keyboard again. D
	OneTime bool `json:"one_time_keyboard,omitempty"`

	// Use this parameter if you want to show the keyboard to specific users only.
	Selective bool `json:"selective,omitempty"`
}

// NewReplyKeyboardMarkup creates ReplyKeyboardMarkup
func NewReplyKeyboardMarkup(rows ...KeyboardRow) ReplyKeyboardMarkup {
	return ReplyKeyboardMarkup{
		Keyboard: rows,
	}
}

// WithResize makes keyboard resizable on client UI.
func (kb ReplyKeyboardMarkup) WithResize(yes bool) ReplyKeyboardMarkup {
	kb.Resize = yes
	return kb
}

// WithOneTime makes keyboard one-time used.
func (kb ReplyKeyboardMarkup) WithOneTime(yes bool) ReplyKeyboardMarkup {
	kb.OneTime = yes
	return kb
}

// WithSelective makes keyboard selective.
func (kb ReplyKeyboardMarkup) WithSelective(yes bool) ReplyKeyboardMarkup {
	kb.Selective = yes
	return kb
}

// EncodeReplyMarkup returns encoded (JSON) representation of InlineKeyboardMarkup.
func (kb ReplyKeyboardMarkup) EncodeReplyMarkup() (string, error) {
	obj, err := json.Marshal(kb)
	return string(obj), err
}

// Upon receiving a message with this object, T
// Telegram clients will display a reply interface to the user
// Act as if the user has selected the bot‘s message and tapped ’Reply'.
// This can be extremely useful if you want to create
// user-friendly step-by-step interfaces
// without having to sacrifice privacy mode.
type ForceReply struct {
	Selective bool `json:"selective"`
}

// NewForceReply creates ForceReply reply markup.
func NewForceReply() ForceReply {
	return ForceReply{}
}

// WithSelective makes keyboard selective.
func (fr ForceReply) WithSelective(yes bool) ForceReply {
	fr.Selective = yes
	return fr
}

func (fr ForceReply) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ForceReply bool `json:"force_reply"`
		Selective  bool `json:"selective"`
	}{
		ForceReply: true,
		Selective:  fr.Selective,
	})
}

// EncodeReplyMarkup returns encoded (JSON) representation of InlineKeyboardMarkup.
func (fr ForceReply) EncodeReplyMarkup() (string, error) {
	obj, err := json.Marshal(fr)
	return string(obj), err
}

// ReplyKeyboardRemove upon receiving a message with this object,
// Telegram clients will remove the current custom keyboard
// and display the default letter-keyboard.
type ReplyKeyboardRemove struct {
	Selective bool `json:"selective"`
}

// NewReplyKeyboardRemove creates ReplyKeyboardRemove reply markup.
func NewReplyKeyboardRemove() ReplyKeyboardRemove {
	return ReplyKeyboardRemove{}
}

// WithSelective makes keyboard remove selective.
func (kr ReplyKeyboardRemove) WithSelective(yes bool) ReplyKeyboardRemove {
	kr.Selective = yes
	return kr
}

// EncodeReplyMarkup returns encoded (JSON) representation of InlineKeyboardMarkup.
func (kr ReplyKeyboardRemove) EncodeReplyMarkup() (string, error) {
	obj, err := json.Marshal(kr)
	return string(obj), err
}

func (kr ReplyKeyboardRemove) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		RemoveKeyboard bool `json:"remove_keyboard"`
		Selective      bool `json:"selective"`
	}{
		RemoveKeyboard: true,
		Selective:      kr.Selective,
	})
}
