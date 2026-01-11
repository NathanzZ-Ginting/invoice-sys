#!/bin/bash

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${BLUE}🚀 Starting Invoice System Microservices (Manual Mode)...${NC}\n"

# Check .env
if [ ! -f .env ]; then
    echo -e "${RED}❌ .env file not found!${NC}"
    exit 1
fi

# Export environment variables
export $(cat .env | grep -v '^#' | xargs)

# Kill any existing processes on these ports
echo -e "${YELLOW}Cleaning up existing processes...${NC}"
lsof -ti:8080,8081,8082,8083,8084,8085 | xargs kill -9 2>/dev/null || true
sleep 2

echo -e "${GREEN}Starting services...${NC}\n"

# Start API Gateway
echo -e "${BLUE}[1/6]${NC} Starting API Gateway on port 8080..."
cd /Users/nathangtg/Invoice-sYS-generator/backend/services/api-gateway/cmd
GOWORK=off go run main.go > /tmp/api-gateway.log 2>&1 &
API_GATEWAY_PID=$!
cd - > /dev/null

sleep 2

# Start Customer Service
echo -e "${BLUE}[2/6]${NC} Starting Customer Service on port 8082..."
cd /Users/nathangtg/Invoice-sYS-generator/backend/services/customer-service/cmd
GOWORK=off go run main.go > /tmp/customer-service.log 2>&1 &
CUSTOMER_PID=$!
cd - > /dev/null

sleep 1

# Start Invoice Service
echo -e "${BLUE}[3/6]${NC} Starting Invoice Service on port 8081..."
cd /Users/nathangtg/Invoice-sYS-generator/backend/services/invoice-service/cmd
GOWORK=off go run main.go > /tmp/invoice-service.log 2>&1 &
INVOICE_PID=$!
cd - > /dev/null

sleep 1

# Start Payment Service
echo -e "${BLUE}[4/6]${NC} Starting Payment Service on port 8083..."
cd /Users/nathangtg/Invoice-sYS-generator/backend/services/payment-service/cmd
GOWORK=off go run main.go > /tmp/payment-service.log 2>&1 &
PAYMENT_PID=$!
cd - > /dev/null

sleep 1

# Start Analytics Service
echo -e "${BLUE}[5/6]${NC} Starting Analytics Service on port 8085..."
cd /Users/nathangtg/Invoice-sYS-generator/backend/services/analytics-service/cmd
GOWORK=off go run main.go > /tmp/analytics-service.log 2>&1 &
ANALYTICS_PID=$!
cd - > /dev/null

sleep 1

# Start Notification Service
echo -e "${BLUE}[6/6]${NC} Starting Notification Service on port 8084..."
cd /Users/nathangtg/Invoice-sYS-generator/backend/services/notification-service/cmd
GOWORK=off go run main.go > /tmp/notification-service.log 2>&1 &
NOTIFICATION_PID=$!
cd - > /dev/null

sleep 3

# Check which services are running
echo -e "\n${YELLOW}Checking services...${NC}"
check_port() {
    local port=$1
    local name=$2
    if lsof -i:$port > /dev/null 2>&1; then
        echo -e "${GREEN}✅ $name (port $port) - RUNNING${NC}"
    else
        echo -e "${RED}❌ $name (port $port) - FAILED${NC}"
        echo -e "${YELLOW}   Check logs: /tmp/${name,,}.log${NC}"
    fi
}

echo ""
check_port 8080 "API-Gateway"
check_port 8082 "Customer-Service"
check_port 8081 "Invoice-Service"
check_port 8083 "Payment-Service"
check_port 8085 "Analytics-Service"
check_port 8084 "Notification-Service"

echo -e "\n${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}API Gateway:${NC}           http://localhost:8080"
echo -e "${GREEN}Customer Service:${NC}      http://localhost:8082"
echo -e "${GREEN}Invoice Service:${NC}       http://localhost:8081"
echo -e "${GREEN}Payment Service:${NC}       http://localhost:8083"
echo -e "${GREEN}Analytics Service:${NC}     http://localhost:8085"
echo -e "${GREEN}Notification Service:${NC}  http://localhost:8084"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

echo -e "\n${YELLOW}Logs location: /tmp/*-service.log${NC}"
echo -e "${YELLOW}To stop all: pkill -f 'go run main.go'${NC}"
echo -e "${YELLOW}To view logs: tail -f /tmp/api-gateway.log${NC}\n"

# Keep script running and handle Ctrl+C
trap "echo -e '\n${YELLOW}Stopping all services...${NC}'; pkill -f 'go run main.go'; exit 0" INT

echo -e "${GREEN}Press Ctrl+C to stop all services${NC}\n"
wait
