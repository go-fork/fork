# Service Provider - Nhà cung cấp Dịch vụ

Package `fork` cung cấp `ServiceProvider` - một component quan trọng trong hệ thống Dependency Injection (DI) của Fork Framework. Service Provider chịu trách nhiệm đăng ký và khởi tạo HTTP framework vào application container.

## Tổng quan

`ServiceProvider` tuân theo interface `ServiceProvider` của package `di` để tích hợp seamlessly với hệ thống service container. Nó đảm nhận việc:

- Đăng ký HTTP WebApp vào DI container
- Cấu hình và khởi tạo WebApp trong quá trình boot
- Load và validate configuration từ config provider
- Thiết lập HTTP adapter từ configuration

## Architecture

```go
type ServiceProvider struct{}
```

Service Provider implement 2 phương thức chính:
- `Register(app interface{})`: Đăng ký bindings vào container
- `Boot(app interface{})`: Khởi tạo và cấu hình services

## API Reference

### NewServiceProvider()

Tạo một instance mới của HTTP service provider:

```go
func NewServiceProvider() di.ServiceProvider {
    return &ServiceProvider{}
}
```

**Returns:** 
- `di.ServiceProvider`: Provider instance để đăng ký HTTP services

### Register(app interface{})

Đăng ký các binding liên quan đến HTTP framework vào container:

```go
func (p *ServiceProvider) Register(app interface{})
```

**Chức năng:**
- Lấy container từ app instance
- Đăng ký factory function để tạo HTTP WebApp
- Đăng ký alias cho WebApp (`http.webapp` → `http`)

**Implementation:**
```go
func (p *ServiceProvider) Register(app interface{}) {
    if appWithContainer, ok := app.(interface {
        Container() *di.Container
    }); ok {
        c := appWithContainer.Container()

        // Đăng ký factory function
        c.Bind("http", func(container *di.Container) interface{} {
            return NewWebApp()
        })

        // Đăng ký alias
        c.Alias("http.webapp", "http")
    }
}
```

### Boot(app interface{})

Được gọi sau khi tất cả service providers đã được đăng ký. Cấu hình HTTP WebApp:

```go
func (p *ServiceProvider) Boot(app interface{})
```

**Chức năng:**
1. Lấy HTTP WebApp từ container
2. Load configuration từ config provider
3. Validate configuration
4. Thiết lập configuration cho WebApp
5. Load và thiết lập HTTP adapter

**Panics:**
- Nếu không tìm thấy adapter trong config
- Nếu có lỗi critical trong việc load WebApp config

## Configuration Loading

### LoadConfigFromProvider()

Service Provider sử dụng function `LoadConfigFromProvider` để load configuration:

```go
appConfig, err := LoadConfigFromProvider(configManager, "http")
if err != nil {
    logger.Error("Failed to load HTTP WebApp config: " + err.Error())
    appConfig = DefaultWebAppConfig()
}
```

**Fallback Strategy:**
- Nếu load config thất bại → sử dụng default config
- Nếu config validation thất bại → sử dụng default config
- Log errors để debugging

### Configuration Validation

```go
if err := appConfig.Validate(); err != nil {
    logger.Error("Invalid HTTP WebApp config: " + err.Error())
    appConfig = DefaultWebAppConfig()
}
```

## Adapter Management

Service Provider quản lý việc load và thiết lập HTTP adapters:

```go
// Lấy tên adapter từ config
adapterName, ok := configManager.GetString("http.adapter")
if !ok {
    logger.Fatal("HTTP adapter not found in config")
}

// Lấy adapter instance từ container
adapter := c.MustMake("http.adapter." + adapterName).(adapter.Adapter)
```

**Supported Adapters:**
- `http` - Standard Go net/http
- `fasthttp` - FastHTTP implementation
- `http2` - HTTP/2 protocol
- `quic` - QUIC/HTTP3 protocol
- `unified` - Multi-protocol adapter

## Usage Examples

### Basic Setup

