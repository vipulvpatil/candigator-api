package s3

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	s3go "github.com/aws/aws-sdk-go/service/s3"
)

type Client interface {
	GetPresignedUploadUrl(bucket, key string) (string, error)
}

type client struct {
	s3Client *s3go.S3
}

type ClientOptions struct {
	Key    string
	Secret string
}

func NewS3Client(opts ClientOptions) (Client, error) {
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(opts.Key, opts.Secret, ""),
		Endpoint:         aws.String("https://candidate-tracker-dev.fra1.digitaloceanspaces.com"),
		Region:           aws.String("us-east-1"),
		S3ForcePathStyle: aws.Bool(false), // // Configures to use subdomain/virtual calling format. Depending on your version, alternatively use o.UsePathStyle = false
	}
	newSession, err := session.NewSession(s3Config)
	if err != nil {
		return nil, err
	}
	return &client{
		s3Client: s3go.New(newSession),
	}, nil
}

func (c *client) GetPresignedUploadUrl(bucket, key string) (string, error) {
	req, _ := c.s3Client.PutObjectRequest(&s3go.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	urlStr, err := req.Presign(5 * time.Minute)
	if err != nil {
		return "", err
	}
	return urlStr, nil
}
