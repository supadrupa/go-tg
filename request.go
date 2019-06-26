package tg

import (
	"io"
	"strconv"
	"time"
)

// RequestFile represents file to be uploaded with request.
type RequestFile struct {
	Body io.Reader
	Name string
}

// Encoder represents request encoder.
type Encoder interface {
	// Writes string argument k to encoder.
	AddString(k string, v string) error

	// Write files argument k to encoder.
	AddFile(k string, file RequestFile) error
}

// Request represents RPC request to Telegram Bot API.
//
// It contains following info:
//  - token
//  - method
//  - args
//  - files
type Request struct {
	token  string
	method string

	args  map[string]string
	files map[string]RequestFile
}

// NewRequest creates request with provided method.
func NewRequest(method string) *Request {
	return &Request{
		method: method,
	}
}

// Token returns request token.
func (r *Request) Token() string {
	return r.token
}

// Method returns request method.
func (r *Request) Method() string {
	return r.method
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

// AddOptInt adds int argument k to request, if value is not 0.
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

// AddOptBool adds bool argument k to request, if v is not false.
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
func (r *Request) AddFile(k string, file RequestFile) *Request {
	// make files map if not exist
	if r.files == nil {
		r.files = make(map[string]RequestFile)
	}

	r.files[k] = file

	return r
}

// AddPeer adds peer argument k to request.
func (r *Request) AddPeer(k string, peer Peer) *Request {
	peer.AddPeerToRequest(k, r)
	return r
}

// AddChatID adds peer as chat_id argument. Just shortcut.
func (r *Request) AddChatID(peer Peer) *Request {
	return r.AddPeer("chat_id", peer)
}

// AddTime adds time to request as unix timestamp.
func (r *Request) AddTime(k string, v time.Time) *Request {
	return r.AddInt64(k, v.Unix())
}

// AddOptTime adds time to request as unix timestamp, if v is not zero.
func (r *Request) AddOptTime(k string, v time.Time) *Request {
	if !v.IsZero() {
		r.AddTime(k, v)
	}

	return r
}

// RequestPart defines interface of object that can be added to the request.
// It's used for add complex structures and
// isolate logic of addeding to request in struct instead method.
type RequestPart interface {
	AddToRequest(r *Request)
}

// AddPart adds complex objects to request.
func (r *Request) AddPart(part RequestPart) *Request {
	part.AddToRequest(r)
	return r
}

// HasFiles returns true if request contains files.
func (r *Request) HasFiles() bool {
	return len(r.files) > 0
}

// Encode request using encoder.
func (r *Request) Encode(encoder Encoder) error {

	// add files
	for k, v := range r.files {
		if err := encoder.AddFile(k, v); err != nil {
			return err
		}
	}

	// add arguments
	for k, v := range r.args {
		if err := encoder.AddString(k, v); err != nil {
			return err
		}
	}

	return nil
}
