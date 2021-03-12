package domain

import (
	"context"
	"io"
)

type Filestore interface {
	Upload(ctx context.Context, file io.Reader, filename string) (string, error)
	GetDirectDownloadURL(id string) (string, error)
}
