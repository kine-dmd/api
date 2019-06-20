package s3Connection

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"log"
)

type S3Uploader interface {
	UploadFile(bucketName string, s3FilePath string, file io.Reader) error
}

type S3Client struct {
	uploader *s3manager.Uploader
}

func MakeS3Connection() *S3Client {
	conn := new(S3Client)
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("eu-west-2")},
	))
	conn.uploader = s3manager.NewUploader(sess)
	return conn
}

func (s3Conn *S3Client) UploadFile(bucketName string, s3FilePath string, file io.Reader) error {
	// Upload the file to S3
	_, err := s3Conn.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(s3FilePath),
		Body:   file,
	})
	if err != nil {
		log.Printf("Unable to upload to file %s in S3 bucket %s. %s", s3FilePath, bucketName, err)
		return err
	}

	// Success - no error to return
	return nil
}
