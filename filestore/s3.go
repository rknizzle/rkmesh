package filestore

import (
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/rknizzle/rkmesh/domain"
	uuid "github.com/satori/go.uuid"
)

type s3Filestore struct {
	bucket   string
	uploader *s3manager.Uploader
	svc      *s3.S3
}

func NewS3Filestore(session *session.Session, bucket string) domain.Filestore {
	uploader := s3manager.NewUploader(session)
	svc := s3.New(session)
	return &s3Filestore{bucket, uploader, svc}
}

func (s *s3Filestore) Upload(ctx context.Context, file io.Reader, filename string) (string, error) {
	u := uuid.NewV4()

	key := filename + "-" + u.String()
	_, err := s.uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String(s.bucket),
		ACL:    aws.String("public-read"),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		return "", err
	}
	return key, nil
}

func (s *s3Filestore) GetDirectDownloadURL(id string) (string, error) {
	req, _ := s.svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(id),
	})
	urlStr, err := req.Presign(30 * time.Minute)
	if err != nil {
		return "", err
	}

	return urlStr, nil
}
