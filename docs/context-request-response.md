# Context, Request & Response - Xử lý HTTP Context

Package `fork/context` cung cấp hệ thống context mạnh mẽ để quản lý HTTP request/response lifecycle trong Fork Framework. Context hoạt động như một container chứa tất cả thông tin và tiện ích cần thiết để xử lý một HTTP request.

## Tổng quan

Context system bao gồm 3 interface chính:

- **Context**: Đối tượng trung tâm quản lý request/response lifecycle
- **Request**: Trừu tượng hóa HTTP request với các tiện ích truy cập
- **Response**: Trừu tượng hóa HTTP response với các tiện ích ghi dữ liệu

## Context Interface

### Core Methods

#### Request/Response Access

```go
// Truy cập request và response objects
Request() Request
Response() Response

// Go context integration
Context() context.Context
WithContext(ctx context.Context) Context
```

#### Middleware Chain Management

```go
// Điều khiển middleware chain
Next()           // Gọi middleware tiếp theo
Abort()          // Dừng middleware chain
IsAborted() bool // Kiểm tra trạng thái abort
```

#### Data Storage

```go
// Lưu trữ dữ liệu trong request lifecycle
Set(key string, value interface{})
Get(key string) (interface{}, bool)
GetString(key string) string
GetBool(key string) bool
GetInt(key string) int
GetFloat64(key string) float64
```

### Request Parameter Access

#### URL Parameters

```go
// Route parameters (/users/:id)
Param(name string) string
Params() map[string]string

// Query parameters (?page=1&limit=10)
Query(name string) string
DefaultQuery(name, defaultValue string) string
QueryArray(name string) []string
QueryMap(prefix string) map[string]string
```

#### Form Data

```go
// Form values
Form(name string) string
DefaultForm(name, defaultValue string) string
FormArray(name string) []string
FormMap(prefix string) map[string]string
```

#### File Uploads

```go
// File handling
FormFile(name string) (*multipart.FileHeader, error)
MultipartForm() (*multipart.Form, error)
SaveUploadedFile(file *multipart.FileHeader, dst string) error
```

### Data Binding

#### JSON/XML Binding

```go
// Automatic data binding
BindJSON(obj interface{}) error
BindXML(obj interface{}) error
BindQuery(obj interface{}) error
BindForm(obj interface{}) error
Bind(obj interface{}) error           // Auto-detect content type
ShouldBind(obj interface{}) error     // Non-validating bind
```

#### With Validation

```go
// Binding with validation
BindAndValidate(obj interface{}) error
ShouldBindAndValidate(obj interface{}) error
ValidateStruct(obj interface{}) error
```

### Response Generation

#### Status and Headers

```go
// HTTP status
Status(code int)

// Headers
Header(key, value string)
GetHeader(key string) string
```

#### Response Body

```go
// JSON response
JSON(code int, obj interface{})

// XML response
XML(code int, obj interface{})

// Plain text
String(code int, format string, values ...interface{})

// Raw data
Data(code int, contentType string, data []byte)

// HTML
HTML(code int, name string, obj interface{})

// File response
File(filepath string)
Attachment(filepath, filename string)
```

#### Cookies

```go
// Cookie management
SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool)
Cookie(name string) (string, error)
Cookies() []*http.Cookie
```

### Utility Methods

```go
// Client information
ClientIP() string
UserAgent() string
ContentType() string
IsWebsocket() bool

// Request data
GetRawData() ([]byte, error)
```

## Request Interface

### Basic Information

```go
// HTTP basics
Method() string                 // GET, POST, PUT, DELETE, etc.
URL() *url.URL                 // Full URL
Path() string                  // URL path
Query() url.Values             // Query parameters
Fragment() string              // URL fragment
```

### Headers and Metadata

```go
// Headers
Header() http.Header           // All headers
GetHeader(key string) string   // Single header
UserAgent() string             // User-Agent header
Referer() string              // Referer header
ContentType() string          // Content-Type header
ContentLength() int64         // Content-Length
```

