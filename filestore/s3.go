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
	Session *session.Session
	Bucket  string
}

func NewS3Filestore(session *session.Session, bucket string) domain.Filestore {
	return &s3Filestore{session, bucket}
}

func (s *s3Filestore) Upload(ctx context.Context, file io.Reader, filename string) (string, error) {
	uploader := s3manager.NewUploader(s.Session)

	u := uuid.NewV4()

	up, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.Bucket),
		ACL:    aws.String("public-read"),
		Key:    aws.String(filename + "-" + u.String()),
		Body:   file,
	})
	if err != nil {
		return "", err
	}
	return up.Location, nil
}
