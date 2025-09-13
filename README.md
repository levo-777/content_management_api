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

### Docker Commands Reference

```bash
# Start services in background
docker-compose up -d

# View running containers
docker-compose ps

# View logs
docker-compose logs -f [service-name]

# Stop services
docker-compose down

# Stop and remove volumes (clean slate)
docker-compose down -v

# Rebuild and start services
docker-compose up --build

# Execute commands in containers
docker-compose exec cms-backend sh
docker-compose exec postgres psql -U postgres -d cms_db
```

---

## 🚀 Production Deployment

### Building and Distributing the Docker Image

#### 1. Build the Production Image

**Option A: Using Docker Registry (Internet Required)**
```bash
# Build the image
docker build -t cms-backend:latest .

# Tag for registry (replace with your registry)
docker tag cms-backend:latest your-registry.com/cms-backend:v1.0.0

# Push to registry
docker push your-registry.com/cms-backend:v1.0.0
```

**Option B: Using Local Image File (No Internet Required)**

Create the image file:
```bash
# Build the image first
docker build -t cms-backend .

# Save the image to a tar file
docker save cms-backend:latest -o cms-backend-docker-image.tar

# Compress the tar file (optional, reduces size)
gzip cms-backend-docker-image.tar
```

#### 2. Deploy on Another Machine

**Option A: Using Docker Registry (Internet Required)**
```bash
# On the target machine
git clone <repository-url>
cd content_management_system_api

# Update docker-compose.yml to use your image
# Change: build: . → image: your-registry.com/cms-backend:v1.0.0

# Start services
docker-compose up -d
```

**Option B: Using Local Image File (No Internet Required)**

Load and run on another machine:
```bash
# Load the pre-built image (copy the .tar file to the target machine first)
docker load -i cms-backend-docker-image.tar

# Start PostgreSQL database
docker run -d --name postgres-cms \
  -e POSTGRES_DB=cms_db \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=your_secure_password \
  -p 5432:5432 \
  -v postgres_data:/var/lib/postgresql/data \
  postgres:15-alpine

# Run the CMS backend container
docker run -d --name cms-backend \
  --link postgres-cms:postgres \
  -e ENV=production \
  -e DB_HOST=postgres \
  -e DB_PORT=5432 \
  -e DB_USER=postgres \
  -e DB_PASSWORD=your_secure_password \
  -e DB_NAME=cms_db \
  -p 8080:8080 \
  --restart unless-stopped \
  cms-backend:latest
```

**Option C: Manual Docker Run with Registry**
```bash
# Start PostgreSQL
docker run -d --name postgres-cms \
  -e POSTGRES_DB=cms_db \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=your_secure_password \
  -p 5432:5432 \
  -v postgres_data:/var/lib/postgresql/data \
  postgres:15-alpine

# Start CMS Backend
docker run -d --name cms-backend \
  --link postgres-cms:postgres \
  -e ENV=production \
  -e DB_HOST=postgres \
  -e DB_PORT=5432 \
  -e DB_USER=postgres \
  -e DB_PASSWORD=your_secure_password \
  -e DB_NAME=cms_db \
  -p 8080:8080 \
  --restart unless-stopped \
  your-registry.com/cms-backend:v1.0.0
```

#### 3. Complete Deployment Package (Recommended for Offline Deployment)

Create a complete deployment package that includes everything needed:

```bash
# Build the image
docker build -t cms-backend .

# Save image to tar file
docker save cms-backend:latest -o cms-backend-docker-image.tar

# Compress the image
gzip cms-backend-docker-image.tar

# Create deployment directory
mkdir cms-deployment
cd cms-deployment

# Copy necessary files
cp ../cms-backend-docker-image.tar.gz .
cp ../docker-compose.yml .
cp ../.env.example .

# Create production docker-compose override
cat > docker-compose.prod.yml << 'EOF'
version: '3.8'
services:
  cms-backend:
    image: cms-backend:latest
    environment:
      ENV: production
EOF

# Create deployment script
cat > deploy.sh << 'EOF'
#!/bin/bash
echo "Loading Docker image..."
docker load -i cms-backend-docker-image.tar.gz

echo "Starting services..."
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

echo "Deployment complete! API available at http://localhost:8080"
EOF

chmod +x deploy.sh

# Create deployment package
cd ..
tar -czf cms-deployment.tar.gz cms-deployment/
```

**Deploy on target machine:**
```bash
# Transfer deployment package
scp cms-deployment.tar.gz user@target-machine:/path/to/destination/

# On target machine
tar -xzf cms-deployment.tar.gz
cd cms-deployment

# Edit environment variables
cp .env.example .env
nano .env

# Run deployment
./deploy.sh
```

#### 4. Production Environment Variables
```env
ENV=production
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_secure_password
DB_NAME=cms_db
```

---

## 🧪 Testing

### Run All Tests
```bash
make test
```

### Run Unit Tests Only
```bash
make test-unit
```

### Run Integration Tests
```bash
make test-integration
```

### Run Integration Tests with Database Setup
```bash
make test-integration-full
```

---

## ⚙️ Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `ENV` | Environment mode (development/production) | development | No |
| `DB_HOST` | Database host | localhost | Yes |
| `DB_PORT` | Database port | 5432 | Yes |
| `DB_USER` | Database username | postgres | Yes |
| `DB_PASSWORD` | Database password | postgres | Yes |
| `DB_NAME` | Database name | cms_db | Yes |

---

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

## 📁 Project Structure

```
content_management_system_api/
├── controllers/          # HTTP request handlers
│   ├── page_controller.go
│   ├── post_controller.go
│   └── media_controller.go
├── models/              # Data models
│   ├── page.go
│   ├── post.go
│   └── media.go
├── routes/              # API route definitions
│   └── routes.go
├── utils/               # Utility functions
│   └── db.go
├── migrations/          # Database migrations
├── tests/               # Test files
│   └── integration/
├── main.go             # Application entry point
├── go.mod              # Go module file (Go 1.22)
├── go.sum              # Go module checksums
├── makefile            # Build and test commands
├── Dockerfile          # Docker configuration
├── docker-compose.yml  # Docker Compose configuration
├── .env.example        # Environment template
└── README.md           # This file
```

---

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality
5. Run the test suite (`make test`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Submit a pull request

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🆘 Support

If you encounter any issues or have questions:

1. Check the [Issues](https://github.com/your-repo/issues) page
2. Create a new issue with detailed information
3. Include your Go version, OS, and error logs

---

**Made with ❤️ using Go 1.22**
