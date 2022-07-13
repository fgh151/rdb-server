package server

import (
	"bytes"
	"context"
	"db-server/utils"
	"github.com/getsentry/sentry-go"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"strconv"
)

func UploadToS3(file io.Reader, path string, name string, contentType string) (string, minio.UploadInfo) {
	name = utils.CleanInputString(name)
	minioClient, err := getClient()
	ctx := context.Background()

	bucketName := os.Getenv("STORAGE_BUCKET")

	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: os.Getenv("STORAGE_LOCATION")})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Debug("We already own %s\n", bucketName)
		} else {
			sentry.CaptureException(err)
		}
	} else {
		log.Debug("Successfully created %s\n", bucketName)
	}

	info, err := minioClient.PutObject(ctx, bucketName, path+"/"+name, file, getSize(file), minio.PutObjectOptions{ContentType: contentType})

	if err != nil {
		sentry.CaptureException(err)
	}

	log.Debug("Successfully uploaded %s of size %d\n", name, info.Size)

	resPath := os.Getenv("STORAGE_PUBLIC_URL") + "/" + os.Getenv("STORAGE_BUCKET") + "/" + path + "/" + name

	return resPath, info
}

func getClient() (*minio.Client, error) {
	endpoint := os.Getenv("STORAGE_ENDPOINT")
	accessKeyID := os.Getenv("STORAGE_ACCESS_KEY")
	secretAccessKey := os.Getenv("STORAGE_SECRET_KEY")
	useSSL := true

	if val, err := strconv.ParseBool(os.Getenv("STORAGE_SSL")); err == nil {
		useSSL = val
	}

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		sentry.CaptureException(err)
	}
	return minioClient, err
}

func getSize(stream io.Reader) int64 {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(stream)

	if err != nil {
		return 0
	}

	return int64(buf.Len())
}
