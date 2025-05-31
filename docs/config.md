# Configuration - Cấu hình Framework

Package `fork` cung cấp hệ thống cấu hình linh hoạt và mạnh mẽ cho Fork HTTP Framework. Cấu hình được quản lý thông qua các struct và có thể được load từ nhiều nguồn khác nhau như YAML files, environment variables, hoặc configuration providers.

## Tổng quan

Framework sử dụng pattern configuration struct để quản lý các cài đặt. Các cấu hình chính bao gồm:

- **WebApp Configuration**: Cấu hình tổng quát cho ứng dụng web
- **Graceful Shutdown**: Cấu hình quá trình shutdown an toàn
- **Adapter Configurations**: Cấu hình cho các HTTP adapters khác nhau

## WebApp Configuration

### Struct WebAppConfig

```go
type WebAppConfig struct {
    GracefulShutdown GracefulShutdownConfig `mapstructure:"graceful_shutdown" yaml:"graceful_shutdown"`
}
```

`WebAppConfig` là struct chính chứa tất cả cấu hình cho WebApp. Hiện tại nó tập trung vào graceful shutdown, trong khi các cấu hình khác đã được chuyển sang các middleware packages riêng biệt.

### Graceful Shutdown Configuration

```go
type GracefulShutdownConfig struct {
    Enabled            bool `mapstructure:"enabled" yaml:"enabled"`
    Timeout            int  `mapstructure:"timeout" yaml:"timeout"`
    WaitForConnections bool `mapstructure:"wait_for_connections" yaml:"wait_for_connections"`
    SignalBufferSize   int  `mapstructure:"signal_buffer_size" yaml:"signal_buffer_size"`
    
    // Callback functions
    OnShutdownStart    func()       `mapstructure:"-" yaml:"-"`
    OnShutdownComplete func()       `mapstructure:"-" yaml:"-"`
    OnShutdownError    func(error)  `mapstructure:"-" yaml:"-"`
}
```

#### Các tham số:

- **Enabled**: Bật/tắt graceful shutdown (mặc định: `true`)
- **Timeout**: Thời gian tối đa chờ shutdown (giây, mặc định: `30`)
- **WaitForConnections**: Có chờ tất cả connections kết thúc không (mặc định: `true`)
- **SignalBufferSize**: Kích thước buffer cho signal channel (mặc định: `1`)

#### Callback Functions:

- **OnShutdownStart**: Được gọi khi bắt đầu quá trình shutdown
- **OnShutdownComplete**: Được gọi khi shutdown hoàn thành
- **OnShutdownError**: Được gọi khi có lỗi trong quá trình shutdown

## Configuration File (YAML)

Framework hỗ trợ file cấu hình YAML với cấu trúc hoàn chỉnh. Tham khảo file [`configs/app.example.yaml`](../configs/app.example.yaml):

### HTTP Framework Configuration

```yaml
http:
  # Graceful shutdown configuration
  graceful_shutdown:
    enabled: true
    timeout: 30
    wait_for_connections: true
    signal_buffer_size: 1
  
  # Adapter configuration
  debug: true
  adapter: "http"  # http, fasthttp, http2, quic, unified
  
  # Adapter-specific configurations
  http:
    addr: "localhost"
    port: 7667
    read_timeout: 10s
    write_timeout: 10s
    # ... more HTTP adapter settings
```

### Supported Adapters

File cấu hình hỗ trợ nhiều adapter khác nhau:

1. **Standard HTTP** (`http`): Go's standard net/http
2. **FastHTTP** (`fasthttp`): High-performance HTTP server
3. **HTTP/2** (`http2`): HTTP/2 protocol support
4. **QUIC/HTTP3** (`quic`): HTTP/3 over QUIC protocol
5. **Unified** (`unified`): Multi-protocol support

Mỗi adapter có các cấu hình riêng biệt được định nghĩa trong file example.

## API Reference

### DefaultWebAppConfig()

Trả về cấu hình mặc định cho WebApp:

```go
func DefaultWebAppConfig() *WebAppConfig {
    return &WebAppConfig{
        GracefulShutdown: GracefulShutdownConfig{
            Enabled:            true,
            Timeout:            30,
            WaitForConnections: true,
            SignalBufferSize:   1,
        },
    }
}
```

### LoadConfigFromProvider()

Load cấu hình từ configuration provider:

```go
func LoadConfigFromProvider(manager config.Manager, key string) (*WebAppConfig, error)
```

- **manager**: Configuration manager instance
- **key**: Configuration key (thường là "http")
- **Returns**: WebAppConfig instance hoặc error

### Validate()

Validate cấu hình trước khi sử dụng:

```go
func (c *WebAppConfig) Validate() error
```

## Usage Examples

### Basic Configuration

```go
// Sử dụng cấu hình mặc định
config := fork.DefaultWebAppConfig()

// Tạo WebApp với cấu hình
app := fork.NewWebApp()
app.SetConfig(config)
```

### Load from Configuration Provider

```go
// Load từ configuration manager
configManager := container.MustMake("config").(config.Manager)
appConfig, err := fork.LoadConfigFromProvider(configManager, "http")
if err != nil {
    log.Fatal("Failed to load config:", err)
}

// Validate config
if err := appConfig.Validate(); err != nil {
    log.Fatal("Invalid config:", err)
}

// Áp dụng cấu hình
app.SetConfig(appConfig)
```

### Custom Graceful Shutdown

```go
config := fork.DefaultWebAppConfig()

// Cấu hình graceful shutdown
config.GracefulShutdown.Timeout = 60 // 60 giây
config.GracefulShutdown.OnShutdownStart = func() {
    log.Println("Starting graceful shutdown...")
}
config.GracefulShutdown.OnShutdownComplete = func() {
    log.Println("Shutdown completed successfully")
}
config.GracefulShutdown.OnShutdownError = func(err error) {
    log.Printf("Shutdown error: %v", err)
}

app.SetConfig(config)
```

## Migration Notes

### Moved Configurations

Các cấu hình sau đây đã được chuyển sang middleware packages riêng biệt:

- **MaxRequestBodySize** → `bodylimit` middleware package
- **AllowedMethods** → `method` middleware package  
- **RequestTimeout** → `timeout` middleware package
- **EnableSecurityHeaders** → `helmet` middleware package

### Backward Compatibility

Framework vẫn hỗ trợ các cấu hình cũ trong file YAML nhưng khuyến nghị sử dụng middleware packages mới để có tính linh hoạt cao hơn.

## Best Practices

1. **Sử dụng configuration files**: Định nghĩa cấu hình trong file YAML để dễ quản lý
2. **Validate configuration**: Luôn validate cấu hình trước khi áp dụng
3. **Environment-specific configs**: Sử dụng config files khác nhau cho dev/staging/production
4. **Graceful shutdown callbacks**: Implement callbacks để cleanup resources properly
5. **Security considerations**: Không commit sensitive data trong config files

## Related Files

- [`config.go`](../config.go) - Configuration structs và functions
- [`configs/app.example.yaml`](../configs/app.example.yaml) - Example configuration file
- [`provider.go`](../provider.go) - Service provider implementation
