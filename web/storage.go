package web

import (
	"context"
	err2 "db-server/err"
	"encoding/json"
	"github.com/getsentry/sentry-go"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
)

func StoragePut(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

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

	file, fileHeader, err := r.FormFile("file")

	defer func() { _ = file.Close() }()

	objectName := fileHeader.Filename
	contentType := fileHeader.Header["Content-Type"][0]

	info, err := minioClient.PutObject(ctx, bucketName, objectName, file, fileHeader.Size, minio.PutObjectOptions{ContentType: contentType})

	if err != nil {
		sentry.CaptureException(err)
	}

	log.Debug("Successfully uploaded %s of size %d\n", objectName, info.Size)

	path := os.Getenv("STORAGE_PUBLIC_URL") + "/" + os.Getenv("STORAGE_BUCKET") + "/" + objectName

	resp := make(map[string]string)

	resp["path"] = path

	wr, _ := json.Marshal(resp)

	w.WriteHeader(200)
	_, err = w.Write(wr)
	err2.DebugErr(err)
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
