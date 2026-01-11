.PHONY: help dev build clean docker-up docker-down

help: ## Show this help
	@echo "Invoice System Microservices"
	@echo ""
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

dev: ## Run all services in development mode
	@echo "🚀 Starting all microservices..."
	@make -j6 dev-gateway dev-customer dev-invoice dev-payment dev-analytics dev-notification

dev-gateway: ## Run API Gateway
	@echo "🌐 Starting API Gateway..."
	@cd services/api-gateway/cmd && go run main.go

dev-customer: ## Run Customer Service
	@echo "👥 Starting Customer Service..."
	@cd services/customer-service/cmd && go run main.go

dev-invoice: ## Run Invoice Service
	@echo "📄 Starting Invoice Service..."
	@cd services/invoice-service/cmd && go run main.go

dev-payment: ## Run Payment Service
	@echo "💰 Starting Payment Service..."
	@cd services/payment-service/cmd && go run main.go

dev-analytics: ## Run Analytics Service
	@echo "📊 Starting Analytics Service..."
	@cd services/analytics-service/cmd && go run main.go

dev-notification: ## Run Notification Service
	@echo "📧 Starting Notification Service..."
	@cd services/notification-service/cmd && go run main.go

build: ## Build all services
	@echo "🔨 Building all services..."
	@cd services/api-gateway/cmd && go build -o ../../../bin/api-gateway
	@cd services/customer-service/cmd && go build -o ../../../bin/customer-service
	@cd services/invoice-service/cmd && go build -o ../../../bin/invoice-service
	@cd services/payment-service/cmd && go build -o ../../../bin/payment-service
	@cd services/analytics-service/cmd && go build -o ../../../bin/analytics-service
	@cd services/notification-service/cmd && go build -o ../../../bin/notification-service
	@echo "✅ All services built successfully!"

docker-up: ## Start all services with Docker Compose
	@echo "🐳 Starting services with Docker..."
	docker-compose up --build

docker-down: ## Stop all Docker services
	@echo "🛑 Stopping Docker services..."
	docker-compose down

clean: ## Clean build artifacts
	@echo "🧹 Cleaning..."
	@rm -rf bin/
	@echo "✅ Cleaned!"

test: ## Run tests for all services
	@echo "🧪 Running tests..."
	@go test ./services/... -v

.DEFAULT_GOAL := help
