package web

import (
	err2 "db-server/err"
	"db-server/server"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func StoragePut(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	file, fileHeader, err := r.FormFile("file")
	defer func() { _ = file.Close() }()
	objectName := fileHeader.Filename
	contentType := fileHeader.Header["Content-Type"][0]
	path, info := server.UploadToS3(file, "", objectName, contentType)

	log.Debug("Successfully uploaded %s of size %d\n", objectName, info.Size)

	resp := make(map[string]string)

	resp["path"] = path

	wr, _ := json.Marshal(resp)

	w.WriteHeader(200)
	_, err = w.Write(wr)
	err2.DebugErr(err)
}
