package context

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
)

// Context đại diện cho một HTTP request/response context.
// Context đóng vai trò như một container chứa tất cả thông tin liên quan đến
// một HTTP request và response, cung cấp các phương thức để truy cập và thao tác
// với dữ liệu, xử lý middleware, quản lý session và thực hiện các chức năng khác.
type Context interface {
	// Request trả về đối tượng Request.
	//
	// Returns:
	//   - Request: Đối tượng Request chứa thông tin về HTTP request hiện tại
	Request() Request

	// Response trả về đối tượng Response.
	//
	// Returns:
	//   - Response: Đối tượng Response cho phép thao tác với HTTP response
	Response() Response

	// Context trả về context.Context gốc.
	//
	// Returns:
	//   - context.Context: Đối tượng context.Context gốc từ request
	Context() context.Context

	// WithContext thiết lập context.Context mới và trả về context cập nhật.
	//
	// Parameters:
	//   - ctx: context.Context mới sẽ được sử dụng
	//
	// Returns:
	//   - Context: Context sau khi được cập nhật context.Context
	WithContext(ctx context.Context) Context

	// Next gọi middleware tiếp theo trong chuỗi.
	// Phương thức này thực thi middleware tiếp theo trong pipeline.
	Next()

	// Abort ngừng thực thi middleware chains.
	// Khi được gọi, các middleware còn lại trong chuỗi sẽ không được thực thi.
	Abort()

	// IsAborted kiểm tra xem context có bị abort không.
	//
	// Returns:
	//   - bool: true nếu context đã bị abort, ngược lại là false
	IsAborted() bool

	// Set thiết lập giá trị cho một khóa trong context.
	//
	// Parameters:
	//   - key: Khóa để lưu trữ giá trị
	//   - value: Giá trị cần lưu trữ
	Set(key string, value interface{})

	// Get lấy giá trị cho một khóa từ context.
	//
	// Parameters:
	//   - key: Khóa cần truy xuất giá trị
	//
	// Returns:
	//   - interface{}: Giá trị được lưu trữ
	//   - bool: true nếu khóa tồn tại, ngược lại là false
	Get(key string) (interface{}, bool)

	// GetString lấy giá trị string cho một khóa từ context.
	//
	// Parameters:
	//   - key: Khóa cần truy xuất giá trị
	//
	// Returns:
	//   - string: Giá trị string được lưu trữ, hoặc chuỗi rỗng nếu khóa không tồn tại hoặc giá trị không phải kiểu string
	GetString(key string) string

	// GetBool lấy giá trị boolean cho một khóa từ context.
	//
	// Parameters:
	//   - key: Khóa cần truy xuất giá trị
	//
	// Returns:
	//   - bool: Giá trị boolean được lưu trữ, hoặc false nếu khóa không tồn tại hoặc giá trị không phải kiểu boolean
	GetBool(key string) bool

	// GetInt lấy giá trị int cho một khóa từ context.
	//
	// Parameters:
	//   - key: Khóa cần truy xuất giá trị
	//
	// Returns:
	//   - int: Giá trị int được lưu trữ, hoặc 0 nếu khóa không tồn tại hoặc giá trị không phải kiểu int
	GetInt(key string) int

	// GetInt64 lấy giá trị int64 cho một khóa từ context.
	//
	// Parameters:
	//   - key: Khóa cần truy xuất giá trị
	//
	// Returns:
	//   - int64: Giá trị int64 được lưu trữ, hoặc 0 nếu khóa không tồn tại hoặc giá trị không phải kiểu int64
	GetInt64(key string) int64

	// GetFloat64 lấy giá trị float64 cho một khóa từ context.
	//
	// Parameters:
	//   - key: Khóa cần truy xuất giá trị
	//
	// Returns:
	//   - float64: Giá trị float64 được lưu trữ, hoặc 0.0 nếu khóa không tồn tại hoặc giá trị không phải kiểu float64
	GetFloat64(key string) float64

	// GetTime lấy giá trị time.Time cho một khóa từ context.
	//
	// Parameters:
	//   - key: Khóa cần truy xuất giá trị
	//
	// Returns:
	//   - time.Time: Giá trị time.Time được lưu trữ, hoặc time.Time zero value nếu khóa không tồn tại hoặc giá trị không phải kiểu time.Time
	GetTime(key string) time.Time

	// GetDuration lấy giá trị duration cho một khóa từ context.
	//
	// Parameters:
	//   - key: Khóa cần truy xuất giá trị
	//
	// Returns:
	//   - time.Duration: Giá trị time.Duration được lưu trữ, hoặc 0 nếu khóa không tồn tại hoặc giá trị không phải kiểu time.Duration
	GetDuration(key string) time.Duration

	// GetStringSlice lấy giá trị []string cho một khóa từ context.
	//
	// Parameters:
	//   - key: Khóa cần truy xuất giá trị
	//
	// Returns:
	//   - []string: Giá trị []string được lưu trữ, hoặc nil nếu khóa không tồn tại hoặc giá trị không phải kiểu []string
	GetStringSlice(key string) []string

	// GetStringMap lấy giá trị map[string]interface{} cho một khóa từ context.
	//
	// Parameters:
	//   - key: Khóa cần truy xuất giá trị
	//
	// Returns:
	//   - map[string]interface{}: Giá trị map[string]interface{} được lưu trữ, hoặc nil nếu khóa không tồn tại hoặc giá trị không phải kiểu map[string]interface{}
	GetStringMap(key string) map[string]interface{}

	// GetStringMapString lấy giá trị map[string]string cho một khóa từ context.
	//
	// Parameters:
	//   - key: Khóa cần truy xuất giá trị
	//
	// Returns:
	//   - map[string]string: Giá trị map[string]string được lưu trữ, hoặc nil nếu khóa không tồn tại hoặc giá trị không phải kiểu map[string]string
	GetStringMapString(key string) map[string]string

	// Method trả về HTTP method của request.
	//
	// Returns:
	//   - string: HTTP method (GET, POST, PUT, DELETE, v.v...)
	Method() string

	// Path trả về đường dẫn HTTP của request.
	// Path là phần của URL sau domain và trước query string.
	//
	// Returns:
	//   - string: Đường dẫn HTTP
	Path() string

	// RawPath trả về đường dẫn gốc HTTP của request.
	// RawPath giữ nguyên tất cả các ký tự encoding trong URL.
	//
	// Returns:
	//   - string: Đường dẫn gốc HTTP chưa được decode
	RawPath() string

	// Param trả về giá trị tham số route.
	// Tham số route là các biến động được định nghĩa trong đường dẫn,
	// ví dụ: "/users/:id" với "id" là tham số route.
	//
	// Parameters:
	//   - name: Tên của tham số route cần truy xuất
	//
	// Returns:
	//   - string: Giá trị của tham số route, hoặc chuỗi rỗng nếu không tìm thấy
	Param(name string) string

	// ParamMap trả về map các tham số route.
	//
	// Returns:
	//   - map[string]string: Map các tham số route với key là tên tham số và value là giá trị
	ParamMap() map[string]string

	// ParamArray trả về mảng các giá trị cho một tham số route.
	// Hữu ích khi một tham số route có thể xuất hiện nhiều lần.
	//
	// Parameters:
	//   - name: Tên của tham số route cần truy xuất
	//
	// Returns:
	//   - []string: Mảng các giá trị của tham số route
	ParamArray(name string) []string

	// Query trả về giá trị tham số query.
	// Tham số query là các tham số được truyền trong URL sau dấu "?".
	//
	// Parameters:
	//   - name: Tên của tham số query cần truy xuất
	//
	// Returns:
	//   - string: Giá trị đầu tiên của tham số query, hoặc chuỗi rỗng nếu không tìm thấy
	Query(name string) string

	// DefaultQuery trả về giá trị tham số query với giá trị mặc định.
	// Sử dụng khi cần một giá trị mặc định nếu tham số query không tồn tại.
	//
	// Parameters:
	//   - name: Tên của tham số query cần truy xuất
	//   - defaultValue: Giá trị mặc định sẽ được trả về nếu tham số không tồn tại
	//
	// Returns:
	//   - string: Giá trị của tham số query hoặc giá trị mặc định nếu không tìm thấy
	DefaultQuery(name, defaultValue string) string

	// QueryArray trả về mảng các giá trị cho một tham số query.
	// Hữu ích khi tham số query xuất hiện nhiều lần trong URL.
	//
	// Parameters:
	//   - name: Tên của tham số query cần truy xuất
	//
	// Returns:
	//   - []string: Mảng các giá trị của tham số query
	QueryArray(name string) []string

	// QueryMap trả về map giá trị cho các tham số query.
	// Tìm tất cả các tham số query bắt đầu bằng prefix và trả về dưới dạng map.
	//
	// Parameters:
	//   - prefix: Tiền tố để tìm các tham số query
	//
	// Returns:
	//   - map[string]string: Map các tham số query với key là phần còn lại sau prefix và value là giá trị
	QueryMap(prefix string) map[string]string

	// Form trả về giá trị form.
	// Truy xuất giá trị từ form data, hỗ trợ cả application/x-www-form-urlencoded và multipart/form-data.
	//
	// Parameters:
	//   - name: Tên của form field cần truy xuất
	//
	// Returns:
	//   - string: Giá trị đầu tiên của form field, hoặc chuỗi rỗng nếu không tìm thấy
	Form(name string) string

	// DefaultForm trả về giá trị form với giá trị mặc định.
	// Sử dụng khi cần một giá trị mặc định nếu form field không tồn tại.
	//
	// Parameters:
	//   - name: Tên của form field cần truy xuất
	//   - defaultValue: Giá trị mặc định sẽ được trả về nếu form field không tồn tại
	//
	// Returns:
	//   - string: Giá trị của form field hoặc giá trị mặc định nếu không tìm thấy
	DefaultForm(name, defaultValue string) string

	// FormArray trả về mảng các giá trị cho một form field.
	// Hữu ích khi form field xuất hiện nhiều lần.
	//
	// Parameters:
	//   - name: Tên của form field cần truy xuất
	//
	// Returns:
	//   - []string: Mảng các giá trị của form field
	FormArray(name string) []string

	// FormMap trả về map giá trị cho các form fields.
	// Tìm tất cả các form fields bắt đầu bằng prefix và trả về dưới dạng map.
	//
	// Parameters:
	//   - prefix: Tiền tố để tìm các form fields
	//
	// Returns:
	//   - map[string]string: Map các form fields với key là phần còn lại sau prefix và value là giá trị
	FormMap(prefix string) map[string]string

	// MultipartForm trả về multipart form.
	// Phân tích multipart form data từ request.
	//
	// Returns:
	//   - *multipart.Form: Đối tượng multipart.Form chứa các form fields và files
	//   - error: Lỗi nếu có trong quá trình phân tích form data
	//
	// Errors:
	//   - http: "Bad Request" nếu không thể phân tích form data
	MultipartForm() (*multipart.Form, error)

	// FormFile trả về file tải lên.
	// Truy xuất file được tải lên từ multipart form.
	//
	// Parameters:
	//   - name: Tên của form field chứa file
	//
	// Returns:
	//   - *multipart.FileHeader: Thông tin về file đã tải lên
	//   - error: Lỗi nếu không thể truy xuất file
	//
	// Errors:
	//   - http: "Bad Request" nếu không thể phân tích form data
	//   - http: "Bad Request" nếu không tìm thấy file
	FormFile(name string) (*multipart.FileHeader, error)

	// SaveUploadedFile lưu file tải lên vào đường dẫn.
	// Lưu file đã được tải lên từ multipart form vào hệ thống tệp.
	//
	// Parameters:
	//   - file: FileHeader chứa thông tin về file cần lưu
	//   - dst: Đường dẫn đích để lưu file
	//
	// Returns:
	//   - error: Lỗi nếu có trong quá trình lưu file
	//
	// Errors:
	//   - io: Các lỗi liên quan đến thao tác file
	SaveUploadedFile(file *multipart.FileHeader, dst string) error

	// BindJSON bind request body vào struct sử dụng JSON.
	// Đọc dữ liệu từ request body và chuyển đổi thành struct thông qua JSON unmarshaling.
	//
	// Parameters:
	//   - obj: Con trỏ đến struct nhận dữ liệu từ JSON
	//
	// Returns:
	//   - error: Lỗi khi parse body hoặc unmarshal JSON
	//
	// Errors:
	//   - io: Lỗi khi đọc request body
	//   - json: Lỗi khi unmarshal dữ liệu JSON
	BindJSON(obj interface{}) error

	// BindXML bind request body vào struct sử dụng XML.
	// Đọc dữ liệu từ request body và chuyển đổi thành struct thông qua XML unmarshaling.
	//
	// Parameters:
	//   - obj: Con trỏ đến struct nhận dữ liệu từ XML
	//
	// Returns:
	//   - error: Lỗi khi parse body hoặc unmarshal XML
	//
	// Errors:
	//   - io: Lỗi khi đọc request body
	//   - xml: Lỗi khi unmarshal dữ liệu XML
	BindXML(obj interface{}) error

	// BindQuery bind query parameters vào struct.
	// Map các query parameters từ URL vào struct sử dụng tag "form" hoặc "json" trên struct fields.
	//
	// Parameters:
	//   - obj: Con trỏ đến struct nhận dữ liệu từ query parameters
	//
	// Returns:
	//   - error: Lỗi khi bind các query parameters vào struct
	//
	// Errors:
	//   - binding: Lỗi khi chuyển đổi kiểu dữ liệu
	BindQuery(obj interface{}) error

	// BindForm bind form values vào struct.
	// Map các giá trị form từ request vào struct sử dụng tag "form" hoặc "json" trên struct fields.
	//
	// Parameters:
	//   - obj: Con trỏ đến struct nhận dữ liệu từ form
	//
	// Returns:
	//   - error: Lỗi khi bind form values vào struct
	//
	// Errors:
	//   - http: Lỗi khi parse form
	//   - binding: Lỗi khi chuyển đổi kiểu dữ liệu
	BindForm(obj interface{}) error

	// Bind bind request body vào struct dựa vào Content-Type.
	// Tự động chọn phương thức binding dựa vào Content-Type của request.
	// Hỗ trợ các định dạng: JSON, XML, form data.
	//
	// Parameters:
	//   - obj: Con trỏ đến struct nhận dữ liệu
	//
	// Returns:
	//   - error: Lỗi khi bind dữ liệu vào struct
	//
	// Errors:
	//   - ErrUnsupportedBinding: Content-Type không được hỗ trợ
	//   - binding: Lỗi từ phương thức binding tương ứng
	Bind(obj interface{}) error

	// ShouldBind bind request body vào struct và trả về lỗi.
	// Hoạt động tương tự như Bind nhưng được thiết kế để sử dụng trong handler mà không tự động trả về lỗi HTTP.
	//
	// Parameters:
	//   - obj: Con trỏ đến struct nhận dữ liệu
	//
	// Returns:
	//   - error: Lỗi khi bind dữ liệu vào struct
	//
	// Errors:
	//   - ErrUnsupportedBinding: Content-Type không được hỗ trợ
	//   - binding: Lỗi từ phương thức binding tương ứng
	ShouldBind(obj interface{}) error

	// Status thiết lập HTTP status code cho response.
	// Đặt status code HTTP cho response được trả về.
	//
	// Parameters:
	//   - code: HTTP status code (ví dụ: 200, 404, 500)
	Status(code int)

	// Header thiết lập header response.
	// Đặt giá trị cho header trong HTTP response.
	//
	// Parameters:
	//   - key: Tên của header
	//   - value: Giá trị của header
	Header(key, value string)

	// GetHeader trả về giá trị của header request theo tên.
	//
	// Phương thức này lấy giá trị của HTTP header từ request hiện tại dựa trên
	// khóa được cung cấp. Nếu header không tồn tại, trả về chuỗi rỗng.
	//
	// Parameters:
	//   - key: Tên của HTTP header cần lấy giá trị
	//
	// Returns:
	//   - string: Giá trị của header, hoặc chuỗi rỗng nếu không tìm thấy
	GetHeader(key string) string

	// Cookie trả về giá trị của cookie từ request dựa theo tên.
	//
	// Phương thức này tìm kiếm HTTP cookie trong request hiện tại bằng cách sử dụng tên
	// được cung cấp và trả về giá trị của nó. Nếu cookie không tồn tại, trả về lỗi.
	//
	// Parameters:
	//   - name: Tên của cookie cần truy xuất
	//
	// Returns:
	//   - string: Giá trị của cookie nếu tìm thấy
	//   - error: Lỗi nếu cookie không tồn tại hoặc có vấn đề khi đọc cookie
	//
	// Errors:
	//   - http.ErrNoCookie: Khi không tìm thấy cookie với tên chỉ định
	Cookie(name string) (string, error)

	// SetCookie thiết lập cookie trong HTTP response.
	//
	// Phương thức này tạo một HTTP cookie mới với các thông số được cung cấp và
	// thêm nó vào header của response. Cookie có thể được cấu hình với tuổi thọ,
	// domain, path và các thuộc tính bảo mật.
	//
	// Parameters:
	//   - name: Tên của cookie
	//   - value: Giá trị của cookie
	//   - maxAge: Thời gian sống tối đa của cookie tính bằng giây, 0 cho session cookie, âm để xóa cookie
	//   - path: Đường dẫn mà cookie có hiệu lực, "/" cho toàn bộ domain
	//   - domain: Domain mà cookie có hiệu lực, rỗng cho host hiện tại
	//   - secure: Chỉ gửi cookie qua kết nối HTTPS nếu là true
	//   - httpOnly: Ngăn JavaScript truy cập cookie nếu là true
	SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool)

	// Cookies trả về tất cả cookies từ request hiện tại.
	//
	// Phương thức này trích xuất tất cả HTTP cookies có trong request và
	// trả về dưới dạng một mảng các đối tượng http.Cookie.
	//
	// Returns:
	//   - []*http.Cookie: Mảng các đối tượng cookie, có thể rỗng nếu không có cookies
	Cookies() []*http.Cookie

	// Render renders một template với dữ liệu và thiết lập HTTP status code.
	//
	// Phương thức này được thiết kế để render template với tên và dữ liệu được cung cấp.
	// Hiện tại chưa có triển khai đầy đủ và chỉ thiết lập status code.
	//
	// Parameters:
	//   - code: HTTP status code cho response
	//   - name: Tên template cần render
	//   - data: Dữ liệu được truyền vào template
	//
	// TODO: Cần triển khai đầy đủ chức năng rendering template
	Render(code int, name string, data interface{})

	// HTML renders nội dung HTML với status code đã cho.
	//
	// Phương thức này thiết lập Content-Type phù hợp cho HTML, thiết lập HTTP status code
	// và ghi chuỗi HTML vào response body.
	//
	// Parameters:
	//   - code: HTTP status code cho response
	//   - html: Chuỗi HTML để render
	HTML(code int, html string)

	// String renders nội dung text với status code cho trước.
	//
	// Phương thức này thiết lập Content-Type cho plain text, thiết lập HTTP status code
	// và ghi chuỗi text vào response body. Nó hỗ trợ định dạng chuỗi tương tự như fmt.Sprintf.
	//
	// Parameters:
	//   - code: HTTP status code cho response
	//   - format: Chuỗi định dạng, tương tự như trong fmt.Sprintf
	//   - values: Các tham số tùy chọn được sử dụng để định dạng chuỗi
	String(code int, format string, values ...interface{})

	// JSON chuyển đổi một đối tượng thành định dạng JSON và ghi vào response.
	//
	// Phương thức này thiết lập Content-Type phù hợp cho JSON, thiết lập HTTP status code
	// và chuyển đổi đối tượng được cung cấp thành JSON rồi ghi vào response body.
	// Nếu quá trình encoding gặp lỗi, lỗi sẽ được xử lý qua phương thức Error.
	//
	// Parameters:
	//   - code: HTTP status code cho response
	//   - obj: Đối tượng cần chuyển đổi thành JSON
	JSON(code int, obj interface{})

	// JSONP chuyển đổi một đối tượng thành định dạng JSONP và ghi vào response.
	//
	// Phương thức này thiết lập Content-Type phù hợp cho JavaScript, thiết lập HTTP status code,
	// và đóng gói đối tượng được cung cấp trong một hàm callback JSONP. Được sử dụng để
	// khắc phục hạn chế same-origin policy khi gọi API từ domain khác thông qua JavaScript.
	//
	// Parameters:
	//   - code: HTTP status code cho response
	//   - callback: Tên hàm JavaScript được sử dụng để đóng gói kết quả JSON
	//   - obj: Đối tượng cần chuyển đổi thành JSON và đóng gói trong callback
	JSONP(code int, callback string, obj interface{})

	// XML renders dữ liệu dạng XML.
	// Chuyển đổi object thành XML và trả về với Content-Type là "application/xml".
	//
	// Parameters:
	//   - code: HTTP status code cho response
	//   - obj: Đối tượng cần chuyển đổi thành XML
	//
	// Errors:
	//   - Không trả về lỗi trực tiếp, nhưng gọi c.Error() nếu encoding thất bại
	XML(code int, obj interface{})

	// File phục vụ một file từ filesystem.
	// Đọc và trả về nội dung của file từ đường dẫn được chỉ định.
	//
	// Parameters:
	//   - filepath: Đường dẫn đến file cần phục vụ
	//
	// Errors:
	//   - Không trả về lỗi trực tiếp, nhưng sẽ trả về HTTP 404 nếu không tìm thấy file
	File(filepath string)

	// FileAttachment phục vụ một file từ filesystem với tên tùy chỉnh.
	// Phục vụ file nhưng với header Content-Disposition để client tải xuống với tên được chỉ định.
	//
	// Parameters:
	//   - filepath: Đường dẫn đến file cần phục vụ
	//   - filename: Tên file được hiển thị khi người dùng tải xuống
	//
	// Errors:
	//   - Không trả về lỗi trực tiếp, nhưng sẽ trả về HTTP 404 nếu không tìm thấy file
	FileAttachment(filepath, filename string)

	// Blob phục vụ dữ liệu nhị phân từ bộ nhớ với content type.
	// Trả về một mảng byte với Content-Type được chỉ định.
	//
	// Parameters:
	//   - code: HTTP status code cho response
	//   - contentType: Kiểu dữ liệu cho Content-Type header
	//   - data: Mảng byte chứa dữ liệu cần trả về
	Blob(code int, contentType string, data []byte)

	// Stream phục vụ dữ liệu từ io.Reader với content type.
	// Trả về dữ liệu từ một stream với Content-Type được chỉ định.
	//
	// Parameters:
	//   - code: HTTP status code cho response
	//   - contentType: Kiểu dữ liệu cho Content-Type header
	//   - r: Reader chứa dữ liệu cần trả về
	//
	// Errors:
	//   - io: Các lỗi từ Reader được truyền vào không được xử lý
	Stream(code int, contentType string, r io.Reader)

	// Redirect thực hiện HTTP redirect.
	// Chuyển hướng client đến một URL mới với status code được chỉ định.
	//
	// Parameters:
	//   - code: HTTP status code cho redirect (thường là 301, 302, 303, 307, 308)
	//   - location: URL đích cho redirect
	Redirect(code int, location string)

	// Error trả về một HTTP error với status code và message.
	// Trả về lỗi HTTP với status code 500 và message từ error.
	//
	// Parameters:
	//   - err: Error object chứa thông tin lỗi
	Error(err error)

	// ClientIP trả về địa chỉ IP của client.
	// Xác định IP của client bằng cách kiểm tra các header X-Forwarded-For và X-Real-IP
	// trước khi sử dụng địa chỉ từ kết nối trực tiếp.
	//
	// Returns:
	//   - string: Địa chỉ IP của client
	ClientIP() string

	// ContentType trả về Content-Type của request.
	// Lấy giá trị của header Content-Type từ request.
	//
	// Returns:
	//   - string: Giá trị của Content-Type header
	ContentType() string

	// IsWebsocket kiểm tra xem request có phải là websocket không.
	// Xác định xem request hiện tại có phải là websocket connection request
	// bằng cách kiểm tra các header Upgrade và Connection.
	//
	// Returns:
	//   - bool: true nếu request là websocket, ngược lại là false
	IsWebsocket() bool

	// GetRawData trả về raw request body.
	// Đọc và trả về toàn bộ nội dung của request body.
	//
	// Returns:
	//   - []byte: Dữ liệu từ request body
	//   - error: Lỗi nếu có khi đọc body
	//
	// Errors:
	//   - io: Lỗi khi đọc từ body
	GetRawData() ([]byte, error)

	// Handlers trả về tất cả các middleware handlers.
	// Lấy danh sách các handlers được đăng ký cho route hiện tại.
	//
	// Returns:
	//   - []func(Context): Mảng các handler functions
	Handlers() []func(Context)

	// SetHandlers thiết lập handlers mới.
	// Thay thế danh sách các middleware handlers hiện tại với handlers mới.
	//
	// Parameters:
	//   - handlers: Mảng các handler functions mới
	SetHandlers(handlers []func(Context))

	// ValidateStruct kiểm tra tính hợp lệ của một struct sử dụng validator.
	// Sử dụng thư viện validator.v10 để kiểm tra struct dựa trên validation tags.
	//
	// Parameters:
	//   - obj: Struct cần validation
	//
	// Returns:
	//   - error: Lỗi validation nếu struct không hợp lệ
	//
	// Errors:
	//   - validator.ValidationErrors: Chi tiết về các trường không hợp lệ
	ValidateStruct(obj interface{}) error

	// ShouldBindAndValidate thực hiện bind và validate struct từ request.
	// Thực hiện binding dữ liệu từ request vào struct và sau đó validate struct.
	// Không tự động trả về lỗi HTTP như BindAndValidate.
	//
	// Parameters:
	//   - obj: Struct cần binding và validation
	//
	// Returns:
	//   - error: Lỗi binding hoặc validation nếu có
	//
	// Errors:
	//   - binding: Lỗi từ quá trình bind request data
	//   - validator.ValidationErrors: Lỗi từ quá trình validate
	ShouldBindAndValidate(obj interface{}) error

	// BindAndValidate thực hiện bind và validate struct từ request.
	// Tương tự ShouldBindAndValidate nhưng tự động trả về lỗi HTTP trong trường hợp thất bại
	// và sẽ thiết lập response status code và body phù hợp.
	//
	// Parameters:
	//   - obj: Struct cần binding và validation
	//
	// Returns:
	//   - error: HTTPError object từ fork/errors nếu binding hoặc validation thất bại
	//
	// Errors:
	//   - forkerrors.BadRequest: Lỗi khi binding request data
	//   - forkerrors.UnprocessableEntity: Lỗi khi validate dữ liệu
	BindAndValidate(obj interface{}) error

	// RegisterValidation đăng ký một hàm validation tùy chỉnh.
	// Thêm một custom validation tag và hàm validation tương ứng vào validator.
	//
	// Parameters:
	//   - tag: Tag name sẽ được sử dụng trong struct tag
	//   - fn: Hàm validation tương ứng với tag
	//
	// Returns:
	//   - error: Lỗi nếu không thể đăng ký validation
	RegisterValidation(tag string, fn validator.Func) error

	// GetValidator trả về validator instance để cấu hình nâng cao.
	// Cho phép truy cập trực tiếp đến validator instance để thực hiện cấu hình nâng cao.
	//
	// Returns:
	//   - *validator.Validate: Instance của validator
	GetValidator() *validator.Validate
}

// ErrUnsupportedBinding là lỗi được trả về khi Content-Type không được hỗ trợ.
// Lỗi này được sử dụng trong phương thức Bind khi không thể xác định
// phương thức binding phù hợp dựa vao Content-Type.
var ErrUnsupportedBinding = errors.New("unsupported binding type")
