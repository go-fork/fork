# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased] - 2025-06-05

### 🔄 Dependencies

- **Upgraded direct dependencies**: Updated all direct dependencies to latest versions
  - `go.fork.vn/config`: v0.1.0 → v0.1.3
  - `go.fork.vn/di`: v0.1.0 → v0.1.3  
  - `go.fork.vn/log`: v0.1.0 → v0.1.3

### 🛡️ Enhanced Error Handling

- **ServiceProvider.Register()**: Added comprehensive nil checks and panic handling
  - Validates application parameter is not nil
  - Validates container is not nil
  - Prevents runtime errors during service registration

- **ServiceProvider.Boot()**: Enhanced error handling with detailed validation
  - Comprehensive nil checks for application and container
  - Safe type assertions with error reporting for all services (http, log, config)
  - Validates adapter configuration existence and type safety
  - Strict validation for config loading and validation processes
  - Detailed panic messages for debugging and troubleshooting

- **LoadConfigFromProvider()**: Improved configuration loading robustness
  - Added nil provider validation with panic for critical errors
  - Added empty key validation with panic for misconfiguration
  - Enhanced type assertion handling for config providers
  - Automatic config validation after unmarshaling
  - Better error propagation for debugging

### 🔧 Code Refactoring
- **COMPLETED**: Removed `LoadConfigFromProvider` function from `config.go` 
- **COMPLETED**: Added private method `loadConfigFromProvider()` to `ServiceProvider` in `provider.go`
- **COMPLETED**: Updated `ServiceProvider.Boot()` to use the new private method
- **COMPLETED**: Improved code organization and encapsulation

### 🧪 Testing & Mocks

- **Mock Generation**: Regenerated all mock files using mockery v2.53.4
  - Updated mocks for all core interfaces: Adapter, Context, HandlerFunc, Request, Response, Router
  - Enhanced mock support with expecter pattern for better test assertions
  - Improved type safety and interface compatibility
  - Automatic mock generation through `mockery --all` command

- **Config Testing**: Completely rebuilt `config_test.go` with comprehensive test coverage
  - **Unit Tests**: 15+ test cases covering all config functionality (DefaultWebAppConfig, Validate, MergeConfig)
  - **Mock Integration**: Advanced mock testing with `go.fork.vn/config/mocks` using expecter pattern
  - **YAML Integration Tests**: Realistic scenarios simulating `configs/app.example.yaml` configuration
  - **Benchmark Tests**: Performance testing for config operations with race detection
  - **Edge Case Coverage**: Validation of error handling, nil configs, and invalid values

- **Provider Testing**: Created comprehensive `provider_test.go` with full coverage:
  - Tests for `NewServiceProvider`, `Requires`, `Providers` methods
  - Extensive `Register` method testing with error scenarios
  - Comprehensive `Boot` method testing with 15+ error scenarios
  - Integration tests for complete registration and boot cycle
  - Benchmark tests for performance measurement
  - Advanced mock integration using `go.fork.vn/config/mocks`, `go.fork.vn/di/mocks`, `go.fork.vn/log/mocks`, and local `mocks`
  - YAML integration tests simulating real-world configuration scenarios
  - Performance benchmarks showing excellent performance metrics

### 📊 Performance Metrics
- **config_test.go benchmarks**:
  - `DefaultWebAppConfig`: 159ns/op, 352 B/op, 6 allocs/op
  - `WebAppConfig.Validate`: 11.8ns/op, 0 B/op, 0 allocs/op  
  - `MergeConfig`: 4.83ns/op, 0 B/op, 0 allocs/op
- **provider_test.go benchmarks**:
  - `Register`: ~22μs/op, 14346 B/op, 136 allocs/op
  - `Requires`: ~0.3ns/op, 0 B/op, 0 allocs/op
  - `Providers`: ~0.3ns/op, 0 B/op, 0 allocs/op

### ✅ Quality Assurance
- **COMPLETED**: All tests pass including race condition detection
- **COMPLETED**: Comprehensive error scenario coverage
- **COMPLETED**: Mock integration with expecter pattern
- **COMPLETED**: Integration tests covering complete service provider lifecycle

