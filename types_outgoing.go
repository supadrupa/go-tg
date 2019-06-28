package tg

// TextMessage represents simple text message.
//
// Related API method: https://core.telegram.org/bots/api#sendmessage
type TextMessage struct {
	// Recipient of the message
	Peer Peer

	// Content
	Text string

	// Text parse mode
	ParseMode ParseMode

	// Pass true if you need to disable web page preview
	DisableWebPagePreview bool

	// Pass true for send message silent.
	DisableNotification bool

	// Reply to message identity.
	ReplyTo MessageIdentity

	// Reply markup of the message.
	ReplyMarkup ReplyMarkup
}

// NewTextMessage creates a simple text message for specified peer and provided text.
func NewTextMessage(to Peer, text string) *TextMessage {
	return &TextMessage{
		Peer: to,
		Text: text,
	}
}

// WithParseMode sets text message parse mode.
func (msg *TextMessage) WithParseMode(pm ParseMode) *TextMessage {
	msg.ParseMode = pm
	return msg
}

// WithWebPagePreview enable or disable message first link web page preview. (default: enabled).
func (msg *TextMessage) WithWebPagePreview(yes bool) *TextMessage {
	msg.DisableWebPagePreview = !yes
	return msg
}

// WithNotification enable or disable notification (default: enabled).
func (msg *TextMessage) WithNotification(yes bool) *TextMessage {
	msg.DisableNotification = !yes
	return msg
}

// WithReplyTo sets ids of original message, if message is reply.
func (msg *TextMessage) WithReplyTo(msgID MessageIdentity) *TextMessage {
	msg.ReplyTo = msgID
	return msg
}

// WithReplyMarkup sets message reply markup.
func (msg *TextMessage) WithReplyMarkup(rm ReplyMarkup) *TextMessage {
	msg.ReplyMarkup = rm
	return msg
}

// BuildSendRequest returns Request for sending message.
func (msg *TextMessage) BuildSendRequest() (*Request, error) {
	r := NewRequest("sendMessage").
		AddChatID(msg.Peer).
		AddString("text", msg.Text).
		AddOptString("parse_mode", msg.ParseMode.String()).
		AddOptBool("disable_web_page_preview", msg.DisableWebPagePreview).
		AddOptBool("disable_notification", msg.DisableNotification)

	addOptMessageIdentityToRequest(r, "reply_to_message_id", msg.ReplyTo)

	return addOptReplyMarkupToRequest(r, "reply_markup", msg.ReplyMarkup)
}

// ForwardMessage represents send forward message.
//
// Example #1: forward recivied message:
//   client.Send(ctx, tg.NewForwardMessage(
//       userID,
//       update.Message,
//   ), nil)
//
// Example #2: forward message by known ids
//   client.Send(ctx, tg.NewForwardMessage(recipient, tg.MessageLocation{
//       Peer: ChatID(12345),
//       Message: MessageID(12345),
//   }, nil)
//
// Related API Method: https://core.telegram.org/bots/api#forwardmessage
type ForwardMessage struct {
	// Recipient of forwarded messages.
	Peer Peer

	// Message to forward.
	Message MessageIdentityFull

	// Set true, if message should be forwarded silent.
	DisableNotification bool
}

// NewForwardMessage creates a forward.
func NewForwardMessage(to Peer, msg MessageIdentityFull) *ForwardMessage {
	return &ForwardMessage{
		Peer:    to,
		Message: msg,
	}
}

// WithNotification on or off notification (default: on).
func (msg *ForwardMessage) WithNotification(yes bool) *ForwardMessage {
	msg.DisableNotification = !yes
	return msg
}

func (msg *ForwardMessage) BuildSendRequest() (*Request, error) {
	srcPeer, srcMessage := msg.Message.GetMessageLocation()

	r := NewRequest("forwardMessage").
		AddChatID(msg.Peer).
		AddPeer("from_chat_id", srcPeer).
		AddOptBool("disable_notification", msg.DisableNotification)

	addOptMessageIdentityToRequest(r, "message_id", srcMessage)

	return r, nil
}
