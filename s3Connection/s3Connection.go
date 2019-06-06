package s3Connection

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"log"
)

type S3Connection struct {
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
	deleter    *s3.S3
}

func MakeS3Connection() *S3Connection {
	conn := new(S3Connection)
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("eu-west-2")},
	))
	conn.uploader = s3manager.NewUploader(sess)
	conn.downloader = s3manager.NewDownloader(sess)
	conn.deleter = s3.New(sess)
	return conn
}

func (s3Conn *S3Connection) UploadFile(bucketName string, s3FilePath string, file io.Reader) error {
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
