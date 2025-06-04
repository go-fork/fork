package fork

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"go.fork.vn/fork/adapter"
	forkCtx "go.fork.vn/fork/context"
	"go.fork.vn/fork/router"
)

// WebApp là đối tượng chính của framework, quản lý HTTP server và routing.
// Nó hoạt động như một wrapper cho router và adapter, cung cấp API đồng nhất
// để xử lý HTTP requests, đăng ký routes và middlewares.
// WebApp khác biệt với DI Application - WebApp chỉ quản lý web layer.
type WebApp struct {
	// adapter là implementation của adapter interface để giao tiếp với server HTTP
	adapter adapter.Adapter

	// router quản lý việc đăng ký và điều hướng HTTP routes
	router router.Router

	// middlewares lưu trữ danh sách các middleware functions của WebApp
	middlewares []router.HandlerFunc

	// config lưu trữ cấu hình bảo mật và hiệu suất của WebApp
	config *WebAppConfig

	// mu bảo vệ truy cập đồng thời vào các thuộc tính của application
	mu sync.RWMutex

	// shutdownCtx cho graceful shutdown handling
	shutdownCtx    context.Context
	shutdownCancel context.CancelFunc

	// activeConnections theo dõi số lượng connections đang hoạt động
	activeConnections int32

	// isShuttingDown đánh dấu trạng thái shutdown
	isShuttingDown bool
}

// NewWebApp tạo một instance mới của WebApp.
// Hàm này khởi tạo router và mảng middlewares rỗng.
//
// Returns:
//   - *WebApp: Một WebApp mới đã được khởi tạo
func NewWebApp() *WebApp {
	ctx, cancel := context.WithCancel(context.Background())
	app := &WebApp{
		router:         router.NewRouter(),
		middlewares:    make([]router.HandlerFunc, 0),
		config:         DefaultWebAppConfig(),
		shutdownCtx:    ctx,
		shutdownCancel: cancel,
	}
	return app
}

// SetAdapter thiết lập adapter cho WebApp và cấu hình handler cho adapter.
// Adapter được sử dụng để giao tiếp với server HTTP.
//
// Parameters:
//   - adapter: Adapter cần thiết lập
func (app *WebApp) SetAdapter(adapter adapter.Adapter) {
	app.mu.Lock()
	defer app.mu.Unlock()

	app.adapter = adapter
	if adapter != nil {
		adapter.SetHandler(app.router)
	}
}

// Use thêm một hoặc nhiều middleware vào WebApp.
// Middleware sẽ được thực thi theo thứ tự đăng ký cho mỗi request.
//
// Parameters:
//   - middleware: Danh sách các middleware functions cần thêm
func (app *WebApp) Use(middleware ...router.HandlerFunc) {
	app.mu.Lock()
	defer app.mu.Unlock()

	app.middlewares = append(app.middlewares, middleware...)
	app.router.Use(middleware...)
}

// Group tạo một router group mới với prefix đường dẫn.
// Group cho phép tổ chức routes theo cấu trúc thư mục và áp dụng middleware cho một nhóm routes.
//
// Parameters:
//   - prefix: Tiền tố đường dẫn cho group
//
// Returns:
//   - router.Router: Router mới đã được tạo với prefix đã chỉ định
func (app *WebApp) Group(prefix string) router.Router {
	return app.router.Group(prefix)
}

// Static đăng ký một thư mục để phục vụ static files.
// Files trong thư mục này sẽ được phục vụ tại đường dẫn có tiền tố được chỉ định.
//
// Parameters:
//   - prefix: Tiền tố URL để phục vụ static files
//   - root: Đường dẫn tới thư mục chứa static files
func (app *WebApp) Static(prefix, root string) {
	app.router.Static(prefix, root)
}

// GET đăng ký handler cho HTTP GET method.
// HTTP GET thường được sử dụng để truy xuất dữ liệu.
//
// Parameters:
//   - path: Đường dẫn URL để đăng ký handler
//   - handlers: Danh sách các handlers xử lý request
func (app *WebApp) GET(path string, handlers ...router.HandlerFunc) {
	app.router.Handle(MethodGet, path, handlers...)
}

// POST đăng ký handler cho HTTP POST method.
// HTTP POST thường được sử dụng để tạo dữ liệu mới.
//
// Parameters:
//   - path: Đường dẫn URL để đăng ký handler
//   - handlers: Danh sách các handlers xử lý request
func (app *WebApp) POST(path string, handlers ...router.HandlerFunc) {
	app.router.Handle(MethodPost, path, handlers...)
}

