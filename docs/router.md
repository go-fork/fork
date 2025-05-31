# Router - Hệ thống Định tuyến

Package `fork/router` cung cấp hệ thống routing mạnh mẽ và hiệu suất cao cho Fork HTTP Framework. Router sử dụng cấu trúc trie được tối ưu để đảm bảo tốc độ tra cứu nhanh và hỗ trợ các pattern phức tạp.

## Tổng quan

Router system bao gồm:

- **Router Interface**: Định nghĩa các phương thức chuẩn cho routing
- **DefaultRouter**: Implementation mặc định với trie structure
- **Route Groups**: Tổ chức routes theo cấu trúc phân cấp
- **Middleware Support**: Tích hợp middleware chain
- **Pattern Matching**: Hỗ trợ parameters, wildcards, và regex

## Router Interface

### Core Methods

#### Route Registration

```go
// Đăng ký route với method và path cụ thể
Handle(method string, path string, handlers ...HandlerFunc)

// Tạo route group với prefix
Group(prefix string) Router

// Thêm middleware vào router
Use(middleware ...HandlerFunc)
```

#### Static Files

```go
// Phục vụ static files
Static(prefix string, root string)
```

#### Route Management

```go
// Lấy tất cả routes đã đăng ký
Routes() []Route

// HTTP handler integration
ServeHTTP(w http.ResponseWriter, req *http.Request)

// Tìm handler cho method và path
Find(method, path string) HandlerFunc
```

## DefaultRouter Implementation

### Architecture

```go
type DefaultRouter struct {
    basePath    string                          // Base path cho router
    routes      map[string]*trieNode           // Route trees theo method
    middleware  []HandlerFunc                  // Global middleware
    staticDirs  map[string]string              // Static directory mappings
    mu          sync.RWMutex                   // Thread safety
}
```

### Features

- **Trie-based Routing**: O(k) lookup time với k là độ dài path
- **Parameter Extraction**: Automatic parameter parsing từ URL
- **Wildcard Support**: Catch-all routes với `*`
- **Middleware Chaining**: Efficient middleware execution
- **Thread Safety**: Concurrent-safe route registration và lookup

## Route Patterns

### Static Routes

```go
router.Handle("GET", "/users", getUsersHandler)
router.Handle("POST", "/users", createUserHandler)
```

### Parameter Routes

```go
// Named parameters với `:name`
router.Handle("GET", "/users/:id", getUserHandler)
router.Handle("PUT", "/users/:id", updateUserHandler)

// Multiple parameters
router.Handle("GET", "/users/:id/posts/:postId", getPostHandler)
```

### Wildcard Routes

```go
// Catch-all với `*name`
router.Handle("GET", "/files/*filepath", serveFileHandler)

// Với prefix
router.Handle("GET", "/api/*any", apiHandler)
```

### Priority và Matching

1. **Static routes** (highest priority)
2. **Parameter routes** 
3. **Wildcard routes** (lowest priority)

## HandlerFunc

### Definition

```go
type HandlerFunc func(ctx forkCtx.Context)
```

### Handler Chain

Mỗi route có thể có multiple handlers tạo thành chain:

```go
router.Handle("GET", "/protected", 
    authMiddleware,
    logMiddleware, 
    actualHandler)
```

## Route Groups

### Creating Groups

```go
// Tạo group với prefix
api := router.Group("/api")
v1 := api.Group("/v1")
v2 := api.Group("/v2")
```

### Group Middleware

```go
// Middleware cho toàn group
api.Use(corsMiddleware)
api.Use(authMiddleware)

// Routes trong group
api.Handle("GET", "/users", getUsersHandler)     // -> /api/users
api.Handle("POST", "/users", createUserHandler)  // -> /api/users
```

### Nested Groups

```go
api := router.Group("/api")
{
    api.Use(corsMiddleware)
    
    v1 := api.Group("/v1")
    {
        v1.Use(legacyMiddleware)
        v1.Handle("GET", "/users", getUsersV1)
    }
    
    v2 := api.Group("/v2")
    {
        v2.Use(modernMiddleware)
        v2.Handle("GET", "/users", getUsersV2)
    }
}
```

