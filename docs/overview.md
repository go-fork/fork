# Getting Started with Fork HTTP Framework

Hướng dẫn này sẽ giúp bạn bắt đầu với Fork HTTP Framework từ cơ bản đến nâng cao. Framework được thiết kế để dễ sử dụng nhưng vẫn mạnh mẽ và linh hoạt cho các ứng dụng production.

## 📋 Mục lục

1. [Cài đặt](#-cài-đặt)
2. [Hello World](#-hello-world)
3. [Configuration](#️-configuration)
4. [Routing](#-routing)
5. [Middleware](#-middleware)
6. [Context & Data Handling](#-context--data-handling)
7. [Error Handling](#-error-handling)
8. [Dependency Injection](#-dependency-injection)
9. [Template Engine](#-template-engine)
10. [Adapters](#-adapters)
11. [Production Setup](#-production-setup)
12. [Best Practices](#-best-practices)

## 🔧 Cài đặt

### Prerequisites

- Go 1.21 hoặc cao hơn
- Git

### Install Framework

```bash
# Khởi tạo Go module
go mod init your-app

# Cài đặt framework
go get go.fork.vn/fork

# Cài đặt các adapter (optional)
go get github.com/Fork/adapter/fasthttp
go get github.com/Fork/adapter/http2
```

### Verify Installation

```go
package main

import (
    "fmt"
    "go.fork.vn/fork"
)

func main() {
    app := fork.New()
    fmt.Println("Fork HTTP Framework installed successfully!")
}
```

## 👋 Hello World

Tạo ứng dụng đầu tiên với Fork:

### main.go

```go
package main

import (
    "go.fork.vn/fork"
)

func main() {
    // Tạo application instance
    app := fork.New()
    
    // Định nghĩa route đơn giản
    app.Get("/", func(c fork.Context) error {
        return c.String(200, "Hello, Fork!")
    })
    
    // Khởi động server
    app.Listen(":3000")
}
```

### Chạy ứng dụng

```bash
go run main.go
```

Truy cập http://localhost:3000 để xem kết quả.

## ⚙️ Configuration

Fork hỗ trợ cấu hình linh hoạt thông qua YAML files.

### Tạo Configuration File

**configs/app.yaml**

```yaml
# Server configuration
server:
  host: "0.0.0.0"
  port: 8080
  read_timeout: "30s"
  write_timeout: "30s"
  idle_timeout: "120s"

# Graceful shutdown
graceful_shutdown:
  enabled: true
  timeout: "30s"
  wait_time: "5s"

# Adapter configuration
adapter:
  type: "fasthttp"
  config:
    max_request_body_size: 4194304
    concurrency: 1000

# Logging
logging:
  level: "info"
  format: "json"
  output: "stdout"

# Development settings
development:
  auto_reload: true
  debug: true
```

### Sử dụng Configuration

```go
package main

import (
    "log"
    "go.fork.vn/fork"
)

func main() {
    // Load configuration
    config, err := fork.LoadConfigFromFile("configs/app.yaml")
    if err != nil {
        log.Fatal("Failed to load config:", err)
    }
    
    // Create app with config
    app := fork.NewWithConfig(config)
    
    // Routes
    app.Get("/", indexHandler)
    app.Get("/config", configHandler)
    
    // Start server với cấu hình
    app.ListenWithConfig()
}

func indexHandler(c fork.Context) error {
    return c.JSON(200, map[string]string{
        "message": "Hello from configured app!",
        "version": "1.0.0",
    })
}

func configHandler(c fork.Context) error {
    config := c.App().Config()
    return c.JSON(200, map[string]interface{}{
        "server_port": config.Server.Port,
        "adapter_type": config.Adapter.Type,
        "debug_mode": config.Development.Debug,
    })
}
```

## 🚦 Routing

Fork cung cấp router mạnh mẽ với hỗ trợ parameters, wildcards và groups.

### Basic Routes

```go
func setupRoutes(app *fork.App) {
    // HTTP methods
    app.Get("/users", listUsers)           // GET /users
    app.Post("/users", createUser)         // POST /users
    app.Put("/users/:id", updateUser)      // PUT /users/123
    app.Delete("/users/:id", deleteUser)   // DELETE /users/123
    app.Patch("/users/:id", patchUser)     // PATCH /users/123
    
    // Route parameters
    app.Get("/users/:id", getUser)                    // /users/123
    app.Get("/users/:id/posts/:postId", getPost)      // /users/123/posts/456
    
    // Wildcard routes
    app.Get("/files/*", serveFiles)        // /files/path/to/file.txt
    app.Static("/static", "./public")      // Serve static files
}
```

### Route Parameters

```go
func getUser(c fork.Context) error {
    // Get route parameter
    userID := c.Param("id")
    
    // Convert to int
    id, err := strconv.Atoi(userID)
    if err != nil {
        return c.Status(400).JSON(map[string]string{
            "error": "Invalid user ID",
        })
    }
    
    // Get user logic here...
    user := getUserByID(id)
    
    return c.JSON(200, user)
}

func getPost(c fork.Context) error {
    userID := c.Param("id")
    postID := c.Param("postId")
    
    return c.JSON(200, map[string]string{
        "user_id": userID,
        "post_id": postID,
    })
}
```

### Route Groups

```go
func setupAPIRoutes(app *fork.App) {
    // API v1 group
    v1 := app.Group("/api/v1")
    {
        // User routes
        users := v1.Group("/users")
        {
            users.Get("", listUsers)
            users.Post("", createUser)
            users.Get("/:id", getUser)
            users.Put("/:id", updateUser)
            users.Delete("/:id", deleteUser)
        }
        
        // Admin routes với middleware
        admin := v1.Group("/admin", adminAuthMiddleware)
        {
            admin.Get("/stats", getStats)
            admin.Post("/settings", updateSettings)
            admin.Get("/users", adminListUsers)
        }
    }
    
    // API v2 group
    v2 := app.Group("/api/v2")
    {
        v2.Get("/health", healthCheck)
        v2.Post("/webhooks", handleWebhook)
    }
}
```

### Query Parameters

```go
func listUsers(c fork.Context) error {
    // Get query parameters
    page := c.QueryInt("page", 1)           // Default: 1
    limit := c.QueryInt("limit", 10)        // Default: 10
    search := c.Query("search", "")         // Default: ""
    active := c.QueryBool("active", true)   // Default: true
    
    // Validation
    if page < 1 {
        page = 1
    }
    if limit > 100 {
        limit = 100
    }
    
    // Get users with pagination
    users, total := getUsersPaginated(page, limit, search, active)
    
    return c.JSON(200, map[string]interface{}{
        "users": users,
        "pagination": map[string]interface{}{
            "page":  page,
            "limit": limit,
            "total": total,
        },
    })
}
```

## 🛡️ Middleware

Middleware là core component của Fork, cho phép xử lý cross-cutting concerns.

### Built-in Middleware

```go
import (
    "github.com/go-fork/middleware/cors"
    "github.com/go-fork/middleware/logger"
    "github.com/go-fork/middleware/recover"
    "github.com/go-fork/middleware/rate"
)

func setupMiddleware(app *fork.App) {
    // Global middleware (chạy cho tất cả routes)
    app.Use(recover.New())  // Panic recovery
    app.Use(logger.New())   // Request logging
    
    // CORS middleware
    app.Use(cors.New(cors.Config{
        AllowOrigins: []string{"http://localhost:3000"},
        AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
        AllowHeaders: []string{"Content-Type", "Authorization"},
    }))
    
    // Rate limiting
    app.Use(rate.New(rate.Config{
        Max:      100,           // 100 requests
        Duration: time.Minute,   // per minute
    }))
}
```

### Custom Middleware

```go
// Authentication middleware
func authMiddleware(c fork.Context) error {
    // Get token from header
    token := c.Get("Authorization")
    if token == "" {
        return c.Status(401).JSON(map[string]string{
            "error": "Missing authorization token",
        })
    }
    
    // Validate token
    user, err := validateToken(token)
    if err != nil {
        return c.Status(401).JSON(map[string]string{
            "error": "Invalid token",
        })
    }
    
    // Store user in context
    c.Set("user", user)
    
    // Continue to next handler
    return c.Next()
}

// Request timing middleware
func timingMiddleware(c fork.Context) error {
    start := time.Now()
    
    // Process request
    err := c.Next()
    
    // Log duration
    duration := time.Since(start)
    log.Printf("Request %s %s took %v", 
        c.Method(), c.Path(), duration)
    
    return err
}

// Usage
func setupRoutes(app *fork.App) {
    // Apply to specific routes
    app.Get("/protected", authMiddleware, protectedHandler)
    
    // Apply to route group
    api := app.Group("/api", authMiddleware, timingMiddleware)
    {
        api.Get("/profile", getProfile)
        api.Post("/logout", logout)
    }
}
```

## 📦 Context & Data Handling

Context là trung tâm của Fork, cung cấp access tới request/response data.

### Request Data Binding

```go
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name" validate:"required,min=2,max=50"`
    Email string `json:"email" validate:"required,email"`
    Age   int    `json:"age" validate:"min=18,max=120"`
}

func createUser(c fork.Context) error {
    var user User
    
    // Parse JSON body
    if err := c.BodyParser(&user); err != nil {
        return c.Status(400).JSON(map[string]string{
            "error": "Invalid JSON format",
        })
    }
    
    // Validate data
    if err := c.Validate(&user); err != nil {
        return c.Status(422).JSON(map[string]string{
            "error": err.Error(),
        })
    }
    
    // Save user
    savedUser, err := saveUser(&user)
    if err != nil {
        return c.Status(500).JSON(map[string]string{
            "error": "Failed to save user",
        })
    }
    
    return c.Status(201).JSON(savedUser)
}
```

### Response Handling

```go
func getUserResponse(c fork.Context) error {
    userID := c.Param("id")
    
    // Different response types
    switch c.Query("format") {
    case "xml":
        return c.XML(200, user)
    case "yaml":
        return c.YAML(200, user)
    default:
        return c.JSON(200, user)
    }
}

func downloadFile(c fork.Context) error {
    filename := c.Param("filename")
    
    // Set download headers
    c.Set("Content-Disposition", "attachment; filename="+filename)
    c.Set("Content-Type", "application/octet-stream")
    
    // Send file
    return c.SendFile("./uploads/" + filename)
}

func streamData(c fork.Context) error {
    // Set streaming headers
    c.Set("Content-Type", "text/plain")
    c.Set("Transfer-Encoding", "chunked")
    
    // Stream data
    for i := 0; i < 10; i++ {
        data := fmt.Sprintf("Chunk %d\n", i)
        if err := c.WriteString(data); err != nil {
            return err
        }
        time.Sleep(time.Second)
    }
    
    return nil
}
```

## ⚠️ Error Handling

Fork cung cấp một hệ thống error handling mạnh mẽ với `HttpError` type.

### Basic Error Handling

```go
func getUserHandler(c fork.Context) error {
    userID := c.Param("id")
    if userID == "" {
        return errors.BadRequest("User ID is required")
    }
    
    user, err := userService.GetUser(userID)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return errors.NotFound("User not found")
        }
        return errors.InternalServerError("Database error")
    }
    
    return c.JSON(200, user)
}
```

### Structured Error Responses

```go
func createUserHandler(c fork.Context) error {
    var req CreateUserRequest
    if err := c.Bind(&req); err != nil {
        return errors.NewBadRequest("Invalid request", map[string]interface{}{
            "error_type": "validation_failed",
            "details": err.Error(),
        }, err)
    }
    
    if validationErrors := validateUser(req); len(validationErrors) > 0 {
        return errors.NewUnprocessableEntity("Validation failed", map[string]interface{}{
            "validation_errors": validationErrors,
        }, nil)
    }
    
    user, err := userService.CreateUser(req)
    if err != nil {
        if isDuplicateKeyError(err) {
            return errors.NewConflict("User already exists", map[string]interface{}{
                "conflicting_field": "email",
                "value": req.Email,
            }, err)
        }
        return errors.NewInternalServerError("Failed to create user", nil, err)
    }
    
    return c.JSON(201, user)
}
```

### HTTP Status Code Helpers

```go
// 4xx Client Errors
errors.BadRequest("Invalid input")
errors.Unauthorized("Authentication required")
errors.Forbidden("Access denied")
errors.NotFound("Resource not found")
errors.MethodNotAllowed("Method not supported")
errors.Conflict("Resource already exists")
errors.UnprocessableEntity("Validation failed")
errors.TooManyRequests("Rate limit exceeded")

// 5xx Server Errors
errors.InternalServerError("Server error")
errors.NotImplemented("Feature not implemented")
errors.BadGateway("Upstream error")
errors.ServiceUnavailable("Service down")
errors.GatewayTimeout("Upstream timeout")
```

### Error Middleware

```go
func errorHandlerMiddleware() fork.Middleware {
    return func(c fork.Context) error {
        err := c.Next()
        if err != nil {
            var httpErr *errors.HttpError
            if errors.As(err, &httpErr) {
                // Log error with context
                logger.WithFields(map[string]interface{}{
                    "status_code": httpErr.StatusCode,
                    "message": httpErr.Message,
                    "path": c.Path(),
                    "method": c.Method(),
                    "details": httpErr.Details,
                }).Error("HTTP error occurred")
                
                return c.JSON(httpErr.StatusCode, httpErr)
            }
            
            // Handle unexpected errors
            logger.WithError(err).Error("Unhandled error")
            return c.JSON(500, errors.InternalServerError("Internal server error"))
        }
        return nil
    }
}

func main() {
    app := fork.New()
    
    // Global error handling
    app.Use(errorHandlerMiddleware())
    
    app.Get("/users/:id", getUserHandler)
    
    app.Listen(":3000")
}
```

## 🔌 Dependency Injection

Fork tích hợp với DI container để quản lý dependencies.

### Setup DI Container

```go
import "go.fork.vn/di"

type UserService struct {
    repo *UserRepository
}

func (s *UserService) GetUser(id int) (*User, error) {
    return s.repo.FindByID(id)
}

type UserRepository struct {
    db *sql.DB
}

func (r *UserRepository) FindByID(id int) (*User, error) {
    // Database logic
    return &User{}, nil
}

func setupDI() *di.Container {
    container := di.New()
    
    // Register services
    container.Register(&UserRepository{})
    container.Register(&UserService{})
    
    return container
}

func main() {
    container := setupDI()
    
    // Create app with DI
    app := fork.NewWithContainer(container)
    
    app.Get("/users/:id", getUserHandler)
    
    app.Listen(":3000")
}
```

### Use DI in Handlers

```go
func getUserHandler(c fork.Context) error {
    // Resolve service from container
    userService := di.Resolve[*UserService](c.Container())
    
    // Get route parameter
    userID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        return c.Status(400).JSON(map[string]string{
            "error": "Invalid user ID",
        })
    }
    
    // Use service
    user, err := userService.GetUser(userID)
    if err != nil {
        return c.Status(404).JSON(map[string]string{
            "error": "User not found",
        })
    }
    
    return c.JSON(200, user)
}
```

## 🎨 Template Engine

Fork hỗ trợ multiple template engines với tích hợp DI container.

### Basic Template Usage

```go
import "go.fork.vn/templates"

func main() {
    app := fork.New()
    
    // Register template provider
    app.RegisterProvider(&templates.Provider{})
    
    // Configure templates
    app.Get("/", homeHandler)
    
    app.Listen(":3000")
}

func homeHandler(c fork.Context) error {
    return c.Render("home", map[string]interface{}{
        "title": "Welcome to Fork",
        "user":  getCurrentUser(c),
    })
}
```

### Multi-Engine Configuration

```yaml
# configs/templates.yaml
templates:
  default_engine: "html"
  auto_reload: true
  cache_enabled: false
  
  engines:
    html:
      enabled: true
      directory: "./views"
      extension: ".html"
      
    handlebars:
      enabled: true
      directory: "./views/hbs"
      extension: ".hbs"
      
    pug:
      enabled: true
      directory: "./views/pug"
      extension: ".pug"
```

### Advanced Template Features

```go
func setupTemplates(app *fork.App) {
    // Configure with layouts and partials
    templateConfig := templates.Config{
        DefaultEngine: "html",
        AutoReload:    true,
        CacheEnabled:  false,
        
        Engines: map[string]templates.EngineConfig{
            "html": {
                Directory: "./views",
                Extension: ".html",
                Layout:    "layouts/main",
                Partials:  []string{"partials/header", "partials/footer"},
            },
        },
        
        Globals: map[string]interface{}{
            "app_name":    "My App",
            "app_version": "1.0.0",
        },
    }
    
    provider := &templates.Provider{Config: templateConfig}
    app.RegisterProvider(provider)
}

func dashboardHandler(c fork.Context) error {
    data := map[string]interface{}{
        "title": "Dashboard",
        "stats": getDashboardStats(),
        "user":  c.Get("user"),
    }
    
    return c.RenderWithLayout("dashboard", "layouts/admin", data)
}
```

## 🔧 Adapters

Fork hỗ trợ nhiều HTTP engines thông qua adapter pattern.

### FastHTTP Adapter

```go
import "github.com/Fork/adapter/fasthttp"

func main() {
    app := fork.New()
    
    // Use FastHTTP adapter for high performance
    adapter := fasthttp.New(fasthttp.Config{
        MaxRequestBodySize: 4 * 1024 * 1024, // 4MB
        Concurrency:        1000,
        ReadTimeout:        30 * time.Second,
        WriteTimeout:       30 * time.Second,
    })
    
    app.SetAdapter(adapter)
    
    // Routes
    app.Get("/", highPerformanceHandler)
    
    app.Listen(":3000")
}
```

### HTTP/2 Adapter

```go
import "github.com/Fork/adapter/http2"

func main() {
    app := fork.New()
    
    // Use HTTP/2 adapter
    adapter := http2.New(http2.Config{
        TLSConfig: &tls.Config{
            Certificates: []tls.Certificate{cert},
        },
    })
    
    app.SetAdapter(adapter)
    
    app.Get("/", http2Handler)
    
    app.ListenTLS(":443", "cert.pem", "key.pem")
}
```

### Adapter Selection via Config

```yaml
# configs/app.yaml
adapter:
  type: "fasthttp"  # net/http, fasthttp, http2, quic
  config:
    max_request_body_size: 4194304
    concurrency: 1000
    read_timeout: "30s"
    write_timeout: "30s"
```

```go
func main() {
    config, _ := fork.LoadConfigFromFile("configs/app.yaml")
    app := fork.NewWithConfig(config)
    
    // Adapter được tự động setup từ config
    
    app.Get("/", handler)
    app.ListenWithConfig()
}
```

## 🚀 Production Setup

### Graceful Shutdown

```go
func main() {
    app := fork.New()
    
    // Configure graceful shutdown
    app.ConfigureGracefulShutdown(fork.GracefulShutdownConfig{
        Timeout:  30 * time.Second,  // Max time to wait for shutdown
        WaitTime: 5 * time.Second,   // Time to wait before starting shutdown
    })
    
    // Setup routes
    setupRoutes(app)
    
    // Start with graceful shutdown
    app.ListenWithGracefulShutdown(":8080")
}
```

### Health Checks

```go
func setupHealthChecks(app *fork.App) {
    app.Get("/health", healthCheck)
    app.Get("/ready", readinessCheck)
    app.Get("/metrics", metricsHandler)
}

func healthCheck(c fork.Context) error {
    return c.JSON(200, map[string]interface{}{
        "status":    "ok",
        "timestamp": time.Now().Unix(),
        "uptime":    getUptime(),
    })
}

func readinessCheck(c fork.Context) error {
    // Check dependencies (database, external services, etc.)
    if !isDatabaseReady() {
        return c.Status(503).JSON(map[string]string{
            "status": "not ready",
            "reason": "database not available",
        })
    }
    
    return c.JSON(200, map[string]string{
        "status": "ready",
    })
}
```

### Docker Deployment

**Dockerfile**

```dockerfile
# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/configs ./configs

EXPOSE 8080

CMD ["./main"]
```

**docker-compose.yml**

```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - ENV=production
    volumes:
      - ./configs:/root/configs
    depends_on:
      - db
      - redis

  db:
    image: postgres:15
    environment:
      POSTGRES_DB: myapp
      POSTGRES_USER: user
      POSTGRES_PASSWORD: password
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    volumes:
      - redis_data:/data

volumes:
  postgres_data:
  redis_data:
```

## 💡 Best Practices

### 1. Project Structure

```
project/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── handlers/
│   │   ├── user.go
│   │   └── auth.go
│   ├── services/
│   │   ├── user.go
│   │   └── email.go
│   ├── models/
│   │   └── user.go
│   └── middleware/
│       └── auth.go
├── configs/
│   ├── app.yaml
│   ├── app.prod.yaml
│   └── app.dev.yaml
├── migrations/
├── docs/
└── README.md
```

### 2. Error Handling Strategy

```go
type APIError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
    TraceID string `json:"trace_id,omitempty"`
}

func globalErrorHandler(c fork.Context, err error) error {
    traceID := c.Get("trace_id").(string)
    
    var apiErr *APIError
    
    switch {
    case errors.As(err, &apiErr):
        apiErr.TraceID = traceID
        return c.Status(apiErr.Code).JSON(apiErr)
    case errors.Is(err, sql.ErrNoRows):
        return c.Status(404).JSON(APIError{
            Code:    404,
            Message: "Resource not found",
            TraceID: traceID,
        })
    default:
        // Log internal errors
        logger.Error("Internal server error", 
            zap.Error(err), 
            zap.String("trace_id", traceID))
        
        return c.Status(500).JSON(APIError{
            Code:    500,
            Message: "Internal server error",
            TraceID: traceID,
        })
    }
}

// Setup global error handler
func setupErrorHandling(app *fork.App) {
    app.SetErrorHandler(globalErrorHandler)
    
    // Error middleware for logging and monitoring
    app.Use(func(c fork.Context) error {
        // Generate trace ID
        traceID := uuid.New().String()
        c.Set("trace_id", traceID)
        
        // Add to response header
        c.Set("X-Trace-ID", traceID)
        
        return c.Next()
    })
}
```

### 3. Request Validation & Binding

```go
type CreateUserRequest struct {
    Name  string `json:"name" validate:"required,min=2,max=50"`
    Email string `json:"email" validate:"required,email"`
    Age   int    `json:"age" validate:"min=18,max=120"`
}

func (r *CreateUserRequest) Validate() error {
    validate := validator.New()
    return validate.Struct(r)
}

// Handler with proper validation
func createUserHandler(c fork.Context) error {
    var req CreateUserRequest
    
    // Parse request body
    if err := c.BodyParser(&req); err != nil {
        return &APIError{
            Code:    400,
            Message: "Invalid request format",
            Details: err.Error(),
        }
    }
    
    // Validate request
    if err := req.Validate(); err != nil {
        return &APIError{
            Code:    422,
            Message: "Validation failed",
            Details: err.Error(),
        }
    }
    
    // Business logic
    user, err := userService.CreateUser(req)
    if err != nil {
        return err // Will be handled by global error handler
    }
    
    return c.Status(201).JSON(user)
}
```

### 4. Logging & Monitoring

```go
import (
    "go.uber.org/zap"
    "github.com/prometheus/client_golang/prometheus"
)

// Setup structured logging
func setupLogging() *zap.Logger {
    config := zap.NewProductionConfig()
    config.OutputPaths = []string{"stdout", "logs/app.log"}
    
    logger, _ := config.Build()
    return logger
}

// Setup metrics
var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
    
    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "Duration of HTTP requests",
        },
        []string{"method", "endpoint"},
    )
)

func setupMonitoring(app *fork.App) {
    // Register metrics
    prometheus.MustRegister(httpRequestsTotal, httpRequestDuration)
    
    // Metrics middleware
    app.Use(func(c fork.Context) error {
        start := time.Now()
        
        err := c.Next()
        
        duration := time.Since(start).Seconds()
        method := c.Method()
        path := c.Path()
        status := strconv.Itoa(c.Response().StatusCode)
        
        httpRequestsTotal.WithLabelValues(method, path, status).Inc()
        httpRequestDuration.WithLabelValues(method, path).Observe(duration)
        
        return err
    })
    
    // Metrics endpoint
    app.Get("/metrics", func(c fork.Context) error {
        return c.Next() // Prometheus handler
    })
}
```

### 5. Security Best Practices

```go
func setupSecurity(app *fork.App) {
    // Enable security middleware through config
    config := map[string]interface{}{
        "http": map[string]interface{}{
            "middleware": map[string]interface{}{
                // CORS
                "cors": map[string]interface{}{
                    "enabled":      true,
                    "allow_origins": []string{"https://yourdomain.com"},
                    "allow_methods": []string{"GET", "POST", "PUT", "DELETE"},
                    "allow_headers": []string{"*"},
                    "max_age":      86400,
                },
                // Rate limiting
                "limiter": map[string]interface{}{
                    "enabled": true,
                    "max":     100,
                    "duration": "1m",
                },
                // Security headers
                "helmet": map[string]interface{}{
                    "enabled": true,
                    "content_type_nosniff": true,
                    "x_frame_options": "DENY",
                    "hsts_max_age": 31536000,
                },
                // Request ID
                "requestid": map[string]interface{}{
                    "enabled": true,
                    "header":  "X-Request-ID",
                },
            },
        },
    }
    
    app.LoadConfigFromMap(config)
}
```

### 6. Database Integration

```go
import (
    "gorm.io/gorm"
    "gorm.io/driver/postgres"
)

type DatabaseService struct {
    db *gorm.DB
}

func NewDatabaseService(config DatabaseConfig) (*DatabaseService, error) {
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
        config.Host, config.User, config.Password, config.Name, config.Port)
    
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        return nil, err
    }
    
    // Configure connection pool
    sqlDB, _ := db.DB()
    sqlDB.SetMaxIdleConns(config.MaxIdleConns)
    sqlDB.SetMaxOpenConns(config.MaxOpenConns)
    sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)
    
    return &DatabaseService{db: db}, nil
}

// Register as service provider
func (d *DatabaseService) Register(container *fork.Container) error {
    container.Singleton("database", func() interface{} {
        return d.db
    })
    return nil
}

func (d *DatabaseService) Boot() error {
    // Run migrations, setup connections etc.
    return nil
}
```

### 7. Testing Strategy

```go
func TestUserAPI(t *testing.T) {
    // Setup test app
    app := fork.New()
    setupTestRoutes(app)
    
    // Test cases
    tests := []struct {
        name           string
        method         string
        url            string
        body           string
        headers        map[string]string
        expectedStatus int
        expectedBody   string
        setup          func()
        cleanup        func()
    }{
        {
            name:   "Create user success",
            method: "POST",
            url:    "/users",
            body:   `{"name":"John","email":"john@test.com","age":25}`,
            headers: map[string]string{
                "Content-Type":  "application/json",
                "Authorization": "Bearer valid-token",
            },
            expectedStatus: 201,
            expectedBody:   `{"id":1,"name":"John","email":"john@test.com"}`,
            setup: func() {
                // Setup test data
                mockDB.Create(&User{ID: 1, Name: "John"})
            },
            cleanup: func() {
                // Cleanup test data
                mockDB.Delete(&User{}, 1)
            },
        },
        {
            name:           "Create user validation error",
            method:         "POST",
            url:            "/users",
            body:           `{"name":"","email":"invalid","age":15}`,
            expectedStatus: 422,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if tt.setup != nil {
                tt.setup()
            }
            defer func() {
                if tt.cleanup != nil {
                    tt.cleanup()
                }
            }()
            
            req := httptest.NewRequest(tt.method, tt.url, strings.NewReader(tt.body))
            
            // Set headers
            for key, value := range tt.headers {
                req.Header.Set(key, value)
            }
            
            resp, err := app.Test(req)
            require.NoError(t, err)
            assert.Equal(t, tt.expectedStatus, resp.StatusCode)
            
            if tt.expectedBody != "" {
                body, _ := io.ReadAll(resp.Body)
                assert.JSONEq(t, tt.expectedBody, string(body))
            }
        })
    }
}

