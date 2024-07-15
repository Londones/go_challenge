package config

import (
	"context"
	"fmt"
	"path/filepath"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)


var firebaseApp *firebase.App

func GetFirebaseApp() *firebase.App {
    if firebaseApp == nil {
        SetupFirebase()
    }
    return firebaseApp
}

func SetupFirebase() (*firebase.App, context.Context, *messaging.Client) {

	ctx := context.Background()

	serviceAccountKeyFilePath, err := filepath.Abs("./internal/config/serviceAccountKey.json")
	if err != nil {
		panic("Unable to load serviceAccountKeys.json file")
	}
    fmt.Println(serviceAccountKeyFilePath)

	opt := option.WithCredentialsFile(serviceAccountKeyFilePath)

	//Firebase admin SDK initialization
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		panic("Firebase load error")
	}

	//Messaging client
	client, _ := app.Messaging(ctx)
    firebaseApp = app

	return app, ctx, client
}
