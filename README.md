# Fork HTTP Framework

Fork HTTP Framework là một framework HTTP linh hoạt và có thể mở rộng cao cho Go, được thiết kế để hỗ trợ nhiều HTTP engine khác nhau thông qua adapter pattern. Framework này cung cấp một API thống nhất để xây dựng RESTful APIs, web services và ứng dụng web hiện đại.

## 🚀 Tính năng chính

- **Multi-Engine Support**: Hỗ trợ nhiều HTTP engine (net/http, fasthttp, http2, quic-h3)
- **Adapter Pattern**: Kiến trúc linh hoạt cho phép chuyển đổi giữa các engine dễ dàng
- **Powerful Router**: Router mạnh mẽ với hỗ trợ parameters, wildcards và trie structure
- **Middleware System**: Hệ thống middleware linh hoạt với middleware groups
- **Dependency Injection**: Tích hợp sẵn với DI container
- **Configuration Management**: Hệ thống cấu hình YAML linh hoạt
- **Context System**: Context mạnh mẽ với data binding và validation
- **Graceful Shutdown**: Hỗ trợ graceful shutdown cho production
- **Template Integration**: Tích hợp với nhiều template engine
- **Production Ready**: Tối ưu hiệu năng và memory usage

## 📦 Cài đặt

```bash
go get go.fork.vn/fork
```

## 🏃 Quick Start

### Basic Application

```go
package main

import (
    "go.fork.vn/fork"
    "go.fork.vn/fork/adapter"
)

func main() {
    // Tạo web application
    app := fork.New()
    
    // Định nghĩa routes
    app.Get("/", func(c fork.Context) error {
        return c.String(200, "Hello, World!")
    })
    
    app.Get("/user/:id", func(c fork.Context) error {
        id := c.Param("id")
        return c.JSON(200, map[string]string{
            "id": id,
            "message": "User found",
        })
    })
    
    // Khởi động server
    app.Listen(":3000")
}
```

### Với Configuration

```go
package main

import (
    "go.fork.vn/fork"
    "go.fork.vn/fork/adapter"
)

func main() {
    // Load configuration từ file
    config, err := fork.LoadConfigFromFile("configs/app.yaml")
    if err != nil {
        panic(err)
    }
    
    // Tạo application với config
    app := fork.NewWithConfig(config)
    
    // Sử dụng fasthttp adapter
    adapter := adapter.NewFastHTTPAdapter()
    app.SetAdapter(adapter)
    
    // Định nghĩa routes
    app.Get("/api/health", healthHandler)
    app.Post("/api/users", createUserHandler)
    
    // Khởi động với graceful shutdown
    app.ListenWithGracefulShutdown(":8080")
}

func healthHandler(c fork.Context) error {
    return c.JSON(200, map[string]string{
        "status": "ok",
        "time": time.Now().Format(time.RFC3339),
    })
}

func createUserHandler(c fork.Context) error {
    var user User
    if err := c.BodyParser(&user); err != nil {
        return c.Status(400).JSON(map[string]string{
            "error": "Invalid request body",
        })
    }
    
    // Validate user data
    if err := c.Validate(&user); err != nil {
        return c.Status(422).JSON(map[string]string{
            "error": err.Error(),
        })
    }
    
    // Save user logic here...
    
    return c.Status(201).JSON(user)
}
```

## 📚 Tài liệu

### Core Components

- **[Configuration](docs/config.md)** - Hệ thống cấu hình và YAML management
- **[Service Provider](docs/service-provider.md)** - Dependency Injection và service management
- **[Web Application](docs/web-application.md)** - Core WebApp object và lifecycle
- **[Context, Request & Response](docs/context-request-response.md)** - HTTP context system
- **[Router](docs/router.md)** - Routing system và middleware
- **[Adapter](docs/adapter.md)** - Adapter pattern và HTTP engines

### Quick Links

- [Getting Started Guide](docs/overview.md)
- [API Reference](docs/)
- [Examples](../examples/)
- [Middleware](../middleware/)
- [Templates](../templates/)

## 🏗️ Kiến trúc

```
┌─────────────────┐
│   Application   │
├─────────────────┤
│   Middleware    │
├─────────────────┤
│     Router      │
├─────────────────┤
│    Context      │
├─────────────────┤
│    Adapter      │
├─────────────────┤
│  HTTP Engine    │
│ (net/http,      │
│  fasthttp,      │
│  http2, quic)   │
└─────────────────┘
```

### Adapter Pattern

Framework sử dụng adapter pattern để hỗ trợ nhiều HTTP engine:

- **net/http**: Standard Go HTTP server
- **fasthttp**: High-performance HTTP server
- **http2**: HTTP/2 support
- **quic**: HTTP/3 over QUIC

## 🔧 Configuration

Framework hỗ trợ cấu hình thông qua YAML file:

```yaml
# configs/app.yaml
server:
  host: "0.0.0.0"
  port: 8080
  read_timeout: "30s"
  write_timeout: "30s"
  idle_timeout: "120s"

graceful_shutdown:
  enabled: true
  timeout: "30s"
  wait_time: "5s"

adapter:
  type: "fasthttp"
  config:
    max_request_body_size: 4194304
    concurrency: 1000
```

## 🛠️ Middleware

Framework cung cấp nhiều middleware có sẵn:

