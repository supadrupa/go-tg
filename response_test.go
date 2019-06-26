package tg

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestResponseParameters_UnmarshalJSON(t *testing.T) {
	for _, tt := range []struct {
		in  string
		out ResponseParameters
	}{
		{`{}`, ResponseParameters{}},
		{`{"migrate_to_chat_id": 12345}`, ResponseParameters{MigrateToChatID: 12345}},
		{`{"retry_after": 60}`, ResponseParameters{RetryAfter: time.Minute}},
		{`{"migrate_to_chat_id": 12345, "retry_after": 60}`, ResponseParameters{MigrateToChatID: 12345, RetryAfter: time.Minute}},
	} {
		tmp := ResponseParameters{}

		err := json.Unmarshal([]byte(tt.in), &tmp)

		if assert.NoError(t, err) {
			assert.Equal(t, tt.out, tmp)
		}
	}
}
