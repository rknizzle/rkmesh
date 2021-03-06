package testFileStore

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const (
	fsHost        = "localhost:9000"
	fsRegion      = "us-east-1"
	fsAccess      = "AKIAIOSFODNN7EXAMPLE"
	fsSecret      = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
	fsModelBucket = "rkmesh"
)

// A TestFilestore adds/removes objects to/from a bucket for use in testing
type TestFilestore struct {
	bucket   string
	uploader *s3manager.Uploader
	svc      *s3.S3
}

// InitTestFileStorage creates a session that the application can use to access a test bucket and
// also returns a TestFilestore that can be used for seeding and clearing that bucket while testing
func InitTestFileStore() (TestFilestore, *session.Session, string, error) {
	sess, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(fsAccess, fsSecret, ""),
		Region:           aws.String(fsRegion),
		Endpoint:         aws.String(fsHost),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		return TestFilestore{}, nil, "", err
	}

	uploader := s3manager.NewUploader(sess)
	svc := s3.New(sess)

	tfs := TestFilestore{uploader: uploader, bucket: fsModelBucket, svc: svc}
	return tfs, sess, fsModelBucket, nil
}

// Seed places test objects into the test bucket for integration tests
func (t *TestFilestore) Seed() {
	// TODO: s.Uplaod a couple test files
}

// Clear deletes all seed objects from the test bucket
func (t *TestFilestore) Clear() error {
	iter := s3manager.NewDeleteListIterator(t.svc, &s3.ListObjectsInput{
		Bucket: aws.String(t.bucket),
	})

	err := s3manager.NewBatchDeleteWithClient(t.svc).Delete(aws.BackgroundContext(), iter)
	if err != nil {
		return err
	}
	return nil
}
