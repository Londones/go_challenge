package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	// "encoding/json"
	// "strings"
	"encoding/base64"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

var firebaseApp *firebase.App

func GetFirebaseApp() *firebase.App {
	if firebaseApp == nil {
		firebaseApp, _, _ = SetupFirebase()
	}
	return firebaseApp
}

func SetupFirebase() (*firebase.App, context.Context, *messaging.Client) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file", err)
	}

	ctx := context.Background()

	sdk, err := base64.StdEncoding.DecodeString(os.Getenv("FIREBASE_SDK"))
	if err != nil {
		panic("Failed to decode FIREBASE_SDK")
	}

	opt := option.WithCredentialsJSON(sdk)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		panic("Firebase load error")
	}

	// Messaging client
	client, err := app.Messaging(ctx)
	if err != nil {
		panic("Firebase messaging client error")
	}

	return app, ctx, client
}

func SetupFirebaseTemp() (*firebase.App, context.Context, *messaging.Client) {

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
