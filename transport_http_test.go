package tg

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMultipartEncoder(t *testing.T) {
	var (
		body        bytes.Buffer
		contentType string
	)

	write := func() {
		encoder := NewMultipartEncoder(&body)

		file := RequestFile{
			Name: "test.txt",
			Body: strings.NewReader("test, test, test"),
		}

		err := encoder.AddString("chat_id", "@channely_updates")
		require.NoError(t, err, "when write string")

		err = encoder.AddFile("document", file)
		require.NoError(t, err, "when write file")

		err = encoder.Close()
		require.NoError(t, err, "close encoder")

		contentType = encoder.ContentType()

		assert.Contains(t, contentType, "multipart/form-data; boundary=")
	}

	read := func() {
		_, params, err := mime.ParseMediaType(contentType)
		require.NoError(t, err, "parse media type: '%s'", contentType)

		reader := multipart.NewReader(&body, params["boundary"])

		form, err := reader.ReadForm(1024)
		require.NoError(t, err, "parse form")

		assert.Equal(t, "@channely_updates", form.Value["chat_id"][0])
		assert.Equal(t, "test.txt", form.File["document"][0].Filename)
	}

	write()
	read()
}

func TestURLEncodedEncoder(t *testing.T) {
	var (
		body        bytes.Buffer
		contentType string
	)

	write := func() {
		encoder := NewURLEncodedEncoder(&body)

		file := RequestFile{
			Name: "test.txt",
			Body: strings.NewReader("test, test, test"),
		}

		err := encoder.AddString("chat_id", "@channely_updates")
		require.NoError(t, err, "when write string")

		err = encoder.AddString("text", "1+1=2")
		require.NoError(t, err, "when write string")

		err = encoder.AddFile("document", file)
		require.EqualError(t, err, "URLEncodedEncoder does not support file uploading")

		err = encoder.Close()
		require.NoError(t, err, "close encoder")

		contentType = encoder.ContentType()

		assert.Equal(t, contentType, "application/x-www-form-urlencoded")
	}

	read := func() {
		vs, err := url.ParseQuery(body.String())
		require.NoError(t, err, "parse query string")

		assert.Equal(t, "@channely_updates", vs.Get("chat_id"))
		assert.Equal(t, "1+1=2", vs.Get("text"))
	}

	write()
	read()
}

func TestDefaultURLBuilders(t *testing.T) {
	assert.Equal(t,
		"https://api.telegram.org/bot1234:secret/getMe",
		defaultBuildCallURL("1234:secret", "getMe"),
	)

	assert.Equal(t,
		"https://api.telegram.org/file/bot1234:secret/photos/user.png",
		defaultBuildFileURL("1234:secret", "photos/user.png"),
	)
}

func TestHTTPTransport_New(t *testing.T) {
	buildEndpointURL := func(token, method string) string {
		return fmt.Sprintf("https://api.telegram.local/bot%s/%s", token, method)
	}

	buildFileURL := func(token, path string) string {
		return fmt.Sprintf("https://api.telegram.local/file/bot%s/%s", token, path)
	}

	doer := &HTTPDoerMock{}

	transport := NewHTTPTransport(
		WithHTTPDoer(doer),
		WithHTTPBuildCallURLFunc(buildEndpointURL),
		WithHTTPBuildFileURLFunc(buildFileURL),
	)

	assert.Equal(t, doer, transport.doer)
	assert.Equal(t,
		"https://api.telegram.local/bottest/getMe",
		transport.buildCallURL("test", "getMe"),
	)

	assert.Equal(t,
		"https://api.telegram.local/file/bottest/photo.png",
		transport.buildDownloadURL("test", "photo.png"),
	)
}

