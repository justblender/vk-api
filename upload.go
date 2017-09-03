package vk_api

import (
	"bytes"
	"mime/multipart"
	"os"
	"io"
	"net/http"
)

// UploadFile explains itself by the name, right?
// Format depends on what file you want to upload, video_file is for videos and so on..
func UploadFile(url, format, path string) error {
	contentType, buffer, err := getFilePart(format, path)
	if err != nil {
		return err
	}
	_, err = http.Post(url, contentType, buffer)
	if err != nil {
		return err
	}
	return nil
}

func getFilePart(format, path string) (contentType string, buffer *bytes.Buffer, err error) {
	buffer = &bytes.Buffer{}
	writer := multipart.NewWriter(buffer)

	fileWriter, err := writer.CreateFormFile(format, path)
	if err != nil {
		return
	}
	fh, err := os.Open(path)
	if err != nil {
		return
	}
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return
	}
	contentType = writer.FormDataContentType()
	writer.Close()
	return
}
