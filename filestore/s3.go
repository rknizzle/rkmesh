package filestore

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/rknizzle/rkmesh/domain"
	uuid "github.com/satori/go.uuid"
)

type s3Filestore struct {
	bucket   string
	uploader *s3manager.Uploader
}

func NewS3Filestore(session *session.Session, bucket string) domain.Filestore {
	uploader := s3manager.NewUploader(session)
	return &s3Filestore{bucket, uploader}
}

func (s *s3Filestore) Upload(ctx context.Context, file io.Reader, filename string) (string, error) {
	u := uuid.NewV4()

	up, err := s.uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String(s.bucket),
		ACL:    aws.String("public-read"),
		Key:    aws.String(filename + "-" + u.String()),
		Body:   file,
	})
	if err != nil {
		return "", err
	}
	return up.Location, nil
}
