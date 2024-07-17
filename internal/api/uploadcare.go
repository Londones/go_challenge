package api

import (
	"context"
	"fmt"
	// "net/http"
	"os"

	"github.com/uploadcare/uploadcare-go/ucare"
	"github.com/uploadcare/uploadcare-go/upload"
)

func CreateUCClient() (ucare.Client, error) {
	creds := ucare.APICreds{
		SecretKey: os.Getenv("UPLOAD_CARE_SECRET_KEY"),
		PublicKey: os.Getenv("UPLOAD_CARE_PUBLIC_KEY"),
	}

	conf := &ucare.Config{
		SignBasedAuthentication: true,
		APIVersion:              ucare.APIv06,
	}

	client, err := ucare.NewClient(creds, conf)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func UploadImage(client ucare.Client, file string) (string, string, error) {
	return uploadFile(client, file)
}

func uploadFile(client ucare.Client, file string) (string, string, error) {
	uploadService := upload.NewService(client)

	f, err := os.Open(file)
	if err != nil {
		return "", "", fmt.Errorf("could not open file: %v", err)
	}

	defer f.Close()

	// Check the file type
	buffer := make([]byte, 512)
	_, err = f.Read(buffer)
	if err != nil {
		return "", "", fmt.Errorf("could not read file: %v", err)
	}

	param := upload.FileParams{
		Data: f,
		Name: f.Name(),
	}

	fileID, err := uploadService.File(context.Background(), param)
	if err != nil {
		return "", "", fmt.Errorf("could not upload file: %v", err)
	}

	fileURL := constructFileURL(fileID)

	return fileURL, fileID, nil
}

func UploadFilePDF(client ucare.Client, file string) (string, string, error) {
	return uploadFile(client, file)
}

func constructFileURL(fileID string) string {
	return "https://ucarecdn.com/" + fileID + "/"
}
