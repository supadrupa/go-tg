package tg

import "time"

// Media define interface files in outgoing message.
//
// Types implementing this interface:
//  - InputFile
//  - FileID
//  - RemoteFile
type Media interface {
	AddFileToRequest(k string, r *Request)
}

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
// Example #1: forward received message:
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

type PhotoMessage struct {
	// Recipient of photo message
	Peer Peer

	// Photo (InputFile, FileID, RemoteFile)
	Photo Media

	// Caption of photo (0-1024)
	Caption string

	// Parse mode of caption
	ParseMode ParseMode

	// Pass true for send message silent.
	DisableNotification bool

	// Reply to message identity.
	ReplyTo MessageIdentity

	// Reply markup of the message.
	ReplyMarkup ReplyMarkup
}

// NewPhotoMessage creates a photo message.
func NewPhotoMessage(to Peer, media Media) *PhotoMessage {
	return &PhotoMessage{
		Peer:  to,
		Photo: media,
	}
}

// WithCaption sets message caption.
func (msg *PhotoMessage) WithCaption(text string) *PhotoMessage {
	msg.Caption = text
	return msg
}

// WithParseMode sets caption parse mode.
func (msg *PhotoMessage) WithParseMode(pm ParseMode) *PhotoMessage {
	msg.ParseMode = pm
	return msg
}

// WithNotification enable or disable notification (default: enabled).
func (msg *PhotoMessage) WithNotification(yes bool) *PhotoMessage {
	msg.DisableNotification = !yes
	return msg
}

// WithReplyTo sets ids of original message, if message is reply.
func (msg *PhotoMessage) WithReplyTo(msgID MessageIdentity) *PhotoMessage {
	msg.ReplyTo = msgID
	return msg
}

// WithReplyMarkup sets message reply markup.
func (msg *PhotoMessage) WithReplyMarkup(rm ReplyMarkup) *PhotoMessage {
	msg.ReplyMarkup = rm
	return msg
}

func (msg *PhotoMessage) BuildSendRequest() (*Request, error) {
	r := NewRequest("sendPhoto").
		AddChatID(msg.Peer).
		AddOptString("caption", msg.Caption).
		AddOptString("parse_mode", msg.ParseMode.String()).
		AddOptBool("disable_notification", msg.DisableNotification)

	addMediaToRequest(r, "photo", msg.Photo)
	addOptMessageIdentityToRequest(r, "reply_to_message_id", msg.ReplyTo)

	return addOptReplyMarkupToRequest(r, "reply_markup", msg.ReplyMarkup)
}

// AudioMessage represents outgoing audio message.
// Audio must be in the .mp3 format.
// Bots can currently send audio files of up to 50 MB in size, this limit may be changed in the future.
type AudioMessage struct {
	// Recipient of photo message.
	Peer Peer

	// Audio media to send (InputFile, FileID, RemoteFile).
	Audio Media

	// Caption of audio (0-1024).
	Caption string

	// Parse mode of caption.
	ParseMode ParseMode

	// Duration of the audio (will be sent in seconds).
	Duration time.Duration

	// Performer of the track.
	Performer string

	// Track name
	Title string

	// Thumbnail of the file sent.
	// Can be ignored if thumbnail generation for the file is supported server-side.
	// The thumbnail should be in JPEG format and less than 200 kB in size.
	// A thumbnail‘s width and height should not exceed 320.
	// Thumbnails can’t be reused and can be only uploaded as a new file.
	Thumb *InputFile

	// Pass true for send message silent.
	DisableNotification bool

	// Reply to message identity.
	ReplyTo MessageIdentity

	// Reply markup of the message.
	ReplyMarkup ReplyMarkup
}

// NewAudioMessage creates outgoing audio message.
func NewAudioMessage(to Peer, audio Media) *AudioMessage {
	return &AudioMessage{
		Peer:  to,
		Audio: audio,
	}
}

// WithCaption sets message caption.
func (msg *AudioMessage) WithCaption(text string) *AudioMessage {
	msg.Caption = text
	return msg
}

// WithTitle sets audio title.
func (msg *AudioMessage) WithTitle(title string) *AudioMessage {
	msg.Title = title
	return msg
}

// WithPerformer sets audio performer.
func (msg *AudioMessage) WithPerformer(performer string) *AudioMessage {
	msg.Performer = performer
	return msg
}

// WithDuration sets audio duration.
func (msg *AudioMessage) WithDuration(d time.Duration) *AudioMessage {
	msg.Duration = d
	return msg
}

// WithThumb sets audio thumb.
func (msg *AudioMessage) WithThumb(thumb InputFile) *AudioMessage {
	msg.Thumb = &thumb
	return msg
}

// WithParseMode sets caption parse mode.
func (msg *AudioMessage) WithParseMode(pm ParseMode) *AudioMessage {
	msg.ParseMode = pm
	return msg
}

// WithNotification enable or disable notification (default: enabled).
func (msg *AudioMessage) WithNotification(yes bool) *AudioMessage {
	msg.DisableNotification = !yes
	return msg
}

// WithReplyTo sets ids of original message, if message is reply.
func (msg *AudioMessage) WithReplyTo(msgID MessageIdentity) *AudioMessage {
	msg.ReplyTo = msgID
	return msg
}

// WithReplyMarkup sets message reply markup.
func (msg *AudioMessage) WithReplyMarkup(rm ReplyMarkup) *AudioMessage {
	msg.ReplyMarkup = rm
	return msg
}

func (msg *AudioMessage) BuildSendRequest() (*Request, error) {
	r := NewRequest("sendAudio").
		AddChatID(msg.Peer).
		AddOptString("caption", msg.Caption).
		AddOptString("parse_mode", msg.ParseMode.String()).
		AddOptString("performer", msg.Performer).
		AddOptString("title", msg.Title).
		AddOptBool("disable_notification", msg.DisableNotification).
		AddOptInt("duration", int(msg.Duration.Seconds())).
		AddOptAttachment("thumb", msg.Thumb)

	addMediaToRequest(r, "audio", msg.Audio)
	addOptMessageIdentityToRequest(r, "reply_to_message_id", msg.ReplyTo)

	return addOptReplyMarkupToRequest(r, "reply_markup", msg.ReplyMarkup)
}
