package filestorage

import (
	"io"
	"log"
	// "github.com/aws/aws-sdk-go/aws"
	// "github.com/aws/aws-sdk-go/aws/session"
	// "github.com/aws/aws-sdk-go/service/s3"
)

const (
// S3_REGION = "us-west-2"
)

func UploadFileToS3(file io.ReadSeeker, contentLength int64, bucket string, fileName string) error {
	// TODO - implement it
	log.Printf("NOT IMPLEMENTED YET: upload to aws s3")
	return nil
	// sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(S3_REGION)}))

	// _, err := s3.New(sess).PutObject(&s3.PutObjectInput{
	// 	Bucket:        aws.String(bucket),
	// 	Key:           aws.String(fileName),
	// 	Body:          file,
	// 	ContentLength: aws.Int64(contentLength),
	// 	ACL:           aws.String("public-read"),
	// })

	// return err
}
