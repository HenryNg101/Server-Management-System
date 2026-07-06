package data_transfer

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
)

type BlobStorage interface {
	Upload(ctx context.Context, objectName string, r io.Reader, size int64) (string, error)
	Download(ctx context.Context, objectName string) (io.ReadCloser, error)
}

type minIORepository struct {
	client     *minio.Client
	bucketName string
}

func NewBlobStorage(client *minio.Client, bucketName string) BlobStorage {
	return &minIORepository{
		client:     client,
		bucketName: bucketName,
	}
}

func (m *minIORepository) Upload(ctx context.Context, objectName string, r io.Reader, size int64) (string, error) {
	_, err := m.client.PutObject(ctx, m.bucketName, objectName, r, size, minio.PutObjectOptions{
		ContentType: "text/csv",
	})
	if err != nil {
		return "", err
	}

	return objectName, nil
}

func (m *minIORepository) Download(ctx context.Context, objectName string) (io.ReadCloser, error) {
	obj, err := m.client.GetObject(ctx, m.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return obj, nil
}