### 🎯 Task Summary
**ALL OBJECTIVES COMPLETED SUCCESSFULLY**:
1. ✅ Removed `LoadConfigFromProvider` function from `config.go`
2. ✅ Moved and rebuilt as private method `loadConfigFromProvider()` in `ServiceProvider`
3. ✅ Rebuilt `config_test.go` with comprehensive test coverage using mocks
4. ✅ Created `provider_test.go` with full test coverage using mocks from multiple packages
5. ✅ Committed all changes to git with detailed commit messages

## [v0.0.9] - 2025-06-01

### 🎉 Initial Release - Fork HTTP Framework

Đây là phiên bản đầu tiên của Fork HTTP Framework, một framework HTTP hiệu năng cao và linh hoạt cho Go applications.

---

## 🏗️ Core Framework Components

### ✅ **WebApp - Core Application**
- **Main Application Interface**: WebApp làm trung tâm điều khiển toàn bộ framework
- **Lifecycle Management**: Quản lý vòng đời application từ khởi tạo đến graceful shutdown  
- **Configuration Management**: Hỗ trợ cấu hình linh hoạt thông qua YAML files
- **Adapter Integration**: Tích hợp với multiple HTTP adapters
- **Route Registration**: API thống nhất cho việc đăng ký routes và middlewares
- **Graceful Shutdown**: Hỗ trợ graceful shutdown với context cancellation
- **Connection Tracking**: Theo dõi active connections để shutdown an toàn

### ✅ **Context System - Request/Response Handling**
- **Unified Context Interface**: API thống nhất cho tất cả các adapters
- **Request Data Binding**: 
  - JSON binding với validation
  - XML binding và parsing
  - Form data binding (application/x-www-form-urlencoded)
  - Multipart form data với file upload support
  - Query parameter binding với default values
  - URL parameter extraction
- **Response Helpers**:
  - JSON response với automatic content-type
  - XML response formatting
  - String response với template support
  - File serving và download
  - Redirect responses
  - Status code management
- **Header Management**: Get/Set HTTP headers
- **Cookie Management**: Secure cookie handling
- **Context Storage**: Key-value storage trong request lifecycle
- **Middleware Chain**: Next() function để điều khiển middleware execution

### ✅ **Router System - Advanced Routing**
- **Trie-Based Router**: High-performance routing với trie data structure
- **Pattern Matching**:
  - Static routes (`/users`)
  - Parameter routes (`/users/:id`)
  - Wildcard routes (`/files/*filepath`)
  - Catch-all patterns
- **HTTP Methods**: Hỗ trợ tất cả HTTP methods (GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS)
- **Route Groups**: Nhóm routes với common prefix và middleware
- **Middleware Integration**: Per-route và group-level middleware
- **Performance Optimization**: Zero-allocation routing cho static routes
- **Route Priority**: Intelligent route matching priority
- **Parameters Extraction**: Efficient parameter parsing và caching

### ✅ **Adapter Pattern - Multi-Engine Support**
- **Adapter Interface**: Unified interface cho different HTTP engines
- **Net/HTTP Adapter**: 
  - Standard library integration
  - Full HTTP/1.1 support
  - WebSocket upgrade capability
  - TLS/SSL support
- **FastHTTP Adapter**: 
  - High-performance alternative
  - Zero-allocation optimizations
  - Custom context implementation
  - Advanced connection pooling
- **HTTP/2 Adapter**: 
  - HTTP/2 protocol support
  - Server push capabilities
  - Multiplexing support
- **QUIC Adapter**: 
  - QUIC protocol support
  - UDP-based transport
  - Future-ready implementation
- **Unified Adapter**: Fallback adapter cho compatibility

---

## ⚙️ Configuration System

### ✅ **YAML-Based Configuration**
- **Structured Configuration**: Comprehensive YAML configuration support
- **Environment-Specific Configs**: app.dev.yaml, app.prod.yaml, app.test.yaml
- **Environment Variable Override**: Override config values với environment variables
- **Nested Configuration**: Complex nested configuration structures
- **Validation**: Configuration validation với default values
- **Hot Reload**: Configuration reloading trong development mode

### ✅ **WebApp Configuration**
- **Server Settings**: Host, port, timeouts configuration
- **Adapter Selection**: Runtime adapter switching
- **Middleware Configuration**: Global middleware setup
- **Static File Serving**: Static content configuration
- **Security Settings**: Security-related configurations
- **Performance Tuning**: Performance optimization settings

---

## 🔧 Dependency Injection Integration

