# 📚 Bookshelf — Mini Book Management System

A full-stack CRUD web app for managing books, built with:
- **Backend:** Go (Gin, GORM, JWT)
- **Frontend:** Vue 3 + Vite
- **Database:** PostgreSQL
- **Docs:** Swagger (auto-generated with Swag)

---

## 🚀 Features
✅ User authentication (JWT) 
✅ Admin can create / edit / delete books  
✅ File upload for book covers  
✅ REST API with Swagger docs  
✅ CORS-ready for Vue frontend  
✅ Auto-migration and admin seeding 
✅ dark-light theme


## 🧠 Backend Setup (Gin + GORM)

### 1️⃣ Requirements
- Go 1.21+
- PostgreSQL
- Swag CLI (`go install github.com/swaggo/swag/cmd/swag@latest`)

### 2️⃣ Environment Variables
Create `.env` inside `bookshelf-backend/` or export manually:

