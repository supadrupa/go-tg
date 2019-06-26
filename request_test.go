package tg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type MapEncoder struct {
	Args  map[string]string
	Files map[string]InputFile
}

var _ Encoder = &MapEncoder{}

func NewMapEncoder() *MapEncoder {
	return &MapEncoder{
		Args:  make(map[string]string),
		Files: make(map[string]InputFile),
	}
}

func (enc *MapEncoder) AddString(k, v string) error {
	enc.Args[k] = v
	return nil
}

func (enc *MapEncoder) AddInputFile(k string, v InputFile) error {
	enc.Files[k] = v
	return nil
}

type MockInputFile struct {
	name string
	n    int
	err  error
}

var _ InputFile = &MockInputFile{}

func (m *MockInputFile) Name() string {
	return m.name
}

func (m *MockInputFile) Read(p []byte) (n int, err error) {
	return m.n, m.err
}

func TestRequest_New(t *testing.T) {
	r := NewRequest("getMe").WithToken("12345:test")

	assert.Equal(t, "getMe", r.Method())
	assert.Equal(t, "12345:test", r.Token())
}

func TestRequest_Add(t *testing.T) {
	r := NewRequest("sendMessage")

	inputFile := &MockInputFile{}

	r.AddString("string", "value").
		AddInt("int", 1).
		AddInt64("int64", 64).
		AddBool("bool", true).
		AddFloat64("float64", 1.23).
		AddInputFile("input-file", inputFile)

	encoder := NewMapEncoder()

	r.Encode(encoder)

	assert.Equal(t, map[string]string{
		"string":  "value",
		"int":     "1",
		"int64":   "64",
		"float64": "1.23",
		"bool":    "true",
	}, encoder.Args)

	assert.Equal(t, map[string]InputFile{
		"input-file": inputFile,
	}, encoder.Files)
}

func TestRequest_HasInputFile(t *testing.T) {
	t.Run("False", func(t *testing.T) {
		r := NewRequest("sendMessage")
		assert.False(t, r.HasInputFile())
	})
	t.Run("True", func(t *testing.T) {
		r := NewRequest("sendMessage")
		r.AddInputFile("test", &MockInputFile{})
		assert.True(t, r.HasInputFile())
	})
}
