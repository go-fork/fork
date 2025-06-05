# Hướng dẫn Di chuyển - v0.1.0

## Tổng quan
Hướng dẫn này giúp bạn di chuyển từ phiên bản trước đó lên v0.1.0.

## Yêu cầu Tiên quyết
- Go 1.23 hoặc mới hơn
- Phiên bản trước đã được cài đặt

## Checklist Di chuyển Nhanh
- [ ] Cập nhật dependencies lên v0.1.3
- [ ] Cập nhật ServiceProvider implementations
- [ ] Cập nhật container access patterns
- [ ] Kiểm tra error handling code
- [ ] Chạy tests để đảm bảo compatibility
- [ ] Cập nhật mock objects nếu sử dụng testing

## Breaking Changes (Thay đổi Phá vỡ)

### API Changes
#### ServiceProvider Interface Changes
```go
// Cách cũ (phiên bản trước)
func (p *ServiceProvider) Register(app interface{}) {
    // app là generic interface{}
    container := app.(*SomeAppType).Container()
}

func (p *ServiceProvider) Boot(app interface{}) {
    // app là generic interface{}
}

// Cách mới (v0.1.0)
func (p *ServiceProvider) Register(app di.Application) {
    // app là di.Application interface với type safety
    container := app.Container() // Trả về di.Container interface
}

func (p *ServiceProvider) Boot(app di.Application) {
    // app là di.Application interface với enhanced validation
}
```

#### Container Access Pattern Changes
```go
// Cách cũ
type ServiceProvider struct {
    container *di.Container // Pointer to interface
}

// Cách mới
type ServiceProvider struct {
    // Không cần store container, lấy từ di.Application
}

// Container operations
// Cũ: container.(*ConcreteType).Method()
// Mới: container.Method() // Direct interface method calls
```

#### Enhanced Error Handling
```go
// Cũ: Silent failures hoặc basic error returns
func someFunction(param interface{}) error {
    if param == nil {
        return errors.New("param is nil")
    }
    // ...
}

// Mới: Panic-based validation cho critical errors
func someFunction(param interface{}) error {
    if param == nil {
        panic("critical error: param cannot be nil - check your configuration")
    }
    // ...
}
```

### Dependency Changes
Dependencies đã được nâng cấp với breaking interface changes:

```go
// go.fork.vn/config v0.1.0 → v0.1.3
// Enhanced interface với better type safety

// go.fork.vn/di v0.1.0 → v0.1.3  
// Container interface redesign, Application interface enhanced

// go.fork.vn/log v0.1.0 → v0.1.3
// Better Fork framework integration
```

## Step-by-Step Migration (Các bước Di chuyển Chi tiết)

### Step 1: Cập nhật Dependencies
```bash
go get go.fork.vn/fork@v0.1.0
go get go.fork.vn/config@v0.1.3
go get go.fork.vn/di@v0.1.3
go get go.fork.vn/log@v0.1.3
go mod tidy
```

### Step 2: Cập nhật ServiceProvider Implementation
```go
// Trước khi update
type MyServiceProvider struct {
    container *di.Container
}

func (p *MyServiceProvider) Register(app interface{}) {
    container := app.(*SomeType).GetContainer()
    // ... registration logic
}

func (p *MyServiceProvider) Boot(app interface{}) {
    // ... boot logic
}

// Sau khi update
type MyServiceProvider struct {
    // Không cần store container
}

func (p *MyServiceProvider) Register(app di.Application) {
    // Validate parameters
    if app == nil {
        panic("MyServiceProvider.Register: application cannot be nil")
    }
    
    container := app.Container()
    if container == nil {
        panic("MyServiceProvider.Register: container cannot be nil")
    }
    // ... registration logic with enhanced error handling
}

func (p *MyServiceProvider) Boot(app di.Application) {
    // Enhanced validation và type-safe operations
    if app == nil {
        panic("MyServiceProvider.Boot: application cannot be nil")
    }
    // ... boot logic
}
```

### Step 3: Cập nhật Container Operations
```go
// Trước
container := app.(*AppType).Container()
service := container.(*ConcreteContainer).Make("service-name")

// Sau
container := app.Container() // Interface method
service := container.Make("service-name") // Interface method
```

### Step 4: Cập nhật Error Handling
```go
// Cải thiện error handling để tận dụng enhanced validation
func yourFunction() {
    // Thay vì silent error returns, sử dụng panic cho critical errors
    config := loadConfig()
    if config == nil {
        panic("critical: configuration failed to load - check your setup")
    }
}
```

### Step 5: Cập nhật Testing Code
```go
// Sử dụng new mock objects
import (
    "go.fork.vn/config/mocks"
    "go.fork.vn/di/mocks"
    "go.fork.vn/log/mocks"
    localMocks "your-project/mocks"
)

func TestYourFunction(t *testing.T) {
    // Sử dụng expecter pattern
    mockApp := mocks.NewMockApplication(t)
    mockContainer := mocks.NewMockContainer(t)
    
    mockApp.EXPECT().Container().Return(mockContainer)
    // ... test logic với enhanced mock support
}
```

### Step 6: Chạy Tests
```bash
go test ./... -race
```

## Common Issues và Solutions (Vấn đề Thường gặp và Giải pháp)

### Issue 1: ServiceProvider Interface Mismatch
**Problem**: `cannot use *MyServiceProvider (type *MyServiceProvider) as type di.ServiceProvider`  
**Solution**: Cập nhật method signatures để sử dụng `di.Application` parameter

### Issue 2: Container Access Error
**Problem**: `container.(*ConcreteType) undefined`  
**Solution**: Sử dụng interface methods trực tiếp thay vì type assertions

### Issue 3: Mock Object Not Found
**Problem**: `package mocks is not imported`  
**Solution**: Import correct mock packages và regenerate mocks với mockery v2.53.4

### Issue 4: Panic Instead of Error
**Problem**: `panic: critical error message`  
**Solution**: Đây là intended behavior cho critical errors. Kiểm tra configuration và setup

## Testing Migration (Di chuyển Testing)

### Mock Generation
```bash
# Regenerate mocks với mockery v2.53.4
go install github.com/vektra/mockery/v2@v2.53.4
mockery --all
```

### Enhanced Test Patterns
```go
// Sử dụng expecter pattern cho cleaner tests
mockService.EXPECT().Method(mock.Anything).Return(expectedResult).Times(1)
```

## Getting Help (Nhận Trợ giúp)
- Kiểm tra [documentation](https://pkg.go.dev/go.fork.vn/fork@v0.1.0)
- Tìm [existing issues](https://github.com/go-fork/fork/issues)
- Tạo [new issue](https://github.com/go-fork/fork/issues/new) nếu cần

## Rollback Instructions (Hướng dẫn Rollback)
Nếu bạn cần rollback:

```bash
go get go.fork.vn/fork@v0.0.9
go get go.fork.vn/config@v0.1.0
go get go.fork.vn/di@v0.1.0
go get go.fork.vn/log@v0.1.0
go mod tidy
```

---
**Cần Trợ giúp?** Đừng ngần ngại tạo issue hoặc discussion trên GitHub.
