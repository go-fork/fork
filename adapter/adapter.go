// Package adapter cung cấp interface và implementations cho việc kết nối
// framework với các HTTP server khác nhau (như net/http standard, fasthttp, v.v.).
// Package này cho phép framework hoạt động trên nhiều server implementations khác nhau.
package adapter

import (
	"net/http"

	"go.fork.vn/fork/context"
)

// Adapter định nghĩa interface cho các HTTP engine adapters.
// Interface này cung cấp các phương thức chuẩn để tương tác với các HTTP server
// implementations khác nhau, cho phép framework hoạt động độc lập với server cụ thể.
type Adapter interface {
	// Name trả về tên của adapter.
	// Thông tin này hữu ích cho việc ghi log và debugging.
	//
	// Returns:
	//   - string: Tên định danh của adapter
	Name() string

	// Serve khởi động HTTP server với cấu hình từ adapter.
	// Phương thức này bắt đầu lắng nghe các HTTP requests trên địa chỉ và cổng đã cấu hình.
	//
	// Returns:
	//   - error: Lỗi nếu có trong quá trình khởi động server hoặc nil nếu thành công
	Serve() error

	// RunTLS khởi động HTTPS server với cấu hình từ adapter.
	// Server sẽ sử dụng TLS/SSL với các tệp chứng chỉ đã cung cấp.
	//
	// Parameters:
	//   - certFile: Đường dẫn đến tệp chứng chỉ SSL/TLS
	//   - keyFile: Đường dẫn đến tệp khóa SSL/TLS
	//
	// Returns:
	//   - error: Lỗi nếu có trong quá trình khởi động server hoặc nil nếu thành công
	RunTLS(certFile, keyFile string) error

	// ServeHTTP xử lý HTTP request.
	// Implements interface http.Handler để tích hợp với Go's standard HTTP library.
	//
	// Parameters:
	//   - w: HTTP response writer để ghi response
	//   - r: HTTP request cần xử lý
	ServeHTTP(w http.ResponseWriter, r *http.Request)

	// HandleFunc đăng ký một handler function với method và path.
	// Phương thức này cho phép đăng ký handler cho các routes cụ thể.
	//
	// Parameters:
	//   - method: HTTP method (GET, POST, PUT, v.v.)
	//   - path: URL path pattern để đăng ký handler
	//   - handler: Function xử lý request với context đã cho
	HandleFunc(method, path string, handler func(ctx context.Context))

	// Use thêm middleware vào adapter.
	// Middleware sẽ được thực thi cho mọi request trước khi đến handler chính.
	//
	// Parameters:
	//   - middleware: Function middleware để thêm vào chain
	Use(middleware func(ctx context.Context))

	// SetHandler thiết lập handler chính cho adapter.
	// Handler này sẽ nhận và xử lý tất cả các requests đến adapter.
	//
	// Parameters:
	//   - handler: HTTP handler để thiết lập
	SetHandler(handler http.Handler)

	// Shutdown đóng HTTP server một cách graceful.
	// Phương thức này chờ các request hiện tại hoàn thành trước khi đóng server.
	//
	// Returns:
	//   - error: Lỗi nếu có trong quá trình đóng server hoặc nil nếu thành công
	Shutdown() error
}
