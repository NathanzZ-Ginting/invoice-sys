package main

import (
	"log"
	"net/http"
	"os"

	"invoice-backend/services/customer-service/internal/handler"
	"invoice-backend/services/customer-service/internal/repository"
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
	repo := repository.NewCustomerRepository(db)
	h := handler.NewCustomerHandler(repo)

	// Setup routes
	r := mux.NewRouter()
	r.HandleFunc("/customers", h.GetAll).Methods("GET")
	r.HandleFunc("/customers/{id}", h.GetByID).Methods("GET")
	r.HandleFunc("/customers", h.Create).Methods("POST")
	r.HandleFunc("/customers/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/customers/{id}", h.Delete).Methods("DELETE")
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"customer-service"}`))
	}).Methods("GET")

	// Apply middleware (no CORS - handled by API Gateway)
	handler := middleware.Recovery(middleware.Logger("CUSTOMER-SERVICE")(r))

	// Start server
	port := os.Getenv("CUSTOMER_SERVICE_PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("🚀 Customer Service running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
