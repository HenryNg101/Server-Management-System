package data_transfer

import (
	"context"

	"github.com/minio/minio-go/v7"
)

type BlobStorage interface {
	Delete(ctx context.Context, objectName string) error
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

func (m *minIORepository) Delete(ctx context.Context, objectName string) error {
	return m.internalClient.RemoveObject(ctx, m.bucketName, objectName, minio.RemoveObjectOptions{})
}