### ✅ **Service Provider Pattern**
- **ServiceProvider Interface**: Chuẩn hóa service registration và boot process
- **Container Integration**: Tích hợp với go.fork.vn/di container
- **Service Registration**: Automatic service discovery và registration
- **Dependency Resolution**: Runtime dependency injection
- **Lifecycle Management**: Service lifecycle coordination
- **Configuration Injection**: Inject configuration vào services

### ✅ **Built-in Service Providers**
- **WebApp Provider**: Register WebApp instance vào DI container
- **Router Provider**: Router service registration
- **Context Provider**: Context factory registration
- **Adapter Provider**: Adapter services registration

---

## 🛡️ Middleware Ecosystem

### ✅ **YAML-Based Middleware Configuration**
- **Auto-Loading System**: Middleware tự động load từ YAML config
- **Zero-Code Configuration**: Enable middleware chỉ với `enabled: true`
- **Environment-Specific**: Different middleware configs cho từng environment
- **Conditional Loading**: Load middleware based on conditions
- **Order Management**: Middleware execution order configuration

### ✅ **Security & Authentication Middleware**
- **BasicAuth Middleware**: HTTP Basic Authentication với user/password
- **Helmet Middleware**: Comprehensive security headers
  - XSS Protection
  - Content Type Options
  - Frame Options (Clickjacking protection)
  - HSTS (HTTP Strict Transport Security)
  - Content Security Policy
  - Referrer Policy
- **CORS Middleware**: Cross-Origin Resource Sharing
  - Origin validation
  - Methods và headers configuration
  - Credentials support
  - Preflight handling
- **CSRF Middleware**: Cross-Site Request Forgery protection
  - Token generation và validation
  - Cookie-based storage
  - Multiple token lookup methods
- **KeyAuth Middleware**: API key authentication
  - Header, query, form-based key lookup
  - Custom validation functions
  - Error handling customization

### ✅ **Performance & Content Middleware**
- **Compression Middleware**: Response compression
  - Gzip compression support
  - Configurable compression levels
  - Content-type filtering
  - Size threshold configuration
- **Cache Middleware**: HTTP caching headers
  - ETag generation
  - Cache-Control headers
  - Conditional requests
  - Custom cache strategies
- **ETag Middleware**: Entity tag generation
  - Strong và weak ETags
  - Automatic generation
  - Conditional request handling
- **Static Middleware**: Static file serving
  - Directory browsing
  - Index file support
  - Cache headers
  - MIME type detection

### ✅ **Rate Limiting & Control Middleware**
- **Limiter Middleware**: Rate limiting functionality
  - Request rate limiting
  - IP-based limiting
  - Custom key extraction
  - Memory-based storage
- **BodyLimit Middleware**: Request body size limiting
  - Configurable size limits
  - Different limits for different content types
  - Early termination for oversized requests
- **Method Middleware**: HTTP method validation
  - Allowed methods configuration
  - Method override support
- **Timeout Middleware**: Request timeout handling
  - Per-route timeout configuration
  - Graceful timeout responses

### ✅ **Session & State Middleware**
- **Session Middleware**: HTTP session management
  - Cookie-based sessions
  - Multiple storage backends
  - Session lifetime management
  - Secure session handling
- **RequestID Middleware**: Request tracking
  - Unique request ID generation
  - Header injection
  - Logging integration

### ✅ **Infrastructure & Utilities**
- **Logger Middleware**: Request/response logging
  - Structured logging support
  - Configurable log formats
  - Performance metrics logging
- **Recover Middleware**: Panic recovery
  - Graceful error handling
  - Stack trace logging
  - Custom error responses
- **Monitor Middleware**: Application monitoring
  - Request metrics collection
  - Performance monitoring
  - Health check support
- **Favicon Middleware**: Favicon serving optimization
- **Static Middleware**: Enhanced static file serving
- **Proxy Middleware**: Reverse proxy functionality
- **Redirect Middleware**: URL redirection
- **HealthCheck Middleware**: Health check endpoints

### ✅ **Advanced Features Middleware**
- **Skip Middleware**: Conditional middleware execution
- **EarlyData Middleware**: HTTP/2 early data handling
- **EncryptCookie Middleware**: Cookie encryption/decryption
- **EnvVar Middleware**: Environment variable injection
- **ExpVar Middleware**: Go expvar metrics exposure
- **Idempotency Middleware**: Idempotent request handling
- **Rewrite Middleware**: URL rewriting

