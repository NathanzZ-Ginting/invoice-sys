#!/bin/bash

# Get script directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$SCRIPT_DIR"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

printf "${BLUE}🚀 Starting Invoice System Microservices...${NC}\n\n"

# Check .env
if [ ! -f .env ]; then
    printf "${YELLOW}⚠️  .env file not found!${NC}\n"
    exit 1
fi

# Load .env variables into current shell
set -a
source .env
set +a

# Export env explicitly for subprocesses
export $(grep -v '^#' .env | xargs)

# Run binaries
printf "${GREEN}Starting api-gateway on port 8080...${NC}\n"
./bin/api-gateway &
sleep 1

printf "${GREEN}Starting customer-service on port 8082...${NC}\n"
./bin/customer-service &
sleep 1

printf "${GREEN}Starting invoice-service on port 8081...${NC}\n"
./bin/invoice-service &
sleep 1

printf "${GREEN}Starting payment-service on port 8083...${NC}\n"
./bin/payment-service &
sleep 1

printf "${GREEN}Starting analytics-service on port 8085...${NC}\n"
./bin/analytics-service &
sleep 1

printf "${GREEN}Starting notification-service on port 8084...${NC}\n"
./bin/notification-service &
sleep 1

printf "\n${GREEN}✅ All services started!${NC}\n"
printf "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}\n"
printf "${GREEN}API Gateway:        ${NC}http://localhost:8080\n"
printf "${GREEN}Customer Service:   ${NC}http://localhost:8082\n"
printf "${GREEN}Invoice Service:    ${NC}http://localhost:8081\n"
printf "${GREEN}Payment Service:    ${NC}http://localhost:8083\n"
printf "${GREEN}Analytics Service:  ${NC}http://localhost:8085\n"
printf "${GREEN}Notification Service: ${NC}http://localhost:8084\n"
printf "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}\n"
printf "\n${YELLOW}Press Ctrl+C to stop all services${NC}\n\n"

wait