// PUT đăng ký handler cho HTTP PUT method.
// HTTP PUT thường được sử dụng để cập nhật toàn bộ dữ liệu.
//
// Parameters:
//   - path: Đường dẫn URL để đăng ký handler
//   - handlers: Danh sách các handlers xử lý request
func (app *WebApp) PUT(path string, handlers ...router.HandlerFunc) {
	app.router.Handle(MethodPut, path, handlers...)
}

// DELETE đăng ký handler cho HTTP DELETE method.
// HTTP DELETE thường được sử dụng để xóa dữ liệu.
//
// Parameters:
//   - path: Đường dẫn URL để đăng ký handler
//   - handlers: Danh sách các handlers xử lý request
func (app *WebApp) DELETE(path string, handlers ...router.HandlerFunc) {
	app.router.Handle(MethodDelete, path, handlers...)
}

// PATCH đăng ký handler cho HTTP PATCH method.
// HTTP PATCH thường được sử dụng để cập nhật một phần dữ liệu.
//
// Parameters:
//   - path: Đường dẫn URL để đăng ký handler
//   - handlers: Danh sách các handlers xử lý request
func (app *WebApp) PATCH(path string, handlers ...router.HandlerFunc) {
	app.router.Handle(MethodPatch, path, handlers...)
}

// HEAD đăng ký handler cho HTTP HEAD method.
// HTTP HEAD tương tự như GET nhưng chỉ trả về headers không có body.
//
// Parameters:
//   - path: Đường dẫn URL để đăng ký handler
//   - handlers: Danh sách các handlers xử lý request
func (app *WebApp) HEAD(path string, handlers ...router.HandlerFunc) {
	app.router.Handle(MethodHead, path, handlers...)
}

// OPTIONS đăng ký handler cho HTTP OPTIONS method.
// HTTP OPTIONS thường được sử dụng để truy vấn thông tin về communication options.
//
// Parameters:
//   - path: Đường dẫn URL để đăng ký handler
//   - handlers: Danh sách các handlers xử lý request
func (app *WebApp) OPTIONS(path string, handlers ...router.HandlerFunc) {
	app.router.Handle(MethodOptions, path, handlers...)
}

// Any đăng ký handler cho tất cả các HTTP methods phổ biến.
// Handler sẽ được đăng ký cho GET, POST, PUT, DELETE, PATCH, HEAD và OPTIONS.
//
// Parameters:
//   - path: Đường dẫn URL để đăng ký handler
//   - handlers: Danh sách các handlers xử lý request
func (app *WebApp) Any(path string, handlers ...router.HandlerFunc) {
	methods := []string{
		MethodGet, MethodPost, MethodPut, MethodDelete,
		MethodPatch, MethodHead, MethodOptions,
	}
	for _, method := range methods {
		app.router.Handle(method, path, handlers...)
	}
}

// Handle đăng ký handler cho một HTTP method cụ thể.
// Đây là phương thức tổng quát cho phép đăng ký handler với bất kỳ HTTP method nào.
//
// Parameters:
//   - method: HTTP method cần đăng ký (GET, POST, PUT, DELETE, v.v.)
//   - path: Đường dẫn URL để đăng ký handler
//   - handlers: Danh sách các handlers xử lý request
func (app *WebApp) Handle(method, path string, handlers ...router.HandlerFunc) {
	app.router.Handle(method, path, handlers...)
}

// Run khởi động HTTP server sử dụng adapter hiện tại.
// Server sẽ lắng nghe và xử lý các HTTP requests theo cấu hình từ adapter.
//
// Returns:
//   - error: Lỗi nếu có trong quá trình khởi động server
//
// Errors:
//   - ErrAdapterNotSet: Trả về khi adapter chưa được thiết lập
func (app *WebApp) Serve() error {
	app.mu.RLock()
	adp := app.adapter
	app.mu.RUnlock()

	if adp == nil {
		return ErrAdapterNotSet
	}

	// Đặt router làm handler cho adapter
	adp.SetHandler(app.router)

	// Chạy server với cấu hình từ adapter
	return adp.Serve()
}

