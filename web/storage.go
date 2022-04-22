package web

import (
	"context"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"net/http"
	"os"
)

func StoragePut(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	endpoint := os.Getenv("STORAGE_ENDPOINT")
	accessKeyID := os.Getenv("STORAGE_ACCESS_KEY")
	secretAccessKey := os.Getenv("STORAGE_SECRET_KEY")
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		sentry.CaptureException(err)
	}

	bucketName := os.Getenv("STORAGE_BUCKET")

	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: os.Getenv("STORAGE_LOCATION")})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			sentry.CaptureException(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}

	file, fileHeader, err := r.FormFile("file")

	defer file.Close()

	//// Upload the zip file
	objectName := fileHeader.Filename
	//filePath := "/tmp/golden-oldies.zip"
	contentType := fileHeader.Header["Content-Type"][0]

	fmt.Println(fileHeader.Header)
	fmt.Println(contentType)

	info, err := minioClient.PutObject(ctx, bucketName, objectName, file, fileHeader.Size, minio.PutObjectOptions{ContentType: contentType})
	//
	//// Upload the zip file with FPutObject
	//info, err := minioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}
	//
	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)

}
