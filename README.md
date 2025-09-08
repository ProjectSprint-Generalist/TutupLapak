# TutupLapak API

A RESTful API built with Go and Gin framework, following clean architecture principles and best practices.

## Features

- 🚀 **Gin Framework** - Fast HTTP web framework
- 📝 **Structured Logging** - Using zerolog for efficient logging
- 🔒 **CORS Support** - Cross-origin resource sharing
- 🏗️ **Clean Architecture** - Organized project structure
- 📊 **Health Checks** - Built-in health and readiness endpoints
- 🔧 **Environment Configuration** - Easy configuration management
- 📋 **Standard REST API** - Following REST conventions

## Project Structure

```
tutuplapak/
├── main.go                 # Application entry point
├── go.mod                  # Go module file
├── .env.example           # Environment variables template
├── README.md              # This file
└── internal/              # Private application code
    ├── config/            # Configuration management
    │   └── config.go
    ├── handlers/          # HTTP request handlers
    │   ├── health.go
    │   └── user.go
    ├── middleware/        # HTTP middleware
    │   ├── cors.go
    │   ├── logger.go
    │   └── recovery.go
    ├── models/            # Data models
    │   ├── response.go
    │   └── user.go
    └── routes/            # Route definitions
        └── routes.go
```

## Getting Started

### Prerequisites

- Go 1.25.0 or higher
- Git

### Installation

1. **Clone the repository**
   ```bash
   git clone <your-repo-url>
   cd tutuplapak
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Run the application**
   ```bash
   go run main.go
   ```

The API will be available at `http://localhost:8080`

## API Endpoints

### Health Check
- `GET /api/v1/health/` - Health status
- `GET /api/v1/health/ready` - Readiness check

### Users
- `GET /api/v1/users/` - Get all users (with pagination)
- `GET /api/v1/users/:id` - Get user by ID
- `POST /api/v1/users/` - Create new user
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user

#### User Model
```json
{
  "id": 1,
  "email": "name@name.com",
  "name": "John Doe",
  "preference": "metric",
  "weightUnit": "kg",
  "heightUnit": "cm", 
  "weight": 75.5,
  "height": 180.0,
  "imageUri": "https://example.com/image.jpg"
}
```

**Note:** All fields except `id` and `email` can be `null` when empty.

### Root
- `GET /` - API information

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `ENVIRONMENT` | Application environment | `development` |
| `PORT` | Server port | `8080` |
| `DATABASE_URL` | Database connection string | - |
| `JWT_SECRET` | JWT signing secret | `your-secret-key` |
| `CORS_ALLOWED_ORIGINS` | Allowed CORS origins | `*` |

## Development

### Adding New Endpoints

1. **Create a new handler** in `internal/handlers/`
2. **Define models** in `internal/models/` if needed
3. **Add routes** in `internal/routes/routes.go`
4. **Update main.go** to initialize the new handler

### Example: Adding a Product Handler

```go
// internal/handlers/product.go
type ProductHandler struct{}

func (h *ProductHandler) GetProducts(c *gin.Context) {
    // Implementation
}
```

```go
// internal/routes/routes.go
products := v1.Group("/products")
{
    products.GET("/", productHandler.GetProducts)
}
```

## Building for Production

```bash
# Build the application
go build -o tutuplapak main.go

# Run the binary
./tutuplapak
```

## Docker Support (Optional)

Create a `Dockerfile`:

```dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o tutuplapak main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/fitbyte .
CMD ["./tutuplapak"]
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License.

## Next Steps

- [ ] Add CI/CD pipeline