## Route Information

### Route Structure

```go
type Route struct {
    Method  string      // HTTP method
    Path    string      // URL path pattern
    Handler HandlerFunc // Handler function
}
```

### Getting Routes

```go
// Lấy tất cả routes
allRoutes := router.Routes()

for _, route := range allRoutes {
    fmt.Printf("%s %s\n", route.Method, route.Path)
}
```

## Static File Serving

### Basic Static

```go
// Phục vụ files từ ./public dưới /static
router.Static("/static", "./public")

// Multiple static directories
router.Static("/css", "./assets/css")
router.Static("/js", "./assets/js")
router.Static("/images", "./assets/images")
```

### Advanced Static Configuration

```go
// Custom static handler
router.Handle("GET", "/uploads/*filepath", func(c forkCtx.Context) {
    filepath := c.Param("filepath")
    
    // Security check
    if strings.Contains(filepath, "..") {
        c.Status(403)
        return
    }
    
    fullPath := "./storage/uploads/" + filepath
    c.File(fullPath)
})
```

## Usage Examples

### Basic Router Setup

```go
func main() {
    router := router.NewDefaultRouter()
    
    // Basic routes
    router.Handle("GET", "/", homeHandler)
    router.Handle("GET", "/about", aboutHandler)
    
    // Parameter routes
    router.Handle("GET", "/users/:id", getUserHandler)
    router.Handle("PUT", "/users/:id", updateUserHandler)
    
    // Start server
    http.ListenAndServe(":8080", router)
}

func homeHandler(c forkCtx.Context) {
    c.JSON(200, map[string]string{
        "message": "Welcome home!",
    })
}

func getUserHandler(c forkCtx.Context) {
    userID := c.Param("id")
    c.JSON(200, map[string]string{
        "user_id": userID,
    })
}
```

### REST API Example

```go
func setupRoutes() router.Router {
    r := router.NewDefaultRouter()
    
    // Middleware
    r.Use(loggerMiddleware)
    r.Use(corsMiddleware)
    
    // API routes
    api := r.Group("/api/v1")
    {
        // Users
        users := api.Group("/users")
        users.Handle("GET", "", listUsers)
        users.Handle("POST", "", createUser)
        users.Handle("GET", "/:id", getUser)
        users.Handle("PUT", "/:id", updateUser)
        users.Handle("DELETE", "/:id", deleteUser)
        
        // Posts
        posts := api.Group("/posts")
        posts.Use(authMiddleware) // Auth required for posts
        posts.Handle("GET", "", listPosts)
        posts.Handle("POST", "", createPost)
        posts.Handle("GET", "/:id", getPost)
        posts.Handle("PUT", "/:id", updatePost)
        posts.Handle("DELETE", "/:id", deletePost)
    }
    
    // Static files
    r.Static("/static", "./public")
    
    return r
}
```

### Middleware Integration

```go
// Logger middleware
func loggerMiddleware(c forkCtx.Context) {
    start := time.Now()
    
    // Process request
    c.Next()
    
    // Log after processing
    duration := time.Since(start)
    log.Printf("%s %s - %v", 
        c.Method(), 
        c.Path(), 
        duration)
}

// Auth middleware
func authMiddleware(c forkCtx.Context) {
    token := c.GetHeader("Authorization")
    
    if !isValidToken(token) {
        c.JSON(401, map[string]string{
            "error": "Unauthorized",
        })
        c.Abort() // Stop chain
        return
    }
    
    // Add user info to context
    user := getUserFromToken(token)
    c.Set("user", user)
    
    c.Next()
}

// Route với middleware
router.Handle("GET", "/protected", 
    authMiddleware,
    func(c forkCtx.Context) {
        user, _ := c.Get("user")
        c.JSON(200, map[string]interface{}{
            "message": "Protected data",
            "user": user,
        })
    })
```

