package aws

import (
	"bytes"
	"capstone/util"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
)

func UploadFileS3(name string, file *multipart.FileHeader, folder string) (string, error) {
	errenv := godotenv.Load()

	if errenv != nil {
		log.Fatal("error load env file")
	}
	src, err := file.Open()
	if err != nil {
		return "", err
	}

	file.Filename = util.GenerateRandomString(name)
	defer src.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, src); err != nil {
		return "", err
	}

	awsRegion := aws.String(os.Getenv("AWS_REGION"))

	sess, err := session.NewSession(&aws.Config{
		Region: awsRegion,
	})
	if err != nil {
		return "", err
	}

	s3Client := s3.New(sess)
	s3Key := fmt.Sprintf("uploads/%s/%s", folder ,file.Filename)
	s3Bucket := os.Getenv("AWS_S3_BUCKET")
	objectInput := &s3.PutObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(s3Key),
		Body:   bytes.NewReader(buf.Bytes()),
	}

	if _, err := s3Client.PutObject(objectInput); err != nil {
		return "", err
	}

	imageURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s3Bucket, *awsRegion, s3Key)
	fmt.Println("i", imageURL)

	return imageURL, nil
}
