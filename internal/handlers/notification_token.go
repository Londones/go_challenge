package handlers

import (
	"encoding/json"
	"go-challenge/internal/database/queries"
	"go-challenge/internal/models"
	"go-challenge/internal/config"
	"net/http"
	"strconv"

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
    notificationTokenIDStr := chi.URLParam(r, "id")
    if notificationTokenIDStr == "" {
        http.Error(w, "notification ID is required", http.StatusBadRequest)
        return
    }

    notificationTokenID, err := strconv.Atoi(notificationTokenIDStr)
    if err != nil {
        http.Error(w, "notification ID must be an integer", http.StatusBadRequest)
        return
    }

    err = h.notificationTokenQueries.DeleteNotificationToken(notificationTokenID)
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

	SendToToken(config.GetFirebaseApp(), notification.Token, notification.Text, notification.Title)

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

func SendToToken(app *firebase.App, fcmToken string, text string, title string) {
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}

	//fcmToken := "d3adVu0tMDyBUTPkgc_l-0:APA91bHjb6-wWkT1ABGSasFqxrsOR3AdfcTjLc8b7f7yukWLt32GS4UA5XdIwZ8p98oOLp-CBcyuYaCYdEPRji_f2WSXO9JKb7XPjotm_3bdkk-7hJyxJS8JuUHt82xzGGJ6Aacy0QWb"

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  text,
		},
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

// func (h *NotificationTokenHandler) DeleteNotificationTokenHandler(w http.ResponseWriter, r *http.Request) {
// 	notificationTokenID := chi.URLParam(r, "id")
// 	if notificationTokenID == "" {
// 		http.Error(w, "notification ID is required", http.StatusBadRequest)
// 		return
// 	}

// 	err := h.notificationTokenQueries.DeleteNotificationToken(notificationTokenID)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusNoContent)
// }

