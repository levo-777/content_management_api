# Content Management API

A **minimalistic** and lightweight content management system API built with Go, Gin, and PostgreSQL. Perfect for small to medium projects that need a simple, fast, and reliable content management solution.

## 🚀 Quick Overview

This API provides a clean and efficient way to manage content through RESTful endpoints. It supports three main content types:

- **Pages**: Static content pages
- **Posts**: Dynamic blog-style posts with author information  
- **Media**: File attachments and media resources

## ✨ Features

- **Minimalistic Design**: Simple, clean API with no bloat
- **RESTful API**: Standard HTTP methods and status codes
- **PostgreSQL Integration**: Robust database with ACID compliance
- **GORM ORM**: Type-safe database operations
- **Auto Migrations**: Database schema management
- **Environment Configuration**: Flexible deployment options
- **Comprehensive Testing**: Unit and integration tests
- **Docker Ready**: Containerized deployment
- **Health Checks**: Built-in monitoring endpoints

## 🛠 Tech Stack

- **Backend**: Go 1.22
- **Web Framework**: Gin
- **Database**: PostgreSQL 15
- **ORM**: GORM
- **Testing**: Go testing framework with SQLMock
- **Containerization**: Docker & Docker Compose

## API Endpoints

### Pages
- `GET /api/v1/pages` - Get all pages
- `GET /api/v1/pages/:id` - Get page by ID
- `POST /api/v1/pages` - Create new page
- `PUT /api/v1/pages/:id` - Update page
- `DELETE /api/v1/pages/:id` - Delete page

### Posts
- `GET /api/v1/posts` - Get all posts
- `GET /api/v1/posts/:id` - Get post by ID
- `POST /api/v1/posts` - Create new post
- `PUT /api/v1/posts/:id` - Update post
- `DELETE /api/v1/posts/:id` - Delete post

### Media
- `GET /api/v1/media` - Get all media
- `GET /api/v1/media/:id` - Get media by ID
- `POST /api/v1/media` - Create new media
- `DELETE /api/v1/media/:id` - Delete media

## Data Models

### Page
```json
{
  "id": 1,
  "title": "Page Title",
  "content": "Page content...",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### Post
```json
{
  "id": 1,
  "title": "Post Title",
  "content": "Post content...",
  "author": "Author Name",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "media": []
}
```

### Media
```json
{
  "id": 1,
  "url": "https://example.com/image.jpg",
  "type": "image",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

## 📋 Prerequisites

### For Local Development (without Docker)
- **Go 1.22** or higher
- **PostgreSQL 15** or higher
- Git

### For Docker Deployment
- **Docker** 20.10+
- **Docker Compose** 2.0+

## 🚀 Getting Started

Choose your preferred setup method:

- [**Local Development**](#-local-development-setup) - Run directly on your machine
- [**Docker Development**](#-docker-development-setup) - Run with Docker Compose
- [**Production Deployment**](#-production-deployment) - Deploy to other machines

---

## 💻 Local Development Setup

### 1. Clone the Repository
```bash
git clone <repository-url>
cd content_management_system_api
```

### 2. Install Go Dependencies
```bash
go mod download
```

### 3. Environment Configuration
```bash
# Copy environment template
cp .env.example .env

# Edit with your database settings
nano .env
```

**Required environment variables:**
```env
ENV=development
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=cms_db
```

### 4. Database Setup
```bash
# Start PostgreSQL service
sudo systemctl start postgresql

# Create database
psql -U postgres -c "CREATE DATABASE cms_db;"
```

### 5. Run the Application
```bash
# Start the API server
go run main.go
```

✅ **API available at:** `http://localhost:8080`

---

## 🐳 Docker Development Setup

### Quick Start (Recommended)

```bash
# Clone the repository
git clone <repository-url>
cd content_management_system_api

# Start all services with Docker Compose
docker-compose up -d

# Check service status
docker-compose ps

# View logs
docker-compose logs -f
```

✅ **API available at:** `http://localhost:8080`

## 📖 API Usage Examples

### Create a Page
```bash
curl -X POST http://localhost:8080/api/v1/pages \
  -H "Content-Type: application/json" \
  -d '{
    "title": "About Us",
    "content": "This is our about page content."
  }'
```

### Create a Post
```bash
curl -X POST http://localhost:8080/api/v1/posts \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Welcome to Our Blog",
    "content": "This is our first blog post!",
    "author": "John Doe"
  }'
```

### Get All Pages
```bash
curl http://localhost:8080/api/v1/pages
```

### Get All Posts
```bash
curl http://localhost:8080/api/v1/posts
```

### Create Media
```bash
curl -X POST http://localhost:8080/api/v1/media \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://example.com/image.jpg",
    "type": "image"
  }'
```

---