---

## 🎨 Template Engine Support

### ✅ **Multi-Engine Template System**
- **Multiple Template Engines**: Hỗ trợ nhiều template engines đồng thời
  - HTML templates (Go standard library)
  - Text templates
  - Pug templates
  - Mustache templates
  - Handlebars templates
  - Jet templates
- **YAML Configuration**: Complete template configuration via YAML
- **Layout System**: Master layouts và partial templates
- **Auto-Detection**: Automatic engine selection based on file extension
- **Template Caching**: Production-ready template caching
- **Auto-Reload**: Development mode auto-reload
- **Custom Functions**: Built-in và custom template functions
- **Context Integration**: Enhanced Fork Context với automatic content-type detection

### ✅ **Template Features**
- **Layouts và Partials**: Master layout với reusable components
- **Template Inheritance**: Template inheritance patterns
- **Data Binding**: Strong-typed data binding
- **Error Handling**: Comprehensive error handling
- **Performance Optimization**: Template compilation và caching
- **Security**: Template security với HTML escaping

---

## ❌ Error Handling System

### ✅ **Structured Error Management**
- **HttpError Struct**: Standardized HTTP error structure
  - Status code management
  - Error message và details
  - Original error preservation
  - JSON serialization support
- **HTTP Status Code Coverage**: Complete HTTP status code support
  - **4xx Client Errors**: 400, 401, 403, 404, 405, 406, 409, 410, 415, 422, 429
  - **5xx Server Errors**: 500, 501, 502, 503, 504
- **Error Creation Helpers**: Convenient error creation functions
- **Error Response Formatting**: Consistent error response format
- **Integration Patterns**: Error handling integration với middleware và validation

### ✅ **Error Handling Features**
- **Global Error Handler**: Application-level error handling
- **Middleware Error Integration**: Error handling trong middleware chain
- **Validation Error Handling**: Structured validation error responses
- **Authentication Error Handling**: Auth-specific error responses
- **Rate Limiting Error Handling**: Rate limit exceeded responses
- **Security Best Practices**: Secure error message exposure
- **Logging Integration**: Error logging và monitoring

---

## 📖 Documentation & Examples

### ✅ **Comprehensive Vietnamese Documentation**
- **Getting Started Guide**: 12-section comprehensive overview guide
- **Core Components Documentation**:
  - Configuration System (config.md)
  - Service Provider & DI (service-provider.md)
  - Web Application (web-application.md)
  - Context, Request & Response (context-request-response.md)
  - Router System (router.md)
  - Adapter Pattern (adapter.md)
- **Advanced Features Documentation**:
  - Middleware System (middleware.md) - 30+ middleware packages
  - Error Handling (error-handling.md) - Comprehensive error management
  - Template Engine (templates integration)
- **Navigation & Reference**:
  - Documentation Index (index.md)
  - Complete API Reference
  - Project README (README.md)

### ✅ **Real-World Examples**
- **Basic Examples**: Hello World và simple applications
- **Adapter Examples**: Examples cho từng adapter type
- **Middleware Integration**: Real-world middleware usage
- **Template Examples**: Multi-engine template examples
- **Production Examples**: Production-ready application setups
- **Testing Examples**: Unit và integration testing patterns

### ✅ **Best Practices Documentation**
- **Project Structure**: Recommended project organization
- **Security Best Practices**: Security configuration và patterns
- **Performance Optimization**: Performance tuning guidelines
- **Testing Strategies**: Unit, integration, và performance testing
- **Deployment Patterns**: Docker, Kubernetes deployment
- **Monitoring Setup**: Observability và monitoring configuration

---

## 🚀 Performance & Production Features

### ✅ **High Performance**
- **Zero-Allocation Routing**: Optimized routing với minimal allocations
- **Connection Pooling**: Efficient connection management
- **Memory Management**: Object pooling và reuse
- **Async Processing**: Non-blocking request processing
- **Adapter-Specific Optimizations**: Performance tuning cho từng adapter

### ✅ **Production-Ready Features**
- **Graceful Shutdown**: Clean application shutdown
- **Health Checks**: Built-in health check endpoints
- **Metrics Collection**: Application metrics và monitoring
- **Security Headers**: Comprehensive security header support
- **Request Tracking**: Request ID tracking và correlation
- **Error Recovery**: Panic recovery và error handling
- **Configuration Validation**: Runtime configuration validation