// RunTLS khởi động HTTPS server với các tệp chứng chỉ SSL/TLS đã chỉ định.
// Server sẽ lắng nghe và xử lý các HTTPS requests theo cấu hình từ adapter.
//
// Parameters:
//   - certFile: Đường dẫn đến tệp chứng chỉ SSL/TLS
//   - keyFile: Đường dẫn đến tệp khóa SSL/TLS
//
// Returns:
//   - error: Lỗi nếu có trong quá trình khởi động server
//
// Errors:
//   - ErrAdapterNotSet: Trả về khi adapter chưa được thiết lập
//   - ErrInvalidCertificate: Trả về khi tệp chứng chỉ hoặc tệp khóa không hợp lệ
func (app *WebApp) RunTLS(certFile, keyFile string) error {
	app.mu.RLock()
	adp := app.adapter
	app.mu.RUnlock()

	if adp == nil {
		return ErrAdapterNotSet
	}

	// Kiểm tra tệp chứng chỉ và khóa
	if certFile == "" || keyFile == "" {
		return ErrInvalidCertificate
	}

	// Đặt router làm handler cho adapter
	adp.SetHandler(app.router)

	// Chạy server với TLS và cấu hình từ adapter
	return adp.RunTLS(certFile, keyFile)
}

// ServeWithGracefulShutdown khởi động server với graceful shutdown tự động
// Phương thức này tự động lắng nghe shutdown signals và thực hiện graceful shutdown
//
// Returns:
//   - error: Lỗi nếu có trong quá trình khởi động hoặc shutdown server
func (app *WebApp) ServeWithGracefulShutdown() error {
	// Bật signal listening nếu graceful shutdown được enable
	if app.config.GracefulShutdown.Enabled {
		app.ListenForShutdownSignals()
	}

	// Chạy server
	return app.Serve()
}

// ServeHTTP xử lý HTTP request và implement interface http.Handler.
// Phương thức này cho phép WebApp hoạt động như một HTTP handler.
//
// Parameters:
//   - w: HTTP response writer để ghi response
//   - r: HTTP request cần xử lý
func (app *WebApp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.router.ServeHTTP(w, r)
}

// Shutdown đóng HTTP server một cách an toàn, chờ các kết nối hiện tại kết thúc.
// Phương thức này nên được gọi khi muốn dừng server một cách graceful.
//
// Returns:
//   - error: Lỗi nếu có trong quá trình đóng server
func (app *WebApp) Shutdown() error {
	app.mu.RLock()
	adp := app.adapter
	app.mu.RUnlock()

	if adp == nil {
		return nil
	}

	return adp.Shutdown()
}

// GracefulShutdown thực hiện graceful shutdown với cấu hình nâng cao
// Phương thức này xử lý signals, timeout và connection tracking
//
// Returns:
//   - error: Lỗi nếu có trong quá trình graceful shutdown
func (app *WebApp) GracefulShutdown() error {
	app.mu.Lock()
	if app.isShuttingDown {
		app.mu.Unlock()
		return nil // Already shutting down
	}
	app.isShuttingDown = true
	config := app.config.GracefulShutdown
	app.mu.Unlock()

	if !config.Enabled {
		return app.Shutdown()
	}

	// Call OnShutdownStart callback
	if config.OnShutdownStart != nil {
		config.OnShutdownStart()
	}

	// Create timeout context
	shutdownCtx, cancel := context.WithTimeout(app.shutdownCtx, time.Duration(config.Timeout)*time.Second)
	defer cancel()

	// Wait for connections if enabled
	if config.WaitForConnections {
		app.waitForConnections(shutdownCtx)
	}

	// Perform actual shutdown
	err := app.Shutdown()

	// Call appropriate callback
	if err != nil && config.OnShutdownError != nil {
		config.OnShutdownError(err)
	} else if err == nil && config.OnShutdownComplete != nil {
		config.OnShutdownComplete()
	}

	return err
}

// ListenForShutdownSignals lắng nghe các signals để thực hiện graceful shutdown
// Hàm này chạy trong goroutine riêng và tự động shutdown khi nhận signals
func (app *WebApp) ListenForShutdownSignals() {
	if !app.config.GracefulShutdown.Enabled {
		return
	}

	sigChan := make(chan os.Signal, app.config.GracefulShutdown.SignalBufferSize)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		app.GracefulShutdown()
	}()
}

