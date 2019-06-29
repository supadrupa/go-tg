package tg

import (
	"encoding/json"
)

// PassportData contains information about Telegram Passport data shared with the bot by the user.
type PassportData struct {
	// Array with information about documents and other Telegram Passport elements that was shared with the bot.
	Data []json.RawMessage `json:"data"`

	// Encrypted credentials required to decrypt the data.
	Credentials json.RawMessage `json:"credentials"`
}
