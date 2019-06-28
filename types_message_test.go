package tg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageIdentity(t *testing.T) {
	for _, tt := range []struct {
		Identity MessageIdentity
		Result   MessageID
	}{
		{MessageID(1), MessageID(1)},
	} {
		assert.Equal(t,
			tt.Result,
			tt.Identity.GetMessageID(),
		)
	}
}
