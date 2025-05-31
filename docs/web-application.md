# Web Application - Ứng dụng Web

Package `fork` cung cấp `WebApp` - đối tượng core của Fork HTTP Framework. WebApp hoạt động như một wrapper cho router và adapter, cung cấp API đồng nhất để xử lý HTTP requests, đăng ký routes và middlewares.

## Tổng quan

`WebApp` là lớp trung tâm quản lý toàn bộ HTTP lifecycle của ứng dụng. Nó khác biệt với DI Application - WebApp chỉ tập trung vào web layer và cung cấp:

- Route registration và management
- Middleware pipeline execution
- HTTP adapter abstraction
- Graceful shutdown handling
- Request/Response context management

## Architecture

```go
type WebApp struct {
    adapter         adapter.Adapter
    router          router.Router
    middlewares     []router.HandlerFunc
    config          *WebAppConfig
    mu              sync.RWMutex
    shutdownCtx     context.Context
    shutdownCancel  context.CancelFunc
    activeConnections int32
    isShuttingDown  bool
}
```

### Core Components

- **Adapter**: HTTP server implementation (net/http, fasthttp, etc.)
- **Router**: Route management và request dispatching
- **Middlewares**: Global middleware chain
- **Config**: Application configuration
- **Shutdown Management**: Graceful shutdown coordination

## API Reference

### Constructor

#### NewWebApp()

Tạo một instance mới của WebApp:

```go
func NewWebApp() *WebApp
```

**Returns:** WebApp instance với router và middlewares được khởi tạo

**Usage:**
```go
app := fork.NewWebApp()
```

### Configuration Management

#### SetConfig()

Thiết lập configuration cho WebApp:

```go
func (app *WebApp) SetConfig(config *WebAppConfig)
```

#### GetConfig()

Lấy configuration hiện tại:

```go
func (app *WebApp) GetConfig() *WebAppConfig
```

### Adapter Management

#### SetAdapter()

Thiết lập HTTP adapter:

```go
func (app *WebApp) SetAdapter(adapter adapter.Adapter)
```

#### GetAdapter()

Lấy adapter hiện tại:

```go
func (app *WebApp) GetAdapter() adapter.Adapter
```

### Route Registration

WebApp cung cấp các phương thức convenient cho việc đăng ký routes:

#### HTTP Methods

```go
func (app *WebApp) GET(path string, handlers ...router.HandlerFunc)
func (app *WebApp) POST(path string, handlers ...router.HandlerFunc)
func (app *WebApp) PUT(path string, handlers ...router.HandlerFunc)
func (app *WebApp) DELETE(path string, handlers ...router.HandlerFunc)
func (app *WebApp) PATCH(path string, handlers ...router.HandlerFunc)
func (app *WebApp) HEAD(path string, handlers ...router.HandlerFunc)
func (app *WebApp) OPTIONS(path string, handlers ...router.HandlerFunc)
```

#### Universal Registration

```go
func (app *WebApp) Any(path string, handlers ...router.HandlerFunc)
```

Đăng ký handler cho tất cả HTTP methods phổ biến.

#### Custom Methods

```go
func (app *WebApp) Handle(method, path string, handlers ...router.HandlerFunc)
```

### Middleware Management

#### Global Middleware

```go
func (app *WebApp) Use(middleware ...router.HandlerFunc)
```

Thêm middleware áp dụng cho tất cả routes.

#### Route Groups

```go
func (app *WebApp) Group(prefix string) router.Router
```

Tạo route group với prefix và có thể áp dụng middleware riêng.

### Static File Serving

```go
func (app *WebApp) Static(prefix string, root string)
```

Phục vụ static files từ filesystem.

### Server Management

#### Start Server

```go
func (app *WebApp) Run() error
func (app *WebApp) RunTLS(certFile, keyFile string) error
```

#### Graceful Shutdown

```go
func (app *WebApp) Shutdown(ctx context.Context) error
func (app *WebApp) SetShutdownTimeout(timeout time.Duration)
```

### Context Management

#### NewContext

```go
func (app *WebApp) NewContext(w http.ResponseWriter, r *http.Request) forkCtx.Context
```

Tạo context mới cho request handling.

### Utility Methods

#### Router Access

```go
func (app *WebApp) Router() router.Router
```

#### Connection Management

```go
func (app *WebApp) GetActiveConnections() int32
func (app *WebApp) IsShuttingDown() bool
```

## Usage Examples

### Basic Setup

```go
func main() {
    // Tạo WebApp
    app := fork.NewWebApp()
    
    // Cấu hình
    config := fork.DefaultWebAppConfig()
    app.SetConfig(config)
    
    // Đăng ký routes
    app.GET("/", func(c forkCtx.Context) {
        c.JSON(200, map[string]string{
            "message": "Hello World!",
        })
    })
    
    // Khởi động server
    log.Fatal(app.Run())
}
```

### With Middleware

