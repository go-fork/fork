# Adapter - HTTP Server Adapters

Package `fork/adapter` cung cấp hệ thống adapter mạnh mẽ cho phép Fork Framework hoạt động với nhiều HTTP server implementations khác nhau. Adapter pattern cho phép framework abstract away các chi tiết implementation cụ thể và cung cấp API đồng nhất.

## Tổng quan

Adapter system bao gồm:

- **Adapter Interface**: Định nghĩa contract chung cho tất cả adapters
- **Multiple Implementations**: net/http, fasthttp, HTTP/2, QUIC/HTTP3
- **Configuration Management**: Cấu hình riêng cho từng adapter
- **Performance Optimization**: Tối ưu cho từng server type
- **Unified API**: API nhất quán cho tất cả adapters

## Adapter Interface

### Core Methods

```go
type Adapter interface {
    // Tên của adapter
    Name() string
    
    // Khởi động server
    Serve() error
    
    // Khởi động HTTPS server
    RunTLS(certFile, keyFile string) error
    
    // HTTP handler integration
    ServeHTTP(w http.ResponseWriter, r *http.Request)
    
    // Đăng ký handler function
    HandleFunc(method, path string, handler HandlerFunc)
    
    // Đăng ký middleware
    Use(middleware HandlerFunc)
    
    // Cấu hình
    SetConfig(config interface{}) error
    GetConfig() interface{}
    
    // Lifecycle management
    Start() error
    Stop() error
    Restart() error
}
```

## Available Adapters

### 1. Standard HTTP Adapter (`http`)

Sử dụng Go's standard `net/http` package:

```go
type HTTPAdapter struct {
    server     *http.Server
    config     *HTTPConfig
    router     router.Router
    middleware []HandlerFunc
}
```

**Features:**
- Stable và well-tested
- Full HTTP/1.1 support
- Built-in HTTPS support
- Graceful shutdown
- Connection pooling

**Configuration:**
```yaml
http:
  addr: "localhost"
  port: 7667
  read_timeout: 10s
  write_timeout: 10s
  read_header_timeout: 5s
  idle_timeout: 120s
  max_header_bytes: 1048576
  tls:
    enabled: false
    cert_file: "./certs/server.crt"
    key_file: "./certs/server.key"
```

### 2. FastHTTP Adapter (`fasthttp`)

High-performance HTTP server implementation:

```go
type FastHTTPAdapter struct {
    server     *fasthttp.Server
    config     *FastHTTPConfig
    router     router.Router
    middleware []HandlerFunc
}
```

**Features:**
- High performance (10x faster than net/http)
- Low memory footprint
- Zero-allocation trong hot paths
- Custom context system
- Built-in compression

**Configuration:**
```yaml
fasthttp:
  addr: "localhost"
  port: 7668
  read_timeout: 10s
  write_timeout: 10s
  max_request_body_size: 4194304
  compression: true
  tls:
    enabled: false
    cert_file: "./certs/server.crt"
    key_file: "./certs/server.key"
```

### 3. HTTP/2 Adapter (`http2`)

Full HTTP/2 protocol support:

```go
type HTTP2Adapter struct {
    server     *http.Server
    config     *HTTP2Config
    router     router.Router
    middleware []HandlerFunc
}
```

**Features:**
- HTTP/2 protocol support
- Server push capabilities
- Multiplexing
- Header compression (HPACK)
- Flow control
- h2c support (HTTP/2 over cleartext)

**Configuration:**
```yaml
http2:
  addr: "localhost"
  port: 7669
  max_concurrent_streams: 250
  initial_window_size: 1048576
  max_frame_size: 16384
  h2c: true
  tls:
    enabled: true
    min_version: "1.2"
    max_version: "1.3"
```

### 4. QUIC/HTTP3 Adapter (`quic`)

Modern HTTP/3 over QUIC protocol:

```go
type QUICAdapter struct {
    server     *http3.Server
    config     *QUICConfig
    router     router.Router
    middleware []HandlerFunc
}
```

**Features:**
- HTTP/3 support
- QUIC transport protocol
- 0-RTT connection resumption
- Built-in encryption
- Improved performance over lossy networks
- UDP-based transport

**Configuration:**
```yaml
quic:
  addr: "localhost"
  port: 7670
  handshake_idle_timeout: 5s
  max_idle_timeout: 30s
  allow_0rtt: false
  enable_datagrams: true
  tls:
    enabled: true
    min_version: "1.3"
    max_version: "1.3"
```

### 5. Unified Adapter (`unified`)

Multi-protocol support trong single adapter:

```go
type UnifiedAdapter struct {
    httpServer   *http.Server
    http2Server  *http.Server
    http3Server  *http3.Server
    config       *UnifiedConfig
    router       router.Router
    middleware   []HandlerFunc
}
```