// Integration tests
func TestIntegration(t *testing.T) {
    // Setup test database
    testDB := setupTestDatabase()
    defer testDB.Close()
    
    // Setup test app with real dependencies
    app := setupAppWithDependencies(testDB)
    
    // Run integration tests
    runAuthenticationTests(t, app)
    runUserManagementTests(t, app)
    runAPITests(t, app)
}
```

### 8. Performance Optimization

```go
func setupPerformanceOptimizations(app *fork.App) {
    // Enable caching middleware
    app.Use(func(c fork.Context) error {
        // Cache static resources
        if strings.HasPrefix(c.Path(), "/static/") {
            c.Set("Cache-Control", "public, max-age=31536000")
        }
        
        // Cache API responses
        if c.Method() == "GET" && shouldCache(c.Path()) {
            cacheKey := generateCacheKey(c)
            if cached := getFromCache(cacheKey); cached != nil {
                return c.JSON(200, cached)
            }
        }
        
        return c.Next()
    })
    
    // Enable compression
    config := map[string]interface{}{
        "http": map[string]interface{}{
            "middleware": map[string]interface{}{
                "compress": map[string]interface{}{
                    "enabled": true,
                    "level":   6,
                    "types": []string{
                        "text/html",
                        "text/css",
                        "text/javascript",
                        "application/json",
                        "application/xml",
                    },
                },
                "etag": map[string]interface{}{
                    "enabled": true,
                    "weak":    false,
                },
            },
        },
    }
    
    app.LoadConfigFromMap(config)
}

