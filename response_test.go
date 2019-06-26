package tg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponseParameters_UnmarshalResult(t *testing.T) {
	dst := struct {
		Test bool `json:"test"`
	}{}

	response := Response{
		Result: []byte(`{"test": true}`),
	}

	err := response.UnmarshalResult(&dst)

	if assert.NoError(t, err) {
		assert.True(t, dst.Test)
	}
}
