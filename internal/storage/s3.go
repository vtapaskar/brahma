package storage

import (
	"bytes"
	"context"
	"fmt"
	"path"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	cfg "github.com/vtapaskar/brahma/internal/config"
)

type S3Client struct {
	client *s3.Client
	bucket string
	prefix string
}

func NewS3Client(ctx context.Context, s3cfg cfg.S3Config) (*S3Client, error) {
	var awsCfg aws.Config
	var err error

	if s3cfg.AccessKeyID != "" && s3cfg.SecretAccessKey != "" {
		awsCfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(s3cfg.Region),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
				s3cfg.AccessKeyID,
				s3cfg.SecretAccessKey,
				"",
			)),
		)
	} else {
		awsCfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(s3cfg.Region),
		)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	var client *s3.Client
	if s3cfg.Endpoint != "" {
		client = s3.NewFromConfig(awsCfg, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(s3cfg.Endpoint)
			o.UsePathStyle = true
		})
	} else {
		client = s3.NewFromConfig(awsCfg)
	}

	return &S3Client{
		client: client,
		bucket: s3cfg.Bucket,
		prefix: s3cfg.Prefix,
	}, nil
}

func (c *S3Client) Upload(key string, data []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	fullKey := key
	if c.prefix != "" {
		fullKey = path.Join(c.prefix, key)
	}

	_, err := c.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(fullKey),
		Body:   bytes.NewReader(data),
	})

	if err != nil {
		return fmt.Errorf("failed to upload to S3: %w", err)
	}

	return nil
}

func (c *S3Client) Download(key string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	fullKey := key
	if c.prefix != "" {
		fullKey = path.Join(c.prefix, key)
	}

	result, err := c.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(fullKey),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to download from S3: %w", err)
	}
	defer result.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(result.Body)

	return buf.Bytes(), nil
}

func (c *S3Client) GenerateKey(deviceID, reportType, reportID, filename string) string {
	timestamp := time.Now().Format("2006/01/02")
	return path.Join(timestamp, deviceID, reportType, fmt.Sprintf("%s_%s", reportID, filename))
}

func (c *S3Client) Delete(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fullKey := key
	if c.prefix != "" {
		fullKey = path.Join(c.prefix, key)
	}

	_, err := c.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(fullKey),
	})

	if err != nil {
		return fmt.Errorf("failed to delete from S3: %w", err)
	}

	return nil
}
