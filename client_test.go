package tg

import (
	"context"
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var (
	ResponseResultTrue = &Response{
		OK:     true,
		Result: []byte("true"),
	}
)

func FakeExecuteRequest(
	do func(ctx context.Context, c *Client) error,
	res *Response,
	err error,
) (*Request, error) {
	var (
		request *Request
	)

	transport := &TransportMock{
		ExecuteFunc: func(ctx context.Context, r *Request) (*Response, error) {
			request = r
			return res, err
		},
	}

	client := NewClient("1234:secret", WithTransport(transport))

	doError := do(context.Background(), client)

	return request, doError
}

func TestClient_New(t *testing.T) {
	transport := &TransportMock{}

	client := NewClient("1234:secret",
		WithTransport(transport),
		WithParseMode(Markdown),
		WithWebPagePreview(false),
	)

	assert.Equal(t,
		transport,
		client.transport,
		"transport is not set",
	)

	assert.Equal(t,
		Markdown,
		client.defaultParseMode,
		"default parse mode is not set",
	)

	assert.False(t,
		client.defaultWebPagePreview,
	)
}

func TestClient_Invoke(t *testing.T) {
	ctx := context.Background()

	t.Run("ExecuteError", func(t *testing.T) {
		exceptedError := errors.New("execute error")

		transport := &TransportMock{
			ExecuteFunc: func(ctx context.Context, r *Request) (*Response, error) {
				assert.Equal(t, "getMe", r.Method())
				assert.Equal(t, "1234:secret", r.Token())

				return nil, exceptedError
			},
		}

		client := NewClient("1234:secret",
			WithTransport(transport),
		)

		req := NewRequest("getMe")

		err := client.Invoke(ctx, req, nil)

		assert.Equal(t, exceptedError, err)
	})

	t.Run("UnmarshalError", func(t *testing.T) {
		transport := &TransportMock{
			ExecuteFunc: func(ctx context.Context, r *Request) (*Response, error) {
				assert.Equal(t, "getMe", r.Method())
				assert.Equal(t, "1234:secret", r.Token())

				return &Response{
					OK: true,
					Result: []byte(`{
						"ok": true,
						"result": {
							"test": 1
						}
					}`),
				}, nil
			},
		}

		client := NewClient("1234:secret",
			WithTransport(transport),
		)

		req := NewRequest("getMe")

		var result int

		err := client.Invoke(ctx, req, &result)

		_, isUnmarshalTypeError := err.(*json.UnmarshalTypeError)

		assert.True(t, isUnmarshalTypeError)
	})

	t.Run("OK", func(t *testing.T) {
		transport := &TransportMock{
			ExecuteFunc: func(ctx context.Context, r *Request) (*Response, error) {
				assert.Equal(t, "getMe", r.Method())
				assert.Equal(t, "1234:secret", r.Token())

				return &Response{
					OK: true,
					Result: []byte(`{
						"ok": true,
						"result": {
							"test": 1
						}
					}`),
				}, nil
			},
		}

		client := NewClient("1234:secret",
			WithTransport(transport),
		)

		req := NewRequest("getMe")

		err := client.Invoke(ctx, req, nil)

		assert.NoError(t, err)
	})
}

func TestProfilePhotosOptions_AddToRequest(t *testing.T) {
	opts := &ProfilePhotosOptions{Offset: 10, Limit: 5}

	args := extractArgs(
		NewRequest("getUserProfilePhotos").AddPart(opts),
	)

	assert.Equal(t, map[string]string{
		"offset": "10",
		"limit":  "5",
	}, args)
}

func TestClient_KickChatMember(t *testing.T) {
	until := time.Now()

	request, _ := FakeExecuteRequest(func(ctx context.Context, client *Client) error {
		return client.KickChatMember(ctx,
			ChatID(-100),
			UserID(1),
			&KickOptions{
				Until: until,
			},
		)
	}, ResponseResultTrue, nil)

	args := extractArgs(request)

	assert.Equal(t, map[string]string{
		"chat_id":    "-100",
		"user_id":    "1",
		"until_date": strconv.FormatInt(until.Unix(), 10),
	}, args)
}

func TestClient_UnbanChatMember(t *testing.T) {
	request, _ := FakeExecuteRequest(func(ctx context.Context, client *Client) error {
		return client.UnbanChatMember(ctx,
			ChatID(-100),
			UserID(1),
		)
	}, ResponseResultTrue, nil)

	args := extractArgs(request)

	assert.Equal(t, map[string]string{
		"chat_id": "-100",
		"user_id": "1",
	}, args)
}

func TestClient_RestrictChatMember(t *testing.T) {
	until := time.Now()

	request, _ := FakeExecuteRequest(func(ctx context.Context, client *Client) error {
		return client.RestrictChatMember(ctx,
			ChatID(-100),
			UserID(1),
			&RestrictOptions{
				Until:                  until,
				CanSendMessages:        true,
				CanSendMediaMessages:   true,
				CanSendOtherMessages:   true,
				CanSendWebPagePreviews: true,
			},
		)
	}, ResponseResultTrue, nil)

	args := extractArgs(request)

	assert.Equal(t, map[string]string{
		"chat_id":                    "-100",
		"user_id":                    "1",
		"until_date":                 strconv.FormatInt(until.Unix(), 10),
		"can_send_messages":          "true",
		"can_send_media_messages":    "true",
		"can_send_other_messages":    "true",
		"can_send_web_page_previews": "true",
	}, args)
}
