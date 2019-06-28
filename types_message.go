package tg

// MessageIdentity is a common interface for everything that identifies a message.
//
// Types implementing this interface:
//  - MessageID
type MessageIdentity interface {
	GetMessageID() MessageID
}

// MessageID represents unique message identifier in chat.
type MessageID int

// GetMessageID it's MessageIdentity implementation.
func (msgID MessageID) GetMessageID() MessageID {
	return msgID
}
