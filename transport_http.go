package tg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

type HTTPEncoder interface {
	Encoder

	ContentType() string
	Close() error
}

type MultipartEncoder struct {
	w *multipart.Writer
}

func NewMultipartEncoder(writer io.Writer) *MultipartEncoder {
	return &MultipartEncoder{
		w: multipart.NewWriter(writer),
	}
}

func (enc *MultipartEncoder) AddString(k string, v string) error {
	return enc.w.WriteField(k, v)
}

func (enc *MultipartEncoder) AddFile(k string, file RequestFile) error {
	writer, err := enc.w.CreateFormFile(k, file.Name)

	if err != nil {
		return errors.Wrapf(err, "create form file '%s'", k)
	}

	if _, err := io.Copy(writer, file.Body); err != nil {
		return errors.Wrapf(err, "copy to form file '%s'", k)
	}

	return nil
}

func (enc *MultipartEncoder) ContentType() string {
	return enc.w.FormDataContentType()
}

func (enc *MultipartEncoder) Close() error {
	return enc.w.Close()
}

type URLEncodedEncoder struct {
	dst   io.Writer
	total int
}

func NewURLEncodedEncoder(dst io.Writer) *URLEncodedEncoder {
	return &URLEncodedEncoder{dst: dst}
}

func (w *URLEncodedEncoder) AddString(k string, v string) error {
	buf := strings.Builder{}

	if w.total > 0 {
		buf.WriteByte('&')
	}

	buf.WriteString(url.QueryEscape(k))
	buf.WriteByte('=')
	buf.WriteString(url.QueryEscape(v))

	n, err := io.WriteString(w.dst, buf.String())
	if err != nil {
		return err
	}
	w.total += n

	return nil
}

func (w *URLEncodedEncoder) AddFile(k string, file RequestFile) error {
	return errors.New("URLEncodedEncoder does not support file uploading")
}

func (w *URLEncodedEncoder) Close() error { return nil }

func (w *URLEncodedEncoder) ContentType() string {
	return "application/x-www-form-urlencoded"
}

type HTTPDoer interface {
	Do(r *http.Request) (*http.Response, error)
}

type HTTPTransport struct {
	doer HTTPDoer

	buildCallURL     func(token string, method string) string
	buildDownloadURL func(token string, path string) string
}

func NewHTTPTransport() *HTTPTransport {
	return &HTTPTransport{
		doer: http.DefaultClient,

		buildCallURL: func(token string, method string) string {
			return fmt.Sprintf("https://api.telegram.org/bot%s/%s", token, method)
		},

		buildDownloadURL: func(token string, path string) string {
			return fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", token, path)
		},
	}
}

func (t *HTTPTransport) newMultipartEncoder(w io.Writer) HTTPEncoder {
	return NewMultipartEncoder(w)
}

func (t *HTTPTransport) newUrlencodedEncoder(w io.Writer) HTTPEncoder {
	return NewURLEncodedEncoder(w)
}

func (t *HTTPTransport) Execute(ctx context.Context, r *Request) (*Response, error) {
	if r.HasFiles() {
		return t.executeStreaming(
			ctx,
			t.newMultipartEncoder,
			r,
		)
	}

	return t.executeSimple(
		ctx,
		t.newUrlencodedEncoder,
		r,
	)
}

func (t *HTTPTransport) Download(ctx context.Context, token string, path string) (io.ReadCloser, error) {
	url := t.buildDownloadURL(token, path)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "build request")
	}

	req = req.WithContext(ctx)

	res, err := t.doer.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "do request")
	}

	return res.Body, nil
}

func (t *HTTPTransport) executeSimple(
	ctx context.Context,
	newEncoder func(io.Writer) HTTPEncoder,
	r *Request,
) (*Response, error) {
	buf := &bytes.Buffer{}

	encoder := newEncoder(buf)

	if err := r.Encode(encoder); err != nil {
		return nil, errors.Wrap(err, "encode")
	}

	if err := encoder.Close(); err != nil {
		return nil, errors.Wrap(err, "encoder close")
	}

	req, err := t.buildHTTPRequest(
		r,
		buf,
		encoder.ContentType(),
	)

	if err != nil {
		return nil, errors.Wrap(err, "build http request")
	}

	res, err := t.executeHTTPRequest(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "execute http request")
	}

	return res, nil
}

func (t *HTTPTransport) executeHTTPRequest(ctx context.Context, r *http.Request) (*Response, error) {
	r = r.WithContext(ctx)

	// execute request
	res, err := t.doer.Do(r)
	if err != nil {
		return nil, errors.Wrap(err, "execute request")
	}
	defer res.Body.Close()

	// read content
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read response")
	}

	response := &Response{
		StatusCode: res.StatusCode,
	}

	// unmarshal content
	if err := json.Unmarshal(content, &response); err != nil {
		return nil, errors.Wrap(err, "unmarshal response")
	}

	return response, nil
}

func (t *HTTPTransport) buildHTTPRequest(
	r *Request,
	body io.Reader,
	contentType string,
) (*http.Request, error) {
	url := t.buildCallURL(r.Token(), r.Method())

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, errors.Wrap(err, "build request")
	}

	// set content type
	req.Header.Set("Content-Type", contentType)

	return req, nil
}

func (t *HTTPTransport) executeStreaming(
	ctx context.Context,
	newEncoder func(io.Writer) HTTPEncoder,
	r *Request,
) (*Response, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	pr, pw := io.Pipe()

	encoder := newEncoder(pw)

	resChan := make(chan *Response)
	errChan := make(chan error)

	// upload
	go func() {
		defer func() {
			if err := pw.Close(); err != nil {
				errChan <- errors.Wrap(err, "close pipe writer")
			}
		}()

		defer func() {
			if err := encoder.Close(); err != nil {
				errChan <- errors.Wrap(err, "close multipart writer")
			}
		}()

		if err := r.Encode(encoder); err != nil {
			errChan <- err
		}
	}()

	go func() {
		req, err := t.buildHTTPRequest(r, pr, encoder.ContentType())
		if err != nil {
			errChan <- errors.Wrap(err, "build http request")
			return
		}

		res, err := t.executeHTTPRequest(ctx, req)
		if err != nil {
			errChan <- errors.Wrap(err, "execute http request")
			return
		}

		resChan <- res
	}()

	select {
	case err := <-errChan:
		return nil, err
	case res := <-resChan:
		return res, nil
	}
}
