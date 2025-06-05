# Nhật ký Thay đổi (Changelog)

Tất cả các thay đổi đáng chú ý của dự án này sẽ được ghi lại trong tệp này.

Định dạng dựa trên [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
và dự án này tuân thủ [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [v0.1.0] - 2025-06-05

### Added

#### docs: Comprehensive documentation restructure
- Professional enterprise-grade README.md with visual design and architecture diagrams
- Enhanced documentation structure from 426 to 959 lines with detailed technical content
- Mermaid architecture diagrams showcasing framework flow and component relationships
- Production deployment patterns with Docker, docker-compose and graceful shutdown examples
- Performance benchmarks and comprehensive testing framework documentation
- Professional Vietnamese technical terminology throughout documentation

#### docs: Enhanced component documentation
- **Router Documentation (docs/router.md)**: Complete restructure (620 lines) based on actual source code
  - Documented actual Router interface with methods: Handle, Group, Use, Static, Routes, ServeHTTP, Find
  - Added DefaultRouter implementation with real struct fields
  - Documented TrieNode and RouteTrie structures from source code
  - Comprehensive route pattern matching and performance optimization
- **Adapter Documentation (docs/adapter.md)**: Complete restructure (1038 lines) based on actual source code
  - Documented actual Adapter interface with methods: Name, Serve, RunTLS, ServeHTTP, HandleFunc, Use, SetHandler, Shutdown
  - Added implementation examples and framework integration patterns

#### test: Comprehensive testing framework
- **Mock Generation**: Regenerated all mock files using mockery v2.53.4
  - Updated mocks for all core interfaces: Adapter, Context, HandlerFunc, Request, Response, Router
  - Enhanced mock support with expecter pattern for better test assertions
  - Improved type safety and interface compatibility
- **WebApp Test Suite**: Added complete test coverage for WebApp functionality (`web_app_test.go`, 746 lines)
  - 30+ test functions covering all WebApp functionality
  - Core features: HTTP methods, middleware, routing, error handling, configuration
  - Advanced features: Router grouping, parameter handling, context management
  - Concurrency tests for thread safety and connection tracking
  - Performance tests with 3 benchmark functions
- **Config Testing**: Complete rebuild of `config_test.go` with comprehensive coverage
  - 15+ test cases covering all config functionality (DefaultWebAppConfig, Validate, MergeConfig)
  - Advanced mock integration with `go.fork.vn/config/mocks` using expecter pattern
  - YAML integration tests mimicking `configs/app.example.yaml` configuration
  - Benchmark tests for performance measurement with race detection
- **Provider Testing**: Comprehensive `provider_test.go` with full coverage
  - Tests for `NewServiceProvider`, `Requires`, `Providers` methods
  - Extensive `Register` method testing with error scenarios
  - Comprehensive `Boot` method testing with 15+ error scenarios
  - Integration tests for complete registration and boot cycle
  - Advanced mock integration using multiple package mocks

### Changed

#### deps: Dependency updates
- `go.fork.vn/config`: v0.1.0 → v0.1.3
- `go.fork.vn/di`: v0.1.0 → v0.1.3  
- `go.fork.vn/log`: v0.1.0 → v0.1.3

#### refactor: Code organization improvements
- Moved `LoadConfigFromProvider` function from `config.go` to private method `loadConfigFromProvider()` in `ServiceProvider`
- Updated `ServiceProvider.Boot()` to use new private method
- Improved code organization and encapsulation

### Fixed

#### fix: Enhanced error handling
- **ServiceProvider.Register()**: Added comprehensive nil checks and panic handling
  - Validate application parameter is not nil
  - Validate container is not nil
  - Prevent runtime errors during service registration
- **ServiceProvider.Boot()**: Enhanced error handling with detailed validation
  - Comprehensive nil checks for application and container
  - Safe type assertions with error reporting for all services (http, log, config)
  - Validate adapter configuration existence and type safety
  - Strict validation for configuration loading and validation process
  - Detailed panic messages for debugging and troubleshooting
- **LoadConfigFromProvider()**: Improved robustness of configuration loading
  - Added nil provider validation with panic for critical errors
  - Added empty key validation with panic for misconfiguration
  - Enhanced type assertion handling for config providers
  - Automatic config validation after unmarshaling
  - Better error propagation for debugging

### Performance (Hiệu suất)

#### perf: Benchmark results
- **config_test.go benchmarks**:
  - `DefaultWebAppConfig`: 159ns/op, 352 B/op, 6 allocs/op
  - `WebAppConfig.Validate`: 11.8ns/op, 0 B/op, 0 allocs/op  
  - `MergeConfig`: 4.83ns/op, 0 B/op, 0 allocs/op
- **provider_test.go benchmarks**:
  - `Register`: ~22μs/op, 14346 B/op, 136 allocs/op
  - `Requires`: ~0.3ns/op, 0 B/op, 0 allocs/op
  - `Providers`: ~0.3ns/op, 0 B/op, 0 allocs/op

## [v0.0.9] - 2025-06-01

### Added

#### feat: Core Framework Components

##### WebApp - Core Application
- Main Application Interface: WebApp as central controller for entire framework
- Lifecycle Management: Application lifecycle from initialization to graceful shutdown
- Configuration Management: Flexible configuration support through YAML files
- Adapter Integration: Integration with multiple HTTP adapters
- Route Registration: Unified API for registering routes and middlewares
- Graceful Shutdown: Support for graceful shutdown with context cancellation
- Connection Tracking: Track active connections for safe shutdown

##### Context System - Request/Response Handling
- Unified Context Interface: Unified API for all adapters
- **Request Data Binding**:
  - JSON binding with validation
  - XML binding and parsing
  - Form data binding (application/x-www-form-urlencoded)
  - Multipart form data with file upload support
  - Query parameter binding with default values
  - URL parameter extraction
- **Response Helpers**:
  - JSON response with automatic content-type
  - XML response formatting
  - String response with template support
  - File serving and download
  - Redirect responses
  - Status code management
- Header Management: Get/Set HTTP headers
- Cookie Management: Secure cookie handling
- Context Storage: Key-value storage in request lifecycle
- Middleware Chain: Next() function to control middleware execution

##### Router System - Advanced Routing
- Trie-Based Router: High-performance routing with trie data structure
- **Pattern Matching**:
  - Static routes (`/users`)
  - Parameter routes (`/users/:id`)
  - Wildcard routes (`/files/*filepath`)
  - Catch-all patterns
- HTTP Methods: Support for all HTTP methods (GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS)
- Route Groups: Group routes with common prefix and middleware
- Middleware Integration: Per-route and group-level middleware
- Performance Optimization: Zero-allocation routing for static routes
- Route Priority: Intelligent route matching priority
- Parameters Extraction: Efficient parameter parsing and caching

##### Adapter Pattern - Multi-Engine Support
- Adapter Interface: Unified interface for different HTTP engines
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
- Unified Adapter: Fallback adapter for compatibility

#### feat: Configuration System

##### YAML-Based Configuration
- Structured Configuration: Comprehensive YAML configuration support
- Environment-Specific Configs: app.dev.yaml, app.prod.yaml, app.test.yaml
- Environment Variable Override: Override config values with environment variables
- Nested Configuration: Complex nested configuration structures
- Validation: Configuration validation with default values
- Hot Reload: Configuration reloading in development mode

##### WebApp Configuration
- Server Settings: Host, port, timeouts configuration
- Adapter Selection: Runtime adapter switching
- Middleware Configuration: Global middleware setup
- Static File Serving: Static content configuration
- Security Settings: Security-related configurations
- Performance Tuning: Performance optimization settings

#### feat: Dependency Injection Integration

##### Service Provider Pattern
- ServiceProvider Interface: Standardize service registration and boot process
- Container Integration: Integration with go.fork.vn/di container
- Service Registration: Automatic service discovery and registration
- Dependency Resolution: Runtime dependency injection
- Lifecycle Management: Service lifecycle coordination
- Configuration Injection: Inject configuration into services

##### Built-in Service Providers
- WebApp Provider: Register WebApp instance into DI container
- Router Provider: Router service registration
- Context Provider: Context factory registration
- Adapter Provider: Adapter services registration

#### feat: Middleware Ecosystem

##### YAML-Based Middleware Configuration
- Auto-Loading System: Middleware automatically load from YAML config
- Zero-Code Configuration: Enable middleware just with `enabled: true`
- Environment-Specific: Different middleware configs for each environment
- Conditional Loading: Load middleware based on conditions
- Order Management: Middleware execution order configuration

##### Security & Authentication Middleware
- **BasicAuth Middleware**: HTTP Basic Authentication with user/password
- **Helmet Middleware**: Comprehensive security headers
  - XSS Protection
  - Content Type Options
  - Frame Options (Clickjacking protection)
  - HSTS (HTTP Strict Transport Security)
  - Content Security Policy
  - Referrer Policy
- **CORS Middleware**: Cross-Origin Resource Sharing
  - Origin validation
  - Methods and headers configuration
  - Credentials support
  - Preflight handling
- **CSRF Middleware**: Cross-Site Request Forgery protection
  - Token generation and validation
  - Cookie-based storage
  - Multiple token lookup methods
- **KeyAuth Middleware**: API key authentication
  - Header, query, form-based key lookup
  - Custom validation functions
  - Error handling customization

##### Performance & Content Middleware
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
  - Strong and weak ETags
  - Automatic generation
  - Conditional request handling
- **Static Middleware**: Static file serving
  - Directory browsing
  - Index file support
  - Cache headers
  - MIME type detection

##### Rate Limiting & Control Middleware
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

##### Session & State Middleware
- **Session Middleware**: HTTP session management
  - Cookie-based sessions
  - Multiple storage backends
  - Session lifetime management
  - Secure session handling
- **RequestID Middleware**: Request tracking
  - Unique request ID generation
  - Header injection
  - Logging integration

##### Infrastructure & Utilities
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
- **Proxy Middleware**: Reverse proxy functionality
- **Redirect Middleware**: URL redirection
- **HealthCheck Middleware**: Health check endpoints

##### Advanced Features Middleware
- **Skip Middleware**: Conditional middleware execution
- **EarlyData Middleware**: HTTP/2 early data handling
- **EncryptCookie Middleware**: Cookie encryption/decryption
- **EnvVar Middleware**: Environment variable injection
- **ExpVar Middleware**: Go expvar metrics exposure
- **Idempotency Middleware**: Idempotent request handling
- **Rewrite Middleware**: URL rewriting

#### feat: Template Engine Support

##### Multi-Engine Template System
- **Multiple Template Engines**: Support for multiple template engines simultaneously
  - HTML templates (Go standard library)
  - Text templates
  - Pug templates
  - Mustache templates
  - Handlebars templates
  - Jet templates
- YAML Configuration: Complete template configuration via YAML
- Layout System: Master layouts and partial templates
- Auto-Detection: Automatic engine selection based on file extension
- Template Caching: Production-ready template caching
- Auto-Reload: Development mode auto-reload
- Custom Functions: Built-in and custom template functions
- Context Integration: Enhanced Fork Context with automatic content-type detection

##### Template Features
- Layouts and Partials: Master layout with reusable components
- Template Inheritance: Template inheritance patterns
- Data Binding: Strong-typed data binding
- Error Handling: Comprehensive error handling
- Performance Optimization: Template compilation and caching
- Security: Template security with HTML escaping

#### feat: Error Handling System

##### Structured Error Management
- **HttpError Struct**: Standardized HTTP error structure
  - Status code management
  - Error message and details
  - Original error preservation
  - JSON serialization support
- **HTTP Status Code Coverage**: Complete HTTP status code support
  - **4xx Client Errors**: 400, 401, 403, 404, 405, 406, 409, 410, 415, 422, 429
  - **5xx Server Errors**: 500, 501, 502, 503, 504
- Error Creation Helpers: Convenient error creation functions
- Error Response Formatting: Consistent error response format
- Integration Patterns: Error handling integration with middleware and validation

##### Error Handling Features
- Global Error Handler: Application-level error handling
- Middleware Error Integration: Error handling in middleware chain
- Validation Error Handling: Structured validation error responses
- Authentication Error Handling: Auth-specific error responses
- Rate Limiting Error Handling: Rate limit exceeded responses
- Security Best Practices: Secure error message exposure
- Logging Integration: Error logging and monitoring

#### docs: Comprehensive Vietnamese Documentation
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

#### docs: Real-World Examples
- Basic Examples: Hello World and simple applications
- Adapter Examples: Examples for each adapter type
- Middleware Integration: Real-world middleware usage
- Template Examples: Multi-engine template examples
- Production Examples: Production-ready application setups
- Testing Examples: Unit and integration testing patterns

#### docs: Best Practices Documentation
- Project Structure: Recommended project organization
- Security Best Practices: Security configuration and patterns
- Performance Optimization: Performance tuning guidelines
- Testing Strategies: Unit, integration, and performance testing
- Deployment Patterns: Docker, Kubernetes deployment
- Monitoring Setup: Observability and monitoring configuration

#### feat: Performance & Production Features

##### High Performance
- Zero-Allocation Routing: Optimized routing with minimal allocations
- Connection Pooling: Efficient connection management
- Memory Management: Object pooling and reuse
- Async Processing: Non-blocking request processing
- Adapter-Specific Optimizations: Performance tuning for each adapter

##### Production-Ready Features
- Graceful Shutdown: Clean application shutdown
- Health Checks: Built-in health check endpoints
- Metrics Collection: Application metrics and monitoring
- Security Headers: Comprehensive security header support
- Request Tracking: Request ID tracking and correlation
- Error Recovery: Panic recovery and error handling
- Configuration Validation: Runtime configuration validation

##### Deployment Support
- Docker Integration: Docker deployment patterns
- Kubernetes Support: K8s deployment configurations
- Environment Management: Multi-environment configuration
- Secret Management: Secure secret handling
- CI/CD Integration: GitHub Actions workflows
- Monitoring Integration: Prometheus metrics, distributed tracing

#### feat: Development Experience

##### Developer-Friendly Features
- Hot Reload: Development mode hot reloading
- Detailed Error Messages: Comprehensive error reporting
- Auto-Configuration: Minimal configuration requirements
- Extensive Logging: Debug and development logging
- Testing Utilities: Built-in testing helpers
- Documentation: Complete API documentation

##### Integration Ecosystem
- Third-Party Integration: Easy integration with external libraries
- Plugin Architecture: Extensible plugin system
- Custom Middleware: Simple custom middleware development
- Custom Adapters: Custom adapter implementation support
- Service Integration: Easy service integration patterns

### Dependencies (Phụ thuộc)

#### Core Dependencies
- **go.fork.vn/config v0.1.2**: YAML configuration management
- **go.fork.vn/di v0.1.0**: Dependency injection container
- **go.fork.vn/log v0.1.0**: Structured logging
- **github.com/go-playground/validator/v10**: Request validation
- **golang.org/x/net**: Network utilities
- **gopkg.in/yaml.v3**: YAML parsing

#### Framework Features
- YAML-First Configuration: All configuration through YAML files
- Auto-Loading Middleware: Zero-code middleware configuration
- Multiple HTTP Adapters: net/http, FastHTTP, HTTP/2, QUIC support
- Comprehensive Middleware: 30+ production-ready middleware packages
- Template Engine Support: Multiple template engines with YAML config
- Structured Error Handling: Complete HTTP error management
- Production-Ready: Graceful shutdown, health checks, monitoring


## 🏁 Next Steps & Roadmap

### Immediate Improvements
- WebSocket support implementation
- GraphQL adapter development  
- Enhanced metrics collection
- Advanced authentication middleware (JWT, OAuth2)
- Content negotiation middleware

### Future Enhancements
- gRPC adapter support
- Distributed tracing integration
- Advanced caching strategies
- Plugin ecosystem expansion
- Performance benchmarking tools

---

## 🤝 Contributing

Chúng tôi hoan nghênh contributions từ community! Vui lòng đọc [CONTRIBUTING.md](CONTRIBUTING.md) để biết cách contribute vào project.

## 📄 License

Project này được phát hành dưới [MIT License](LICENSE).

## 🙏 Acknowledgments

Cảm ơn tất cả contributors và community đã hỗ trợ phát triển Fork HTTP Framework!

---

**Fork HTTP Framework** - Build powerful, scalable web applications with Go! 🚀
