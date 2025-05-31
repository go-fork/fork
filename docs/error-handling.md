# Error Handling System - Fork HTTP Framework

Fork HTTP Framework cung cấp một hệ thống error handling mạnh mẽ và linh hoạt với `HttpError` type, cho phép xử lý lỗi HTTP một cách structured và consistent. System hỗ trợ đầy đủ các HTTP status codes phổ biến với các helper functions tiện lợi.

## 📋 Mục lục

1. [Tổng quan Error System](#-tổng-quan-error-system)
2. [HttpError Structure](#-httperror-structure)
3. [Error Creation Methods](#-error-creation-methods)
4. [HTTP Status Code Coverage](#-http-status-code-coverage)
5. [Usage Patterns](#-usage-patterns)
6. [Error Handling Best Practices](#-error-handling-best-practices)
7. [Integration với Middleware](#-integration-với-middleware)
8. [Custom Error Types](#️-custom-error-types)

## 🏗️ Tổng quan Error System

### HttpError Type

Fork sử dụng `HttpError` struct để đại diện cho tất cả HTTP errors. Type này implement Go's `error` interface và cung cấp thêm metadata về HTTP status codes, messages, và details.

```go
type HttpError struct {
    StatusCode int                     `json:"status_code"`
    Message    string                  `json:"message"`
    Details    map[string]interface{}  `json:"details,omitempty"`
    Err        error                   `json:"-"`
}
```

### Key Features

- **Structured Error Information**: Status code, message, details và underlying error
- **JSON Serializable**: Tự động serialize thành JSON response
- **Error Chain Support**: Implement `Unwrap()` method cho error chaining
- **Type Safety**: Strong typing với compile-time safety
- **Rich Context**: Chi tiết metadata cho debugging và logging
- **Performance Optimized**: Minimal allocation designs

## 🔧 HttpError Structure

### Fields Overview

```go
type HttpError struct {
    // StatusCode - HTTP status code (400, 404, 500, etc.)
    StatusCode int `json:"status_code"`
    
    // Message - Human-readable error message
    Message string `json:"message"`
    
    // Details - Additional metadata (validation errors, context, etc.)
    Details map[string]interface{} `json:"details,omitempty"`
    
    // Err - Underlying error (not serialized to prevent info leakage)
    Err error `json:"-"`
}
```

### Core Methods

#### Error() - Implement Go's error interface
```go
func (e *HttpError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("HTTP Error %d: %s - %s", e.StatusCode, e.Message, e.Err.Error())
    }
    return fmt.Sprintf("HTTP Error %d: %s", e.StatusCode, e.Message)
}
```

#### Unwrap() - Support error chaining
```go
func (e *HttpError) Unwrap() error {
    return e.Err
}
```

### Usage Examples

```go
// Basic usage
err := errors.BadRequest("Invalid input data")
fmt.Println(err.Error()) // "HTTP Error 400: Invalid input data"

// With details
err = errors.NewBadRequest("Validation failed", map[string]interface{}{
    "field": "email",
    "value": "invalid-email",
    "constraint": "must be valid email format",
}, nil)

// With underlying error
originalErr := fmt.Errorf("database connection failed")
err = errors.NewInternalServerError("Database error", nil, originalErr)
fmt.Println(err.Error()) // "HTTP Error 500: Database error - database connection failed"

// Error unwrapping
if errors.Is(err, originalErr) {
    // Handle specific error type
}
```

## 🛠️ Error Creation Methods

### Constructor Patterns

Fork cung cấp hai patterns cho mỗi HTTP status code:

1. **Simple Constructor**: `StatusName(message)`
2. **Full Constructor**: `NewStatusName(message, details, err)`

### Universal Constructors

#### NewHttpError - Complete control
```go
func NewHttpError(statusCode int, message string, details map[string]interface{}, err error) *HttpError
```

#### SimpleHttpError - Basic error
```go
func SimpleHttpError(statusCode int, message string) *HttpError
```

### Examples

```go
// Universal constructors
err1 := errors.NewHttpError(418, "I'm a teapot", nil, nil)
err2 := errors.SimpleHttpError(418, "I'm a teapot")

// Using specific constructors
err3 := errors.BadRequest("Invalid request format")
err4 := errors.NewBadRequest("Validation failed", map[string]interface{}{
    "errors": []string{"email is required", "password too short"},
}, nil)
```

## 📊 HTTP Status Code Coverage

### 4xx Client Errors

#### 400 Bad Request
```go
// Simple form
err := errors.BadRequest("Invalid request format")

// Full form
err := errors.NewBadRequest("Validation failed", map[string]interface{}{
    "validation_errors": []map[string]string{
        {"field": "email", "message": "invalid format"},
        {"field": "password", "message": "too short"},
    },
}, validationErr)
```

#### 401 Unauthorized
```go
// Authentication required
err := errors.Unauthorized("Authentication required")

// With details
err := errors.NewUnauthorized("Invalid token", map[string]interface{}{
    "token_type": "Bearer",
    "expired_at": "2023-12-01T10:00:00Z",
}, nil)
```

#### 403 Forbidden
```go
// Access denied
err := errors.Forbidden("Access denied")

// With permission details
err := errors.NewForbidden("Insufficient permissions", map[string]interface{}{
    "required_permission": "admin:read",
    "user_permissions": []string{"user:read", "user:write"},
}, nil)
```

#### 404 Not Found
```go
// Resource not found
err := errors.NotFound("User not found")

// With resource details
err := errors.NewNotFound("Resource not found", map[string]interface{}{
    "resource_type": "user",
    "resource_id": "12345",
    "available_endpoints": []string{"/api/v1/users", "/api/v1/users/{id}"},
}, nil)
```

#### 405 Method Not Allowed
```go
// Method not supported
err := errors.MethodNotAllowed("Method not allowed")

// With allowed methods
err := errors.NewMethodNotAllowed("Method not allowed", map[string]interface{}{
    "allowed_methods": []string{"GET", "POST", "PUT"},
    "requested_method": "DELETE",
}, nil)
```

#### 406 Not Acceptable
```go
// Content negotiation failed
err := errors.NotAcceptable("Not acceptable")

// With supported formats
err := errors.NewNotAcceptable("Content type not supported", map[string]interface{}{
    "supported_types": []string{"application/json", "application/xml"},
    "requested_type": "text/plain",
}, nil)
```

#### 409 Conflict
```go
// Resource conflict
err := errors.Conflict("Resource already exists")

// With conflict details
err := errors.NewConflict("Email already registered", map[string]interface{}{
    "conflicting_field": "email",
    "conflicting_value": "user@example.com",
    "existing_resource_id": "67890",
}, nil)
```

#### 410 Gone
```go
// Resource permanently deleted
err := errors.Gone("Resource no longer available")

// With historical info
err := errors.NewGone("API version deprecated", map[string]interface{}{
    "deprecated_version": "v1",
    "current_version": "v2",
    "migration_guide": "https://docs.example.com/migration",
}, nil)
```

#### 415 Unsupported Media Type
```go
// Media type not supported
err := errors.UnsupportedMediaType("Unsupported media type")

// With supported types
err := errors.NewUnsupportedMediaType("Media type not supported", map[string]interface{}{
    "supported_types": []string{"application/json", "multipart/form-data"},
    "received_type": "application/xml",
}, nil)
```

#### 422 Unprocessable Entity
```go
// Semantic validation failed
err := errors.UnprocessableEntity("Validation failed")

// With validation details
err := errors.NewUnprocessableEntity("Business rule violation", map[string]interface{}{
    "rule": "start_date_must_be_before_end_date",
    "start_date": "2023-12-31",
    "end_date": "2023-01-01",
}, nil)
```

#### 429 Too Many Requests
```go
// Rate limit exceeded
err := errors.TooManyRequests("Rate limit exceeded")

// With rate limit info
err := errors.NewTooManyRequests("Rate limit exceeded", map[string]interface{}{
    "limit": 100,
    "window": "1h",
    "reset_at": "2023-12-01T11:00:00Z",
    "retry_after": 3600,
}, nil)
```

### 5xx Server Errors

#### 500 Internal Server Error
```go
// Generic server error
err := errors.InternalServerError("Internal server error")

// With error details (for logging, not exposed to client)
err := errors.NewInternalServerError("Database connection failed", map[string]interface{}{
    "error_id": "ERR-2023-001",
    "timestamp": time.Now().UTC(),
}, dbErr)
```

#### 501 Not Implemented
```go
// Feature not implemented
err := errors.NotImplemented("Feature not implemented")

// With implementation timeline
err := errors.NewNotImplemented("Feature under development", map[string]interface{}{
    "feature": "advanced_search",
    "planned_release": "v2.1.0",
    "eta": "Q2 2024",
}, nil)
```

#### 502 Bad Gateway
```go
// Upstream server error
err := errors.BadGateway("Bad gateway")

// With upstream details
err := errors.NewBadGateway("Upstream service unavailable", map[string]interface{}{
    "upstream_service": "payment-service",
    "upstream_status": "503",
    "retry_after": 30,
}, nil)
```

#### 503 Service Unavailable
```go
// Service temporarily unavailable
err := errors.ServiceUnavailable("Service unavailable")

// With maintenance info
err := errors.NewServiceUnavailable("Scheduled maintenance", map[string]interface{}{
    "maintenance_start": "2023-12-01T02:00:00Z",
    "maintenance_end": "2023-12-01T04:00:00Z",
    "retry_after": 7200,
}, nil)
```

#### 504 Gateway Timeout
```go
// Upstream timeout
err := errors.GatewayTimeout("Gateway timeout")

// With timeout details
err := errors.NewGatewayTimeout("Upstream service timeout", map[string]interface{}{
    "timeout_duration": "30s",
    "upstream_service": "analytics-service",
    "retry_recommended": true,
}, nil)
```

## 🎯 Usage Patterns

### Handler Error Handling

```go
func getUserHandler(c fork.Context) error {
    userID := c.Param("id")
    if userID == "" {
        return errors.BadRequest("User ID is required")
    }
    
    user, err := userService.GetUser(userID)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return errors.NotFound("User not found")
        }
        return errors.NewInternalServerError("Failed to fetch user", map[string]interface{}{
            "user_id": userID,
            "error_id": generateErrorID(),
        }, err)
    }
    
    return c.JSON(200, user)
}
```

### Validation Error Handling

```go
func createUserHandler(c fork.Context) error {
    var req CreateUserRequest
    if err := c.Bind(&req); err != nil {
        return errors.NewBadRequest("Invalid request format", map[string]interface{}{
            "error_type": "json_parse_error",
            "details": err.Error(),
        }, err)
    }
    
    if validationErrors := validateCreateUserRequest(req); len(validationErrors) > 0 {
        return errors.NewUnprocessableEntity("Validation failed", map[string]interface{}{
            "validation_errors": validationErrors,
        }, nil)
    }
    
    user, err := userService.CreateUser(req)
    if err != nil {
        if isDuplicateKeyError(err) {
            return errors.NewConflict("User already exists", map[string]interface{}{
                "conflicting_field": "email",
                "value": req.Email,
            }, err)
        }
        return errors.NewInternalServerError("Failed to create user", nil, err)
    }
    
    return c.JSON(201, user)
}

func validateCreateUserRequest(req CreateUserRequest) []map[string]string {
    var errors []map[string]string
    
    if req.Email == "" {
        errors = append(errors, map[string]string{
            "field": "email",
            "message": "email is required",
        })
    } else if !isValidEmail(req.Email) {
        errors = append(errors, map[string]string{
            "field": "email", 
            "message": "invalid email format",
        })
    }
    
    if len(req.Password) < 8 {
        errors = append(errors, map[string]string{
            "field": "password",
            "message": "password must be at least 8 characters",
        })
    }
    
    return errors
}
```

### Authentication & Authorization

```go
func authMiddleware(c fork.Context) error {
    token := c.Get("Authorization")
    if token == "" {
        return errors.NewUnauthorized("Authentication required", map[string]interface{}{
            "auth_methods": []string{"Bearer token", "API key"},
        }, nil)
    }
    
    claims, err := validateToken(token)
    if err != nil {
        return errors.NewUnauthorized("Invalid token", map[string]interface{}{
            "error": "token_validation_failed",
            "hint": "Check token format and expiration",
        }, err)
    }
    
    c.Set("user", claims.User)
    return c.Next()
}

func adminRequiredMiddleware(c fork.Context) error {
    user := c.Get("user").(User)
    if !user.IsAdmin {
        return errors.NewForbidden("Admin access required", map[string]interface{}{
            "required_role": "admin",
            "user_role": user.Role,
        }, nil)
    }
    
    return c.Next()
}
```

### Rate Limiting

```go
func rateLimitMiddleware(limiter *RateLimiter) fork.Middleware {
    return func(c fork.Context) error {
        clientID := c.IP()
        
        allowed, resetTime, err := limiter.Allow(clientID)
        if err != nil {
            return errors.NewInternalServerError("Rate limiter error", nil, err)
        }
        
        if !allowed {
            return errors.NewTooManyRequests("Rate limit exceeded", map[string]interface{}{
                "limit": limiter.Limit,
                "window": limiter.Window.String(),
                "reset_at": resetTime.Unix(),
                "retry_after": int(time.Until(resetTime).Seconds()),
            }, nil)
        }
        
        return c.Next()
    }
}
```

## 🔧 Error Handling Best Practices

### 1. Consistent Error Format

```go
// Good: Consistent error structure
func getUser(id string) (*User, error) {
    if id == "" {
        return nil, errors.BadRequest("User ID is required")
    }
    
    user, err := db.FindUser(id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, errors.NotFound("User not found")
        }
        return nil, errors.NewInternalServerError("Database error", nil, err)
    }
    
    return user, nil
}
```

### 2. Error Context & Details

```go
// Good: Rich error context
func processPayment(paymentID string) error {
    payment, err := getPayment(paymentID)
    if err != nil {
        return errors.NewBadRequest("Invalid payment", map[string]interface{}{
            "payment_id": paymentID,
            "error_type": "payment_not_found",
            "suggested_action": "verify payment ID and try again",
        }, err)
    }
    
    if payment.Status != "pending" {
        return errors.NewConflict("Payment already processed", map[string]interface{}{
            "payment_id": paymentID,
            "current_status": payment.Status,
            "allowed_status": "pending",
        }, nil)
    }
    
    return nil
}
```

### 3. Error Logging

```go
func errorHandler(err error, c fork.Context) {
    var httpErr *errors.HttpError
    if errors.As(err, &httpErr) {
        // Log error với context
        logger.WithFields(map[string]interface{}{
            "status_code": httpErr.StatusCode,
            "message": httpErr.Message,
            "path": c.Path(),
            "method": c.Method(),
            "user_id": c.Get("user_id"),
            "request_id": c.Get("request_id"),
            "details": httpErr.Details,
        }).Error("HTTP error occurred")
        
        // Separate logging cho underlying error
        if httpErr.Err != nil {
            logger.WithError(httpErr.Err).Error("Underlying error")
        }
        
        // Response to client
        c.JSON(httpErr.StatusCode, httpErr)
    } else {
        // Generic error handling
        logger.WithError(err).Error("Unhandled error")
        c.JSON(500, errors.InternalServerError("Internal server error"))
    }
}
```

### 4. Security Considerations

```go
func secureErrorHandler(err error, c fork.Context) {
    var httpErr *errors.HttpError
    if errors.As(err, &httpErr) {
        // Remove sensitive info for production
        if isProduction() && httpErr.StatusCode >= 500 {
            // Log full error internally
            logger.Error("Internal error", "error", err, "details", httpErr.Details)
            
            // Return sanitized error to client
            sanitized := errors.InternalServerError("Internal server error")
            c.JSON(sanitized.StatusCode, sanitized)
            return
        }
        
        c.JSON(httpErr.StatusCode, httpErr)
    }
}
```

## 🔌 Integration với Middleware

### Error Handling Middleware

```go
func ErrorHandlingMiddleware() fork.Middleware {
    return func(c fork.Context) error {
        err := c.Next()
        if err != nil {
            var httpErr *errors.HttpError
            if errors.As(err, &httpErr) {
                return c.JSON(httpErr.StatusCode, httpErr)
            }
            
            // Convert generic error to HttpError
            httpErr = errors.NewInternalServerError("Internal server error", map[string]interface{}{
                "error_id": generateErrorID(),
                "timestamp": time.Now().UTC(),
            }, err)
            
            return c.JSON(httpErr.StatusCode, httpErr)
        }
        return nil
    }
}
```

### Recovery Middleware Integration

```go
func RecoveryWithErrorHandling() fork.Middleware {
    return func(c fork.Context) error {
        defer func() {
            if r := recover(); r != nil {
                var err error
                switch x := r.(type) {
                case string:
                    err = errors.NewInternalServerError("Panic occurred", map[string]interface{}{
                        "panic_message": x,
                        "stack_trace": string(debug.Stack()),
                    }, nil)
                case error:
                    err = errors.NewInternalServerError("Panic occurred", map[string]interface{}{
                        "panic_error": x.Error(),
                        "stack_trace": string(debug.Stack()),
                    }, x)
                default:
                    err = errors.InternalServerError("Unknown panic occurred")
                }
                
                c.JSON(500, err)
            }
        }()
        
        return c.Next()
    }
}
```

## 🛠️ Custom Error Types

### Domain-Specific Errors

```go
// Business logic errors
func NewBusinessRuleViolation(rule string, context map[string]interface{}) *errors.HttpError {
    return errors.NewUnprocessableEntity("Business rule violation", map[string]interface{}{
        "error_type": "business_rule_violation",
        "rule": rule,
        "context": context,
    }, nil)
}

// Validation errors
func NewValidationError(field, message string) *errors.HttpError {
    return errors.NewBadRequest("Validation failed", map[string]interface{}{
        "error_type": "validation_error",
        "field": field,
        "message": message,
    }, nil)
}

// Service availability errors
func NewServiceUnavailableError(service string, retryAfter int) *errors.HttpError {
    return errors.NewServiceUnavailable("Service temporarily unavailable", map[string]interface{}{
        "service": service,
        "retry_after": retryAfter,
        "error_type": "service_unavailable",
    }, nil)
}
```

### Error Factory Pattern

```go
type ErrorFactory struct {
    service string
    version string
}

func NewErrorFactory(service, version string) *ErrorFactory {
    return &ErrorFactory{
        service: service,
        version: version,
    }
}

func (f *ErrorFactory) BadRequest(message string, details map[string]interface{}) *errors.HttpError {
    if details == nil {
        details = make(map[string]interface{})
    }
    details["service"] = f.service
    details["version"] = f.version
    details["timestamp"] = time.Now().UTC()
    
    return errors.NewBadRequest(message, details, nil)
}

func (f *ErrorFactory) InternalError(message string, err error) *errors.HttpError {
    return errors.NewInternalServerError(message, map[string]interface{}{
        "service": f.service,
        "version": f.version,
        "error_id": generateErrorID(),
        "timestamp": time.Now().UTC(),
    }, err)
}

// Usage
errFactory := NewErrorFactory("user-service", "v1.2.0")
err := errFactory.BadRequest("Invalid user data", map[string]interface{}{
    "validation_errors": validationErrors,
})
```

## 🔗 Tài liệu liên quan

- **[Context System](context-request-response.md)** - Context error handling patterns
- **[Middleware System](middleware.md)** - Error handling middleware
- **[Web Application](web-application.md)** - Application-level error handling
- **[Configuration](config.md)** - Error logging và monitoring config

---

**Fork HTTP Framework Error Handling** - Build robust applications with comprehensive error handling! 🚀
