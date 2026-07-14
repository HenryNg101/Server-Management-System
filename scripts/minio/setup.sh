#!/bin/sh

mc alias set myminio http://${MINIO_INTERNAL_ENDPOINT} ${MINIO_ACCESS_KEY} ${MINIO_SECRET_KEY};

# Create bucket (idempotent by, if exist, just returns as true instead)
mc mb myminio/${MINIO_BUCKET} || true;

# Set public access (idempotent by, if exist, just returns as true instead)
mc anonymous set public myminio/${MINIO_BUCKET} || true;

echo "MinIO setup completed!"

exit 0;