### Parameter Extraction

```go
// Route với multiple parameters
router.Handle("GET", "/users/:userId/posts/:postId/comments/:commentId", 
    func(c forkCtx.Context) {
        userID := c.Param("userId")
        postID := c.Param("postId")
        commentID := c.Param("commentId")
        
        c.JSON(200, map[string]string{
            "user_id": userID,
            "post_id": postID,
            "comment_id": commentID,
        })
    })

// Wildcard route
router.Handle("GET", "/files/*path", func(c forkCtx.Context) {
    filePath := c.Param("path")
    
    // Serve file
    fullPath := "./storage/" + filePath
    c.File(fullPath)
})
```

### Advanced Route Organization

```go
type RouteGroup struct {
    router router.Router
}

func NewRouteGroup(r router.Router) *RouteGroup {
    return &RouteGroup{router: r}
}

func (rg *RouteGroup) SetupUserRoutes() {
    users := rg.router.Group("/users")
    users.Use(validateJSONMiddleware)
    
    users.Handle("GET", "", rg.listUsers)
    users.Handle("POST", "", rg.createUser)
    users.Handle("GET", "/:id", rg.getUser)
    users.Handle("PUT", "/:id", rg.updateUser)
    users.Handle("DELETE", "/:id", rg.deleteUser)
    
    // User profile routes
    profile := users.Group("/:id/profile")
    profile.Handle("GET", "", rg.getUserProfile)
    profile.Handle("PUT", "", rg.updateUserProfile)
    profile.Handle("POST", "/avatar", rg.uploadAvatar)
}

func (rg *RouteGroup) listUsers(c forkCtx.Context) {
    // Implementation...
}
```

## Performance Optimization

### Trie Structure

Router sử dụng trie (prefix tree) để tối ưu hiệu suất:

```go
type trieNode struct {
    path     string           // Path segment
    children map[string]*trieNode // Child nodes
    handler  HandlerFunc      // Handler for this node
    isParam  bool            // Is parameter node
    isWild   bool            // Is wildcard node
    paramName string         // Parameter name
}
```

### Lookup Complexity

- **Static routes**: O(1) average case
- **Parameter routes**: O(k) với k là số segments
- **Wildcard routes**: O(k) với fallback matching

### Memory Optimization

- Shared path segments
- Lazy node creation
- Efficient string handling

## Thread Safety

Router implementation là thread-safe:

```go
type DefaultRouter struct {
    // ... other fields
    mu sync.RWMutex // Read-write mutex
}

func (r *DefaultRouter) Handle(method, path string, handlers ...HandlerFunc) {
    r.mu.Lock()
    defer r.mu.Unlock()
    // ... implementation
}

func (r *DefaultRouter) Find(method, path string) HandlerFunc {
    r.mu.RLock()
    defer r.mu.RUnlock()
    // ... implementation
}
```

## Best Practices

1. **Route Organization**: Sử dụng groups để tổ chức routes logic
2. **Middleware Order**: Đặt global middleware trước specific middleware
3. **Parameter Validation**: Validate parameters trong handlers
4. **Error Handling**: Implement proper error handling cho route not found
5. **Static Files**: Sử dụng dedicated static file server cho production
6. **Performance**: Minimize middleware overhead cho high-traffic routes
7. **Security**: Implement proper authentication và authorization

## Integration với WebApp

WebApp sử dụng router thông qua wrapper methods:

```go
// WebApp methods delegate to router
func (app *WebApp) GET(path string, handlers ...router.HandlerFunc) {
    app.router.Handle("GET", path, handlers...)
}

func (app *WebApp) Group(prefix string) router.Router {
    return app.router.Group(prefix)
}
```

## Related Files

- [`router/router.go`](../router/router.go) - Router implementation
- [`router/trie.go`](../router/trie.go) - Trie data structure
- [`router/router_test.go`](../router/router_test.go) - Tests và examples
- [`constants.go`](../constants.go) - HTTP method constants
