package di

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

// @inject
func ProvideFCMClient() *messaging.Client {
	credentialsFile := "firebase-adminsdk.json"

	// If credentials file does not exist, return nil without crashing the app
	if _, err := os.Stat(credentialsFile); os.IsNotExist(err) {
		log.Println("⚠️ firebase-adminsdk.json not found. FCM Client disabled.")
		return nil
	}

	opt := option.WithCredentialsFile(credentialsFile)
	firebaseApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Printf("⚠️ Failed to initialize Firebase app: %v\n", err)
		return nil
	}

	client, err := firebaseApp.Messaging(context.Background())
	if err != nil {
		log.Printf("⚠️ Failed to initialize FCM Client: %v\n", err)
		return nil
	}

	log.Println("✅ Firebase FCM Client initialized successfully")
	return client
}
