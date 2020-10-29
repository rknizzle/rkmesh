package main

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/rknizzle/rkmesh/domain"
	"github.com/satori/go.uuid"
	"io"
)

type s3FileRepository struct {
	Session *session.Session
	Bucket  string
}

func NewS3FileRepository(session *session.Session, bucket string) domain.FileRepository {
	return &s3FileRepository{session, bucket}
}

func (s *s3FileRepository) Upload(ctx context.Context, file io.Reader, filename string) (string, error) {
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