```go
import (
    "github.com/go-fork/middleware/cors"
    "github.com/go-fork/middleware/logger"
    "github.com/go-fork/middleware/recover"
)

app := fork.New()

// Global middleware
app.Use(recover.New())
app.Use(logger.New())
app.Use(cors.New())

// Route-specific middleware
app.Get("/api/admin/*", adminAuth, adminHandler)
```

## 🎯 Context & Data Binding

```go
type User struct {
    ID    int    `json:"id" validate:"required"`
    Name  string `json:"name" validate:"required,min=2,max=50"`
    Email string `json:"email" validate:"required,email"`
}

func createUser(c fork.Context) error {
    var user User
    
    // Parse request body
    if err := c.BodyParser(&user); err != nil {
        return c.Status(400).JSON(ErrorResponse{
            Error: "Invalid JSON",
        })
    }
    
    // Validate data
    if err := c.Validate(&user); err != nil {
        return c.Status(422).JSON(ErrorResponse{
            Error: err.Error(),
        })
    }
    
    // Get route parameters
    userID := c.Param("id")
    
    // Get query parameters
    filter := c.Query("filter", "all")
    
    // Set response headers
    c.Set("X-User-ID", strconv.Itoa(user.ID))
    
    return c.JSON(201, user)
}
```

## 🚦 Router Features

```go
app := fork.New()

// Basic routes
app.Get("/", homeHandler)
app.Post("/users", createUserHandler)
app.Put("/users/:id", updateUserHandler)
app.Delete("/users/:id", deleteUserHandler)

// Route parameters
app.Get("/users/:id", getUserHandler)
app.Get("/users/:id/posts/:postId", getPostHandler)

// Wildcard routes
app.Get("/files/*", fileHandler)

// Route groups
api := app.Group("/api/v1")
{
    api.Get("/users", listUsersHandler)
    api.Post("/users", createUserHandler)
    
    admin := api.Group("/admin", adminMiddleware)
    {
        admin.Get("/stats", statsHandler)
        admin.Post("/settings", updateSettingsHandler)
    }
}
```

## 🔌 Dependency Injection

```go
import "go.fork.vn/di"

// Register services
container := di.NewContainer()
container.Register(&UserService{})
container.Register(&EmailService{})

// Create app with DI
app := fork.NewWithContainer(container)

func getUserHandler(c fork.Context) error {
    // Resolve service from container
    userService := di.Resolve[*UserService](c.Container())
    
    users, err := userService.GetAll()
    if err != nil {
        return c.Status(500).JSON(ErrorResponse{
            Error: "Failed to get users",
        })
    }
    
    return c.JSON(200, users)
}
```

## 📊 Performance

Framework được tối ưu cho hiệu năng cao:

- **Zero-allocation routing** với trie structure
- **Memory pooling** cho objects có thể tái sử dụng
- **Efficient middleware chain** với minimal overhead
- **Adapter-based engine selection** cho performance tuning

### Benchmarks

```
BenchmarkRouter-8           5000000    240 ns/op     0 allocs/op
BenchmarkContext-8          3000000    450 ns/op     1 allocs/op
BenchmarkMiddleware-8       2000000    680 ns/op     2 allocs/op
```

## 🔒 Security Features

- **CORS middleware** với configurable options
- **CSRF protection** middleware
- **Rate limiting** middleware
- **Input validation** với struct tags
- **Secure headers** middleware
- **Authentication middleware** support

## 🧪 Testing

Framework cung cấp testing utilities:

```go
func TestUserAPI(t *testing.T) {
    app := fork.New()
    app.Post("/users", createUserHandler)
    
    // Test request
    req := httptest.NewRequest("POST", "/users", strings.NewReader(`{"name":"John"}`))
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := app.Test(req)
    assert.NoError(t, err)
    assert.Equal(t, 201, resp.StatusCode)
}
```

## 🚀 Production Deployment

### Docker

```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/configs ./configs
CMD ["./main"]
```

### Graceful Shutdown

```go
app := fork.New()

// Configure graceful shutdown
app.ConfigureGracefulShutdown(fork.GracefulShutdownConfig{
    Timeout:  30 * time.Second,
    WaitTime: 5 * time.Second,
})

// Start with graceful shutdown
app.ListenWithGracefulShutdown(":8080")
```

## 📝 Examples

Xem thêm examples trong thư mục [examples/](../examples/):

- [Basic HTTP Server](../examples/http/simple-config-example/)
- [FastHTTP Example](../examples/adapter/fasthttp_example/)
- [HTTP/2 Example](../examples/adapter/http2_sample/)
- [Graceful Shutdown](../examples/http/graceful-shutdown-example/)
- [Configuration Example](../examples/http/config-provider-example/)

## 🤝 Contributing

Chúng tôi hoan nghênh các contributions! Vui lòng:

1. Fork repository
2. Tạo feature branch
3. Commit changes
4. Push branch
5. Tạo Pull Request

## 📄 License

MIT License - xem file [LICENSE](LICENSE) để biết thêm chi tiết.

## 🔗 Links

- [Documentation](docs/)
- [Examples](../examples/)
- [Middleware](../middleware/)
- [Templates](../templates/)
- [Issues](https://github.com/go-fork/http/issues)

---

**Fork HTTP Framework** - Build fast, scalable web applications in Go 🚀
