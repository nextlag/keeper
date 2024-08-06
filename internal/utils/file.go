package utils

import (
	"io"
	"mime/multipart"
	"os"
)

// SaveUploadedFile saves an uploaded file to a specified destination.
// It creates the destination directory if it does not exist, then copies
// the contents of the uploaded file to a new file at the destination path.
func SaveUploadedFile(file *multipart.FileHeader, fileName, dst string) error {
	err := os.MkdirAll(dst, os.ModePerm)
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst + "/" + fileName)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