func shouldCache(path string) bool {
    cachePaths := []string{"/api/public/", "/api/static/"}
    for _, p := range cachePaths {
        if strings.HasPrefix(path, p) {
            return true
        }
    }
    return false
}
```

### 9. Deployment Automation

**GitHub Actions Workflow**

```yaml
# .github/workflows/deploy.yml
name: Deploy to Production

on:
  push:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: 1.21
      
      - name: Run tests
        run: |
          go mod download
          go test -v ./...
          go test -race -coverprofile=coverage.out ./...
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3

  build:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Build Docker image
        run: |
          docker build -t myapp:${{ github.sha }} .
          docker tag myapp:${{ github.sha }} myapp:latest
      
      - name: Push to registry
        run: |
          echo ${{ secrets.DOCKER_PASSWORD }} | docker login -u ${{ secrets.DOCKER_USERNAME }} --password-stdin
          docker push myapp:${{ github.sha }}
          docker push myapp:latest

  deploy:
    needs: build
    runs-on: ubuntu-latest
    environment: production
    steps:
      - name: Deploy to production
        run: |
          # Update Kubernetes deployment
          kubectl set image deployment/myapp myapp=myapp:${{ github.sha }}
          kubectl rollout status deployment/myapp
```

### 10. Monitoring & Observability

```go
// Distributed tracing with OpenTelemetry
func setupTracing(app *fork.App) {
    tracer := otel.Tracer("myapp")
    
    app.Use(func(c fork.Context) error {
        ctx, span := tracer.Start(c.Context(), c.Method()+" "+c.Path())
        defer span.End()
        
        // Add span context to fork context
        c.SetUserContext(ctx)
        
        // Add trace attributes
        span.SetAttributes(
            attribute.String("http.method", c.Method()),
            attribute.String("http.path", c.Path()),
            attribute.String("http.user_agent", c.Get("User-Agent")),
        )
        
        err := c.Next()
        
        if err != nil {
            span.RecordError(err)
            span.SetStatus(codes.Error, err.Error())
        }
        
        span.SetAttributes(
            attribute.Int("http.status_code", c.Response().StatusCode),
        )
        
        return err
    })
}
```

### 4. Environment Configuration & Secrets Management

```go
type Config struct {
    Environment string         `yaml:"environment"`
    Database    DatabaseConfig `yaml:"database"`
    Redis       RedisConfig    `yaml:"redis"`
    JWT         JWTConfig      `yaml:"jwt"`
    External    ExternalConfig `yaml:"external"`
}

