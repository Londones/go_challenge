package api

import (
	"os"

	"github.com/uploadcare/uploadcare-go/ucare"
)

func CreateUCClient() (ucare.Client, error) {
	creds := ucare.APICreds{
		SecretKey: os.Getenv("UPLOADCARE_SECRET_KEY"),
		PublicKey: os.Getenv("UPLOADCARE_PUBLIC_KEY"),
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
