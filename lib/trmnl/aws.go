package trmnl

import (
	"bytes"
	"context"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	bucketName         = os.Getenv("BUCKET_NAME")
	credentialFilename = os.Getenv("CREDENTIAL_FILENAME")
	writeToFilePath    = os.Getenv("WRITE_TO_FILE_PATH")
)

func UploadToS3(pngBytes []byte) {
	if len(writeToFilePath) > 0 {
		if err := os.WriteFile(writeToFilePath, pngBytes, 0644); err != nil {
			log.Fatal(err)
		}
	}

	s3Session, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion("eu-west-1"),
		config.WithSharedConfigProfile("terminal"),
		config.WithSharedCredentialsFiles([]string{credentialFilename}),
	)
	if err != nil {
		log.Println(err)
		return
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancelFunc()
	client := s3.NewFromConfig(s3Session)
	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Body:   bytes.NewReader(pngBytes),
		Bucket: aws.String(bucketName),
		Key:    aws.String("temperature-graph.png"),
	})
	if err != nil {
		log.Println(err)
	}
}