type DatabaseConfig struct {
    Host             string        `yaml:"host"`
    Port             int           `yaml:"port"`
    Name             string        `yaml:"name"`
    User             string        `yaml:"user"`
    Password         string        `yaml:"password"`
    MaxIdleConns     int           `yaml:"max_idle_conns"`
    MaxOpenConns     int           `yaml:"max_open_conns"`
    ConnMaxLifetime  time.Duration `yaml:"conn_max_lifetime"`
}

func LoadConfig() (*Config, error) {
    env := os.Getenv("ENV")
    if env == "" {
        env = "development"
    }
    
    configFile := fmt.Sprintf("configs/app.%s.yaml", env)
    
    config, err := fork.LoadConfigFromFile(configFile)
    if err != nil {
        return nil, err
    }
    
    // Override with environment variables
    if dbPassword := os.Getenv("DB_PASSWORD"); dbPassword != "" {
        config.Database.Password = dbPassword
    }
    if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
        config.JWT.Secret = jwtSecret
    }
    
    return config, nil
}

// Docker Secrets support
func loadDockerSecret(secretName string) (string, error) {
    secretPath := fmt.Sprintf("/run/secrets/%s", secretName)
    if _, err := os.Stat(secretPath); os.IsNotExist(err) {
        return "", nil
    }
    
    data, err := os.ReadFile(secretPath)
    if err != nil {
        return "", err
    }
    
    return strings.TrimSpace(string(data)), nil
}
```

### 5. Comprehensive Testing Strategy

```go
package main

