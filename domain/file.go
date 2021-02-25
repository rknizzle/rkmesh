package domain

import (
	"context"
	"io"
)

type FileRepository interface {
	Upload(ctx context.Context, file io.Reader, filename string) (string, error)
}
