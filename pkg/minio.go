package pkg

import (
	"github.com/minio/minio-go"
	"log"
	"os"
)

func NewMinioClient() *minio.Client {

	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKeyID := os.Getenv("MINIO_ACCESS_KEY")
	secretAccessKey := os.Getenv("MINIO_SECRET_KEY")
	useSSL := false

	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil { log.Fatalln(err) }
	return minioClient
}
