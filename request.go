package tg

import (
	"io"
	"strconv"
	"time"
)

// Encoder represents request encoder.
type Encoder interface {
	AddString(k string, v string) error
	AddInputFile(k string, name string, content io.Reader) error
}

// Request represents Telegram Bot API request.
type Request struct {
	token  string
	method string
	args   map[string]string
	files  map[string]requestFile
}

type requestFile struct {
	Body io.Reader
	Name string
}

// Token returns request token
func (r *Request) Token() string {
	return r.token
}

// Method returns request method
func (r *Request) Method() string {
	return r.method
}

// NewRequest creates a request object with provided token and method
func NewRequest(method string) *Request {
	return &Request{
		method: method,
	}
}

// WithToken sets request token.
func (r *Request) WithToken(token string) *Request {
	r.token = token
	return r
}

// AddString adds string argument k to request.
func (r *Request) AddString(k string, v string) *Request {
	if r.args == nil {
		r.args = make(map[string]string)
	}

	r.args[k] = v

	return r
}

// AddInt adds int argument k to request.
func (r *Request) AddInt(k string, v int) *Request {
	return r.AddString(k, strconv.Itoa(v))
}

func (r *Request) AddOptInt(k string, v int) *Request {
	if v != 0 {
		r.AddInt(k, v)
	}

	return r
}

// AddInt64 adds int64 argument k to request.
func (r *Request) AddInt64(k string, v int64) *Request {
	return r.AddString(k, strconv.FormatInt(v, 10))
}

// AddBool adds bool argument k to request.
func (r *Request) AddBool(k string, v bool) *Request {
	return r.AddString(k, strconv.FormatBool(v))
}

// AddOptBool adds bool argument k to request.
func (r *Request) AddOptBool(k string, v bool) *Request {
	if v {
		r.AddBool(k, v)
	}
	return r
}

// AddBool adds float64 argument k to request.
func (r *Request) AddFloat64(k string, v float64) *Request {
	return r.AddString(k, strconv.FormatFloat(v, 'f', -1, 64))
}

// AddInputFile adds InputFile to request.
func (r *Request) AddInputFile(k string, name string, content io.Reader) *Request {
	if r.files == nil {
		r.files = make(map[string]requestFile)
	}

	r.files[k] = requestFile{Name: name, Body: content}
	return r
}

func (r *Request) AddPeer(k string, peer Peer) *Request {
	peer.AsPeer(k, r)
	return r
}

func (r *Request) AddChatID(peer Peer) *Request {
	return r.AddPeer("chat_id", peer)
}

func (r *Request) AddTime(k string, v time.Time) *Request {
	return r.AddInt64(k, v.Unix())
}

func (r *Request) AddOptTime(k string, v time.Time) *Request {
	if !v.IsZero() {
		r.AddTime(k, v)
	}

	return r
}

type Addedable interface {
	AddToRequest(r *Request)
}

func (r *Request) Add(add Addedable) *Request {
	add.AddToRequest(r)
	return r
}

// HasInputFile returns true if request contains InputFile.
func (r *Request) HasInputFile() bool {
	return len(r.files) > 0
}

// Encode request via encoder.
func (r *Request) Encode(encoder Encoder) error {
	for k, v := range r.files {
		if err := encoder.AddInputFile(k, v.Name, v.Body); err != nil {
			return err
		}
	}

	for k, v := range r.args {
		if err := encoder.AddString(k, v); err != nil {
			return err
		}
	}

	return nil
}
