package tg

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewInputFile(t *testing.T) {
	body := strings.NewReader("test")

	file := NewInputFile(
		"test.txt",
		body,
	)

	assert.Equal(t, "test.txt", file.Name)
	assert.Equal(t, body, file.Body)
}

func TestNewInputFileBytes(t *testing.T) {
	file := NewInputFileBytes(
		"test.txt",
		[]byte("test"),
	)

	body, err := ioutil.ReadAll(file.Body)
	require.NoError(t, err)

	assert.Equal(t, "test.txt", file.Name)
	assert.Equal(t, []byte("test"), body)
}

func TestNewInputFileLocal(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		file, err := NewInputFileLocal("./README.md")
		require.NoError(t, err)
		defer file.Close()

		assert.Equal(t, "README.md", file.Name)
	})

	t.Run("Failed", func(t *testing.T) {
		_, err := NewInputFileLocal("./SHOULD-NOT-EXIST.txt")
		require.Error(t, err)
	})
}

func TestNewInputFileLocalBuffer(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		file, err := NewInputFileLocalBuffer("./README.md")
		require.NoError(t, err)
		defer file.Close()

		assert.Equal(t, "README.md", file.Name)
	})

	t.Run("Failed", func(t *testing.T) {
		_, err := NewInputFileLocalBuffer("./SHOULD-NOT-EXIST.txt")
		require.Error(t, err)
	})
}

func TestInputFile_AddFileToRequest(t *testing.T) {
	r := NewRequest("test")

	file := NewInputFileBytes(
		"test.txt",
		[]byte("test"),
	)

	file.AddFileToRequest("document", r)

	assertInputFileEqual(t,
		r,
		"document",
		file,
	)
}