**Features:**
- Simultaneous HTTP/1.1, HTTP/2, HTTP/3
- Protocol detection
- Shared routing và middleware
- Advanced load balancing
- WebSocket support
- Server push cho HTTP/2

## Configuration System

### Base Configuration

Mỗi adapter có base configuration struct:

```go
type BaseConfig struct {
    Addr           string        `yaml:"addr"`
    Port           int           `yaml:"port"`
    ReadTimeout    time.Duration `yaml:"read_timeout"`
    WriteTimeout   time.Duration `yaml:"write_timeout"`
    TLS            TLSConfig     `yaml:"tls"`
}

type TLSConfig struct {
    Enabled     bool     `yaml:"enabled"`
    CertFile    string   `yaml:"cert_file"`
    KeyFile     string   `yaml:"key_file"`
    MinVersion  string   `yaml:"min_version"`
    MaxVersion  string   `yaml:"max_version"`
    CipherSuites []string `yaml:"cipher_suites"`
}
```

### Adapter-Specific Configurations

Each adapter extends base config với specific fields:

```go
// FastHTTP specific
type FastHTTPConfig struct {
    BaseConfig           `yaml:",inline"`
    MaxRequestBodySize   int  `yaml:"max_request_body_size"`
    Compression          bool `yaml:"compression"`
}

// HTTP/2 specific  
type HTTP2Config struct {
    BaseConfig            `yaml:",inline"`
    MaxConcurrentStreams  int `yaml:"max_concurrent_streams"`
    InitialWindowSize     int `yaml:"initial_window_size"`
    MaxFrameSize          int `yaml:"max_frame_size"`
    H2C                   bool `yaml:"h2c"`
}
```

## Usage Examples

### Basic HTTP Adapter

```go
func main() {
    // Tạo HTTP adapter
    config := &adapter.HTTPConfig{
        Addr: "localhost",
        Port: 8080,
        ReadTimeout: 10 * time.Second,
        WriteTimeout: 10 * time.Second,
    }
    
    httpAdapter := http.NewAdapter(config)
    
    // Tạo WebApp với adapter
    app := fork.NewWebApp()
    app.SetAdapter(httpAdapter)
    
    // Routes
    app.GET("/", func(c forkCtx.Context) {
        c.JSON(200, map[string]string{
            "message": "Hello from HTTP adapter",
        })
    })
    
    // Start server
    log.Fatal(app.Run())
}
```

### FastHTTP Adapter

```go
func main() {
    // FastHTTP configuration
    config := &adapter.FastHTTPConfig{
        BaseConfig: adapter.BaseConfig{
            Addr: "localhost",
            Port: 8080,
        },
        MaxRequestBodySize: 4 * 1024 * 1024, // 4MB
        Compression: true,
    }
    
    fastAdapter := fasthttp.NewAdapter(config)
    
    app := fork.NewWebApp()
    app.SetAdapter(fastAdapter)
    
    app.GET("/fast", func(c forkCtx.Context) {
        c.JSON(200, map[string]string{
            "message": "Hello from FastHTTP",
            "adapter": fastAdapter.Name(),
        })
    })
    
    log.Fatal(app.Run())
}
```

### HTTP/2 với TLS

```go
func main() {
    config := &adapter.HTTP2Config{
        BaseConfig: adapter.BaseConfig{
            Addr: "localhost",
            Port: 8443,
            TLS: adapter.TLSConfig{
                Enabled:  true,
                CertFile: "./certs/server.crt",
                KeyFile:  "./certs/server.key",
                MinVersion: "1.2",
                MaxVersion: "1.3",
            },
        },
        MaxConcurrentStreams: 250,
        InitialWindowSize: 1024 * 1024, // 1MB
        H2C: false, // Require TLS
    }
    
    http2Adapter := http2.NewAdapter(config)
    
    app := fork.NewWebApp()
    app.SetAdapter(http2Adapter)
    
    app.GET("/", func(c forkCtx.Context) {
        c.JSON(200, map[string]interface{}{
            "protocol": "HTTP/2",
            "tls": true,
        })
    })
    
    log.Fatal(app.RunTLS("", "")) // Uses config certs
}
```

### Unified Multi-Protocol

```go
func main() {
    config := &adapter.UnifiedConfig{
        PrimaryAddr: "localhost",
        HTTP: adapter.HTTPProtocolConfig{
            Enabled: true,
            Port: 8080,
        },
        HTTP2: adapter.HTTP2ProtocolConfig{
            Enabled: true,
            Port: 8443,
            H2CPort: 8081,
        },
        HTTP3: adapter.HTTP3ProtocolConfig{
            Enabled: true,
            Port: 9443,
        },
        TLS: adapter.TLSConfig{
            Enabled: true,
            CertFile: "./certs/unified.crt",
            KeyFile: "./certs/unified.key",
        },
    }
    
    unifiedAdapter := unified.NewAdapter(config)
    
    app := fork.NewWebApp()
    app.SetAdapter(unifiedAdapter)
    
    app.GET("/", func(c forkCtx.Context) {
        protocol := c.Request().Protocol()
        c.JSON(200, map[string]interface{}{
            "message": "Multi-protocol support",
            "protocol": protocol,
            "secure": c.Request().IsSecure(),
        })
    })
    
    log.Fatal(app.Run())
}
```

