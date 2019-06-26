package tg

import (
	"context"
	"io"
)

// Transport define interface for execute requests.
type Transport interface {
	Execute(ctx context.Context, r *Request) (*Response, error)
	Download(ctx context.Context, token string, path string) (io.ReadCloser, error)
}
