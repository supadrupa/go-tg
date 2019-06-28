package tg

import (
	"bytes"
	"strconv"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type mapEncoder struct {
	Args  map[string]string
	Files map[string]RequestFile
}

var _ Encoder = &mapEncoder{}

func newMapEncoder() *mapEncoder {
	return &mapEncoder{
		Args:  make(map[string]string),
		Files: make(map[string]RequestFile),
	}
}

func (enc *mapEncoder) AddString(k, v string) error {
	enc.Args[k] = v
	return nil
}

func (enc *mapEncoder) AddFile(k string, v RequestFile) error {
	enc.Files[k] = v
	return nil
}

func extractArgs(r *Request) map[string]string {
	enc := newMapEncoder()
	if err := r.Encode(enc); err != nil {
		panic(err)
	}
	return enc.Args
}

func extractFiles(r *Request) map[string]RequestFile {
	enc := newMapEncoder()
	if err := r.Encode(enc); err != nil {
		panic(err)
	}
	return enc.Files
}

func TestRequest_New(t *testing.T) {
	r := NewRequest("getMe").WithToken("12345:test")

	assert.Equal(t, "getMe", r.Method())
	assert.Equal(t, "12345:test", r.Token())
}

func assertRequestArgEqual(
	t *testing.T,
	r *Request,
	key string,
	excepted string,
) bool {
	args := extractArgs(r)
	v := args[key]

	return assert.Equalf(t,
		excepted,
		v,
		"request argument '%s' equals '%s', excepted '%s'",
		key,
		v,
		excepted,
	)
}

func assertRequestFileEqual(
	t *testing.T,
	r *Request,
	key string,
	excepted RequestFile,
) bool {
	args := extractFiles(r)
	v := args[key]

	return assert.Equalf(t,
		excepted,
		v,
		"request file '%s' equals '%+v', excepted '%+v'",
		key,
		v,
		excepted,
	)
}

type peerMock string

func (m peerMock) AddPeerToRequest(k string, r *Request) {
	r.AddString(k, string(m))
}

type simplePart struct {
	K string
	V string
}

func (part simplePart) AddToRequest(r *Request) {
	r.AddString(part.K, part.V)
}

func TestRequest_Add(t *testing.T) {

	t.Run("String", func(t *testing.T) {
		assertRequestArgEqual(t,
			NewRequest("test").AddString("v", "123"),
			"v",
			"123",
		)
	})

	t.Run("OptString", func(t *testing.T) {
		assertRequestArgEqual(t,
			NewRequest("test").AddOptString("v", "123"),
			"v",
			"123",
		)

		assertRequestArgEqual(t,
			NewRequest("test").AddOptString("v", ""),
			"v",
			"",
		)
	})

	t.Run("Int", func(t *testing.T) {
		assertRequestArgEqual(t,
			NewRequest("test").AddInt("v", -1003413),
			"v",
			"-1003413",
		)
	})

	t.Run("OptInt", func(t *testing.T) {
		assertRequestArgEqual(t,
			NewRequest("test").AddOptInt("v", -1003413),
			"v",
			"-1003413",
		)

		assertRequestArgEqual(t,
			NewRequest("test").AddOptInt("v", 0),
			"v",
			"",
		)
	})

	t.Run("Int64", func(t *testing.T) {
		assertRequestArgEqual(t,
			NewRequest("test").AddInt64("v", -1003411203043),
			"v",
			"-1003411203043",
		)
	})

	t.Run("Bool", func(t *testing.T) {
		assertRequestArgEqual(t,
			NewRequest("test").AddBool("v", true),
			"v",
			"true",
		)

		assertRequestArgEqual(t,
			NewRequest("test").AddBool("v", false),
			"v",
			"false",
		)
	})

	t.Run("Bool", func(t *testing.T) {
		assertRequestArgEqual(t,
			NewRequest("test").AddBool("v", true),
			"v",
			"true",
		)

		assertRequestArgEqual(t,
			NewRequest("test").AddBool("v", false),
			"v",
			"false",
		)
	})

	t.Run("OptBool", func(t *testing.T) {
		assertRequestArgEqual(t,
			NewRequest("test").AddOptBool("v", true),
			"v",
			"true",
		)

		assertRequestArgEqual(t,
			NewRequest("test").AddOptBool("v", false),
			"v",
			"",
		)
	})

	t.Run("Float64", func(t *testing.T) {
		assertRequestArgEqual(t,
			NewRequest("test").AddFloat64("v", 25.941481),
			"v",
			"25.941481",
		)
	})

	t.Run("File", func(t *testing.T) {
		file := RequestFile{
			Body: &bytes.Buffer{},
			Name: "test.png",
		}

		assertRequestFileEqual(t,
			NewRequest("test").AddFile("v", file),
			"v",
			file,
		)
	})

	t.Run("Peer", func(t *testing.T) {
		peer := peerMock("test")

		assertRequestArgEqual(t,
			NewRequest("test").AddPeer("v", peer),
			"v",
			string(peer),
		)
	})

	t.Run("ChatID", func(t *testing.T) {
		peer := peerMock("chat_id_value")

		assertRequestArgEqual(t,
			NewRequest("test").AddChatID(peer),
			"chat_id",
			string(peer),
		)
	})

	t.Run("Time", func(t *testing.T) {
		now := time.Now()

		assertRequestArgEqual(t,
			NewRequest("test").AddTime("t", now),
			"t",
			strconv.FormatInt(now.Unix(), 10),
		)
	})

	t.Run("OptTime", func(t *testing.T) {
		{
			var now time.Time

			assertRequestArgEqual(t,
				NewRequest("test").AddOptTime("t", now),
				"t",
				"",
			)
		}

		{
			now := time.Now()

			assertRequestArgEqual(t,
				NewRequest("test").AddOptTime("t", now),
				"t",
				strconv.FormatInt(now.Unix(), 10),
			)
		}
	})

	t.Run("Part", func(t *testing.T) {
		part := simplePart{"v", "test"}

		assertRequestArgEqual(t,
			NewRequest("test").AddPart(part),
			"v",
			"test",
		)
	})

}

func TestRequest_HasFiles(t *testing.T) {
	t.Run("WithoutFiles", func(t *testing.T) {
		r := NewRequest("test")

		assert.False(t, r.HasFiles())
	})
	t.Run("WithFiles", func(t *testing.T) {
		r := NewRequest("test").AddFile("mkk", RequestFile{
			Name: "Test.png",
			Body: &bytes.Buffer{},
		})

		assert.True(t, r.HasFiles())
	})
}

type encoderMock struct {
	args           map[string]string
	addStringError error

	files        map[string]RequestFile
	addFileError error
}

func (mock *encoderMock) AddFile(k string, v RequestFile) error {
	if mock.addFileError != nil {
		return mock.addFileError
	}

	if mock.files == nil {
		mock.files = map[string]RequestFile{
			k: v,
		}

		return nil
	}

	mock.files[k] = v

	return nil
}

func (mock *encoderMock) AddString(k string, v string) error {
	if mock.addStringError != nil {
		return mock.addStringError
	}

	if mock.args == nil {
		mock.args = map[string]string{
			k: v,
		}

		return nil
	}

	mock.args[k] = v

	return nil
}

func TestRequest_Encode(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		enc := &encoderMock{}

		req := NewRequest("test").
			AddString("test", "v").
			AddFile("f", RequestFile{Name: "test.png"})

		err := req.Encode(enc)

		assert.NoError(t, err)
	})

	t.Run("AddStringError", func(t *testing.T) {
		enc := &encoderMock{
			addStringError: errors.New("test"),
		}

		req := NewRequest("test").
			AddString("test", "v").
			AddFile("f", RequestFile{Name: "test.png"})

		err := req.Encode(enc)

		assert.Error(t, err)
	})

	t.Run("AddFileError", func(t *testing.T) {
		enc := &encoderMock{
			addFileError: errors.New("test"),
		}

		req := NewRequest("test").
			AddString("test", "v").
			AddFile("f", RequestFile{Name: "test.png"})

		err := req.Encode(enc)

		assert.Error(t, err)
	})
}
