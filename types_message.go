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
	Sender *User `json:"sender"`

	// Date the message was sent in Unix time
	Date int64 `json:"date,omitempty"`

	// Conversation the message belongs to
	Chat Chat `json:"chat,omitempty"`
}

// GetMessageID for compatibility with the MessageIdentity interface.
func (msg Message) GetMessageID() MessageID {
	return msg.ID
}

// GetMessageID for compatibility with the MessageIdentityFull interface.
func (msg Message) GetMessageLocation() (Peer, MessageIdentity) {
	return msg.Chat.ID, msg.ID
}