func TestHTTPTransport_Download(t *testing.T) {
	ctx := context.Background()

	t.Run("Error", func(t *testing.T) {
		doer := &HTTPDoerMock{
			DoFunc: func(r *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, r.Method)

				return &http.Response{
					StatusCode: http.StatusOK,
					Header: http.Header{
						"Content-Type": []string{"application/octet-stream"},
					},
					Body: ioutil.NopCloser(
						strings.NewReader("test, test, test..."),
					),
				}, nil
			},
		}

		trans := NewHTTPTransport(
			WithHTTPDoer(doer),
		)

		res, err := trans.Download(ctx, "123:secret", "image.png")
		assert.NoError(t, err)
		if assert.NotNil(t, res) {
			body, err := ioutil.ReadAll(res)
			if assert.NoError(t, err) {
				assert.Equal(t, "test, test, test...", string(body))
			}
		}

	})

	t.Run("Success", func(t *testing.T) {
		doError := errors.New("do test error")

		doer := &HTTPDoerMock{
			DoFunc: func(r *http.Request) (*http.Response, error) {
				return nil, doError
			},
		}

		trans := NewHTTPTransport(
			WithHTTPDoer(doer),
		)

		result, err := trans.Download(ctx, "123:secret", "image.png")
		assert.Nil(t, result)
		assert.EqualError(t, errors.Cause(err), "do test error")
	})
}

func TestHTTPTransport_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("Streaming", func(t *testing.T) {
		file := RequestFile{
			Name: "test.txt",
			Body: strings.NewReader("test, test, test"),
		}

		request := NewRequest("sendPhoto").
			WithToken("12345:secret").
			AddString("chat_id", "@channely_updates").
			AddFile("document", file)

		t.Run("BadRequest", func(t *testing.T) {
			doError := errors.New("do test error")

			doer := &HTTPDoerMock{
				DoFunc: func(r *http.Request) (*http.Response, error) {
					return nil, doError
				},
			}

			transport := NewHTTPTransport(
				WithHTTPDoer(doer),
			)

			response, err := transport.Execute(ctx, request)

			assert.EqualError(t, errors.Cause(err), "do test error")
			assert.Nil(t, response, "response should be nil if error")
		})

		t.Run("OK", func(t *testing.T) {
			doer := &HTTPDoerMock{
				DoFunc: func(r *http.Request) (*http.Response, error) {
					assert.Equal(t, http.MethodPost, r.Method)

					// read body
					_, err := io.Copy(ioutil.Discard, r.Body)
					require.NoError(t, err)

					return &http.Response{
						StatusCode: http.StatusOK,
						Header: http.Header{
							"Content-Type": []string{"application/json"},
						},
						Body: ioutil.NopCloser(
							strings.NewReader(`{
								"ok": true,
								"result": {
									"photo": {}
								}
							}`),
						),
					}, nil
				},
			}

			transport := NewHTTPTransport(
				WithHTTPDoer(doer),
			)

			response, err := transport.Execute(ctx, request)
			assert.NoError(t, err)
			if assert.NotNil(t, response) {
				assert.Equal(t, http.StatusOK, response.StatusCode)
				assert.True(t, response.OK)
				assert.JSONEq(t, `{"photo": {}}`, string(response.Result))
			}
		})

	})

	t.Run("Simple", func(t *testing.T) {
		request := NewRequest("sendPhoto").
			WithToken("12345:secret").
			AddString("chat_id", "@channely_updates").
			AddString("text", "test")

		t.Run("BadRequest", func(t *testing.T) {
			doError := errors.New("do test error")

			doer := &HTTPDoerMock{
				DoFunc: func(r *http.Request) (*http.Response, error) {
					return nil, doError
				},
			}

			transport := NewHTTPTransport(
				WithHTTPDoer(doer),
			)

			response, err := transport.Execute(ctx, request)

			assert.EqualError(t, errors.Cause(err), "do test error")
			assert.Nil(t, response, "response should be nil if error")
		})

		t.Run("OK", func(t *testing.T) {
			doer := &HTTPDoerMock{
				DoFunc: func(r *http.Request) (*http.Response, error) {
					assert.Equal(t, http.MethodPost, r.Method)

					// read body
					_, err := io.Copy(ioutil.Discard, r.Body)
					require.NoError(t, err)

					return &http.Response{
						StatusCode: http.StatusOK,
						Header: http.Header{
							"Content-Type": []string{"application/json"},
						},
						Body: ioutil.NopCloser(
							strings.NewReader(`{
								"ok": true,
								"result": {
									"text": "test"
								}
							}`),
						),
					}, nil
				},
			}

			transport := NewHTTPTransport(
				WithHTTPDoer(doer),
			)

			response, err := transport.Execute(ctx, request)
			assert.NoError(t, err)
			if assert.NotNil(t, response) {
				assert.Equal(t, http.StatusOK, response.StatusCode)
				assert.True(t, response.OK)
				assert.JSONEq(t, `{"text": "test"}`, string(response.Result))
			}
		})
	})
}
