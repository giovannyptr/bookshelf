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

---
 Project Structure


bookshelf/
├─ bookshelf-backend/
│ ├─ cmd/server/main.go # main entry
│ ├─ internal/auth/ # login/register handlers
│ ├─ internal/books/ # book CRUD handlers
│ ├─ internal/users/ # user repository
│ ├─ models/ # GORM models
│ ├─ docs/ # auto-generated swagger files
│ └─ go.mod / go.sum
│
├─ bookshelf-frontend/
│ ├─ src/views/BooksPage.vue # list view
│ ├─ src/views/BookDetail.vue # detail/edit view
│ ├─ src/components/BookForm.vue # shared form
│ ├─ src/lib/api.js # axios instance
│ └─ .env
│
└─ uploads/ # image storage folder


---

