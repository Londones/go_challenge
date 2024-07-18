package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"go-challenge/internal/database/queries"
	"go-challenge/internal/models"
	"go-challenge/internal/utils"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
)

type ModifiedReport []map[string]interface{}

type ReportsHandler struct {
	reportsQueries *queries.DatabaseService
	clients        map[*websocket.Conn]bool
	broadcast      chan []byte
	upgrader       websocket.Upgrader
	mu             sync.Mutex
	lastPosition   int64
}

type LogEntry struct {
	Level     string          `json:"level"`
	Title     string          `json:"title"`
	Message   string          `json:"msg"`
	Time      time.Time       `json:"time"`
	ExtraData json.RawMessage `json:"extra_data,omitempty"`
}

func NewReportsHandler(reportsQueries *queries.DatabaseService) *ReportsHandler {
	handler := &ReportsHandler{
		reportsQueries: reportsQueries,
		clients:        make(map[*websocket.Conn]bool),
		broadcast:      make(chan []byte),
		upgrader:       websocket.Upgrader{},
		lastPosition:   0,
	}

	go handler.broadcastReports()
	go handler.watchLogFile()

	return handler
}

func (h *ReportsHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	h.mu.Lock()
	h.clients[conn] = true
	h.mu.Unlock()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			h.mu.Lock()
			delete(h.clients, conn)
			h.mu.Unlock()
			break
		}
	}
}

func (h *ReportsHandler) watchLogFile() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		utils.Logger("error", "File Watcher", "Error creating watcher", err.Error())
		return
	}
	defer watcher.Close()

	err = watcher.Add("logs.log")
	if err != nil {
		utils.Logger("error", "File Watcher", "Error adding file to watcher", err.Error())
		return
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				h.processNewLogEntries()
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			utils.Logger("error", "File Watcher", "Error watching file", err.Error())
		}
	}
}

func (h *ReportsHandler) processNewLogEntries() {
	file, err := os.Open("logs.log")
	if err != nil {
		utils.Logger("error", "File Processing", "Error opening log file", err.Error())
		return
	}
	defer file.Close()

	_, err = file.Seek(h.lastPosition, 0)
	if err != nil {
		utils.Logger("error", "File Processing", "Error seeking in file", err.Error())
		return
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var logEntry LogEntry
		err := json.Unmarshal([]byte(line), &logEntry)
		if err != nil {
			utils.Logger("error", "File Processing", "Error unmarshaling log entry", err.Error())
			continue
		}

		if logEntry.Title == "Reported Message" || logEntry.Title == "Reported Annonce" {
			h.broadcast <- []byte(logEntry.Message)
		}
	}

	h.lastPosition, _ = file.Seek(0, 1)
}

func (h *ReportsHandler) broadcastReports() {
	for {
		message := <-h.broadcast
		h.mu.Lock()
		for client := range h.clients {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				utils.Logger("error", "WebSocket", "Error sending message", err.Error())
				client.Close()
				delete(h.clients, client)
			}
		}
		h.mu.Unlock()
	}
}

