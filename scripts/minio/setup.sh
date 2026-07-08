#!/bin/sh

mc alias set myminio http://${MINIO_INTERNAL_ENDPOINT} ${MINIO_ACCESS_KEY} ${MINIO_SECRET_KEY};
mc mb myminio/${MINIO_BUCKET};
mc anonymous set public myminio/${MINIO_BUCKET};
exit 0;