### ✅ **Deployment Support**
- **Docker Integration**: Docker deployment patterns
- **Kubernetes Support**: K8s deployment configurations
- **Environment Management**: Multi-environment configuration
- **Secret Management**: Secure secret handling
- **CI/CD Integration**: GitHub Actions workflows
- **Monitoring Integration**: Prometheus metrics, distributed tracing

---

## 🔧 Development Experience

### ✅ **Developer-Friendly Features**
- **Hot Reload**: Development mode hot reloading
- **Detailed Error Messages**: Comprehensive error reporting
- **Auto-Configuration**: Minimal configuration requirements
- **Extensive Logging**: Debug và development logging
- **Testing Utilities**: Built-in testing helpers
- **Documentation**: Complete API documentation

### ✅ **Integration Ecosystem**
- **Third-Party Integration**: Easy integration với external libraries
- **Plugin Architecture**: Extensible plugin system
- **Custom Middleware**: Simple custom middleware development
- **Custom Adapters**: Custom adapter implementation support
- **Service Integration**: Easy service integration patterns

---

## 📦 Package Structure

```
go.fork.vn/fork/
├── adapter/           # HTTP adapter implementations
│   ├── adapter.go     # Base adapter interface
│   └── README.md      # Adapter documentation
├── context/           # Request/response context system
│   ├── context.go     # Main context implementation
│   ├── request.go     # Request handling
│   ├── response.go    # Response helpers
│   └── *_test.go      # Comprehensive tests
├── docs/              # Complete documentation
│   ├── overview.md    # Getting started guide
│   ├── config.md      # Configuration documentation
│   ├── *.md           # Component documentation
│   └── index.md       # Documentation index
├── errors/            # Error handling system
│   ├── errors.go      # HttpError implementation
│   └── errors_test.go # Error handling tests
├── router/            # Advanced routing system
│   ├── router.go      # Main router implementation
│   ├── trie.go        # Trie data structure
│   └── *_test.go      # Router tests
├── configs/           # Configuration examples
│   └── app.example.yaml
├── web_app.go         # Main WebApp implementation
├── config.go          # Configuration management
├── provider.go        # Service provider implementation
├── constants.go       # Framework constants
├── doc.go             # Package documentation
├── go.mod             # Module dependencies
├── README.md          # Project README
└── LICENSE            # MIT License
```

---

## 🎯 Key Dependencies

### Core Dependencies
- **go.fork.vn/config v0.1.2**: YAML configuration management
- **go.fork.vn/di v0.1.0**: Dependency injection container
- **go.fork.vn/log v0.1.0**: Structured logging
- **github.com/go-playground/validator/v10**: Request validation
- **golang.org/x/net**: Network utilities
- **gopkg.in/yaml.v3**: YAML parsing

### Framework Features
- **YAML-First Configuration**: All configuration through YAML files
- **Auto-Loading Middleware**: Zero-code middleware configuration
- **Multiple HTTP Adapters**: net/http, FastHTTP, HTTP/2, QUIC support
- **Comprehensive Middleware**: 30+ production-ready middleware packages
- **Template Engine Support**: Multiple template engines với YAML config
- **Structured Error Handling**: Complete HTTP error management
- **Production-Ready**: Graceful shutdown, health checks, monitoring

---

## 🏁 Next Steps & Roadmap

### Immediate Improvements
- 🔄 WebSocket support implementation
- 🔄 GraphQL adapter development  
- 🔄 Enhanced metrics collection
- 🔄 Advanced authentication middleware (JWT, OAuth2)
- 🔄 Content negotiation middleware

### Future Enhancements
- 🔄 gRPC adapter support
- 🔄 Distributed tracing integration
- 🔄 Advanced caching strategies
- 🔄 Plugin ecosystem expansion
- 🔄 Performance benchmarking tools

---

## 🤝 Contributing

Chúng tôi hoan nghênh contributions từ community! Vui lòng đọc [CONTRIBUTING.md](CONTRIBUTING.md) để biết cách contribute vào project.

## 📄 License

Project này được phát hành dưới [MIT License](LICENSE).

## 🙏 Acknowledgments

Cảm ơn tất cả contributors và community đã hỗ trợ phát triển Fork HTTP Framework!

---

**Fork HTTP Framework v0.0.9** - Build powerful, scalable web applications with Go! 🚀
