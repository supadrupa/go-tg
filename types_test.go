package tg

import (
	"context"
	"io"
	"testing"

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
