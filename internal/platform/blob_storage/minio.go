package blob_storage

import (
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinIOSession(endpoint, accessKey, secretKey string, useSSL bool) *minio.Client {
	conn, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatal("Failed to initiate a MinIO session:", err)
	}
	return conn
}
