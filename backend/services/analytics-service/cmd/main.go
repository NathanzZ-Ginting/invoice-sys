package main

import (
	"log"
	"net/http"
	"os"

	"invoice-backend/services/analytics-service/internal/handler"
	"invoice-backend/services/analytics-service/internal/repository"
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
	repo := repository.NewAnalyticsRepository(db)
	h := handler.NewAnalyticsHandler(repo)

	// Setup routes
	r := mux.NewRouter()
	r.HandleFunc("/dashboard/stats", h.GetDashboardStats).Methods("GET")
	r.HandleFunc("/dashboard/revenue", h.GetRevenueByPeriod).Methods("GET")
	r.HandleFunc("/dashboard/top-customers", h.GetTopCustomers).Methods("GET")
	r.HandleFunc("/dashboard/overdue", h.GetOverdueInvoices).Methods("GET")
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"analytics-service"}`))
	}).Methods("GET")

	// Apply middleware (no CORS - handled by API Gateway)
	handler := middleware.Recovery(middleware.Logger("ANALYTICS-SERVICE")(r))

	// Start server
	port := os.Getenv("ANALYTICS_SERVICE_PORT")
	if port == "" {
		port = "8085"
	}

	log.Printf("🚀 Analytics Service running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
