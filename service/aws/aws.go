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

type S3Object struct {
	Bucket string
	Key    string
}

func CreateObject(name, folder string, file *multipart.FileHeader) S3Object {
	errenv := godotenv.Load()
	if errenv != nil {
		log.Fatal("error load env file")
	}

	file.Filename = util.GenerateRandomString(name)
	bucket := os.Getenv("AWS_S3_BUCKET")
	key := fmt.Sprintf("uploads/%s/%s", folder, file.Filename)
	object := S3Object{
		Bucket: bucket,
		Key:    key,
	}

	return object
}

func UploadFileS3(awsObject S3Object, file *multipart.FileHeader) (string, error) {
	// awsS3 := S3Object{}
	// errenv := godotenv.Load()

	// if errenv != nil {
	// 	log.Fatal("error load env file")
	// }
	src, err := file.Open()
	if err != nil {
		return "", err
	}

	// file.Filename = util.GenerateRandomString(name)
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
	// s3Key := fmt.Sprintf("uploads/%s/%s", folder, file.Filename)
	// s3Bucket := os.Getenv("AWS_S3_BUCKET")
	objectInput := &s3.PutObjectInput{
		Bucket: aws.String(awsObject.Bucket),
		Key:    aws.String(awsObject.Key),
		Body:   bytes.NewReader(buf.Bytes()),
	}

	if _, err := s3Client.PutObject(objectInput); err != nil {
		return "", err
	}

	imageURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", awsObject.Bucket, *awsRegion, awsObject.Key)

	return imageURL, nil
}

func DeleteObject(object ...S3Object) {
	errenv := godotenv.Load()
	if errenv != nil {
		log.Fatal("error load env file")
	}
	region := aws.String(os.Getenv("AWS_REGION"))
	sess, err := session.NewSession(&aws.Config{
		Region: region,
	})
	if err != nil {
		log.Println(err)
	}

	s3Client := s3.New(sess)
	for i := range object {
		params := &s3.DeleteObjectInput{
			Bucket: &object[i].Bucket,
			Key:    &object[i].Key,
		}

		_, err = s3Client.DeleteObject(params)
		if err != nil {
			log.Println(err)
		}

	}
}