### Dynamic Adapter Selection

```go
func createAdapter(adapterType string) adapter.Adapter {
    switch adapterType {
    case "fasthttp":
        return fasthttp.NewAdapter(&adapter.FastHTTPConfig{
            BaseConfig: adapter.BaseConfig{
                Addr: "localhost",
                Port: 8080,
            },
            Compression: true,
        })
        
    case "http2":
        return http2.NewAdapter(&adapter.HTTP2Config{
            BaseConfig: adapter.BaseConfig{
                Addr: "localhost", 
                Port: 8443,
                TLS: adapter.TLSConfig{Enabled: true},
            },
            H2C: true,
        })
        
    default:
        return http.NewAdapter(&adapter.HTTPConfig{
            Addr: "localhost",
            Port: 8080,
        })
    }
}

func main() {
    adapterType := os.Getenv("HTTP_ADAPTER")
    if adapterType == "" {
        adapterType = "http"
    }
    
    selectedAdapter := createAdapter(adapterType)
    
    app := fork.NewWebApp()
    app.SetAdapter(selectedAdapter)
    
    log.Printf("Using %s adapter", selectedAdapter.Name())
    log.Fatal(app.Run())
}
```

## Performance Comparison

### Benchmarks

| Adapter | Requests/sec | Memory Usage | Features |
|---------|-------------|--------------|----------|
| HTTP | 50K | Medium | Standard, Stable |
| FastHTTP | 500K | Low | High Performance |
| HTTP/2 | 100K | Medium | Multiplexing, Push |
| QUIC/HTTP3 | 80K | Medium | Modern, 0-RTT |
| Unified | 200K | High | Multi-protocol |

### Use Cases

- **HTTP**: General purpose, development, legacy systems
- **FastHTTP**: High-performance APIs, microservices
- **HTTP/2**: Modern web applications, server push
- **QUIC/HTTP3**: Mobile apps, poor network conditions
- **Unified**: Enterprise applications, protocol migration

## Advanced Features

### Server Push (HTTP/2)

```go
app.GET("/", func(c forkCtx.Context) {
    // Check for HTTP/2 push support
    if pusher, ok := c.Response().Pusher(); ok {
        // Push CSS file
        pusher.Push("/static/style.css", &http.PushOptions{
            Method: "GET",
            Header: http.Header{
                "Content-Type": []string{"text/css"},
            },
        })
    }
    
    c.HTML(200, "index.html", nil)
})
```

### Connection Upgrading

```go
// WebSocket upgrade
app.GET("/ws", func(c forkCtx.Context) {
    if c.IsWebsocket() {
        // Handle WebSocket connection
        handleWebSocket(c)
    } else {
        c.JSON(400, map[string]string{
            "error": "WebSocket upgrade required",
        })
    }
})
```

### Custom Context Adaptation

```go
type FastHTTPContext struct {
    forkCtx.Context
    fastCtx *fasthttp.RequestCtx
}

func (c *FastHTTPContext) FastHTTPContext() *fasthttp.RequestCtx {
    return c.fastCtx
}
```

## Best Practices

1. **Adapter Selection**: Choose adapter based on requirements
2. **Configuration**: Use environment-specific configs
3. **TLS Setup**: Proper certificate management
4. **Performance Tuning**: Adjust timeouts và buffer sizes
5. **Error Handling**: Implement proper error handling for each adapter
6. **Resource Management**: Monitor memory và connection usage
7. **Security**: Configure TLS properly, validate inputs

## Monitoring và Metrics

### Health Checks

```go
app.GET("/health", func(c forkCtx.Context) {
    adapter := app.GetAdapter()
    
    c.JSON(200, map[string]interface{}{
        "status": "healthy",
        "adapter": adapter.Name(),
        "timestamp": time.Now(),
    })
})
```

### Performance Metrics

```go
app.Use(func(c forkCtx.Context) {
    start := time.Now()
    c.Next()
    
    duration := time.Since(start)
    log.Printf("Request processed in %v by %s adapter", 
        duration, 
        app.GetAdapter().Name())
})
```

## Related Files

- [`adapter/adapter.go`](../adapter/adapter.go) - Adapter interface
- [`adapter/http/`](../adapter/http/) - Standard HTTP adapter
- [`adapter/fasthttp/`](../adapter/fasthttp/) - FastHTTP adapter  
- [`adapter/http2/`](../adapter/http2/) - HTTP/2 adapter
- [`adapter/quic/`](../adapter/quic/) - QUIC/HTTP3 adapter
- [`adapter/unified/`](../adapter/unified/) - Unified adapter