```go
func main() {
    app := fork.NewWebApp()
    
    // Global middleware
    app.Use(func(c forkCtx.Context) {
        log.Printf("%s %s", c.Method(), c.Path())
        c.Next()
    })
    
    // Routes
    app.GET("/users", getUsers)
    app.POST("/users", createUser)
    
    app.Run()
}
```

### Route Groups

```go
func main() {
    app := fork.NewWebApp()
    
    // API v1 group
    v1 := app.Group("/api/v1")
    v1.Use(authMiddleware) // Middleware cho group
    v1.GET("/users", getUsers)
    v1.POST("/users", createUser)
    
    // API v2 group
    v2 := app.Group("/api/v2")
    v2.Use(authMiddleware, rateLimitMiddleware)
    v2.GET("/users", getUsersV2)
    v2.POST("/users", createUserV2)
    
    app.Run()
}
```

### Static Files

```go
func main() {
    app := fork.NewWebApp()
    
    // Phục vụ static files
    app.Static("/static", "./public")
    app.Static("/uploads", "./storage/uploads")
    
    // API routes
    app.GET("/api/health", healthCheck)
    
    app.Run()
}
```

### Custom Configuration

```go
func main() {
    app := fork.NewWebApp()
    
    // Custom config
    config := &fork.WebAppConfig{
        GracefulShutdown: fork.GracefulShutdownConfig{
            Enabled:            true,
            Timeout:            60, // 60 seconds
            WaitForConnections: true,
            OnShutdownStart: func() {
                log.Println("Starting graceful shutdown...")
            },
            OnShutdownComplete: func() {
                log.Println("Shutdown completed")
            },
        },
    }
    
    app.SetConfig(config)
    app.Run()
}
```

### With Custom Adapter

```go
func main() {
    app := fork.NewWebApp()
    
    // Sử dụng FastHTTP adapter
    fastAdapter := fasthttp.NewAdapter(&fasthttp.Config{
        Addr: "localhost:8080",
        ReadTimeout: 10 * time.Second,
        WriteTimeout: 10 * time.Second,
    })
    
    app.SetAdapter(fastAdapter)
    
    app.GET("/", func(c forkCtx.Context) {
        c.String(200, "Hello from FastHTTP!")
    })
    
    app.Run()
}
```

### Graceful Shutdown

```go
func main() {
    app := fork.NewWebApp()
    
    // Setup routes
    app.GET("/", homeHandler)
    
    // Setup graceful shutdown
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    
    go func() {
        <-c
        log.Println("Shutting down server...")
        
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()
        
        if err := app.Shutdown(ctx); err != nil {
            log.Printf("Shutdown error: %v", err)
        }
    }()
    
    log.Fatal(app.Run())
}
```

### Multiple HTTP Methods

```go
func main() {
    app := fork.NewWebApp()
    
    // Single route, multiple methods
    app.Handle("GET", "/users", listUsers)
    app.Handle("POST", "/users", createUser)
    app.Handle("PUT", "/users/:id", updateUser)
    app.Handle("DELETE", "/users/:id", deleteUser)
    
    // Hoặc sử dụng Any cho tất cả methods
    app.Any("/health", func(c forkCtx.Context) {
        c.JSON(200, map[string]string{"status": "ok"})
    })
    
    app.Run()
}
```

## Lifecycle Management

### Application Lifecycle

1. **Creation**: `NewWebApp()` tạo instance mới
2. **Configuration**: `SetConfig()` thiết lập cấu hình
3. **Route Registration**: Đăng ký routes và middlewares
4. **Adapter Setup**: `SetAdapter()` thiết lập HTTP adapter
5. **Server Start**: `Run()` hoặc `RunTLS()` khởi động server
6. **Request Processing**: Xử lý incoming requests
7. **Graceful Shutdown**: `Shutdown()` dọn dẹp resources

### Request Lifecycle

1. **Request Received**: HTTP adapter nhận request
2. **Context Creation**: `NewContext()` tạo request context
3. **Route Matching**: Router tìm handler phù hợp
4. **Middleware Execution**: Thực thi middleware chain
5. **Handler Execution**: Thực thi route handler
6. **Response Generation**: Tạo và gửi response

## Thread Safety

WebApp được thiết kế thread-safe:

- **Read-Write Mutex**: Bảo vệ truy cập đồng thời
- **Atomic Operations**: Quản lý connection counter
- **Immutable State**: Configuration không thay đổi sau khi set

## Best Practices

1. **Single Instance**: Sử dụng một WebApp instance cho toàn ứng dụng
2. **Configuration First**: Thiết lập config trước khi đăng ký routes
3. **Middleware Order**: Đăng ký global middleware trước routes
4. **Error Handling**: Implement proper error handling trong handlers
5. **Graceful Shutdown**: Luôn implement graceful shutdown
6. **Resource Cleanup**: Cleanup resources trong shutdown callbacks

## Related Files

- [`web_app.go`](../web_app.go) - WebApp implementation
- [`config.go`](../config.go) - Configuration structures
- [`constants.go`](../constants.go) - HTTP constants và definitions
- [`router/`](../router/) - Router implementation
- [`adapter/`](../adapter/) - Adapter interface và implementations
