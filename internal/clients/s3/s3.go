package s3

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	s3go "github.com/aws/aws-sdk-go/service/s3"
)

type Client interface {
	GetPresignedUploadUrl(path, fileName string) (string, error)
}

type client struct {
	s3Client *s3go.S3
	s3Bucket string
}

type ClientOptions struct {
	Key      string
	Secret   string
	Endpoint string
	Bucket   string
}

func NewS3Client(opts ClientOptions) (Client, error) {
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(opts.Key, opts.Secret, ""),
		Endpoint:         aws.String(opts.Endpoint),
		Region:           aws.String("us-east-1"),
		S3ForcePathStyle: aws.Bool(false),
	}
	newSession, err := session.NewSession(s3Config)
	if err != nil {
		return nil, err
	}
	return &client{
		s3Client: s3go.New(newSession),
		s3Bucket: opts.Bucket,
	}, nil
}

func (c *client) GetPresignedUploadUrl(path, fileName string) (string, error) {
	req, _ := c.s3Client.PutObjectRequest(&s3go.PutObjectInput{
		Bucket: aws.String(c.s3Bucket),
		Key:    aws.String(path + "/" + fileName),
	})
	urlStr, err := req.Presign(15 * time.Minute)
	if err != nil {
		return "", err
	}
	return urlStr, nil
}
