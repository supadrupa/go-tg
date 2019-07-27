package tg

import (
	"context"
	"encoding/json"
	"io"
	"regexp"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestParseMode_String(t *testing.T) {
	for k, v := range map[ParseMode]string{
		Plain:    "",
		Markdown: "markdown",
		HTML:     "HTML",
	} {
		assert.Equal(t,
			v,
			k.String(),
		)
	}
}

func TestParsePeer(t *testing.T) {
	for _, tt := range []struct {
		Input  string
		Result Peer
		Error  bool
	}{
		{"@channely_updates", Username("channely_updates"), false},
		{"-1001072262979", ChatID(-1001072262979), false},
		{"bad", nil, true},
	} {
		if tt.Error {
			_, err := ParsePeer(tt.Input)
			assert.Error(t, err)
		} else {
			v, err := ParsePeer(tt.Input)
			if assert.NoError(t, err) {
				assert.Equal(t, tt.Result, v)
			}
		}
	}
}

func TestFileID_AddFileToRequest(t *testing.T) {
	r := NewRequest("test")

	id := FileID("xxx")

	id.AddFileToRequest("test", r)

	assertRequestArgEqual(t,
		r,
		"test",
		"xxx",
	)
}

func TestPeers(t *testing.T) {
	for _, tt := range []struct {
		Peer  Peer
		Value string
	}{
		{
			Peer:  UserID(1),
			Value: "1",
		},
		{
			Peer:  Username("mr-linch"),
			Value: "@mr-linch",
		},
		{
			Peer: User{
				ID: UserID(1),
			},
			Value: "1",
		},
		{
			Peer:  ChatID(1234),
			Value: "1234",
		},
		{
			Peer: Chat{
				ID: ChatID(123),
			},
			Value: "123",
		},
	} {
		r := NewRequest("test")

		tt.Peer.AddPeerToRequest("chat_id", r)

		assert.Equal(t, tt.Value, r.args["chat_id"])
	}
}

func TestUserProfilesPhotos_First(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		photos := UserProfilePhotos{}

		assert.Nil(t, photos.First())
	})
	t.Run("NotEmpty", func(t *testing.T) {
		photos := UserProfilePhotos{
			Items: []PhotoSizeSlice{
				PhotoSizeSlice{
					PhotoSize{FileID: FileID("test1")},
				},
				PhotoSizeSlice{
					PhotoSize{FileID: FileID("test2")},
				},
			},
		}

		assert.Equal(t, FileID("test1"), photos.First()[0].FileID)
	})
}

func TestUserProfilesPhotos_Last(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		photos := UserProfilePhotos{}

		assert.Nil(t, photos.Last())
	})

	t.Run("NotEmpty", func(t *testing.T) {
		photos := UserProfilePhotos{
			Items: []PhotoSizeSlice{
				PhotoSizeSlice{
					PhotoSize{FileID: FileID("test1")},
				},
				PhotoSizeSlice{
					PhotoSize{FileID: FileID("test2")},
				},
			},
		}

		assert.Equal(t, FileID("test2"), photos.Last()[0].FileID)
	})
}

func TestFile_NewReader(t *testing.T) {
	var (
		path string
	)

	transport := &TransportMock{
		DownloadFunc: func(
			ctx context.Context,
			token string,
			p string,
		) (io.ReadCloser, error) {
			path = p
			return nil, errors.New("test error")
		},
	}

	client := NewClient("test",
		WithTransport(transport),
	)

	file := File{
		ID:   FileID("random"),
		Path: "images/photo_01.png",

		client: client,
	}

	body, err := file.NewReader(context.Background())

	assert.Nil(t, body)
	assert.EqualError(t, err, "test error")
	assert.Equal(t, file.Path, path)
}

func TestWebhookError_Error(t *testing.T) {

	err := WebhookError{"invalid response", time.Now()}

	assert.Regexp(t,
		regexp.MustCompile(`invalid response at (.+) \((.+)\)`),
		err.Error(),
	)
}

func TestWebhookInfo_IsSet(t *testing.T) {
	wh := WebhookInfo{}
	assert.False(t, wh.IsSet())

	wh = WebhookInfo{URL: "https://google.com"}
	assert.True(t, wh.IsSet())
}

func TestWebhookInfo_HasError(t *testing.T) {
	wh := WebhookInfo{}
	assert.False(t, wh.HasError())

	wh = WebhookInfo{Error: &WebhookError{Message: "Test"}}
	assert.True(t, wh.HasError())
}

func TestWebhookInfo_UnmarshalJSON(t *testing.T) {
	t.Run("WithoutError", func(t *testing.T) {
		webhookInfo := WebhookInfo{}

		err := json.Unmarshal([]byte(`{
		   "url":"http://test.com",
		   "has_custom_certificate":true,
		   "pending_update_count":42,
		   "max_connections":40
		}`), &webhookInfo)

		assert.Equal(t, WebhookInfo{
			URL:                  "http://test.com",
			HasCustomCertificate: true,
			PendingUpdateCount:   42,
			MaxConnections:       40,
		}, webhookInfo)

		assert.NoError(t, err)
	})

	t.Run("WithError", func(t *testing.T) {
		webhookInfo := WebhookInfo{}

		err := json.Unmarshal([]byte(`{
		   "url":"http://test.com",
		   "has_custom_certificate":true,
		   "pending_update_count":42,
           "last_error_message": "trouble",
           "last_error_date": 1563980388,
		   "max_connections":40
		}`), &webhookInfo)

		assert.Equal(t, WebhookInfo{
			URL:                  "http://test.com",
			HasCustomCertificate: true,
			PendingUpdateCount:   42,
			MaxConnections:       40,
			Error: &WebhookError{
				Message: "trouble",
				Date:    time.Unix(1563980388, 0),
			},
		}, webhookInfo)

		assert.NoError(t, err)
	})

	t.Run("WithInvalidTypeError", func(t *testing.T) {
		webhookInfo := WebhookInfo{}

		err := json.Unmarshal([]byte(`{
		   "url":"http://test.com",
		   "has_custom_certificate":true,
		   "pending_update_count": "x",
           "last_error_message": "trouble",
           "last_error_date": 1563980388,
		   "max_connections":40
		}`), &webhookInfo)

		assert.Error(t, err)
	})
}
