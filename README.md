# 🚀 Savannah POS - Modern E-commerce Backend

A high-performance backend system built with Go that handles everything from secure authentication to real-time order notifications. Designed to demonstrate modern backend engineering practices.

## ✨ What Makes This Special

**🔥 Real Business Logic** - Not just another CRUD app. Features hierarchical categories, OIDC auth, and dual notification systems.

**🛡️ Production-Ready Security** - OpenID Connect authentication with proper state validation and JWT management.

**📱 Smart Notifications** - Orders trigger instant SMS to customers + email alerts to admins via Africa's Talking API.

**🌳 Hierarchical Categories** - Products organized in unlimited-depth categories using efficient closure table pattern.

## 🏗️ Tech Stack: Go 1.21+ • Gin Framework • PostgreSQL • Docker • Kubernetes • OIDC • Africa's Talking SMS

# 1. Clone and run
git clone https://github.com/Doris-Mwito5/savannah-pos.git

cd savannah-pos

# 2. Start database
docker run -d -p 5432:5432 --name savannah-db \
  -e POSTGRES_PASSWORD=postgres postgres:17-alpine
  
Quick Start 
# 3. Launch the API
go run cmd/server/main.go

# 🎉 API running at http://localhost:8080

💡 Experience the Features
# 1. Start OIDC authentication flow
curl -X POST http://localhost:8080/v1/auth/login

# 2. Create an order (triggers real SMS + email)
curl -X POST http://localhost:8080/v1/orders \
  -H "Authorization: Bearer {token}" \
  -d '{"items": [{"product_id": 1, "quantity": 2}]}'

# 3. Get average price by category hierarchy
curl http://localhost:8080/v1/categories/3/average-price