```go
package main

import (
    "go.fork.vn/app"
    "go.fork.vn/fork"
)

func main() {
    // Tạo application instance
    application := app.NewApplication()
    
    // Đăng ký HTTP service provider
    httpProvider := fork.NewServiceProvider()
    application.RegisterProvider(httpProvider)
    
    // Boot application
    application.Boot()
    
    // Lấy HTTP WebApp từ container
    httpApp := application.Container().MustMake("http").(*fork.WebApp)
    
    // Sử dụng WebApp...
}
```

### Custom Configuration

```go
func main() {
    application := app.NewApplication()
    
    // Load custom configuration
    configManager := application.Container().MustMake("config").(config.Manager)
    configManager.SetString("http.adapter", "fasthttp")
    configManager.SetInt("http.graceful_shutdown.timeout", 60)
    
    // Đăng ký provider
    application.RegisterProvider(fork.NewServiceProvider())
    application.Boot()
    
    httpApp := application.Container().MustMake("http").(*fork.WebApp)
}
```

### With Middleware Registration

```go
func main() {
    application := app.NewApplication()
    
    // Đăng ký HTTP provider
    application.RegisterProvider(fork.NewServiceProvider())
    
    // Đăng ký middleware providers
    application.RegisterProvider(cors.NewServiceProvider())
    application.RegisterProvider(logger.NewServiceProvider())
    
    application.Boot()
    
    httpApp := application.Container().MustMake("http").(*fork.WebApp)
    
    // Apply middlewares
    httpApp.Use(
        container.MustMake("middleware.cors").(func(ctx forkCtx.Context)),
        container.MustMake("middleware.logger").(func(ctx forkCtx.Context)),
    )
}
```

## Integration with DI Container

### Container Interface Requirements

Service Provider expect app to implement:

```go
interface {
    Container() *di.Container
}
```

### Service Bindings

Provider tạo các bindings sau trong container:

| Key | Type | Description |
|-----|------|-------------|
| `http` | `*fork.WebApp` | HTTP WebApp instance |
| `http.webapp` | `*fork.WebApp` | Alias cho `http` |

### Dependencies

Service Provider yêu cầu các services sau trong container:

| Service | Type | Purpose |
|---------|------|---------|
| `log` | `log.Manager` | Logging operations |
| `config` | `config.Manager` | Configuration management |
| `http.adapter.{name}` | `adapter.Adapter` | HTTP adapter instances |

## Logging

Service Provider thực hiện comprehensive logging:

```go
// Success logging
logger.Info("HTTP WebApp config loaded successfully",
    "graceful_shutdown_enabled", appConfig.GracefulShutdown.Enabled,
    "graceful_shutdown_timeout", appConfig.GracefulShutdown.Timeout,
)

// Error logging
logger.Error("Failed to load HTTP WebApp config: " + err.Error())
logger.Error("Invalid HTTP WebApp config: " + err.Error())

// Fatal logging
logger.Fatal("HTTP adapter not found in config")
```

## Error Handling

### Graceful Degradation

- **Config load failure**: Fall back to default configuration
- **Config validation failure**: Fall back to default configuration
- **Adapter missing**: Fatal error (application cannot continue)

### Error Scenarios

1. **Configuration Provider Unavailable**: Uses default config
2. **Invalid Configuration**: Uses default config và logs warning
3. **Missing Adapter Configuration**: Fatal error
4. **Adapter Load Failure**: Fatal error

## Best Practices

1. **Register Early**: Đăng ký HTTP provider trước các middleware providers
2. **Configuration Validation**: Luôn validate config trong custom providers
3. **Error Handling**: Implement graceful degradation cho non-critical failures
4. **Logging**: Comprehensive logging cho debugging và monitoring
5. **Dependency Management**: Đảm bảo required services available trong container

## Related Files

- [`provider.go`](../provider.go) - Service provider implementation
- [`config.go`](../config.go) - Configuration management
- [`web_app.go`](../web_app.go) - WebApp implementation
- Package `di` - Dependency injection system
- Package `config` - Configuration management system
