package trmnl

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"log"
	"os"
)

var (
	bucketName         = os.Getenv("BUCKET_NAME")
	credentialFilename = os.Getenv("CREDENTIAL_FILENAME")
)

func UploadToS3(pngBytes []byte) {
	s3Session := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("eu-west-1"),
		Credentials: credentials.NewSharedCredentials(credentialFilename, "terminal"),
	}))
	uploader := s3manager.NewUploader(s3Session)
	_, err := uploader.Upload(&s3manager.UploadInput{
		Body:   bytes.NewBuffer(pngBytes),
		Bucket: aws.String(bucketName),
		Key:    aws.String("temperature-graph.png"),
	})
	if err != nil {
		log.Println(err)
	}
}
