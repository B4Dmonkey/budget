package main

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestUpload(t *testing.T) {
	mockEnv := new(MockEnvConfig)
	mockEnv.On("IsDev").Return(true)
	mockEnv.On("Addr").Return("test_address")
	ctx := context.Background()
	app := CreateApp(ctx, mockEnv)
	app.InitializeServer()
	ts := httptest.NewServer(app.mux)
	defer ts.Close()

	filePath := "/Users/appstack/Downloads/Chase Activity Sept 27.CSV"

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	// Create a new multipart form file
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", file.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Copy the file content to the multipart form
	if _, err := io.Copy(part, file); err != nil {
		t.Fatal(err)
	}
	writer.Close()

	// Create a new HTTP request
	req, err := http.NewRequest(http.MethodPost, ts.URL+"/upload", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Perform the request
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

}
