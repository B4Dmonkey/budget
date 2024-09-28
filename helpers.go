package main

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
)

func SaveFileToDisk(header *multipart.FileHeader, file multipart.File) error {
	// todo: check if the folder is there. Not important since this is for me
	out, err := os.Create("uploads/" + header.Filename)
	if err != nil {
		return errors.New("Unable to create the file: " + err.Error())
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return errors.New("Error saving file: " + err.Error())
	}

	return nil
}
