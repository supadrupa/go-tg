package tg

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddOptMessageIdentityToRequest(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		var identity MessageIdentity

		r := NewRequest("test")

		addOptMessageIdentityToRequest(r, "test", identity)

		args := extractArgs(r)

		assert.Empty(t, args)
	})

	t.Run("NotNil", func(t *testing.T) {
		identity := MessageID(1)

		r := NewRequest("test")

		addOptMessageIdentityToRequest(r, "test", identity)

		args := extractArgs(r)

		assert.Equal(t, map[string]string{
			"test": "1",
		}, args)
	})
}

func TestAddOptReplyMarkup(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		var rm ReplyMarkup

		r := NewRequest("test")

		_, err := addOptReplyMarkupToRequest(r, "test", rm)
		require.NoError(t, err)

		args := extractArgs(r)

		assert.Empty(t, args)
	})

	t.Run("NotNil", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			rm := &ReplyMarkupMock{
				EncodeReplyMarkupFunc: func() (string, error) {
					return "{}", nil
				},
			}

			r := NewRequest("test")
			_, err := addOptReplyMarkupToRequest(r, "test", rm)
			require.NoError(t, err)

			args := extractArgs(r)

			assert.Equal(t, map[string]string{
				"test": "{}",
			}, args)
		})
		t.Run("Failed", func(t *testing.T) {
			testErr := errors.New("test")

			rm := &ReplyMarkupMock{
				EncodeReplyMarkupFunc: func() (string, error) {
					return "", testErr
				},
			}

			r := NewRequest("test")
			_, err := addOptReplyMarkupToRequest(r, "test", rm)
			require.Equal(t, testErr, err)
		})
	})
}
