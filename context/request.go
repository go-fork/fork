package context

import (
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

// Request interface định nghĩa các phương thức để truy cập thông tin HTTP request.
// Interface này bọc http.Request và cung cấp các phương thức thuận tiện để truy cập
// và xử lý các thành phần của HTTP request.
type Request interface {
	// Method trả về HTTP method của request (GET, POST, PUT, DELETE, v.v.).
	//
	// Returns:
	//   - string: HTTP method của request
	Method() string

	// URL trả về URL đầy đủ của request.
	// URL này chứa thông tin về scheme, host, path, và query parameters.
	//
	// Returns:
	//   - *url.URL: URL đầy đủ của request
	URL() *url.URL

	// Header trả về tất cả HTTP headers của request.
	// Headers cho phép truy cập thông tin như Content-Type, Authorization, v.v.
	//
	// Returns:
	//   - http.Header: Map các HTTP headers của request
	Header() http.Header

	// Cookies trả về tất cả cookies được gửi từ client.
	//
	// Returns:
	//   - []*http.Cookie: Slice chứa tất cả cookies của request
	Cookies() []*http.Cookie

	// Cookie trả về cookie với tên cụ thể từ request.
	//
	// Parameters:
	//   - name: Tên của cookie cần lấy
	//
	// Returns:
	//   - *http.Cookie: Cookie được tìm thấy
	//   - error: Lỗi nếu không tìm thấy cookie
	Cookie(name string) (*http.Cookie, error)

	// Body trả về body của request dưới dạng io.ReadCloser.
	// Lưu ý: Body chỉ có thể đọc một lần, nên cần xử lý cẩn thận.
	//
	// Returns:
	//   - io.ReadCloser: Body của request
	Body() io.ReadCloser

	// Form trả về tất cả các form values của request, bao gồm cả query string và post form.
	// Phương thức này sẽ tự động gọi ParseForm nếu cần.
	//
	// Returns:
	//   - url.Values: Map các form values của request
	Form() url.Values

	// PostForm trả về các form values từ post request.
	// Chỉ bao gồm các values từ application/x-www-form-urlencoded form.
	//
	// Returns:
	//   - url.Values: Map các post form values của request
	PostForm() url.Values

	// FormValue trả về form value đầu tiên cho name từ form data,
	// bao gồm cả query string và post form.
	//
	// Parameters:
	//   - name: Tên của form field cần lấy
	//
	// Returns:
	//   - string: Giá trị của form field hoặc chuỗi rỗng nếu không tồn tại
	FormValue(name string) string

	// PostFormValue trả về form value đầu tiên cho name từ post form data.
	// Chỉ tìm kiếm trong application/x-www-form-urlencoded form.
	//
	// Parameters:
	//   - name: Tên của post form field cần lấy
	//
	// Returns:
	//   - string: Giá trị của post form field hoặc chuỗi rỗng nếu không tồn tại
	PostFormValue(name string) string

	// MultipartForm trả về multipart form của request.
	// Phương thức này sẽ tự động gọi ParseMultipartForm nếu cần.
	//
	// Returns:
	//   - *multipart.Form: Multipart form của request
	//   - error: Lỗi nếu không thể parse multipart form
	MultipartForm() (*multipart.Form, error)

	// FormFile trả về file upload đầu tiên cho form field name.
	// Phương thức này hữu ích cho việc xử lý file uploads.
	//
	// Parameters:
	//   - name: Tên của form field chứa file upload
	//
	// Returns:
	//   - *multipart.FileHeader: Thông tin về file upload
	//   - error: Lỗi nếu không tìm thấy file hoặc không thể parse multipart form
	FormFile(name string) (*multipart.FileHeader, error)

	// RemoteAddr trả về địa chỉ IP và port của client.
	//
	// Returns:
	//   - string: Địa chỉ IP và port của client (format: "IP:port")
	RemoteAddr() string

	// ParseForm parse form data của request.
	// Phương thức này parse query string và application/x-www-form-urlencoded form.
	//
	// Returns:
	//   - error: Lỗi nếu không thể parse form data
	ParseForm() error

	// ParseMultipartForm parse multipart form của request.
	// Phương thức này parse multipart form data với giới hạn bộ nhớ tối đa.
	//
	// Parameters:
	//   - maxMemory: Giới hạn bộ nhớ tối đa (bytes) được sử dụng để parse form
	//
	// Returns:
	//   - error: Lỗi nếu không thể parse multipart form
	ParseMultipartForm(maxMemory int64) error

	// UserAgent trả về user agent của request.
	// Thông tin này thường được sử dụng để xác định browser hoặc client đang được sử dụng.
	//
	// Returns:
	//   - string: User agent của request hoặc chuỗi rỗng nếu không có
	UserAgent() string

	// Referer trả về referer của request.
	// Referer là URL của trang mà từ đó request hiện tại được gửi.
	//
	// Returns:
	//   - string: URL của referer hoặc chuỗi rỗng nếu không có
	Referer() string

	// ContentLength trả về content length của request.
	// Content length là kích thước của request body tính bằng bytes.
	//
	// Returns:
	//   - int64: Content length của request body
	ContentLength() int64

	// Host trả về host của request.
	// Host là tên miền và cổng (nếu có) mà request được gửi đến.
	//
	// Returns:
	//   - string: Host của request (ví dụ: "example.com:8080")
	Host() string

	// RequestURI trả về URI của request.
	// URI bao gồm path và query string của request.
	//
	// Returns:
	//   - string: URI của request (ví dụ: "/users?page=1")
	RequestURI() string

	// Scheme trả về scheme của request (http hoặc https).
	// Scheme xác định giao thức được sử dụng để gửi request.
	//
	// Returns:
	//   - string: "https" nếu request là secure, ngược lại là "http"
	Scheme() string

	// IsSecure kiểm tra xem request có an toàn không (HTTPS).
	// Phương thức này xác định nếu connection được mã hóa bằng TLS/SSL.
	//
	// Returns:
	//   - bool: true nếu request sử dụng HTTPS, ngược lại là false
	IsSecure() bool

	// Protocol trả về protocol của request.
	// Trả về phiên bản HTTP protocol được sử dụng (ví dụ: "HTTP/1.1").
	//
	// Returns:
	//   - string: Protocol của request
	Protocol() string

	// Request trả về http.Request gốc.
	// Phương thức này cho phép truy cập trực tiếp đối tượng http.Request của Go.
	//
	// Returns:
	//   - *http.Request: Đối tượng http.Request gốc
	Request() *http.Request
}

// forkRequest là implementation private cho Request interface.
// Được thiết kế để chỉ tạo được thông qua hàm NewRequest để đảm bảo an toàn.
type forkRequest struct {
	request *http.Request
}

// NewRequest tạo một request mới từ http.Request.
// Phương thức này tạo một đối tượng Request mới bọc http.Request gốc.
//
// Parameters:
//   - r: http.Request gốc
//
// Returns:
//   - Request: Đối tượng Request mới đã được khởi tạo
func NewRequest(r *http.Request) Request {
	return &forkRequest{
		request: r,
	}
}

// Method trả về HTTP method của request.
// Triển khai phương thức Method của Request interface.
//
// Returns:
//   - string: HTTP method của request (GET, POST, PUT, DELETE, v.v.)
func (r *forkRequest) Method() string {
	return r.request.Method
}

// URL trả về URL đầy đủ của request.
// Triển khai phương thức URL của Request interface.
//
// Returns:
//   - *url.URL: URL đầy đủ của request
func (r *forkRequest) URL() *url.URL {
	return r.request.URL
}

// Header trả về tất cả HTTP headers của request.
// Triển khai phương thức Header của Request interface.
//
// Returns:
//   - http.Header: Map các HTTP headers của request
func (r *forkRequest) Header() http.Header {
	return r.request.Header
}

// Cookies trả về tất cả cookies được gửi từ client.
// Triển khai phương thức Cookies của Request interface.
//
// Returns:
//   - []*http.Cookie: Slice chứa tất cả cookies của request
func (r *forkRequest) Cookies() []*http.Cookie {
	return r.request.Cookies()
}

// Cookie trả về cookie với tên cụ thể từ request.
// Triển khai phương thức Cookie của Request interface.
//
// Parameters:
//   - name: Tên của cookie cần lấy
//
// Returns:
//   - *http.Cookie: Cookie được tìm thấy
//   - error: Lỗi nếu không tìm thấy cookie
func (r *forkRequest) Cookie(name string) (*http.Cookie, error) {
	return r.request.Cookie(name)
}

// Body trả về body của request dưới dạng io.ReadCloser.
// Triển khai phương thức Body của Request interface.
//
// Returns:
//   - io.ReadCloser: Body của request
func (r *forkRequest) Body() io.ReadCloser {
	return r.request.Body
}

// Form trả về tất cả các form values của request.
// Triển khai phương thức Form của Request interface.
//
// Returns:
//   - url.Values: Map các form values của request
func (r *forkRequest) Form() url.Values {
	return r.request.Form
}

// PostForm trả về các form values từ post request.
// Triển khai phương thức PostForm của Request interface.
//
// Returns:
//   - url.Values: Map các post form values của request
func (r *forkRequest) PostForm() url.Values {
	return r.request.PostForm
}

// FormValue trả về form value đầu tiên cho name từ form data.
// Triển khai phương thức FormValue của Request interface.
//
// Parameters:
//   - name: Tên của form field cần lấy
//
// Returns:
//   - string: Giá trị của form field hoặc chuỗi rỗng nếu không tồn tại
func (r *forkRequest) FormValue(name string) string {
	return r.request.FormValue(name)
}

// PostFormValue trả về form value đầu tiên cho name từ post form data.
// Triển khai phương thức PostFormValue của Request interface.
//
// Parameters:
//   - name: Tên của post form field cần lấy
//
// Returns:
//   - string: Giá trị của post form field hoặc chuỗi rỗng nếu không tồn tại
func (r *forkRequest) PostFormValue(name string) string {
	return r.request.PostFormValue(name)
}

// MultipartForm trả về multipart form của request.
// Triển khai phương thức MultipartForm của Request interface.
// Nếu multipart form chưa được parse, phương thức sẽ tự động parse với giới hạn mặc định là 32MB.
//
// Returns:
//   - *multipart.Form: Multipart form của request
//   - error: Lỗi nếu không thể parse multipart form
func (r *forkRequest) MultipartForm() (*multipart.Form, error) {
	if r.request.MultipartForm == nil {
		if err := r.request.ParseMultipartForm(32 << 20); err != nil {
			return nil, err
		}
	}
	return r.request.MultipartForm, nil
}

// FormFile trả về file upload đầu tiên cho form field name.
// Triển khai phương thức FormFile của Request interface.
//
// Parameters:
//   - name: Tên của form field chứa file upload
//
// Returns:
//   - *multipart.FileHeader: Thông tin về file upload
//   - error: Lỗi nếu không tìm thấy file hoặc không thể parse multipart form
func (r *forkRequest) FormFile(name string) (*multipart.FileHeader, error) {
	_, fileHeader, err := r.request.FormFile(name)
	return fileHeader, err
}

// RemoteAddr trả về địa chỉ IP và port của client.
// Triển khai phương thức RemoteAddr của Request interface.
//
// Returns:
//   - string: Địa chỉ IP và port của client (format: "IP:port")
func (r *forkRequest) RemoteAddr() string {
	return r.request.RemoteAddr
}

// ParseForm parse form data của request.
// Triển khai phương thức ParseForm của Request interface.
//
// Returns:
//   - error: Lỗi nếu không thể parse form data
func (r *forkRequest) ParseForm() error {
	return r.request.ParseForm()
}

// ParseMultipartForm parse multipart form của request.
// Triển khai phương thức ParseMultipartForm của Request interface.
//
// Parameters:
//   - maxMemory: Giới hạn bộ nhớ tối đa (bytes) được sử dụng để parse form
//
// Returns:
//   - error: Lỗi nếu không thể parse multipart form
func (r *forkRequest) ParseMultipartForm(maxMemory int64) error {
	return r.request.ParseMultipartForm(maxMemory)
}

// UserAgent trả về user agent của request.
// Triển khai phương thức UserAgent của Request interface.
//
// Returns:
//   - string: User agent của request hoặc chuỗi rỗng nếu không có
func (r *forkRequest) UserAgent() string {
	return r.request.UserAgent()
}

// Referer trả về referer của request.
// Triển khai phương thức Referer của Request interface.
//
// Returns:
//   - string: URL của referer hoặc chuỗi rỗng nếu không có
func (r *forkRequest) Referer() string {
	return r.request.Referer()
}

// ContentLength trả về content length của request.
// Triển khai phương thức ContentLength của Request interface.
//
// Returns:
//   - int64: Content length của request body
func (r *forkRequest) ContentLength() int64 {
	return r.request.ContentLength
}

// Host trả về host của request.
// Triển khai phương thức Host của Request interface.
//
// Returns:
//   - string: Host của request (ví dụ: "example.com:8080")
func (r *forkRequest) Host() string {
	return r.request.Host
}

// RequestURI trả về URI của request.
// Triển khai phương thức RequestURI của Request interface.
//
// Returns:
//   - string: URI của request (ví dụ: "/users?page=1")
func (r *forkRequest) RequestURI() string {
	return r.request.RequestURI
}

// Scheme trả về scheme của request (http hoặc https).
// Triển khai phương thức Scheme của Request interface.
//
// Returns:
//   - string: "https" nếu request là secure, ngược lại là "http"
func (r *forkRequest) Scheme() string {
	if r.request.TLS != nil {
		return "https"
	}
	return "http"
}

// IsSecure kiểm tra xem request có an toàn không (HTTPS).
// Triển khai phương thức IsSecure của Request interface.
//
// Returns:
//   - bool: true nếu request sử dụng HTTPS, ngược lại là false
func (r *forkRequest) IsSecure() bool {
	return r.Scheme() == "https"
}

// Protocol trả về protocol của request.
// Triển khai phương thức Protocol của Request interface.
// Phương thức này cải tiến cách xác định protocol để hỗ trợ load balancing giữa các giao thức
// HTTP/1.1, HTTP/2, và HTTP/3, bằng cách kiểm tra cả ProtoMajor, Proto header và các header đặc thù.
//
// Returns:
//   - string: Protocol của request ("HTTP/1.1", "HTTP/2", "HTTP/3", hoặc "h2c-upgrade")
func (r *forkRequest) Protocol() string {
	req := r.request

	// Xác định protocol dựa trên thông tin trong request
	switch {
	case req.ProtoMajor == 3 || strings.HasPrefix(req.Proto, "HTTP/3"):
		return "HTTP/3"
	case req.Header.Get("Alt-Used") != "": // Client đã nâng cấp lên HTTP/3
		return "HTTP/3"
	case req.ProtoMajor == 2 || strings.HasPrefix(req.Proto, "HTTP/2"):
		return "HTTP/2"
	case req.Header.Get("Upgrade") == "h2c": // Client yêu cầu nâng cấp lên HTTP/2 Cleartext
		return "h2c-upgrade"
	default:
		return "HTTP/1.1"
	}
}

// Request trả về http.Request gốc.
// Triển khai phương thức Request của Request interface.
//
// Returns:
//   - *http.Request: Đối tượng http.Request gốc
func (r *forkRequest) Request() *http.Request {
	return r.request
}
