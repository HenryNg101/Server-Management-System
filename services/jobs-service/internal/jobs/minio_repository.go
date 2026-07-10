package jobs

import (
	"context"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
)

type BlobStorage interface {
	Upload(ctx context.Context, objectName string, r io.Reader, size int64) (string, error)
	Download(ctx context.Context, objectName string) (io.ReadCloser, error)
	Delete(ctx context.Context, objectName string) error
	GetPresignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error)
}

type minIORepository struct {
	internalClient *minio.Client
	publicClient   *minio.Client
	bucketName     string
}

func NewBlobStorage(internalClient *minio.Client, publicClient *minio.Client, bucketName string) BlobStorage {
	return &minIORepository{
		internalClient: internalClient,
		publicClient:   publicClient,
		bucketName:     bucketName,
	}
}

func (m *minIORepository) Upload(ctx context.Context, objectName string, r io.Reader, size int64) (string, error) {
	_, err := m.internalClient.PutObject(ctx, m.bucketName, objectName, r, size, minio.PutObjectOptions{
		ContentType: "text/csv",
	})
	if err != nil {
		return "", err
	}

	return objectName, nil
}

func (m *minIORepository) Download(ctx context.Context, objectName string) (io.ReadCloser, error) {
	obj, err := m.internalClient.GetObject(ctx, m.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (m *minIORepository) Delete(ctx context.Context, objectName string) error {
	return m.internalClient.RemoveObject(ctx, m.bucketName, objectName, minio.RemoveObjectOptions{})
}

// Create pre-signed URL for access
func (m *minIORepository) GetPresignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	params := make(url.Values)

	// Force download + clean filename
	params.Set("response-content-disposition", "attachment")

	url, err := m.publicClient.PresignedGetObject(
		ctx, m.bucketName, objectName, expiry, params, // request params
	)
	if err != nil {
		return "", err
	}

	return url.String(), nil
}