import (
    "testing"
    "net/http/httptest"
    "strings"
    "io"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "go.fork.vn/fork"
)

// Unit tests for handlers
func TestUserHandlers(t *testing.T) {
    app := fork.New()
    setupTestRoutes(app)
    
    tests := []struct {
        name           string
        method         string
        url            string
        body           string
        headers        map[string]string
        expectedStatus int
        expectedBody   string
        setup          func()
        cleanup        func()
    }{
        {
            name:   "Get user success",
            method: "GET",
            url:    "/users/1",
            headers: map[string]string{
                "Authorization": "Bearer valid-token",
            },
            expectedStatus: 200,
            expectedBody:   `{"id":1,"name":"John Doe","email":"john@example.com"}`,
            setup: func() {
                mockUserService.EXPECT().GetUser(1).Return(&User{
                    ID: 1, Name: "John Doe", Email: "john@example.com",
                }, nil)
            },
        },
        {
            name:           "Get user not found",
            method:         "GET",
            url:            "/users/999",
            expectedStatus: 404,
            setup: func() {
                mockUserService.EXPECT().GetUser(999).Return(nil, sql.ErrNoRows)
            },
        },
        {
            name:   "Create user success",
            method: "POST",
            url:    "/users",
            body:   `{"name":"Jane Doe","email":"jane@example.com","age":25}`,
            headers: map[string]string{
                "Content-Type":  "application/json",
                "Authorization": "Bearer valid-token",
            },
            expectedStatus: 201,
            setup: func() {
                mockUserService.EXPECT().CreateUser(gomock.Any()).Return(&User{
                    ID: 2, Name: "Jane Doe", Email: "jane@example.com",
                }, nil)
            },
        },
        {
            name:           "Create user validation error",
            method:         "POST",
            url:            "/users",
            body:           `{"name":"","email":"invalid-email","age":15}`,
            headers:        map[string]string{"Content-Type": "application/json"},
            expectedStatus: 422,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if tt.setup != nil {
                tt.setup()
            }
            
            req := httptest.NewRequest(tt.method, tt.url, strings.NewReader(tt.body))
            
            // Set headers
            for key, value := range tt.headers {
                req.Header.Set(key, value)
            }
            
            resp, err := app.Test(req)
            require.NoError(t, err)
            defer resp.Body.Close()
            
            assert.Equal(t, tt.expectedStatus, resp.StatusCode)
            
            if tt.expectedBody != "" {
                body, _ := io.ReadAll(resp.Body)
                assert.JSONEq(t, tt.expectedBody, string(body))
            }
            
            if tt.cleanup != nil {
                tt.cleanup()
            }
        })
    }
}

