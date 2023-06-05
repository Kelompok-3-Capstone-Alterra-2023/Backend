package aws

import (
	"bytes"
	"capstone/util"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func UploadFileS3(name string, file *multipart.FileHeader) (string, error) {
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

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})
	if err != nil {
		return "", err
	}

	s3Client := s3.New(sess)
	// change the path
	s3Key := fmt.Sprintf("uploads/thumbnail/%s", file.Filename)
	s3Bucket := "capstone-prevent"
	objectInput := &s3.PutObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(s3Key),
		Body:   bytes.NewReader(buf.Bytes()),
	}

	if _, err := s3Client.PutObject(objectInput); err != nil {
		return "", err
	}

	imagePath := s3Bucket + "/" + s3Key

	return imagePath, nil
}
