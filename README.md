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

🧠 COMPLETE RUN GUIDE
1️⃣ UNZIP AND GO INTO PROJECT
cd bookshelf

2️⃣ RUN EVERYTHING WITH DOCKER (RECOMMENDED)
---
cd bookshelf-backend
docker-compose up --build
---


✅ This command does all of the following:

### Builds your Go backend container

### Spins up PostgreSQL container (port 5433)

Auto-migrates database (users, books)

Auto-seeds admin user:

admin@mail.com / adminbookshelf


Serves your API and Swagger at:

http://localhost:8080

http://localhost:8080/swagger/index.html

When you see:

🚀 Server running on http://localhost:8080


it’s ready 🎉

3️⃣ (OPTIONAL) RUN DATABASE MANUALLY

If recruiter prefers to inspect DB directly:

docker exec -it bookshelf-db psql -U bookshelf -d bookshelf


To see data:

SELECT * FROM users;
SELECT * FROM books;


Exit:

\q

4️⃣ (OPTIONAL) RUN BACKEND LOCALLY WITHOUT DOCKER

If Go and PostgreSQL already installed locally:

cd bookshelf-backend/cmd/server
go mod tidy
go run main.go

5️⃣ RUN FRONTEND

In another terminal:

cd bookshelf-frontend
npm install
npm run dev


Frontend → http://localhost:5173

(Ensure .env contains VITE_API_BASE=http://localhost:8080)

✅ QUICK CHECKLIST
Service	Command	URL
Backend (Docker)	docker-compose up --build	http://localhost:8080

Swagger Docs	—	http://localhost:8080/swagger/index.html

Frontend (Vue)	npm run dev	http://localhost:5173

Database (Inspect)	docker exec -it bookshelf-db psql -U bookshelf -d bookshelf	—