### Connection Information

```go
// Connection details
Host() string                  // Host header
RemoteAddr() string           // Client address
RequestURI() string           // Original request URI
Scheme() string               // http or https
IsSecure() bool               // HTTPS check
Protocol() string             // HTTP protocol version
```

### Cookies and Form Data

```go
// Cookies
Cookies() []*http.Cookie
Cookie(name string) (*http.Cookie, error)

// Form handling
Form() url.Values              // All form values
PostForm() url.Values          // POST form values
FormValue(key string) string   // Single form value
PostFormValue(key string) string
```

### File Operations

```go
// File uploads
MultipartForm() (*multipart.Form, error)
FormFile(key string) (multipart.File, *multipart.FileHeader, error)

// Request body
Body() io.ReadCloser
```

## Response Interface

### Core Writing

```go
// Write operations
Header() http.Header           // Access headers
Write(data []byte) (int, error) // Write body data
WriteHeader(code int)          // Set status code
WriteString(s string) (int, error) // Write string
```

### Status Management

```go
// Status tracking
Status() int                   // Current status code
Size() int                     // Bytes written
Written() bool                 // Has response been written
```

### Advanced Features

```go
// Connection management
Flush()                        // Flush buffered data
Hijack() (net.Conn, *bufio.ReadWriter, error) // Hijack connection

// HTTP/2 features
Pusher() (http.Pusher, bool)   // Server push support

// Utilities
ResponseWriter() http.ResponseWriter // Original writer
Reset(w http.ResponseWriter)   // Reset for reuse
```

## Usage Examples

### Basic Request Handling

```go
app.GET("/users/:id", func(c forkCtx.Context) {
    // Get route parameter
    userID := c.Param("id")
    
    // Get query parameters
    format := c.DefaultQuery("format", "json")
    
    // Get headers
    authHeader := c.GetHeader("Authorization")
    
    // Response
    c.JSON(200, map[string]interface{}{
        "user_id": userID,
        "format":  format,
        "auth":    authHeader != "",
    })
})
```

### Data Binding

```go
type CreateUserRequest struct {
    Name  string `json:"name" binding:"required"`
    Email string `json:"email" binding:"required,email"`
    Age   int    `json:"age" binding:"min=18"`
}

app.POST("/users", func(c forkCtx.Context) {
    var req CreateUserRequest
    
    // Bind and validate
    if err := c.BindAndValidate(&req); err != nil {
        c.JSON(400, map[string]interface{}{
            "error": "Invalid request",
            "details": err.Error(),
        })
        return
    }
    
    // Create user...
    user := createUser(req)
    
    c.JSON(201, user)
})
```

### File Upload

```go
app.POST("/upload", func(c forkCtx.Context) {
    // Get uploaded file
    file, err := c.FormFile("upload")
    if err != nil {
        c.JSON(400, map[string]string{
            "error": "No file uploaded",
        })
        return
    }
    
    // Save file
    filename := fmt.Sprintf("./uploads/%s", file.Filename)
    if err := c.SaveUploadedFile(file, filename); err != nil {
        c.JSON(500, map[string]string{
            "error": "Failed to save file",
        })
        return
    }
    
    c.JSON(200, map[string]interface{}{
        "message": "File uploaded successfully",
        "filename": file.Filename,
        "size": file.Size,
    })
})
```

### Cookie Management

```go
app.GET("/login", func(c forkCtx.Context) {
    // Set authentication cookie
    c.SetCookie("session_id", "abc123", 3600, "/", "", false, true)
    
    c.JSON(200, map[string]string{
        "message": "Logged in successfully",
    })
})

app.GET("/profile", func(c forkCtx.Context) {
    // Get cookie
    sessionID, err := c.Cookie("session_id")
    if err != nil {
        c.JSON(401, map[string]string{
            "error": "Not authenticated",
        })
        return
    }
    
    // Validate session...
    user := getUserBySession(sessionID)
    c.JSON(200, user)
})
```

