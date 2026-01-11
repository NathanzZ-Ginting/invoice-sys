package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"invoice-backend/services/api-gateway/internal/proxy"
	"invoice-backend/services/shared/pkg/middleware"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	godotenv.Load("../../.env")

	// Setup router
	r := mux.NewRouter()

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"api-gateway"}`))
	}).Methods("GET")

	// Route to services
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		
		// Determine which service to route to based on path
		var serviceURL string
		
		if strings.HasPrefix(path, "/customers") {
			serviceURL = proxy.GetServiceURL("CUSTOMER_SERVICE")
		} else if strings.HasPrefix(path, "/invoices") || strings.HasPrefix(path, "/currency-rates") {
			serviceURL = proxy.GetServiceURL("INVOICE_SERVICE")
		} else if strings.HasPrefix(path, "/payments") {
			serviceURL = proxy.GetServiceURL("PAYMENT_SERVICE")
		} else if strings.HasPrefix(path, "/dashboard") || strings.HasPrefix(path, "/analytics") {
			serviceURL = proxy.GetServiceURL("ANALYTICS_SERVICE")
		} else if strings.HasPrefix(path, "/notifications") {
			serviceURL = proxy.GetServiceURL("NOTIFICATION_SERVICE")
		} else {
			http.Error(w, "Service not found", http.StatusNotFound)
			return
		}
		
		proxy.ProxyRequest(w, r, serviceURL)
	})

	// Apply middleware
	handler := middleware.Recovery(middleware.CORS(middleware.Logger("API-GATEWAY")(r)))

	// Start server
	port := os.Getenv("API_GATEWAY_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🌐 API Gateway running on port %s", port)
	log.Printf("📡 Routing:")
	log.Printf("  /customers/*     -> Customer Service  (port %s)", os.Getenv("CUSTOMER_SERVICE_PORT"))
	log.Printf("  /invoices/*      -> Invoice Service   (port %s)", os.Getenv("INVOICE_SERVICE_PORT"))
	log.Printf("  /payments/*      -> Payment Service   (port %s)", os.Getenv("PAYMENT_SERVICE_PORT"))
	log.Printf("  /dashboard/*     -> Analytics Service (port %s)", os.Getenv("ANALYTICS_SERVICE_PORT"))
	log.Printf("  /notifications/* -> Notification Svc  (port %s)", os.Getenv("NOTIFICATION_SERVICE_PORT"))
	
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