// TrackConnection tăng số lượng active connections
func (app *WebApp) TrackConnection() {
	atomic.AddInt32(&app.activeConnections, 1)
}

// UntrackConnection giảm số lượng active connections
func (app *WebApp) UntrackConnection() {
	atomic.AddInt32(&app.activeConnections, -1)
}

// GetActiveConnections trả về số lượng connections đang hoạt động
func (app *WebApp) GetActiveConnections() int32 {
	return atomic.LoadInt32(&app.activeConnections)
}

// waitForConnections chờ tất cả connections kết thúc hoặc timeout
func (app *WebApp) waitForConnections(ctx context.Context) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return // Timeout
		case <-ticker.C:
			if app.GetActiveConnections() == 0 {
				return // All connections closed
			}
		}
	}
}

// SetShutdownTimeout thiết lập timeout cho graceful shutdown
func (app *WebApp) SetShutdownTimeout(timeout time.Duration) {
	app.mu.Lock()
	defer app.mu.Unlock()
	app.config.GracefulShutdown.Timeout = int(timeout.Seconds())
}

// Router trả về router hiện tại của WebApp.
// Router chịu trách nhiệm quản lý routes và điều hướng requests.
//
// Returns:
//   - router.Router: Router hiện tại của WebApp
func (app *WebApp) Router() router.Router {
	return app.router
}

// NewContext tạo một context mới để xử lý HTTP request/response.
// Context cung cấp các tiện ích để truy cập request và xử lý response.
//
// Parameters:
//   - w: HTTP response writer
//   - r: HTTP request
//
// Returns:
//   - forkCtx.Context: Context mới đã được khởi tạo
func (app *WebApp) NewContext(w http.ResponseWriter, r *http.Request) forkCtx.Context {
	return forkCtx.NewContext(w, r)
}

// GetAdapter trả về adapter hiện tại của WebApp.
// Adapter là interface giao tiếp với HTTP server cơ bản.
//
// Returns:
//   - adapter.Adapter: Adapter hiện tại của WebApp
func (app *WebApp) GetAdapter() adapter.Adapter {
	app.mu.RLock()
	defer app.mu.RUnlock()

	return app.adapter
}

// SetConfig thiết lập cấu hình cho WebApp.
//
// Parameters:
//   - config: Cấu hình mới cho WebApp
func (app *WebApp) SetConfig(config *WebAppConfig) {
	app.mu.Lock()
	defer app.mu.Unlock()

	if config != nil {
		app.config = config
	}
}

// GetConfig trả về cấu hình hiện tại của WebApp.
//
// Returns:
//   - *WebAppConfig: Cấu hình hiện tại của WebApp
func (app *WebApp) GetConfig() *WebAppConfig {
	app.mu.RLock()
	defer app.mu.RUnlock()

	return app.config
}

// CleanupResources dọn dẹp tài nguyên không sử dụng để tránh memory leaks
func (app *WebApp) CleanupResources() {
	app.mu.Lock()
	defer app.mu.Unlock()

	// Cleanup router groups nếu có method tương ứng
	if cleaner, ok := app.router.(interface{ Cleanup() }); ok {
		cleaner.Cleanup()
	}

	// Cancel shutdown context if not already cancelled
	if app.shutdownCancel != nil {
		app.shutdownCancel()
	}
}

// IsShuttingDown kiểm tra xem WebApp có đang trong quá trình shutdown không
func (app *WebApp) IsShuttingDown() bool {
	app.mu.RLock()
	defer app.mu.RUnlock()
	return app.isShuttingDown
}

// EnableSecurityMiddleware bật các middleware bảo mật tự động
// Note: Security headers are now handled by the helmet middleware package.
// Request size, method validation, and timeout are handled by their respective middleware packages:
// - bodylimit: handles max_request_body_size
// - method: handles allowed_methods
// - timeout: handles request_timeout
// - helmet: handles enable_security_headers
func (app *WebApp) EnableSecurityMiddleware() {
	app.Use(app.createConnectionTrackingMiddleware())
}

// createConnectionTrackingMiddleware tạo middleware để theo dõi active connections
func (app *WebApp) createConnectionTrackingMiddleware() router.HandlerFunc {
	return func(c forkCtx.Context) {
		if app.config.GracefulShutdown.Enabled && app.config.GracefulShutdown.WaitForConnections {
			app.TrackConnection()
			defer app.UntrackConnection()
		}
		c.Next()
	}
}
