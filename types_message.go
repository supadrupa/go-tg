package tg

// MessageIdentity is a common interface for everything that identifies the message.
//
// Types implementing this interface:
//  - MessageID
type MessageIdentity interface {
	GetMessageID() MessageID
}

// MessageIdentityFull is a common interface for everything that identifies the message in a particular chat.
//
// Types implementing this interface:
//  - MessageLocation
//  - Message
type MessageIdentityFull interface {
	GetMessageLocation() (Peer, MessageIdentity)
}

// MessageLocation it's implementation of MessageIdentityFull.
type MessageLocation struct {
	Chat    Peer
	Message MessageIdentity
}

func (ml MessageLocation) GetMessageLocation() (Peer, MessageIdentity) {
	return ml.Chat, ml.Message
}

// MessageID represents unique message identifier in chat.
type MessageID int

// GetMessageID it's MessageIdentity implementation.
func (msgID MessageID) GetMessageID() MessageID {
	return msgID
}

// Message represents incoming message.
type Message struct {
	// Unique message identifier inside this chat
	ID MessageID `json:"message_id"`

	// Optional. Sender, empty for messages sent to channels
	From *User `json:"from"`

	// Date the message was sent in Unix time
	Date int64 `json:"date"`

	// Conversation the message belongs to
	Chat Chat `json:"chat"`

	// Optional. For forwarded messages, sender of the original message.
	ForwardFrom *User `json:"forward_from,omitempty"`

	// Optional. For messages forwarded from channels,
	// information about the original channel.
	ForwardFromChat *Chat `json:"forward_from_chat,omitempty"`

	// Optional. For messages forwarded from channels,
	// identifier of the original message in the channel.
	ForwardFromMessageID MessageID `json:"forward_from_message_id,omitempty"`

	// Optional. For messages forwarded from channels,
	// signature of the post author if present.
	ForwardSignature string `json:"forward_signature,omitempty"`

	// Optional. Sender's name for messages forwarded
	// from users who disallow adding a link to their account
	// in forwarded messages.
	ForwardSenderName string `json:"forward_sender_name,omitempty"`

	// Optional. For forwarded messages, date the original message was sent in Unix time.
	ForwardDate int64 `json:"forward_date,omitempty"`

	// Optional. For replies, the original message.
	// Note that the Message object in this field will not contain
	// further reply_to_message fields even if it itself is a reply.
	ReplyToMessage *Message `json:"reply_to_message,omitempty"`

	// Optional. Date the message was last edited in Unix time
	EditDate int64 `json:"edit_date,omitempty"`

	// Optional. The unique identifier of a media message group this message belongs to
	MediaGroupID string `json:"media_group_id,omitempty"`

	// Optional. Signature of the post author for messages in channels
	AuthorSignature string `json:"author_signature,omitempty"`

	// Optional. For text messages, the actual UTF-8 text of the message, 0-4096 characters.
	Text string `json:"text,omitempty"`

	// Optional. For text messages, special entities like usernames, URLs, bot commands, etc. that appear in the text.
	Entities MessageEntitySlice `json:"entities,omitempty"`

	// Optional. For messages with a caption, special entities like usernames, URLs, bot commands, etc. that appear in the caption.
	CaptionEntities MessageEntitySlice `json:"caption_entities,omitempty"`

	// Optional. Message is an audio file, information about the file
	Audio *Audio `json:"audio,omitempty"`

	// Optional. Message is a general file, information about the file
	Document *Document `json:"document,omitempty"`

	// Optional. Message is an animation, information about the animation. For backward compatibility, when this field is set, the document field will also be set
	Animation *Animation `json:"animation,omitempty"`

	// Optional. Message is a game, information about the game.
	Game *Game `json:"game,omitempty"`

	// Optional. Message is a photo, available sizes of the photo
	Photo PhotoSizeSlice `json:"photo,omitempty"`

	// Optional. Message is a sticker, information about the sticker
	Sticker *Sticker `json:"sticker,omitempty"`

	// Optional. Message is a video, information about the video
	Video *Video `json:"video,omitempty"`

	// Optional. Message is a voice message, information about the file
	Voice *Voice `json:"voice,omitempty"`

	// Optional. Message is a video note, information about the video message
	VideoNote *VideoNote `json:"video_note,omitempty"`

	// Optional. Caption for the animation, audio, document, photo, video or voice, 0-1024 characters
	Caption string `json:"caption,omitempty"`

	// Optional. Message is a shared contact, information about the contact
	Contact *Contact `json:"contact,omitempty"`

	// Optional. Message is a shared location, information about the location
	Location *Location `json:"location,omitempty"`

	// Optional. Message is a venue, information about the venue
	Venue *Venue `json:"venue,omitempty"`

	// Optional. Message is a native poll, information about the poll
	Poll *Poll `json:"poll,omitempty"`

	// Optional. New members that were added to the group or supergroup
	// and information about them (the bot itself may be one of these members).
	NewChatMembers UserSlice `json:"new_chat_members,omitempty"`

	// Optional. A member was removed from the group, information about them (this member may be the bot itself)
	LeftChatMember *User `json:"left_chat_member,omitempty"`

	// Optional. A chat title was changed to this value.
	NewChatTitle string `json:"new_chat_title,omitempty"`

	// Optional. A chat photo was change to this value.
	NewChatPhoto string `json:"new_chat_photo,omitempty"`

	// Optional. Service message: the chat photo was deleted
	DeleteChatPhoto bool `json:"delete_chat_photo,omitempty"`

	// Optional. Service message: the group has been created
	GroupChatCreated bool `json:"group_chat_created,omitempty"`

	// Optional. Service message: the supergroup has been created. This field can‘t be received in a message coming through updates, because bot can’t be a member of a supergroup when it is created. It can only be found in reply_to_message if someone replies to a very first message in a directly created supergroup.
	SupergroupChatCreated bool `json:"supergroup_chat_created,omitempty"`

	// Optional. Service message: the channel has been created.
	// This field can‘t be received in a message coming through updates,
	// because bot can’t be a member of a channel when it is created.
	// It can only be found in reply_to_message if someone
	// replies to a very first message in a channel.
	ChannelChatCreated bool `json:"channel_chat_created,omitempty"`

	// Optional. The group has been migrated to a supergroup with the specified identifier.
	MigrateToChatID ChatID `json:"migrate_to_chat_id,omitempty"`

	// Optional. The supergroup has been migrated from a group with the specified identifier.
	MigrateFromChatID ChatID `json:"migrate_from_chat_id,omitempty"`

	// Optional. Specified message was pinned.
	// Note that the Message object in this field will not contain
	// further reply_to_message fields even if it is itself a reply.
	PinnedMessage *Message `json:"pinned_message,omitempty"`

	// Optional. Message is an invoice for a payment, information about the invoice.
	Invoice *Invoice `json:"invoice,omitempty"`

	// Optional. Message is a service message about a successful payment,
	// information about the payment.
	SuccessfulPayment *SuccessfulPayment `json:"successful_payment,omitempty"`

	// Optional. The domain name of the website on which the user has logged in.
	ConnectedWebsite string `json:"connected_website,omitempty"`

	// Optional. Telegram Passport data.
	PassportData string `json:"passport_data,omitempty"`

	// Optional. Inline keyboard attached to the message.
	// LoginURL buttons are represented as ordinary url buttons.
	ReplyMarkup InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

// GetMessageID for compatibility with the MessageIdentity interface.
func (msg Message) GetMessageID() MessageID {
	return msg.ID
}

// GetMessageID for compatibility with the MessageIdentityFull interface.
func (msg Message) GetMessageLocation() (Peer, MessageIdentity) {
	return msg.Chat.ID, msg.ID
}

// MessageEntity represents one special entity in a text message.
// For example, hashtags, usernames, URLs, etc.
type MessageEntity struct {
	// Type of the entity. Can be mention (@username), hashtag, cashtag, bot_command,
	// url, email, phone_number, bold (bold text), italic (italic text), code (monowidth string),
	// pre (monowidth block), text_link (for clickable text URLs),
	// text_mention (for users without usernames)
	Type string `json:"type"`

	// Offset in UTF-16 code units to the start of the entity
	Offset int `json:"offset"`

	// Length of the entity in UTF-16 code units
	Length int `json:"length"`

	// Optional. For “text_link” only, url that will be opened after user taps on the text
	URL string `json:"url,omitempty"`

	// Optional. For “text_mention” only, the mentioned user
	User *User `json:"user,omitempty"`
}

// MessageEntitySlice it's slice of message entity alias.
type MessageEntitySlice []MessageEntity

// Audio represents an audio file to be treated as music by the Telegram clients.
type Audio struct {
	// Unique identifier for this file
	FileID FileID `json:"file_id"`

	// Duration of the audio in seconds as defined by sender.
	DurationSeconds int `json:"duration"`

	// Optional. Performer of the audio as defined by sender or by audio tags
	Performer string `json:"performer,omitempty"`

	// Optional. Title of the audio as defined by sender or by audio tags
	Title string `json:"title,omitempty"`

	// Optional. MIME type of the file as defined by sender
	MIMEType string `json:"mime_type,omitempty"`

	// Optional. File size
	FileSize int `json:"file_size,omitempty"`

	// Optional. Thumbnail of the album cover to which the music file belongs
	Thumb *PhotoSize `json:"thumb,omitempty"`
}

// Document object represents a general file (as opposed to photos, voice messages and audio files).
type Document struct {
	// Unique file identifier.
	FileID string `json:"file_id"`

	// Optional. Document thumbnail as defined by sender
	Thumb *PhotoSize `json:"thumb,omitempty"`

	// Optional. Original filename as defined by sender
	FileName string `json:"file_name,omitempty"`

	// Optional. MIME type of the file as defined by sender
	MIMEType string `json:"mime_type,omitempty"`

	// Optional. File size
	FileSize int `json:"file_size,omitempty"`
}

// Video object represents a video file
type Video struct {
	// Unique file identifier.
	FileID string `json:"file_id"`

	// Video width as defined by sender
	Width int `json:"width"`

	// Video height as defined by sender
	Height int `json:"height"`

	// Duration of the video in seconds as defined by sender
	DurationSeconds int `json:"duration"`

	// Optional. Video thumbnail as defined by sender
	Thumb *PhotoSize `json:"thumb,omitempty"`

	// Optional. MIME type of the file as defined by sender
	MIMEType string `json:"mime_type,omitempty"`

	// Optional. File size
	FileSize int `json:"file_size,omitempty"`
}

// Animation object represents an animation file (GIF or H.264/MPEG-4 AVC video without sound).
type Animation struct {
	// Unique file identifier.
	FileID string `json:"file_id"`

	// Animation width as defined by sender
	Width int `json:"width"`

	// Animation height as defined by sender
	Height int `json:"height"`

	// Duration of the Animation in seconds as defined by sender
	DurationSeconds int `json:"duration"`

	// Optional. Animation thumbnail as defined by sender
	Thumb *PhotoSize `json:"thumb,omitempty"`

	// Optional. Original animation filename as defined by sender
	FileName string `json:"file_name,omitempty"`

	// Optional. MIME type of the file as defined by sender
	MIMEType string `json:"mime_type,omitempty"`

	// Optional. File size
	FileSize int `json:"file_size,omitempty"`
}

// Voice object represents a voice note.
type Voice struct {
	// Unique file identifier.
	FileID string `json:"file_id"`

	// Duration of the audio in seconds as defined by sender
	DurationSeconds int `json:"duration"`

	// Optional. MIME type of the file as defined by sender
	MIMEType string `json:"mime_type,omitempty"`

	// Optional. File size
	FileSize int `json:"file_size,omitempty"`
}

// VideoNote object represents a video message (available in Telegram apps as of v.4.0).
type VideoNote struct {
	// Unique file identifier.
	FileID string `json:"file_id"`

	// Video width and height (diameter of the video message) as defined by sender.
	Length int `json:"length"`

	// Duration of the video in seconds as defined by sender
	DurationSeconds int `json:"duration"`

	// Optional. Video thumbnail as defined by sender
	Thumb *PhotoSize `json:"thumb,omitempty"`

	// Optional. File size
	FileSize int `json:"file_size,omitempty"`
}

// Contact object represents a phone contact.
type Contact struct {
	// Contact's phone number
	PhoneNumber string `json:"phone_number"`

	// Contact's first name
	FirstName string `json:"first_name"`

	// Optional. Contact's last name
	LastName string `json:"last_name,omitempty"`

	// Optional. Contact's user identifier in Telegram
	VCard string `json:"v_card,omitempty"`
}

// Location object represents a point on the map.
type Location struct {
	// Longitude as defined by sender.
	Longitude float64 `json:"longitude"`

	// Latitude as defined by sender.
	Latitude float64 `json:"latitude"`
}

// Venue object represents a venue.
type Venue struct {
	// Venue location
	Location Location `json:"location"`

	// Name of the venue.
	Title string `json:"title"`

	// Address of the venue
	Address string `json:"address"`

	// Optional. Foursquare identifier of the venue
	FoursquareID string `json:"foursquare_id,omitempty"`

	// Optional. Foursquare type of the venue.
	// For example: "arts_entertainment/default", "arts_entertainment/aquarium" or "food/icecream"
	FoursquareType string `json:"foursquare_type,omitempty"`
}

// PollOption object contains information about one answer option in a poll.
type PollOption struct {
	// Option text, 1-100 characters.
	Text string `json:"text"`

	// Number of users that voted for this option.
	VoterCount int `json:"voter_count"`
}

// Poll object contains information about a poll.
type Poll struct {
	// Unique poll identifier.
	ID string `json:"id"`

	// Poll question, 1-255 characters.
	Question string `json:"question"`

	// List of poll options.
	Options []PollOption `json:"options"`

	// True, if the poll is closed.
	IsClosed bool `json:"is_closed"`
}