// Integration tests
func TestIntegration(t *testing.T) {
    // Setup test database
    testDB := setupTestDatabase(t)
    defer testDB.Close()
    
    // Setup test app with real dependencies
    app := setupAppWithTestDependencies(testDB)
    
    t.Run("User Registration Flow", func(t *testing.T) {
        // Register user
        registerReq := httptest.NewRequest("POST", "/auth/register", 
            strings.NewReader(`{"name":"Test User","email":"test@example.com","password":"password123"}`))
        registerReq.Header.Set("Content-Type", "application/json")
        
        registerResp, err := app.Test(registerReq)
        require.NoError(t, err)
        assert.Equal(t, 201, registerResp.StatusCode)
        
        // Login
        loginReq := httptest.NewRequest("POST", "/auth/login",
            strings.NewReader(`{"email":"test@example.com","password":"password123"}`))
        loginReq.Header.Set("Content-Type", "application/json")
        
        loginResp, err := app.Test(loginReq)
        require.NoError(t, err)
        assert.Equal(t, 200, loginResp.StatusCode)
        
        // Extract token from response
        var loginResult map[string]interface{}
        json.NewDecoder(loginResp.Body).Decode(&loginResult)
        token := loginResult["token"].(string)
        
        // Access protected endpoint
        profileReq := httptest.NewRequest("GET", "/users/profile", nil)
        profileReq.Header.Set("Authorization", "Bearer "+token)
        
        profileResp, err := app.Test(profileReq)
        require.NoError(t, err)
        assert.Equal(t, 200, profileResp.StatusCode)
    })
}