func (h *ReportsHandler) CreateReportedMessage(w http.ResponseWriter, r *http.Request) {
	var reportedMessage *models.ReportedMessage
	if err := json.NewDecoder(r.Body).Decode(&reportedMessage); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	reportedMessage, err := h.reportsQueries.CreateReportedMessage(reportedMessage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	message, err := h.reportsQueries.GetMessageByID(reportedMessage.MessageID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	reason, err := h.reportsQueries.GetReportReasonById(reportedMessage.ReasonID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	modifiedReport := ModifiedReport{
		{
			"id":             reportedMessage.ID,
			"message":        message.Content,
			"reporterUserId": reportedMessage.ReporterUserID,
			"reportedUserId": reportedMessage.ReportedUserID,
			"createdAt":      reportedMessage.CreatedAt,
			"reason":         reason.Reason,
			"isHandled":      reportedMessage.IsHandled,
			"type":           "message",
		},
	}

	modifiedReportToJSON, err := json.Marshal(modifiedReport)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.Logger("warn", "Reported Message", "A user reported a message", string(modifiedReportToJSON))

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
}

func (h *ReportsHandler) CreateReportedAnnonce(w http.ResponseWriter, r *http.Request) {
	var reportedAnnonce *models.ReportedAnnonce
	if err := json.NewDecoder(r.Body).Decode(&reportedAnnonce); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(reportedAnnonce)

	reportedAnnonce, err := h.reportsQueries.CreateReportedAnnonce(reportedAnnonce)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	annonce, err := h.reportsQueries.GetAnnonceByID(reportedAnnonce.AnnonceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	reason, err := h.reportsQueries.GetReportReasonById(reportedAnnonce.ReasonID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	modifiedReport := ModifiedReport{
		{
			"id": reportedAnnonce.ID,
			"annonce": map[string]interface{}{
				"ID":          annonce.ID,
				"Title":       annonce.Title,
				"Description": annonce.Description,
			},
			"reporterUserId": reportedAnnonce.ReporterUserID,
			"reportedUserId": reportedAnnonce.ReportedUserID,
			"createdAt":      reportedAnnonce.CreatedAt,
			"reason":         reason.Reason,
			"isHandled":      reportedAnnonce.IsHandled,
			"type":           "annonce",
		},
	}

	modifiedReportToJSON, err := json.Marshal(modifiedReport)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.Logger("warn", "Reported Annonce", "A user reported an annonce", string(modifiedReportToJSON))

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
}

func (h *ReportsHandler) GetReportedMessages(w http.ResponseWriter, r *http.Request) {
	reportedMessages, err := h.reportsQueries.GetReportedMessages()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	modifiedReport := make(ModifiedReport, 0)
	for _, reportedMessage := range reportedMessages {
		message, err := h.reportsQueries.GetMessageByID(reportedMessage.MessageID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		reason, err := h.reportsQueries.GetReportReasonById(reportedMessage.ReasonID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		modifiedReport = append(modifiedReport, map[string]interface{}{
			"id":             reportedMessage.ID,
			"message":        message.Content,
			"reporterUserId": reportedMessage.ReporterUserID,
			"reportedUserId": reportedMessage.ReportedUserID,
			"createdAt":      reportedMessage.CreatedAt,
			"reason":         reason.Reason,
			"isHandled":      reportedMessage.IsHandled,
			"type":           "message",
		})
	}

	modifiedReportToJSON, err := json.Marshal(modifiedReport)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(modifiedReportToJSON)
}

func (h *ReportsHandler) GetReportedAnnonces(w http.ResponseWriter, r *http.Request) {
	reportedAnnonces, err := h.reportsQueries.GetReportedAnnonces()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	modifiedReport := make(ModifiedReport, 0)
	for _, reportedAnnonce := range reportedAnnonces {
		annonce, err := h.reportsQueries.GetAnnonceByID(reportedAnnonce.AnnonceID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		reason, err := h.reportsQueries.GetReportReasonById(reportedAnnonce.ReasonID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		modifiedReport = append(modifiedReport, map[string]interface{}{
			"id": reportedAnnonce.ID,
			"annonce": map[string]interface{}{
				"ID":          annonce.ID,
				"Title":       annonce.Title,
				"Description": annonce.Description,
			},
			"reporterUserId": reportedAnnonce.ReporterUserID,
			"reportedUserId": reportedAnnonce.ReportedUserID,
			"createdAt":      reportedAnnonce.CreatedAt,
			"reason":         reason.Reason,
			"isHandled":      reportedAnnonce.IsHandled,
			"type":           "annonce",
		})
	}

	modifiedReportToJSON, err := json.Marshal(modifiedReport)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(modifiedReportToJSON)
}

func (h *ReportsHandler) GetAllReports(w http.ResponseWriter, r *http.Request) {
	reportedMessages, err := h.reportsQueries.GetReportedMessages()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	reportedAnnonces, err := h.reportsQueries.GetReportedAnnonces()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	modifiedReport := make(ModifiedReport, 0)
	for _, reportedMessage := range reportedMessages {
		message, err := h.reportsQueries.GetMessageByID(reportedMessage.MessageID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		reason, err := h.reportsQueries.GetReportReasonById(reportedMessage.ReasonID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		modifiedReport = append(modifiedReport, map[string]interface{}{
			"id":             reportedMessage.ID,
			"message":        message.Content,
			"reporterUserId": reportedMessage.ReporterUserID,
			"reportedUserId": reportedMessage.ReportedUserID,
			"createdAt":      reportedMessage.CreatedAt,
			"reason":         reason.Reason,
			"isHandled":      reportedMessage.IsHandled,
			"type":           "message",
		})
	}

	for _, reportedAnnonce := range reportedAnnonces {
		annonce, err := h.reportsQueries.GetAnnonceByID(reportedAnnonce.AnnonceID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		reason, err := h.reportsQueries.GetReportReasonById(reportedAnnonce.ReasonID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		modifiedReport = append(modifiedReport, map[string]interface{}{
			"id": reportedAnnonce.ID,
			"annonce": map[string]interface{}{
				"ID":          annonce.ID,
				"Title":       annonce.Title,
				"Description": annonce.Description,
			},
			"reporterUserId": reportedAnnonce.ReporterUserID,
			"reportedUserId": reportedAnnonce.ReportedUserID,
			"createdAt":      reportedAnnonce.CreatedAt,
			"reason":         reason.Reason,
			"isHandled":      reportedAnnonce.IsHandled,
			"type":           "annonce",
		})
	}

	modifiedReportToJSON, err := json.Marshal(modifiedReport)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(modifiedReportToJSON)
}

func (h *ReportsHandler) GetReportReasons(w http.ResponseWriter, r *http.Request) {
	reasons, err := h.reportsQueries.GetReasons()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	reasonsModifiedJSON := make([]map[string]interface{}, 0)

	for _, reason := range reasons {
		reasonsModifiedJSON = append(reasonsModifiedJSON, map[string]interface{}{
			"id":     reason.ID,
			"reason": reason.Reason,
		})
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(reasonsModifiedJSON); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *ReportsHandler) GetReasonByID(w http.ResponseWriter, r *http.Request) {
	reasonID := r.URL.Query().Get("id")
	reasonIDtoUint, err := strconv.ParseUint(reasonID, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	reason, err := h.reportsQueries.GetReportReasonById(uint(reasonIDtoUint))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	reasonToJSON, err := json.Marshal(reason)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(reasonToJSON)
}
