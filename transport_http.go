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

// HTTPEncoder extends Encoder interface
// with HTTP request specific methods.
type HTTPEncoder interface {
	Encoder

	// Returns content type of request
	ContentType() string

	// End request.
	Close() error
}

// MultipartEncoder encodes the request using multipart encoding.
type MultipartEncoder struct {
	w *multipart.Writer
}

// NewMultipartEncoder creates multipart encoder.
func NewMultipartEncoder(writer io.Writer) *MultipartEncoder {
	return &MultipartEncoder{
		w: multipart.NewWriter(writer),
	}
}

// AddString encodes string value
func (enc *MultipartEncoder) AddString(k string, v string) error {
	return enc.w.WriteField(k, v)
}

// AddFile encodes file value.
func (enc *MultipartEncoder) AddFile(k string, file InputFile) error {
	writer, err := enc.w.CreateFormFile(k, file.Name)
	if err != nil {
		return errors.Wrapf(err, "create form file '%s'", k)
	}

	if _, err := io.Copy(writer, file.Body); err != nil {
		return errors.Wrapf(err, "copy to form file '%s'", k)
	}

	return nil
}

// ContentType returns HTTP request content type.
func (enc *MultipartEncoder) ContentType() string {
	return enc.w.FormDataContentType()
}

// Close multipart encoder.
func (enc *MultipartEncoder) Close() error {
	return enc.w.Close()
}

// URLEncodedEncoder encodes request using urlencoded.
type URLEncodedEncoder struct {
	dst   io.Writer
	total int
}

// NewURLEncodedEncoder creates urlencoded encoder.
func NewURLEncodedEncoder(dst io.Writer) *URLEncodedEncoder {
	return &URLEncodedEncoder{dst: dst}
}

// AddString encodes string value.
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

// AddFile not supported by this encoder, returns error.
func (w *URLEncodedEncoder) AddFile(k string, file InputFile) error {
	return errors.New("URLEncodedEncoder does not support file uploading")
}

// Close encoder. no-op.
func (w *URLEncodedEncoder) Close() error { return nil }

// Returns content type
func (w *URLEncodedEncoder) ContentType() string {
	return "application/x-www-form-urlencoded"
}

// HTTPDoer define interface of used HTTP client.
type HTTPDoer interface {
	Do(r *http.Request) (*http.Response, error)
}

// HTTPTransport default transport for requests.
// Uses multipart for requests containing files, and urlencoded for others.
type HTTPTransport struct {
	doer HTTPDoer

	buildCallURL     func(token string, method string) string
	buildDownloadURL func(token string, path string) string
}

var (
	defaultBuildCallURL = func(token, method string) string {
		return fmt.Sprintf("https://api.telegram.org/bot%s/%s", token, method)
	}

	defaultBuildFileURL = func(token, path string) string {
		return fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", token, path)
	}
)

// HTTPTransportOption use this for configure transport.
type HTTPTransportOption func(t *HTTPTransport)

// WithHTTPDoer sets HTTPTransport request executor.
func WithHTTPDoer(doer HTTPDoer) HTTPTransportOption {
	return func(t *HTTPTransport) {
		t.doer = doer
	}
}

// WithHTTPBuildCallFunc sets function used for build API call URLs.
func WithHTTPBuildCallURLFunc(f func(token, method string) string) HTTPTransportOption {
	return func(t *HTTPTransport) {
		t.buildCallURL = f
	}
}

// WithHTTPBuildFileURLFunc sets function used for build download URLs.
func WithHTTPBuildFileURLFunc(f func(token, path string) string) HTTPTransportOption {
	return func(t *HTTPTransport) {
		t.buildDownloadURL = f
	}
}

// NewHTTPTransport creates HTTPTransport with default configuration.
func NewHTTPTransport(opts ...HTTPTransportOption) *HTTPTransport {
	t := &HTTPTransport{
		doer:             http.DefaultClient,
		buildCallURL:     defaultBuildCallURL,
		buildDownloadURL: defaultBuildFileURL,
	}

	for _, opt := range opts {
		opt(t)
	}

	return t
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

	// TODO: handle status and content type

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
		defer pw.Close()
		defer encoder.Close()

		if err := r.Encode(encoder); err != nil {
			errChan <- err
		}
	}()

	// send
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