// Performance tests
func BenchmarkAPIEndpoints(b *testing.B) {
    app := fork.New()
    setupRoutes(app)
    
    b.Run("GetUser", func(b *testing.B) {
        req := httptest.NewRequest("GET", "/users/1", nil)
        
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            resp, _ := app.Test(req)
            resp.Body.Close()
        }
    })
    
    b.Run("CreateUser", func(b *testing.B) {
        body := `{"name":"Test","email":"test@example.com","age":25}`
        
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            req := httptest.NewRequest("POST", "/users", strings.NewReader(body))
            req.Header.Set("Content-Type", "application/json")
            
            resp, _ := app.Test(req)
            resp.Body.Close()
        }
    })
}

// Test helpers
func setupTestDatabase(t *testing.T) *sql.DB {
    db, err := sql.Open("sqlite3", ":memory:")
    require.NoError(t, err)
    
    // Run test migrations
    runTestMigrations(db)
    
    return db
}

func setupAppWithTestDependencies(db *sql.DB) *fork.App {
    app := fork.New()
    
    // Register test services
    container := app.Container()
    container.Singleton("database", func() interface{} { return db })
    container.Singleton("user_service", func() interface{} { 
        return NewUserService(db) 
    })
    
    setupRoutes(app)
    return app
}
```

## 🎯 Next Steps

Sau khi đã nắm được các kiến thức cơ bản, đây là các bước tiếp theo để thành thạo Fork Framework:

### 1. Khám phá Documentation Chi tiết

- **[Configuration Guide](config.md)** - Quản lý cấu hình YAML và environment
- **[Service Provider & DI](service-provider.md)** - Dependency injection và service registration
- **[Web Application](web-application.md)** - Core WebApp API và lifecycle
- **[Context, Request & Response](context-request-response.md)** - Data binding và response handling
- **[Router](router.md)** - Advanced routing với trie structure
- **[Adapters](adapter.md)** - HTTP adapters cho performance tối ưu
- **[Middleware](middleware.md)** - 30+ middleware packages với YAML config
- **[Error Handling](error-handling.md)** - Comprehensive error management system

### 2. Thực hành với Examples

```bash
# Clone examples repository
git clone github.com/go-fork/examples
cd examples

