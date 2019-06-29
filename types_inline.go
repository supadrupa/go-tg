package tg

// InlineQueryID it's unique identifier of InlineQuery
type InlineQueryID string

// InlineMessageID unique inline message ID.
type InlineMessageID string

// InlineQuery object represents an incoming inline query.
// When the user sends an empty query, your bot could return some default or trending results.
type InlineQuery struct {
	// Unique identifier for this query
	ID InlineQueryID `json:"id"`

	// Sender
	From User `json:"from"`

	// Optional. Sender location, only for bots that request user location
	Location *Location `json:"location,omitempty"`

	// Text of the query (up to 512 characters)
	Query string `json:"query"`

	// Offset of the results to be returned, can be controlled by the bot.
	Offset string `json:"offset"`
}

// ChosenInlineResult represents a result of an inline query that was chosen by the user and sent to their chat partner.
type ChosenInlineResult struct {
	// The unique identifier for the result that was chosen
	ResultID string `json:"result_id"`

	// The user that chose the result.
	From User `json:"from"`

	// Optional. Sender location, only for bots that require user location
	Location *Location `json:"location,omitempty"`

	// Optional. Identifier of the sent inline message.
	// Available only if there is an inline keyboard attached to the message.
	// Will be also received in callback queries and can be used to edit the message.
	InlineMessageID InlineMessageID `json:"inline_message_id,omitempty"`

	// The query that was used to obtain the result
	Query string `json:"query"`
}
