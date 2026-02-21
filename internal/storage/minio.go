package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

type MinIOClient struct {
	client *minio.Client
	bucket string
	logger *zap.Logger
}

func NewMinIOClient(endpoint, accessKey, secretKey, bucket string, useSSL bool, logger *zap.Logger) (*MinIOClient, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	return &MinIOClient{
		client: client,
		bucket: bucket,
		logger: logger,
	}, nil
}

func (m *MinIOClient) EnsureBucket(ctx context.Context) error {
	exists, err := m.client.BucketExists(ctx, m.bucket)
	if err != nil {
		return fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		err = m.client.MakeBucket(ctx, m.bucket, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
		m.logger.Info("bucket created", zap.String("bucket", m.bucket))
	}

	return nil
}

func (m *MinIOClient) Upload(ctx context.Context, objectName string, reader io.Reader, objectSize int64, contentType string) (string, error) {
	_, err := m.client.PutObject(ctx, m.bucket, objectName, reader, objectSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload object: %w", err)
	}

	return fmt.Sprintf("/%s/%s", m.bucket, objectName), nil
}

func (m *MinIOClient) Download(ctx context.Context, objectName string) (*minio.Object, error) {
	object, err := m.client.GetObject(ctx, m.bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to download object: %w", err)
	}
	return object, nil
}

func (m *MinIOClient) Delete(ctx context.Context, objectName string) error {
	err := m.client.RemoveObject(ctx, m.bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}
	return nil
}

func (m *MinIOClient) GetPresignedURL(ctx context.Context, objectName string) (string, error) {
	url, err := m.client.PresignedGetObject(ctx, m.bucket, objectName, 3600, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get presigned url: %w", err)
	}
	return url.String(), nil
}
