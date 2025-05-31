# Middleware System - Fork HTTP Framework

Middleware là một thành phần cốt lõi của Fork HTTP Framework, cung cấp các chức năng cross-cutting như authentication, logging, compression, security headers và nhiều tính năng khác. Framework cung cấp một hệ thống middleware mạnh mẽ và linh hoạt với 30 middleware packages đa dạng, **tự động apply thông qua YAML configuration**.

## 📋 Mục lục

1. [Tổng quan Middleware System](#-tổng-quan-middleware-system)
2. [YAML-Based Configuration](#️-yaml-based-configuration)
3. [Middleware Categories](#️-middleware-categories)
4. [Auto-Loading Patterns](#-auto-loading-patterns)
5. [Advanced Configuration](#️-advanced-configuration)
6. [Custom Middleware](#️-custom-middleware)
7. [Performance & Best Practices](#-performance--best-practices)
8. [Middleware Reference](#-middleware-reference)

## 🏗️ Tổng quan Middleware System

### Auto-Configuration Pattern

Fork sử dụng **YAML-based auto-configuration** để tự động apply middleware. Thay vì manual configuration, bạn chỉ cần enable/disable middleware trong YAML config file và framework sẽ tự động load và apply middleware theo đúng thứ tự.

```yaml
# configs/app.yaml
http:
  middleware:
    # Security middleware
    recover:
      enabled: true
    helmet:
      enabled: true
    cors:
      enabled: true
      
    # Monitoring middleware  
    logger:
      enabled: true
    monitor:
      enabled: false
      
    # Performance middleware
    compression:
      enabled: true
    cache:
      enabled: true
```

### Key Features

- **Auto-Loading**: Middleware tự động apply khi `enabled: true`
- **YAML Configuration**: Tất cả config trong YAML files
- **Dependency Injection**: Tích hợp sâu với DI container
- **Execution Order**: Framework tự động sắp xếp middleware theo priority
- **Type Safety**: Strong typing với validation
- **Zero Code**: Không cần manual `app.Use()` calls
- **Environment Aware**: Config khác nhau cho dev/prod/test

## ⚙️ YAML-Based Configuration

### Basic Configuration Structure

Fork middleware được cấu hình hoàn toàn thông qua YAML files. Framework tự động detect và apply middleware dựa trên configuration.

```yaml
# configs/app.yaml
http:
  middleware:
    # Enable/disable middleware
    middleware_name:
      enabled: true|false
      # Middleware-specific config
      config_key: config_value
```

### Complete Configuration Example

```yaml
# configs/app.yaml
http:
  middleware:
    # 🔒 Security & Authentication
    recover:
      enabled: true
      enable_stack_trace: true
      
    helmet:
      enabled: true
      xss_protection: "1; mode=block"
      content_type_nosniff: "nosniff"
      x_frame_options: "DENY"
      hsts_max_age: 31536000
      content_security_policy: "default-src 'self'"
      
    cors:
      enabled: true
      allow_origins: ["https://app.example.com", "*.mydomain.com"]
      allow_methods: ["GET", "POST", "PUT", "DELETE"]
      allow_headers: ["Content-Type", "Authorization"]
      allow_credentials: true
      max_age: 86400
      
    csrf:
      enabled: true
      token_lookup: "header:X-CSRF-Token"
      cookie_name: "_csrf"
      cookie_same_site: "Strict"
      expiration: "24h"
      
    basicauth:
      enabled: false
      users:
        admin: "password"
        user: "secret"
      realm: "Restricted Area"
      
    keyauth:
      enabled: false
      key_lookup: "header:X-API-Key"
      
    # 📊 Monitoring & Debugging  
    logger:
      enabled: true
      format: "${time} | ${status} | ${latency} | ${ip} | ${method} ${path}"
      time_format: "2006/01/02 - 15:04:05"
      output: "stdout"
      enable_colors: true
      skip_paths: ["/health", "/metrics"]
      
    monitor:
      enabled: false
      title: "Fork Metrics"
      refresh: "5s"
      api_only: false
      
    requestid:
      enabled: true
      header: "X-Request-ID"
      
    pprof:
      enabled: false
      prefix: "/debug/pprof"
      
    # ⚡ Performance & Content
    compression:
      enabled: true
      level: 6
      types: ["text/html", "text/css", "text/javascript", "application/json"]
      
    cache:
      enabled: true
      duration: "5m"
      cache_header: "X-Cache"
      methods: ["GET", "HEAD"]
      status_codes: [200, 301, 404]
      
    etag:
      enabled: true
      weak: false
      
    static:
      enabled: true
      root: "./public"
      index_names: ["index.html", "index.htm"]
      browse: false
      max_age: 3600
      
    bodylimit:
      enabled: true
      max_bytes: 4194304  # 4MB
      
    favicon:
      enabled: true
      file: "./static/favicon.ico"
      cache_control: "public, max-age=31536000"
      
    # 🚦 Rate Limiting & Control
    limiter:
      enabled: true
      max: 100
      duration: "1m"
      key_generator: "ip"
      
    timeout:
      enabled: true
      timeout: "30s"
      
    method:
      enabled: false
      methods: ["GET", "POST", "PUT", "DELETE"]
      
    # 🔄 Session & State
    session:
      enabled: false
      store: "memory"
      cookie_name: "session_id"
      cookie_secure: true
      cookie_http_only: true
      expiration: "24h"
      
    # 🌐 Infrastructure & Utilities
    healthcheck:
      enabled: true
      path: "/health"
      
    proxy:
      enabled: false
      targets: ["http://server1:8080", "http://server2:8080"]
      strategy: "round_robin"
      
    redirect:
      enabled: false
      rules:
        "/old-path": "/new-path"
        "/api/v1/*": "/api/v2/$1"
      status_code: 301
      
    # 🔧 Advanced Features
    skip:
      enabled: false
      
    earlydata:
      enabled: false
      
    encryptcookie:
      enabled: false
      
    envvar:
      enabled: false
      
    expvar:
      enabled: false
      
    idempotency:
      enabled: false
      
    rewrite:
      enabled: false
```

### Environment-Specific Configuration

```yaml
# configs/app.dev.yaml
http:
  middleware:
    logger:
      enabled: true
      enable_colors: true
      output: "stdout"
    cors:
      enabled: true
      allow_origins: ["*"]  # Relaxed for development
    recover:
      enabled: true
      enable_stack_trace: true
    helmet:
      enabled: false  # Disabled in development

---
# configs/app.prod.yaml  
http:
  middleware:
    logger:
      enabled: true
      enable_colors: false
      output: "file"
      level: "info"
    cors:
      enabled: true
      allow_origins: ["https://myapp.com"]
    recover:
      enabled: true
      enable_stack_trace: false
    helmet:
      enabled: true
    limiter:
      enabled: true
      max: 1000

---
# configs/app.test.yaml
http:
  middleware:
    logger:
      enabled: false  # No noise in tests
    recover:
      enabled: true
    # Minimal middleware for fast tests
```

## 🗂️ Middleware Categories

### 🔒 Security & Authentication

#### 1. **BasicAuth** ✅
HTTP Basic Authentication với username/password.

```yaml
# YAML Configuration
http:
  middleware:
    basicauth:
      enabled: true
      users:
        admin: "password"
        user: "secret"
      realm: "Restricted Area"
```

#### 2. **Helmet** ✅  
Security headers (XSS, CSRF, HSTS, Content Security Policy).

```yaml
# YAML Configuration
http:
  middleware:
    helmet:
      enabled: true
      xss_protection: "1; mode=block"
      content_type_nosniff: "nosniff"
      x_frame_options: "DENY"
      hsts_max_age: 31536000
      content_security_policy: "default-src 'self'"
```

#### 3. **CORS** ✅
Cross-Origin Resource Sharing với full CORS support.

```yaml
# YAML Configuration
http:
  middleware:
    cors:
      enabled: true
      allow_origins: ["https://example.com", "*.mydomain.com"]
      allow_methods: ["GET", "POST", "PUT", "DELETE"]
      allow_headers: ["Content-Type", "Authorization"]
      allow_credentials: true
      max_age: 86400
```

#### 4. **CSRF** ✅
Cross-Site Request Forgery protection.

```yaml
# YAML Configuration
http:
  middleware:
    csrf:
      enabled: true
      token_lookup: "header:X-CSRF-Token"
      cookie_name: "_csrf"
      cookie_same_site: "Strict"
      expiration: "24h"
```

#### 5. **Keyauth** 🔧
API key authentication.

```yaml
# YAML Configuration
http:
  middleware:
    keyauth:
      enabled: true
      key_lookup: "header:X-API-Key"
      valid_keys: ["key1", "key2", "key3"]
```

### 📊 Monitoring & Debugging

#### 6. **Logger** ✅
HTTP request logging với customizable format.

```yaml
# YAML Configuration
http:
  middleware:
    logger:
      enabled: true
      format: "${time} | ${status} | ${latency} | ${ip} | ${method} ${path}"
      time_format: "2006/01/02 - 15:04:05"
      output: "stdout"
      enable_colors: true
      skip_paths: ["/health", "/metrics"]
```

#### 7. **Monitor** ✅
Real-time metrics và monitoring dashboard.

```yaml
# YAML Configuration
http:
  middleware:
    monitor:
      enabled: true
      title: "Fork Metrics"
      refresh: "5s"
      api_only: false
```

#### 8. **Recover** ✅
Panic recovery với stack traces.

```yaml
# YAML Configuration
http:
  middleware:
    recover:
      enabled: true
      enable_stack_trace: true
```

#### 9. **Pprof** 🔧
Go profiling endpoints cho performance analysis.

```yaml
# YAML Configuration
http:
  middleware:
    pprof:
      enabled: false  # Only enable in development
      prefix: "/debug/pprof"
```

#### 10. **RequestID** ✅
Unique request ID generation cho tracing.

```yaml
# YAML Configuration
http:
  middleware:
    requestid:
      enabled: true
      header: "X-Request-ID"
```

### ⚡ Performance & Content

#### 11. **Compression** ✅
Gzip/Deflate response compression.

```yaml
# YAML Configuration
http:
  middleware:
    compression:
      enabled: true
      level: 6  # 1-9 compression level
      types: ["text/html", "text/css", "text/javascript", "application/json"]
```

#### 12. **Cache** ✅
HTTP caching với TTL và invalidation.

```yaml
# YAML Configuration
http:
  middleware:
    cache:
      enabled: true
      duration: "5m"
      cache_header: "X-Cache"
      methods: ["GET", "HEAD"]
      status_codes: [200, 301, 404]
```

#### 13. **ETag** ✅
Entity tag cho cache validation.

```yaml
# YAML Configuration
http:
  middleware:
    etag:
      enabled: true
      weak: false
```

#### 14. **Static** ✅
Static file serving với optimization.

```yaml
# YAML Configuration
http:
  middleware:
    static:
      enabled: true
      root: "./public"
      index_names: ["index.html", "index.htm"]
      browse: false
      max_age: 3600
```

#### 15. **BodyLimit** ✅
Request body size limitation.

```yaml
# YAML Configuration
http:
  middleware:
    bodylimit:
      enabled: true
      max_bytes: 4194304  # 4MB
```

### 🚦 Rate Limiting & Control

#### 16. **Limiter** ✅
Rate limiting với token bucket algorithm.

```yaml
# YAML Configuration
http:
  middleware:
    limiter:
      enabled: true
      max: 100          # 100 requests
      duration: "1m"    # per minute
      key_generator: "ip"  # Limit by IP
```

#### 17. **Timeout** ✅
Request timeout handling.

```yaml
# YAML Configuration
http:
  middleware:
    timeout:
      enabled: true
      timeout: "30s"
```

#### 18. **Method** ✅
HTTP method validation.

```yaml
# YAML Configuration
http:
  middleware:
    method:
      enabled: true
      methods: ["GET", "POST", "PUT", "DELETE"]
```

### 🔄 Session & State

#### 19. **Session** ✅
Session management với multiple stores.

```yaml
# YAML Configuration
http:
  middleware:
    session:
      enabled: true
      store: "memory"  # memory, redis, file
      cookie_name: "session_id"
      cookie_secure: true
      cookie_http_only: true
      expiration: "24h"
```

### 🌐 Infrastructure & Utilities

#### 20. **HealthCheck** 🔧
Health check endpoints.

```yaml
# YAML Configuration
http:
  middleware:
    healthcheck:
      enabled: true
      path: "/health"
```

#### 21. **Proxy** 🔧
Reverse proxy với load balancing.

```yaml
# YAML Configuration
http:
  middleware:
    proxy:
      enabled: false
      targets: ["http://server1:8080", "http://server2:8080"]
      strategy: "round_robin"
```

#### 22. **Redirect** 🔧
URL redirection với pattern matching.

```yaml
# YAML Configuration
http:
  middleware:
    redirect:
      enabled: false
      rules:
        "/old-path": "/new-path"
        "/api/v1/*": "/api/v2/$1"
      status_code: 301
```

#### 23. **Favicon** ✅
Favicon serving optimization.

```yaml
# YAML Configuration
http:
  middleware:
    favicon:
      enabled: true
      file: "./static/favicon.ico"
      cache_control: "public, max-age=31536000"
```

### 🔧 Advanced Features

#### 24-30. **Advanced Middleware**
- **Skip**: Conditional middleware execution
- **EarlyData**: HTTP/2 early data handling
- **EncryptCookie**: Cookie encryption/decryption  
- **EnvVar**: Environment variable injection
- **ExpVar**: Go expvar metrics exposure
- **Idempotency**: Idempotent request handling
- **Rewrite**: URL rewriting middleware

```yaml
# YAML Configuration for Advanced Middleware
http:
  middleware:
    skip:
      enabled: false
    earlydata:
      enabled: false
    encryptcookie:
      enabled: false
    envvar:
      enabled: false
    expvar:
      enabled: false
    idempotency:
      enabled: false
    rewrite:
      enabled: false
```

## 🔄 Auto-Loading System

### Framework Auto-Detection

Fork tự động scan và apply middleware dựa trên YAML configuration:

```go
// Framework tự động thực hiện điều này:
func (app *WebApp) loadMiddlewareFromConfig() {
    config := app.Config.HTTP.Middleware
    
    // Auto-apply middleware theo thứ tự priority
    if config.Recover.Enabled {
        app.Use(recover.New(config.Recover))
    }
    if config.Logger.Enabled {
        app.Use(logger.New(config.Logger))
    }
    if config.CORS.Enabled {
        app.Use(cors.New(config.CORS))
    }
    // ... tự động apply tất cả middleware được enable
}
```

### Service Provider Auto-Registration

```go
// Framework tự động register middleware providers
func (app *WebApp) registerMiddlewareProviders() {
    providers := []ServiceProvider{
        &middleware.RecoverProvider{},
        &middleware.LoggerProvider{},
        &middleware.CORSProvider{},
        &middleware.HelmetProvider{},
        // ... all middleware providers
    }
    
    for _, provider := range providers {
        app.RegisterProvider(provider)
    }
}
```

## ⚙️ YAML-Based Configuration

### Basic Configuration Structure

Fork middleware được cấu hình hoàn toàn thông qua YAML files. Framework tự động detect và apply middleware dựa trên configuration.

```yaml
# configs/app.yaml
http:
  middleware:
    # Enable/disable middleware
    middleware_name:
      enabled: true|false
      # Middleware-specific config
      config_key: config_value
```

### Complete Configuration Example

```yaml
# configs/app.yaml
http:
  middleware:
    # 🔒 Security & Authentication
    recover:
      enabled: true
      enable_stack_trace: true
      
    helmet:
      enabled: true
      xss_protection: "1; mode=block"
      content_type_nosniff: "nosniff"
      x_frame_options: "DENY"
      hsts_max_age: 31536000
      content_security_policy: "default-src 'self'"
      
    cors:
      enabled: true
      allow_origins: ["https://app.example.com", "*.mydomain.com"]
      allow_methods: ["GET", "POST", "PUT", "DELETE"]
      allow_headers: ["Content-Type", "Authorization"]
      allow_credentials: true
      max_age: 86400
      
    csrf:
      enabled: true
      token_lookup: "header:X-CSRF-Token"
      cookie_name: "_csrf"
      cookie_same_site: "Strict"
      expiration: "24h"
      
    basicauth:
      enabled: false
      users:
        admin: "password"
        user: "secret"
      realm: "Restricted Area"
      
    keyauth:
      enabled: false
      key_lookup: "header:X-API-Key"
      
    # 📊 Monitoring & Debugging  
    logger:
      enabled: true
      format: "${time} | ${status} | ${latency} | ${ip} | ${method} ${path}"
      time_format: "2006/01/02 - 15:04:05"
      output: "stdout"
      enable_colors: true
      skip_paths: ["/health", "/metrics"]
      
    monitor:
      enabled: false
      title: "Fork Metrics"
      refresh: "5s"
      api_only: false
      
    requestid:
      enabled: true
      header: "X-Request-ID"
      
    pprof:
      enabled: false
      prefix: "/debug/pprof"
      
    # ⚡ Performance & Content
    compression:
      enabled: true
      level: 6
      types: ["text/html", "text/css", "text/javascript", "application/json"]
      
    cache:
      enabled: true
      duration: "5m"
      cache_header: "X-Cache"
      methods: ["GET", "HEAD"]
      status_codes: [200, 301, 404]
      
    etag:
      enabled: true
      weak: false
      
    static:
      enabled: true
      root: "./public"
      index_names: ["index.html", "index.htm"]
      browse: false
      max_age: 3600
      
    bodylimit:
      enabled: true
      max_bytes: 4194304  # 4MB
      
    favicon:
      enabled: true
      file: "./static/favicon.ico"
      cache_control: "public, max-age=31536000"
      
    # 🚦 Rate Limiting & Control
    limiter:
      enabled: true
      max: 100
      duration: "1m"
      key_generator: "ip"
      
    timeout:
      enabled: true
      timeout: "30s"
      
    method:
      enabled: false
      methods: ["GET", "POST", "PUT", "DELETE"]
      
    # 🔄 Session & State
    session:
      enabled: false
      store: "memory"
      cookie_name: "session_id"
      cookie_secure: true
      cookie_http_only: true
      expiration: "24h"
      
    # 🌐 Infrastructure & Utilities
    healthcheck:
      enabled: true
      path: "/health"
      
    proxy:
      enabled: false
      targets: ["http://server1:8080", "http://server2:8080"]
      strategy: "round_robin"
      
    redirect:
      enabled: false
      rules:
        "/old-path": "/new-path"
        "/api/v1/*": "/api/v2/$1"
      status_code: 301
      
    # 🔧 Advanced Features
    skip:
      enabled: false
      
    earlydata:
      enabled: false
      
    encryptcookie:
      enabled: false
      
    envvar:
      enabled: false
      
    expvar:
      enabled: false
      
    idempotency:
      enabled: false
      
    rewrite:
      enabled: false
```

### Environment-Specific Configuration

```yaml
# configs/app.dev.yaml
http:
  middleware:
    logger:
      enabled: true
      enable_colors: true
      output: "stdout"
    cors:
      enabled: true
      allow_origins: ["*"]  # Relaxed for development
    recover:
      enabled: true
      enable_stack_trace: true
    helmet:
      enabled: false  # Disabled in development

---
# configs/app.prod.yaml  
http:
  middleware:
    logger:
      enabled: true
      enable_colors: false
      output: "file"
      level: "info"
    cors:
      enabled: true
      allow_origins: ["https://myapp.com"]
    recover:
      enabled: true
      enable_stack_trace: false
    helmet:
      enabled: true
    limiter:
      enabled: true
      max: 1000

---
# configs/app.test.yaml
http:
  middleware:
    logger:
      enabled: false  # No noise in tests
    recover:
      enabled: true
    # Minimal middleware for fast tests
```

## 🗂️ Middleware Categories

### 🔒 Security & Authentication

#### 1. **BasicAuth** ✅
HTTP Basic Authentication với username/password.

```go
import "github.com/go-fork/middleware/basicauth"

app.Use(basicauth.New(basicauth.Config{
    Users: map[string]string{
        "admin": "password",
        "user":  "secret",
    },
    Realm: "Restricted Area",
}))
```

#### 2. **Helmet** ✅  
Security headers (XSS, CSRF, HSTS, Content Security Policy).

```go
import "github.com/go-fork/middleware/helmet"

app.Use(helmet.New(helmet.Config{
    XSSProtection:         "1; mode=block",
    ContentTypeNosniff:    "nosniff",
    XFrameOptions:         "DENY",
    HSTSMaxAge:           31536000,
    ContentSecurityPolicy: "default-src 'self'",
}))
```

#### 3. **CORS** ✅
Cross-Origin Resource Sharing với full CORS support.

```go
import "github.com/go-fork/middleware/cors"

app.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"https://example.com", "*.mydomain.com"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders:     []string{"Content-Type", "Authorization"},
    AllowCredentials: true,
    MaxAge:          86400,
}))
```

#### 4. **CSRF** ✅
Cross-Site Request Forgery protection.

```go
import "github.com/go-fork/middleware/csrf"

app.Use(csrf.New(csrf.Config{
    TokenLookup:    "header:X-CSRF-Token",
    CookieName:     "_csrf",
    CookieSameSite: "Strict",
    Expiration:     24 * time.Hour,
}))
```

#### 5. **Keyauth** 🔧
API key authentication.

```go
import "github.com/go-fork/middleware/keyauth"

app.Use(keyauth.New(keyauth.Config{
    KeyLookup: "header:X-API-Key",
    Validator: func(key string) bool {
        return validateAPIKey(key)
    },
}))
```

### 📊 Monitoring & Debugging

#### 6. **Logger** ✅
HTTP request logging với customizable format.

```go
import "github.com/go-fork/middleware/logger"

app.Use(logger.New(logger.Config{
    Format: "${time} | ${status} | ${latency} | ${ip} | ${method} ${path}",
    TimeFormat: "2006/01/02 - 15:04:05",
    Output: "stdout",
    EnableColors: true,
    SkipPaths: []string{"/health", "/metrics"},
}))
```

#### 7. **Monitor** ✅
Real-time metrics và monitoring dashboard.

```go
import "github.com/go-fork/middleware/monitor"

app.Use(monitor.New(monitor.Config{
    Title: "Fork Metrics",
    Refresh: time.Second * 5,
    APIOnly: false,
    Next: func(c Context) bool {
        return c.Path() == "/private"
    },
}))
```

#### 8. **Recover** ✅
Panic recovery với stack traces.

```go
import "github.com/go-fork/middleware/recover"

app.Use(recover.New(recover.Config{
    EnableStackTrace: true,
    StackTraceHandler: func(c Context, e interface{}) {
        log.Printf("Panic recovered: %v", e)
    },
}))
```

#### 9. **Pprof** 🔧
Go profiling endpoints cho performance analysis.

```go
import "github.com/go-fork/middleware/pprof"

app.Use(pprof.New(pprof.Config{
    Prefix: "/debug/pprof",
}))
```

#### 10. **RequestID** ✅
Unique request ID generation cho tracing.

```go
import "github.com/go-fork/middleware/requestid"

app.Use(requestid.New(requestid.Config{
    Header:    "X-Request-ID",
    Generator: func() string {
        return uuid.New().String()
    },
}))
```

### ⚡ Performance & Content

#### 11. **Compression** ✅
Gzip/Deflate response compression.

```go
import "github.com/go-fork/middleware/compression"

app.Use(compression.New(compression.Config{
    Level: 6, // 1-9 compression level
    Types: []string{
        "text/html",
        "text/css", 
        "text/javascript",
        "application/json",
    },
}))
```

#### 12. **Cache** ✅
HTTP caching với TTL và invalidation.

```go
import "github.com/go-fork/middleware/cache"

app.Use(cache.New(cache.Config{
    Duration: 5 * time.Minute,
    CacheHeader: "X-Cache",
    Methods: []string{"GET", "HEAD"},
    StatusCodes: []int{200, 301, 404},
}))
```

#### 13. **ETag** ✅
Entity tag cho cache validation.

```go
import "github.com/go-fork/middleware/etag"

app.Use(etag.New(etag.Config{
    Weak: false,
    Generator: func(body []byte) string {
        return fmt.Sprintf(`"%x"`, md5.Sum(body))
    },
}))
```

#### 14. **Static** ✅
Static file serving với optimization.

```go
import "github.com/go-fork/middleware/static"

app.Use(static.New(static.Config{
    Root:       "./public",
    IndexNames: []string{"index.html", "index.htm"},
    Browse:     false,
    MaxAge:     3600,
}))
```

#### 15. **BodyLimit** ✅
Request body size limitation.

```go
import "github.com/go-fork/middleware/bodylimit"

app.Use(bodylimit.New(bodylimit.Config{
    MaxBytes: 4 * 1024 * 1024, // 4MB
}))
```

### 🚦 Rate Limiting & Control

#### 16. **Limiter** ✅
Rate limiting với token bucket algorithm.

```go
import "github.com/go-fork/middleware/limiter"

app.Use(limiter.New(limiter.Config{
    Max:        100,                  // 100 requests
    Duration:   time.Minute,          // per minute
    KeyGenerator: func(c Context) string {
        return c.IP()  // Limit by IP
    },
    LimitReached: func(c Context) error {
        return c.Status(429).JSON(map[string]string{
            "error": "Too many requests",
        })
    },
}))
```

#### 17. **Timeout** ✅
Request timeout handling.

```go
import "github.com/go-fork/middleware/timeout"

app.Use(timeout.New(timeout.Config{
    Timeout: 30 * time.Second,
    TimeoutHandler: func(c Context) error {
        return c.Status(408).JSON(map[string]string{
            "error": "Request timeout",
        })
    },
}))
```

#### 18. **Method** ✅
HTTP method validation.

```go
import "github.com/go-fork/middleware/method"

app.Use(method.New(method.Config{
    Methods: []string{"GET", "POST", "PUT", "DELETE"},
}))
```

### 🔄 Session & State

#### 19. **Session** ✅
Session management với multiple stores.

```go
import "github.com/go-fork/middleware/session"

app.Use(session.New(session.Config{
    Store: session.NewMemoryStore(),
    CookieName: "session_id",
    CookieSecure: true,
    CookieHTTPOnly: true,
    Expiration: 24 * time.Hour,
}))
```

### 🌐 Infrastructure & Utilities

#### 20. **HealthCheck** 🔧
Health check endpoints.

```go
import "github.com/go-fork/middleware/healthcheck"

app.Use(healthcheck.New(healthcheck.Config{
    Path: "/health",
    Checker: func() error {
        // Check database, external services, etc.
        return checkDependencies()
    },
}))
```

#### 21. **Proxy** 🔧
Reverse proxy với load balancing.

```go
import "github.com/go-fork/middleware/proxy"

app.Use(proxy.New(proxy.Config{
    Targets: []string{
        "http://server1:8080",
        "http://server2:8080",
    },
    Strategy: "round_robin",
}))
```

#### 22. **Redirect** 🔧
URL redirection với pattern matching.

```go
import "github.com/go-fork/middleware/redirect"

app.Use(redirect.New(redirect.Config{
    Rules: map[string]string{
        "/old-path": "/new-path",
        "/api/v1/*": "/api/v2/$1",
    },
    StatusCode: 301,
}))
```

#### 23. **Favicon** ✅
Favicon serving optimization.

```go
import "github.com/go-fork/middleware/favicon"

app.Use(favicon.New(favicon.Config{
    File: "./static/favicon.ico",
    CacheControl: "public, max-age=31536000",
}))
```

### 🔧 Advanced Features

#### 24. **Skip** 🔧
Conditional middleware execution.

```go
import "github.com/go-fork/middleware/skip"

app.Use(skip.New(authMiddleware, skip.Config{
    Skipper: func(c Context) bool {
        return c.Path() == "/public" || 
               strings.HasPrefix(c.Path(), "/api/public/")
    },
}))
```

#### 25-30. **Advanced Middleware** ❌
- **EarlyData**: HTTP/2 early data handling
- **EncryptCookie**: Cookie encryption/decryption  
- **EnvVar**: Environment variable injection
- **ExpVar**: Go expvar metrics exposure
- **Idempotency**: Idempotent request handling
- **Rewrite**: URL rewriting middleware

## ⚙️ Advanced Configuration

### Conditional Middleware Loading

```yaml
# configs/app.yaml
http:
  middleware:
    # Conditional loading based on environment
    logger:
      enabled: ${HTTP_LOGGER_ENABLED:true}
      level: ${LOG_LEVEL:info}
      
    cors:
      enabled: ${CORS_ENABLED:false}
      allow_origins: ${CORS_ORIGINS:["*"]}
      
    # Feature flags
    monitor:
      enabled: ${FEATURE_MONITORING:false}
      
    pprof:
      enabled: ${DEBUG_MODE:false}
```

### Route-Specific Middleware

```yaml
# configs/routes.yaml
routes:
  api:
    prefix: "/api"
    middleware:
      - keyauth
      - limiter
      - logger
    routes:
      - path: "/users"
        method: "GET"
        handler: "getUsersHandler"
        
  admin:
    prefix: "/admin"
    middleware:
      - basicauth
      - csrf
      - helmet
    routes:
      - path: "/dashboard"
        method: "GET"
        handler: "dashboardHandler"
```

### Middleware Groups

```yaml
# configs/middleware-groups.yaml
middleware_groups:
  security:
    - recover
    - helmet
    - cors
    - csrf
    
  performance:
    - compression
    - cache
    - etag
    
  monitoring:
    - logger
    - monitor
    - requestid
    
  api:
    - keyauth
    - limiter
    - bodylimit

# Apply groups
http:
  middleware_groups:
    - security
    - performance
    - monitoring
```

### Custom Configuration Loading

```go
// Custom config loader
type MiddlewareConfig struct {
    HTTP struct {
        Middleware map[string]interface{} `yaml:"middleware"`
    } `yaml:"http"`
}

func LoadCustomMiddlewareConfig(file string) (*MiddlewareConfig, error) {
    data, err := ioutil.ReadFile(file)
    if err != nil {
        return nil, err
    }
    
    var config MiddlewareConfig
    err = yaml.Unmarshal(data, &config)
    return &config, err
}
```

## 🛠️ Custom Middleware

### YAML-Compatible Custom Middleware

```go
// Define custom middleware với YAML support
type CustomMiddlewareConfig struct {
    Enabled bool          `yaml:"enabled"`
    Prefix  string        `yaml:"prefix"`
    Timeout time.Duration `yaml:"timeout"`
}

type CustomMiddlewareProvider struct{}

func (p *CustomMiddlewareProvider) Register(container *di.Container) error {
    container.Register(func(config CustomMiddlewareConfig) fork.Middleware {
        return func(c fork.Context) error {
            if !config.Enabled {
                return c.Next()
            }
            
            // Add custom logic
            c.Set("X-Custom-Prefix", config.Prefix)
            
            return c.Next()
        }
    })
    return nil
}

func (p *CustomMiddlewareProvider) Boot(container *di.Container) error {
    return nil
}
```

### Register Custom Middleware

```yaml
# configs/app.yaml
http:
  middleware:
    # Built-in middleware
    logger:
      enabled: true
      
    # Custom middleware
    custom:
      enabled: true
      prefix: "MyApp"
      timeout: "30s"
```

```go
// Register custom middleware provider
func main() {
    app := fork.New()
    
    // Register custom provider
    app.RegisterProvider(&CustomMiddlewareProvider{})
    
    // Framework sẽ tự động load custom middleware từ YAML
    app.Start(":8080")
}
```

### Middleware Factory Pattern

```go
// Factory cho dynamic middleware creation
type MiddlewareFactory struct {
    registry map[string]func(config interface{}) fork.Middleware
}

func NewMiddlewareFactory() *MiddlewareFactory {
    return &MiddlewareFactory{
        registry: make(map[string]func(config interface{}) fork.Middleware),
    }
}

func (f *MiddlewareFactory) Register(name string, factory func(config interface{}) fork.Middleware) {
    f.registry[name] = factory
}

func (f *MiddlewareFactory) Create(name string, config interface{}) fork.Middleware {
    if factory, exists := f.registry[name]; exists {
        return factory(config)
    }
    return nil
}

// Usage
factory := NewMiddlewareFactory()
factory.Register("custom", func(config interface{}) fork.Middleware {
    cfg := config.(CustomMiddlewareConfig)
    return NewCustomMiddleware(cfg)
})
```

## 🚀 Performance & Best Practices

### YAML-Based Performance Optimization

```yaml
# configs/performance.yaml
http:
  middleware:
    # Optimal order được framework tự động handle
    recover:
      enabled: true
      order: 1  # Framework priority
      
    logger:
      enabled: true
      order: 2
      skip_paths: ["/health", "/metrics", "/static/*"]
      
    cors:
      enabled: true
      order: 3
      
    compression:
      enabled: true
      order: 4
      level: 6  # Balance between speed and compression ratio
      
    cache:
      enabled: true
      order: 5
      duration: "5m"
      
    # Performance monitoring
    monitor:
      enabled: false  # Disable in production
      
    pprof:
      enabled: false  # Only enable for debugging
```

### Environment-Based Performance Tuning

```yaml
# configs/app.prod.yaml
http:
  middleware:
    # Production optimizations
    logger:
      enabled: true
      output: "file"
      enable_colors: false
      buffer_size: 1024
      
    compression:
      enabled: true
      level: 9  # Maximum compression for production
      
    cache:
      enabled: true
      duration: "30m"  # Longer cache in production
      
    limiter:
      enabled: true
      max: 10000  # Higher limits for production
      
    # Disable debugging middleware
    monitor:
      enabled: false
    pprof:
      enabled: false

---
# configs/app.dev.yaml
http:
  middleware:
    # Development optimizations
    logger:
      enabled: true
      output: "stdout"
      enable_colors: true
      
    compression:
      enabled: false  # Disable for faster development
      
    cache:
      enabled: false  # Disable for fresh data
      
    limiter:
      enabled: false  # No limits in development
      
    # Enable debugging
    monitor:
      enabled: true
    pprof:
      enabled: true
```

### Resource Management

```yaml
# configs/resources.yaml
http:
  middleware:
    bodylimit:
      enabled: true
      max_bytes: 10485760  # 10MB
      
    timeout:
      enabled: true
      timeout: "30s"
      
    session:
      enabled: true
      store: "redis"  # Use Redis for scalability
      expiration: "1h"
      cleanup_interval: "10m"
      
    cache:
      enabled: true
      store: "redis"
      max_memory: "100mb"
      eviction_policy: "lru"
```

### Security Best Practices

```yaml
# configs/security.yaml
http:
  middleware:
    # Security layer 1: Recovery
    recover:
      enabled: true
      enable_stack_trace: false  # Never expose in production
      
    # Security layer 2: Headers
    helmet:
      enabled: true
      xss_protection: "1; mode=block"
      content_type_nosniff: "nosniff"
      x_frame_options: "DENY"
      hsts_max_age: 31536000
      content_security_policy: "default-src 'self'; script-src 'self' 'unsafe-inline'"
      
    # Security layer 3: CORS
    cors:
      enabled: true
      allow_origins: ["https://myapp.com"]  # Specific domains only
      allow_credentials: false  # Disable if not needed
      
    # Security layer 4: CSRF
    csrf:
      enabled: true
      token_lookup: "header:X-CSRF-Token"
      cookie_same_site: "Strict"
      
    # Security layer 5: Rate limiting
    limiter:
      enabled: true
      max: 100
      duration: "1m"
      
    # Security layer 6: Request validation
    bodylimit:
      enabled: true
      max_bytes: 1048576  # 1MB limit
```

## 📚 Middleware Reference

### Complete Package List

| Name | Module Path | Category | Description |
|------|-------------|----------|-------------|
| basicauth | `github.com/go-fork/middleware/basicauth` | Security | HTTP Basic Authentication |
| helmet | `github.com/go-fork/middleware/helmet` | Security | Security headers |
| cors | `github.com/go-fork/middleware/cors` | Security | Cross-Origin Resource Sharing |
| csrf | `github.com/go-fork/middleware/csrf` | Security | CSRF protection |
| keyauth | `github.com/go-fork/middleware/keyauth` | Security | API key authentication |
| limiter | `github.com/go-fork/middleware/limiter` | Control | Rate limiting |
| bodylimit | `github.com/go-fork/middleware/bodylimit` | Control | Request body size limit |
| method | `github.com/go-fork/middleware/method` | Control | HTTP method validation |
| timeout | `github.com/go-fork/middleware/timeout` | Control | Request timeout |
| skip | `github.com/go-fork/middleware/skip` | Control | Conditional middleware execution |
| logger | `github.com/go-fork/middleware/logger` | Monitoring | HTTP request logging |
| monitor | `github.com/go-fork/middleware/monitor` | Monitoring | Real-time metrics dashboard |
| recover | `github.com/go-fork/middleware/recover` | Monitoring | Panic recovery |
| requestid | `github.com/go-fork/middleware/requestid` | Monitoring | Request ID generation |
| pprof | `github.com/go-fork/middleware/pprof` | Monitoring | Go profiling endpoints |
| compression | `github.com/go-fork/middleware/compression` | Performance | Response compression |
| cache | `github.com/go-fork/middleware/cache` | Performance | HTTP caching |
| etag | `github.com/go-fork/middleware/etag` | Performance | Entity tag validation |
| static | `github.com/go-fork/middleware/static` | Content | Static file serving |
| favicon | `github.com/go-fork/middleware/favicon` | Content | Favicon optimization |
| session | `github.com/go-fork/middleware/session` | State | Session management |
| proxy | `github.com/go-fork/middleware/proxy` | Infrastructure | Reverse proxy |
| redirect | `github.com/go-fork/middleware/redirect` | Infrastructure | URL redirection |
| healthcheck | `github.com/go-fork/middleware/healthcheck` | Infrastructure | Health check endpoints |
| earlydata | `github.com/go-fork/middleware/earlydata` | Advanced | HTTP/2 early data |
| encryptcookie | `github.com/go-fork/middleware/encryptcookie` | Advanced | Cookie encryption |
| envvar | `github.com/go-fork/middleware/envvar` | Advanced | Environment variable injection |
| expvar | `github.com/go-fork/middleware/expvar` | Advanced | Go expvar metrics |
| idempotency | `github.com/go-fork/middleware/idempotency` | Advanced | Idempotent requests |
| rewrite | `github.com/go-fork/middleware/rewrite` | Advanced | URL rewriting |

### Quick Setup Examples

#### Development Stack
```yaml
http:
  middleware:
    recover:
      enabled: true
      enable_stack_trace: true
    logger:
      enabled: true
      enable_colors: true
    cors:
      enabled: true
      allow_origins: ["*"]
    monitor:
      enabled: true
```

#### Production Stack  
```yaml
http:
  middleware:
    recover:
      enabled: true
      enable_stack_trace: false
    helmet:
      enabled: true
    cors:
      enabled: true
      allow_origins: ["https://myapp.com"]
    logger:
      enabled: true
      output: "file"
    compression:
      enabled: true
    cache:
      enabled: true
    limiter:
      enabled: true
```

#### API Stack
```yaml
http:
  middleware:
    recover:
      enabled: true
    keyauth:
      enabled: true
    limiter:
      enabled: true
    bodylimit:
      enabled: true
    compression:
      enabled: true
    requestid:
      enabled: true
```

#### Security-First Stack
```yaml
http:
  middleware:
    recover:
      enabled: true
    helmet:
      enabled: true
    cors:
      enabled: true
    csrf:
      enabled: true
    basicauth:
      enabled: true
    limiter:
      enabled: true
      max: 50
```

## 🔗 Tài liệu liên quan

- **[Configuration System](config.md)** - YAML configuration và validation
- **[Service Provider](service-provider.md)** - DI integration patterns
- **[Router System](router.md)** - Router middleware integration
- **[Context System](context-request-response.md)** - Context manipulation trong middleware
- **[Web Application](web-application.md)** - Application-level middleware setup

---

**Fork HTTP Framework Middleware System** - Build powerful, composable web applications with zero-code middleware configuration! 🚀
