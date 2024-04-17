package s3client

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	conf "github.com/muhwyndhamhp/tigerhall-kittens/utils/config"
)

type S3Client struct {
	client *s3.Client
}

const defaultBucketName = "tigerhall-kittens"

func NewS3Client() *S3Client {
	accountId := conf.Get(conf.CF_ACCOUNT_ID)
	accessKeyId := conf.Get(conf.CF_R2_ACCESS_KEY_ID)
	accessKeySecret := conf.Get(conf.CF_R2_SECRET_ACCESS_KEY)

	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId),
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		log.Fatal(err)
	}

	client := s3.NewFromConfig(cfg)

	return &S3Client{client}
}

func (c *S3Client) UploadImage(ctx context.Context, r *bytes.Reader, filename, contentType string, size int64) (string, error) {
	obj := &s3.PutObjectInput{
		Bucket:        aws.String(defaultBucketName),
		Key:           aws.String(AppendTimestamp(filename)),
		ContentType:   aws.String("image/jpeg"),
		Body:          r,
		ContentLength: aws.Int64(size),
	}

	_, err := c.client.PutObject(ctx, obj)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://tigerhall-kittens.mwyndham.dev/%s", filename), nil
}

func AppendTimestamp(fileName string) string {
	extension := filepath.Ext(fileName)
	name := fileName[0 : len(fileName)-len(extension)]
	fileName = fmt.Sprintf("%s-%s%s", name, time.Now().Format("20060102150405"), extension)
	return fileName
}
