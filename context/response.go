package context

import (
	"bufio"
	"errors"
	"net"
	"net/http"
)

// Response interface định nghĩa các phương thức để truy cập và điều khiển HTTP response.
// Interface này mở rộng http.ResponseWriter và cung cấp các phương thức thuận tiện
// để tạo và quản lý HTTP responses.
type Response interface {
	// Header trả về HTTP headers của response.
	// Headers này có thể được cập nhật trước khi response được gửi đi.
	//
	// Returns:
	//   - http.Header: Map các HTTP headers của response
	Header() http.Header

	// Write viết dữ liệu vào response body.
	// Nếu WriteHeader chưa được gọi, Write sẽ thiết lập status code là 200 OK.
	//
	// Parameters:
	//   - data: Dữ liệu để viết vào response body
	//
	// Returns:
	//   - int: Số bytes đã được viết
	//   - error: Lỗi nếu có trong quá trình viết dữ liệu
	Write(data []byte) (int, error)

	// WriteHeader thiết lập HTTP status code cho response.
	// Phương thức này chỉ nên được gọi một lần trước khi viết body.
	//
	// Parameters:
	//   - code: HTTP status code (200, 404, 500, v.v.)
	WriteHeader(code int)

	// Flush ghi dữ liệu từ buffer vào network.
	// Hữu ích cho streaming responses và server-sent events.
	Flush()

	// Status trả về HTTP status code hiện tại của response.
	//
	// Returns:
	//   - int: HTTP status code đã được thiết lập hoặc 200 nếu chưa được thiết lập
	Status() int

	// Size trả về kích thước của response body đã được viết.
	//
	// Returns:
	//   - int: Kích thước tính bằng bytes của response body
	Size() int

	// Written kiểm tra xem response đã được viết chưa.
	//
	// Returns:
	//   - bool: true nếu response đã được viết, ngược lại là false
	Written() bool

	// Hijack cho phép lấy control của connection từ HTTP server.
	// Hữu ích cho việc nâng cấp protocol (ví dụ: từ HTTP sang WebSocket).
	//
	// Returns:
	//   - net.Conn: Network connection
	//   - *bufio.ReadWriter: Buffer để đọc và viết dữ liệu
	//   - error: Lỗi nếu không thể hijack connection
	Hijack() (net.Conn, *bufio.ReadWriter, error)

	// ResponseWriter trả về http.ResponseWriter gốc.
	// Hữu ích khi cần truy cập các phương thức đặc biệt của implementation cụ thể.
	//
	// Returns:
	//   - http.ResponseWriter: Đối tượng ResponseWriter gốc
	ResponseWriter() http.ResponseWriter

	// Reset reset response writer về trạng thái ban đầu.
	// Phương thức này thường được sử dụng khi tái sử dụng response writer.
	//
	// Parameters:
	//   - w: http.ResponseWriter mới để reset
	Reset(w http.ResponseWriter)

	// Pusher trả về http.Pusher nếu server hỗ trợ HTTP/2 server push.
	// Server push cho phép server chủ động gửi resources đến client.
	//
	// Returns:
	//   - http.Pusher: Đối tượng Pusher nếu có hỗ trợ
	//   - bool: true nếu HTTP/2 server push được hỗ trợ, ngược lại là false
	Pusher() (http.Pusher, bool)
}

// forkResponse là implementation private cho Response interface.
// Được thiết kế để chỉ tạo được thông qua hàm NewResponse để đảm bảo an toàn.
type forkResponse struct {
	// writer là http.ResponseWriter gốc
	writer http.ResponseWriter

	// statusCode lưu trữ HTTP status code của response
	statusCode int

	// size lưu trữ kích thước của response body đã được viết
	size int

	// written kiểm tra xem response đã được viết chưa.
	written bool
}

// NewResponse tạo một response mới từ http.ResponseWriter.
// Phương thức này khởi tạo một đối tượng Response với các giá trị mặc định.
//
// Parameters:
//   - w: http.ResponseWriter gốc
//
// Returns:
//   - Response: Đối tượng Response mới đã được khởi tạo
func NewResponse(w http.ResponseWriter) Response {
	return &forkResponse{
		writer:     w,
		statusCode: http.StatusOK,
		size:       0,
		written:    false,
	}
}

