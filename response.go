package tg

import (
	"encoding/json"
	"time"
)

// ResponseParameters contains information about why a request was unsuccessful.
type ResponseParameters struct {
	// The group has been migrated to a supergroup with the specified identifier.
	MigrateToChatID int

	// Optional. In case of exceeding flood control,
	// the time left to wait before request can be repeated.
	RetryAfter time.Duration
}

func (params *ResponseParameters) UnmarshalJSON(data []byte) error {
	helper := struct {
		MigrateToChatID int `json:"migrate_to_chat_id"`
		RetryAfter      int `json:"retry_after"`
	}{}

	if err := json.Unmarshal(data, &helper); err != nil {
		return err
	}

	params.MigrateToChatID = helper.MigrateToChatID
	params.RetryAfter = time.Second * time.Duration(helper.RetryAfter)

	return nil
}

// Response represents Telegram Bot API response.
type Response struct {
	// If equals true, the request was successful
	// and the result of the query can be found in the Result field.
	OK bool `json:"ok"`

	// Telegram Bot API method.
	Method string `json:"-"`

	// Result of request in case of success.
	Result json.RawMessage `json:"result"`

	// Description of response/error
	Description string `json:"description"`

	// StatusCode contains HTTP status code of response.
	StatusCode int `json:"-"`

	// Error code from Telegram.
	ErrorCode int `json:"error_code"`

	// Contains information about why a request was unsuccessful.
	Parameters *ResponseParameters `json:"parameters"`
}

func (response Response) UnmarshalResult(dst interface{}) error {
	return json.Unmarshal(response.Result, dst)
}