# Basic web application
cd basic-web-app
go run main.go

# REST API with database
cd rest-api
go run main.go

# Microservice with gRPC
cd microservice
go run main.go

# Full-stack application
cd fullstack-app
go run main.go
```

### 3. Middleware Ecosystem

Khám phá hơn 30 middleware packages:

```yaml
# configs/app.yaml
http:
  middleware:
    # Security
    cors.enabled: true
    helmet.enabled: true
    csrf.enabled: true
    
    # Performance  
    compress.enabled: true
    etag.enabled: true
    cache.enabled: true
    
    # Monitoring
    logger.enabled: true
    recover.enabled: true
    metrics.enabled: true
    
    # Rate limiting
    limiter.enabled: true
    timeout.enabled: true
```

### 4. Template Engines

Tích hợp với multiple template engines:

```go
// Register multiple template engines
app.RegisterTemplateEngine("html", html.New())
app.RegisterTemplateEngine("pug", pug.New())
app.RegisterTemplateEngine("handlebars", handlebars.New())

// Use in handlers
func indexHandler(c fork.Context) error {
    return c.Render("index.html", map[string]interface{}{
        "title": "Fork Framework",
        "users": userService.GetUsers(),
    })
}
```

### 5. Production Deployment

```bash
# Build for production
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Docker deployment
docker build -t myapp .
docker run -p 8080:8080 myapp

# Kubernetes deployment
kubectl apply -f k8s/
kubectl get pods
```

### 6. Performance Tuning

Tối ưu performance cho production:

- **Adapter Selection**: Chọn FastHTTP cho high-performance applications
- **Middleware Optimization**: Chỉ enable middleware cần thiết
- **Connection Pooling**: Cấu hình database connection pool
- **Caching Strategy**: Implement Redis caching cho API responses
- **Monitoring**: Setup Prometheus metrics và distributed tracing

### 7. Community & Resources

- **GitHub**: [github.com/go-fork/http](github.com/go-fork/http)
- **Documentation**: [https://docs.Fork.vn](https://docs.Fork.vn)
- **Examples**: [github.com/go-fork/examples](github.com/go-fork/examples)
- **Discord**: [https://discord.gg/Fork](https://discord.gg/Fork)
- **Blog**: [https://blog.Fork.vn](https://blog.Fork.vn)

## 📚 API Reference Complete

### Core Components

| Component | Description | Documentation |
|-----------|-------------|---------------|
| **App** | Main application instance | [web-application.md](web-application.md) |
| **Context** | Request/Response context | [context-request-response.md](context-request-response.md) |
| **Router** | URL routing system | [router.md](router.md) |
| **Adapter** | HTTP server adapters | [adapter.md](adapter.md) |
| **Config** | Configuration management | [config.md](config.md) |
| **ServiceProvider** | Dependency injection | [service-provider.md](service-provider.md) |
| **Middleware** | Request processing pipeline | [middleware.md](middleware.md) |
| **Errors** | Error handling system | [error-handling.md](error-handling.md) |

### Framework Features

- ✅ **High Performance**: Multiple HTTP adapters (net/http, FastHTTP, HTTP/2)
- ✅ **YAML Configuration**: Comprehensive config management
- ✅ **Dependency Injection**: Built-in DI container
- ✅ **30+ Middleware**: Auto-configured middleware ecosystem
- ✅ **Template Engines**: Multiple template engine support
- ✅ **Error Handling**: Structured error management
- ✅ **Testing**: Comprehensive testing utilities
- ✅ **Production Ready**: Graceful shutdown, health checks, monitoring
- ✅ **Documentation**: Complete Vietnamese documentation

### Quick Reference Commands

```bash
# Create new project
go mod init myapp
go get go.fork.vn/fork

# Run development server
go run main.go

# Run tests
go test -v ./...

# Build for production
go build -o myapp

# Run with config
ENV=production ./myapp

# Docker build
docker build -t myapp .

# Deploy to production
kubectl apply -f deployment.yaml
```

---

Happy coding với Fork! 🚀