// Header trả về HTTP headers của response.
// Triển khai phương thức Header của Response interface.
//
// Returns:
//   - http.Header: Map các HTTP headers của response
func (r *forkResponse) Header() http.Header {
	return r.writer.Header()
}

// Write viết dữ liệu vào response body.
// Nếu WriteHeader chưa được gọi, Write sẽ gọi WriteHeader(http.StatusOK) ngầm định.
// Triển khai phương thức Write của Response interface.
//
// Parameters:
//   - data: Dữ liệu để viết vào response body
//
// Returns:
//   - int: Số bytes đã được viết
//   - error: Lỗi nếu có trong quá trình viết dữ liệu
func (r *forkResponse) Write(data []byte) (int, error) {
	r.written = true
	n, err := r.writer.Write(data)
	r.size += n
	return n, err
}

// WriteHeader thiết lập HTTP status code cho response.
// Phương thức này chỉ có hiệu lực nếu response chưa được viết.
// Triển khai phương thức WriteHeader của Response interface.
//
// Parameters:
//   - code: HTTP status code (200, 404, 500, v.v.)
func (r *forkResponse) WriteHeader(code int) {
	if r.written {
		return
	}
	r.statusCode = code
	r.writer.WriteHeader(code)
	r.written = true
}

// Flush ghi dữ liệu từ buffer vào network.
// Phương thức này sẽ gọi Flush của http.Flusher nếu writer hỗ trợ.
// Triển khai phương thức Flush của Response interface.
func (r *forkResponse) Flush() {
	if flusher, ok := r.writer.(http.Flusher); ok {
		flusher.Flush()
	}
}

// Status trả về HTTP status code hiện tại của response.
// Triển khai phương thức Status của Response interface.
//
// Returns:
//   - int: HTTP status code đã được thiết lập
func (r *forkResponse) Status() int {
	return r.statusCode
}

// Size trả về kích thước của response body đã được viết.
// Triển khai phương thức Size của Response interface.
//
// Returns:
//   - int: Kích thước tính bằng bytes của response body
func (r *forkResponse) Size() int {
	return r.size
}

// Written kiểm tra xem response đã được viết chưa.
// Triển khai phương thức Written của Response interface.
//
// Returns:
//   - bool: true nếu response đã được viết, ngược lại là false
func (r *forkResponse) Written() bool {
	return r.written
}

// Hijack cho phép lấy control của connection từ HTTP server.
// Triển khai phương thức Hijack của Response interface.
// Phương thức này chỉ hoạt động nếu writer gốc hỗ trợ http.Hijacker.
//
// Returns:
//   - net.Conn: Network connection
//   - *bufio.ReadWriter: Buffer để đọc và viết dữ liệu
//   - error: Lỗi nếu writer không hỗ trợ http.Hijacker
func (r *forkResponse) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := r.writer.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, errors.New("http: response does not implement http.Hijacker")
}

// ResponseWriter trả về http.ResponseWriter gốc.
// Triển khai phương thức ResponseWriter của Response interface.
//
// Returns:
//   - http.ResponseWriter: Đối tượng ResponseWriter gốc
func (r *forkResponse) ResponseWriter() http.ResponseWriter {
	return r.writer
}

// Reset reset response writer về trạng thái ban đầu.
// Triển khai phương thức Reset của Response interface.
//
// Parameters:
//   - w: http.ResponseWriter mới để reset
func (r *forkResponse) Reset(w http.ResponseWriter) {
	r.writer = w
	r.statusCode = http.StatusOK
	r.size = 0
	r.written = false
}

// Pusher trả về http.Pusher nếu server hỗ trợ HTTP/2 server push.
// Triển khai phương thức Pusher của Response interface.
//
// Returns:
//   - http.Pusher: Đối tượng Pusher nếu writer hỗ trợ HTTP/2 server push
//   - bool: true nếu HTTP/2 server push được hỗ trợ, ngược lại là false
func (r *forkResponse) Pusher() (http.Pusher, bool) {
	pusher, ok := r.writer.(http.Pusher)
	return pusher, ok
}
