module invoice-backend/services/notification-service

go 1.25.4

replace invoice-backend/services/shared => ../shared

require (
	github.com/gorilla/mux v1.8.1
	github.com/joho/godotenv v1.5.1
	invoice-backend/services/shared v0.0.0-00010101000000-000000000000
)
