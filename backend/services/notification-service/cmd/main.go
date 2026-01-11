package main

import (
	"log"
	"net/http"
	"os"

	"invoice-backend/services/notification-service/internal/handler"
	"invoice-backend/services/shared/pkg/middleware"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	godotenv.Load("../../../.env")

	// Initialize handler
	h := handler.NewNotificationHandler()

	// Setup routes
	r := mux.NewRouter()
	r.HandleFunc("/notifications/send", h.SendEmail).Methods("POST")
	r.HandleFunc("/notifications/reminder", h.SendPaymentReminder).Methods("POST")
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"notification-service"}`))
	}).Methods("GET")

	// Apply middleware (no CORS - handled by API Gateway)
	handler := middleware.Recovery(middleware.Logger("NOTIFICATION-SERVICE")(r))

	// Start server
	port := os.Getenv("NOTIFICATION_SERVICE_PORT")
	if port == "" {
		port = "8084"
	}

	log.Printf("🚀 Notification Service running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