### Middleware Context Usage

```go
// Authentication middleware
func AuthMiddleware() func(c forkCtx.Context) {
    return func(c forkCtx.Context) {
        token := c.GetHeader("Authorization")
        
        user, err := validateToken(token)
        if err != nil {
            c.JSON(401, map[string]string{
                "error": "Invalid token",
            })
            c.Abort() // Stop middleware chain
            return
        }
        
        // Store user in context
        c.Set("user", user)
        c.Next() // Continue to next middleware
    }
}

// Route handler
app.GET("/protected", AuthMiddleware(), func(c forkCtx.Context) {
    // Get user from context
    user, exists := c.Get("user")
    if !exists {
        c.JSON(500, map[string]string{
            "error": "User not found in context",
        })
        return
    }
    
    c.JSON(200, map[string]interface{}{
        "message": "Protected data",
        "user": user,
    })
})
```

### Stream Response

```go
app.GET("/stream", func(c forkCtx.Context) {
    // Set headers for streaming
    c.Header("Content-Type", "text/plain")
    c.Header("Transfer-Encoding", "chunked")
    
    // Write status
    c.Status(200)
    
    // Stream data
    for i := 0; i < 10; i++ {
        data := fmt.Sprintf("Chunk %d\n", i)
        c.Response().Write([]byte(data))
        c.Response().Flush()
        
        time.Sleep(100 * time.Millisecond)
    }
})
```

### Custom Response

```go
app.GET("/download", func(c forkCtx.Context) {
    // Set download headers
    c.Header("Content-Disposition", "attachment; filename=data.txt")
    c.Header("Content-Type", "application/octet-stream")
    
    // Generate and write data
    data := generateData()
    c.Data(200, "application/octet-stream", data)
})
```

## Validation System

### Built-in Validation

Context sử dụng `github.com/go-playground/validator/v10` cho validation:

```go
type User struct {
    Name  string `json:"name" validate:"required,min=3,max=50"`
    Email string `json:"email" validate:"required,email"`
    Age   int    `json:"age" validate:"min=18,max=100"`
}

func (c *forkContext) ValidateStruct(obj interface{}) error {
    return c.validator.Struct(obj)
}
```

### Custom Validation

```go
// Register custom validation
app.GET("/setup", func(c forkCtx.Context) {
    validator := c.GetValidator()
    
    validator.RegisterValidation("username", func(fl validator.FieldLevel) bool {
        return len(fl.Field().String()) >= 3
    })
})
```

## Error Handling

### Centralized Error Handling

```go
// Error middleware
func ErrorMiddleware() func(c forkCtx.Context) {
    return func(c forkCtx.Context) {
        defer func() {
            if err := recover(); err != nil {
                c.JSON(500, map[string]interface{}{
                    "error": "Internal server error",
                    "details": err,
                })
            }
        }()
        
        c.Next()
    }
}
```

### Context Error Management

```go
func (c *forkContext) Error(err error) {
    // Store error for later handling
    c.errors = append(c.errors, err)
}
```

## Best Practices

1. **Use Context Storage**: Store request-scoped data trong context
2. **Validate Input**: Luôn validate input data trước khi xử lý
3. **Handle Errors**: Implement proper error handling cho all scenarios
4. **Security Headers**: Set appropriate security headers
5. **Content Type**: Specify correct content types
6. **Resource Cleanup**: Cleanup resources sau khi sử dụng
7. **Middleware Order**: Đặt middleware theo thứ tự logic phù hợp

## Related Files

- [`context/context.go`](../context/context.go) - Context implementation
- [`context/context_interface.go`](../context/context_interface.go) - Context interface definitions
- [`context/request.go`](../context/request.go) - Request implementation
- [`context/response.go`](../context/response.go) - Response implementation
- [`mocks/`](../mocks/) - Mock implementations for testing
