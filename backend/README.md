# Invoice System - Microservices

Backend sistem invoice dengan arsitektur microservices menggunakan Go + Supabase.

## 🏗️ Arsitektur

```
Frontend → API Gateway :8080 → Customer Service    :8082
                              → Invoice Service     :8081  
                              → Payment Service     :8083
                              → Analytics Service   :8085
                              → Notification Service :8084
```

## 🚀 Quick Start

### 1. Setup Environment
```bash
cp .env.example .env
# Edit .env dengan Supabase credentials
```

### 2. Jalankan Services

**Cara termudah (1 command):**
```bash
./start.sh
```

**Atau pakai Make:**
```bash
make dev
```

**Atau pakai Docker:**
```bash
docker-compose up --build
```

### 3. Test
```bash
curl http://localhost:8080/health
curl http://localhost:8080/customers
curl http://localhost:8080/invoices
```

### 4. Jalankan Frontend
```bash
cd ../frontend
npm install
npm start
```

## 📦 Services

| Service | Port | Endpoint |
|---------|------|----------|
| API Gateway | 8080 | All routes |
| Customer | 8082 | `/customers` |
| Invoice | 8081 | `/invoices` |
| Payment | 8083 | `/payments` |
| Analytics | 8085 | `/dashboard/*` |
| Notification | 8084 | `/notifications/*` |

## 🔧 Troubleshooting

**Port sudah dipakai:**
```bash
lsof -ti:8080 | xargs kill -9
```

**Service error:**
- Cek `.env` sudah diisi
- Cek Go version: `go version` (perlu 1.25+)
- Run `go mod tidy` di folder cmd masing-masing service

## 📝 Environment Variables

```env
SUPABASE_URL=your_url
SUPABASE_KEY=your_key

API_GATEWAY_PORT=8080
CUSTOMER_SERVICE_PORT=8082
INVOICE_SERVICE_PORT=8081
PAYMENT_SERVICE_PORT=8083
ANALYTICS_SERVICE_PORT=8085
NOTIFICATION_SERVICE_PORT=8084

SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your_email
SMTP_PASSWORD=your_password
```

## ⚡ Make Commands

```bash
make dev       # Run all services
make build     # Build all services
make clean     # Clean build files
make help      # Show all commands
```

Done! 🎉
