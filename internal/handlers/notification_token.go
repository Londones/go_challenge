package handlers

import (
	"encoding/json"
	"go-challenge/internal/database/queries"
	"go-challenge/internal/models"
	"go-challenge/internal/config"
	"net/http"

	"github.com/go-chi/chi/v5"

	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
)

type NotificationTokenHandler struct {
	notificationTokenQueries      *queries.DatabaseService
}

func NewNotificationTokenHandler(notificationTokenQueries *queries.DatabaseService) *NotificationTokenHandler {
	return &NotificationTokenHandler{notificationTokenQueries: notificationTokenQueries}
}

func (h *NotificationTokenHandler) CreateNotificationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var notificationToken models.NotificationToken
	err := json.NewDecoder(r.Body).Decode(&notificationToken)
	if err != nil {
		http.Error(w, "error decoding notification", http.StatusBadRequest)
		return
	}

	err = h.notificationTokenQueries.CreateNotificationToken(&notificationToken)
	if err != nil {
		http.Error(w, "error saving notification", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(notificationToken)
}

func (h *NotificationTokenHandler) DeleteNotificationTokenHandler(w http.ResponseWriter, r *http.Request) {
    userID := chi.URLParam(r, "userID")
    if userID == "" {
        http.Error(w, "userID is required", http.StatusBadRequest)
        return
    }

    err := h.notificationTokenQueries.DeleteNotificationTokenForUser(userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

func (h *NotificationTokenHandler) SendNotificationHandler(w http.ResponseWriter, r *http.Request) {
	var notification models.Notification
	err := json.NewDecoder(r.Body).Decode(&notification)
	if err != nil {
		http.Error(w, "error decoding notification", http.StatusBadRequest)
		return
	}

	payload := make(map[string]string)
	payload["RoomID"] = notification.RoomID
	SendToToken(config.GetFirebaseApp(), notification.Token, notification.Text, notification.Title, payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(notification)

}

func FindUserToken(notificationTokenQueries *queries.DatabaseService, userID string) (string, error) {
	notificationToken, err := notificationTokenQueries.GetNotificationTokenByUserID(userID)
	if err != nil {
		return "", err
	}
	return notificationToken.Token, nil

}

func SendToToken(app *firebase.App, fcmToken string, text string, title string, payload map[string]string) {
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  text,
		},
		Data: payload,
		Token: fcmToken,
	}

	response, err := client.Send(ctx, message)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Successfully sent message:", response)
}

func (h *NotificationTokenHandler) TestSendNotificationHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	client, err := config.GetFirebaseApp().Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}

	// query token from database
	fcmToken := "dqLX69OBQnexTV5QJEdRTj:APA91bHCmnH8MKNpYZuTaxONWwyAOs2mK1hUxOhKqS2MK5Nk-1WUCZdpC1BVQcFnbDGW6wcrBeS6c67nmOUuYdieXOCRp2_7zPYt8XTcYjTb8rZYmdn4EjSmh45VkYo9x8WgYw7BThxt"

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: "Title",
			Body:  "Message",
		},
		Token: fcmToken,
	}

	response, err := client.Send(ctx, message)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Successfully sent message:", response)
}

