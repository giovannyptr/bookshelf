# 📚 Bookshelf — Mini Book Management System

A full-stack CRUD web application for managing books, built with modern technologies:

- **Backend:** Go (Gin, GORM, JWT)
- **Frontend:** Vue 3 + Vite
- **Database:** PostgreSQL
- **Documentation:** Swagger (auto-generated)

---

## 🚀 Features

✅ **User Authentication** - JWT-based login/register system  
✅ **Book Management** - Full CRUD operations for books  
✅ **File Upload** - Book cover image upload support  
✅ **REST API** - RESTful API with Swagger documentation  
✅ **Admin Panel** - Admin users can manage all books  
✅ **Responsive UI** - Dark/light theme support  
✅ **CORS Support** - Frontend-backend communication ready  
✅ **Auto Migration** - Database schema auto-setup  

---

## 📋 Prerequisites

Before you begin, ensure you have the following installed:

- **Go 1.21+** - [Download Go](https://golang.org/dl/)
- **Node.js 18+** - [Download Node.js](https://nodejs.org/)
- **PostgreSQL 12+** - [Download PostgreSQL](https://www.postgresql.org/download/)
- **Docker & Docker Compose** (optional) - [Download Docker](https://www.docker.com/get-started)

### Additional Tools

Install these tools for development:

```bash
# Swagger CLI for API documentation generation
go install github.com/swaggo/swag/cmd/swag@latest

# Air for Go hot reload (optional)
go install github.com/air-verse/air@latest
```

---

## 🚀 Quick Start

### Option 1: Using Docker (Recommended)

1. **Clone the repository**
   ```bash
   git clone https://github.com/giovannyptr/bookshelf.git
   cd bookshelf
   ```

2. **Start the database**
   ```bash
   cd bookshelf-backend
   docker-compose up -d
   ```

3. **Start the backend**
   ```bash
   # Install dependencies
   go mod download
   
   # Generate Swagger docs
   swag init -g cmd/server/main.go
   
   # Run the server
   go run cmd/server/main.go
   ```

4. **Start the frontend**
   ```bash
   cd ../bookshelf-frontend
   
   # Install dependencies
   npm install
   
   # Start development server
   npm run dev
   ```

5. **Access the application**
   - Frontend: http://localhost:5173
   - Backend API: http://localhost:8080
   - Swagger Docs: http://localhost:8080/swagger/index.html

### Option 2: Manual Setup

If you prefer not to use Docker, see the detailed setup instructions below.

---

## 🔧 Detailed Setup

### Backend Setup (Go + Gin + GORM)

#### 1. Navigate to Backend Directory
```bash
cd bookshelf-backend
```

#### 2. Setup Database

**Option A: Using Docker**
```bash
# Start PostgreSQL container
docker-compose up -d

# Verify database is running
docker ps
```

**Option B: Manual PostgreSQL Setup**
```bash
# Create database and user
psql -U postgres
CREATE DATABASE bookshelf;
CREATE USER bookshelf WITH ENCRYPTED PASSWORD 'bookshelf';
GRANT ALL PRIVILEGES ON DATABASE bookshelf TO bookshelf;
\q
```

#### 3. Environment Configuration

Create a `.env` file or set environment variables:

```bash
# Database configuration
DB_HOST=localhost
DB_PORT=5433  # 5432 if using manual setup
DB_USER=bookshelf
DB_PASSWORD=bookshelf
DB_NAME=bookshelf

# JWT Secret (change in production)
JWT_SECRET=your-super-secret-jwt-key

# Server configuration
PORT=8080
GIN_MODE=debug  # release for production
```

#### 4. Install Dependencies & Run

```bash
# Download Go modules
go mod download

# Generate Swagger documentation
swag init -g cmd/server/main.go

# Run the application
go run cmd/server/main.go

# Or use Air for hot reload (if installed)
air
```

The backend will be available at `http://localhost:8080`

---

### Frontend Setup (Vue 3 + Vite)

#### 1. Navigate to Frontend Directory
```bash
cd bookshelf-frontend
```

#### 2. Install Dependencies
```bash
# Using npm
npm install

# Or using yarn
yarn install

# Or using pnpm
pnpm install
```

#### 3. Environment Configuration

Create a `.env` file:
```bash
# API Base URL
VITE_API_BASE_URL=http://localhost:8080
```

#### 4. Start Development Server
```bash
# Start development server
npm run dev

# Or using yarn
yarn dev

# Or using pnpm
pnpm dev
```

The frontend will be available at `http://localhost:5173`

#### 5. Build for Production
```bash
# Build for production
npm run build

# Preview production build
npm run preview
```

---

## 📖 Usage

### Default Admin Account

The application automatically creates a default admin account:

- **Email:** `admin@mail.com`
- **Password:** `adminbookshelf`

### API Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/auth/register` | Register new user | No |
| POST | `/auth/login` | User login | No |
| GET | `/books` | Get all books (paginated) | No |
| GET | `/books/:id` | Get book by ID | No |
| POST | `/books` | Create new book | Yes (Admin) |
| PUT | `/books/:id` | Update book | Yes (Admin) |
| DELETE | `/books/:id` | Delete book | Yes (Admin) |
| POST | `/upload` | Upload book cover | Yes |

### Swagger Documentation

Interactive API documentation is available at:
`http://localhost:8080/swagger/index.html`

---

## 🗄️ Database Schema

### Books Table
- `id` (UUID, Primary Key)
- `title` (String, Required)
- `author` (String, Required)
- `isbn` (String, Unique)
- `published_year` (Integer)
- `genre` (String)
- `description` (Text)
- `cover_image_url` (String)
- `created_at` (Timestamp)
- `updated_at` (Timestamp)

### Users Table
- `id` (UUID, Primary Key)
- `name` (String, Required)
- `email` (String, Unique, Required)
- `password_hash` (String, Required)
- `is_admin` (Boolean, Default: false)
- `created_at` (Timestamp)
- `updated_at` (Timestamp)

---

## 🚀 Deployment

### Backend Deployment

1. **Build the application**
   ```bash
   go build -o bookshelf cmd/server/main.go
   ```

2. **Set production environment variables**
   ```bash
   export GIN_MODE=release
   export DB_HOST=your-production-db-host
   export JWT_SECRET=your-production-jwt-secret
   ```

3. **Run the binary**
   ```bash
   ./bookshelf
   ```

### Frontend Deployment

1. **Build for production**
   ```bash
   npm run build
   ```

2. **Deploy the `dist` folder** to your static hosting service (Netlify, Vercel, AWS S3, etc.)

---

## 🧪 Testing

### Backend Tests
```bash
cd bookshelf-backend
go test ./...
```

### Frontend Tests
```bash
cd bookshelf-frontend
npm run test
```

---

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🐛 Troubleshooting

### Common Issues

**Database Connection Error**
- Ensure PostgreSQL is running
- Check database credentials in environment variables
- Verify database port (5433 for Docker, 5432 for manual setup)

**CORS Errors**
- Ensure frontend is running on `http://localhost:5173`
- Check CORS configuration in backend

**Swagger Not Loading**
- Run `swag init -g cmd/server/main.go` to regenerate docs
- Ensure `docs` package is imported in main.go

**JWT Token Issues**
- Ensure JWT_SECRET is set and consistent
- Check token expiration time
- Verify Authorization header format: `Bearer <token>`

### Getting Help

- Check the [Issues](https://github.com/giovannyptr/bookshelf/issues) page
- Create a new issue with detailed error information
- Include logs, environment details, and steps to reproduce

---

## 🏗️ Project Structure

```
bookshelf/
├── README.md
├── bookshelf-backend/
│   ├── cmd/server/main.go          # Application entry point
│   ├── internal/
│   │   ├── api/response.go         # API response helpers
│   │   ├── auth/                   # Authentication handlers
│   │   ├── books/                  # Book management
│   │   ├── platform/               # Database & CORS setup
│   │   └── users/                  # User management
│   ├── models/                     # Database models
│   ├── docs/                       # Swagger documentation
│   ├── docker-compose.yml          # PostgreSQL container
│   └── go.mod                      # Go dependencies
├── bookshelf-frontend/
│   ├── src/
│   │   ├── components/             # Vue components
│   │   ├── views/                  # Page views
│   │   ├── lib/                    # Utilities (API, auth, etc.)
│   │   └── App.vue                 # Main app component
│   ├── public/                     # Static assets
│   └── package.json                # Node.js dependencies
└── uploads/                        # File upload directory
```
