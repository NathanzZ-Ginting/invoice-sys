package main

import (
	"log"
	"net/http"
	"os"

	"invoice-backend/services/invoice-service/internal/handler"
	"invoice-backend/services/invoice-service/internal/repository"
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
	repo := repository.NewInvoiceRepository(db)
	h := handler.NewInvoiceHandler(repo)

	// Setup routes
	r := mux.NewRouter()
	r.HandleFunc("/invoices", h.GetAll).Methods("GET")
	r.HandleFunc("/invoices/{id}", h.GetByID).Methods("GET")
	r.HandleFunc("/invoices", h.Create).Methods("POST")
	r.HandleFunc("/invoices/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/invoices/{id}", h.Delete).Methods("DELETE")
	r.HandleFunc("/invoices/{id}/pdf", h.GeneratePDF).Methods("GET")
	r.HandleFunc("/currency-rates", h.GetCurrencyRates).Methods("GET")
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"invoice-service"}`))
	}).Methods("GET")

	// Apply middleware (no CORS - handled by API Gateway)
	handler := middleware.Recovery(middleware.Logger("INVOICE-SERVICE")(r))

	// Start server
	port := os.Getenv("INVOICE_SERVICE_PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("🚀 Invoice Service running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
