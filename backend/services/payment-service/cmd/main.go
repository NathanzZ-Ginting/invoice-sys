package main

import (
	"log"
	"net/http"
	"os"

	"invoice-backend/services/payment-service/internal/handler"
	"invoice-backend/services/payment-service/internal/repository"
	"invoice-backend/services/shared/pkg/database"
	"invoice-backend/services/shared/pkg/middleware"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	godotenv.Load("../../../.env")

	// Connect to database
	db, err := database.NewClient()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize repository and handler
	repo := repository.NewPaymentRepository(db)
	h := handler.NewPaymentHandler(repo)

	// Setup routes
	r := mux.NewRouter()
	r.HandleFunc("/payments", h.Create).Methods("POST")
	r.HandleFunc("/payments/invoice/{id}", h.GetByInvoiceID).Methods("GET")
	r.HandleFunc("/payments", h.GetAll).Methods("GET")
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"payment-service"}`))
	}).Methods("GET")

	// Apply middleware (no CORS - handled by API Gateway)
	handler := middleware.Recovery(middleware.Logger("PAYMENT-SERVICE")(r))

	// Start server
	port := os.Getenv("PAYMENT_SERVICE_PORT")
	if port == "" {
		port = "8083"
	}

	log.Printf("🚀 Payment Service running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
