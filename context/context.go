package context

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/go-playground/validator/v10"
	forkerrors "go.fork.vn/fork/errors"
)

// forkContext là implementation private cho Context interface.
//
// Triển khai giao diện Context, chứa tất cả trạng thái và hành vi cần thiết cho một HTTP request lifecycle.
// Chỉ được khởi tạo thông qua hàm NewContext để đảm bảo an toàn và khởi tạo đúng.
type forkContext struct {
	// request chứa thông tin về HTTP request hiện tại
	request Request

	// response cho phép thao tác với HTTP response hiện tại
	response Response

	// ctx là context.Context gốc từ request, dùng để kiểm soát timeout, hủy bỏ, truyền dữ liệu giữa các goroutine
	ctx context.Context

	// params chứa các tham số từ URL path (route parameters)
	params map[string]string

	// handlers là mảng các middleware functions cho request hiện tại
	handlers []func(Context)

	// index là vị trí hiện tại trong chuỗi handlers, dùng để điều khiển flow middleware
	index int

	// aborted đánh dấu trạng thái đã dừng thực thi handlers hay chưa
	aborted bool

	// store lưu trữ dữ liệu tùy chỉnh trong phạm vi của request (key-value)
	store map[string]interface{}

	// validator dùng để xác thực struct theo validation tags
	validator *validator.Validate
}

// NewContext tạo một context mới cho mỗi HTTP request.
//
// Hàm này khởi tạo và trả về một Context mới từ HTTP request và response.
// Thiết lập validator và các cấu hình mặc định cho context.
//
// Params:
//   - w: http.ResponseWriter để ghi HTTP response
//   - r: *http.Request chứa thông tin request
//
// Returns:
//   - Context: Context mới đã được khởi tạo
func NewContext(w http.ResponseWriter, r *http.Request) Context {
	// Khởi tạo validator với cấu hình mặc định
	validate := validator.New()

	// Đăng ký hàm định dạng lỗi tùy chỉnh
	// Ưu tiên sử dụng tên từ tag json, sau đó là form, cuối cùng là tên trường
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("json")
		if name == "" {
			name = fld.Tag.Get("form")
		}
		if name == "" {
			name = fld.Name
		}
		return name
	})

	// Lưu ý: Không còn sử dụng translator nữa

	return &forkContext{
		request:   NewRequest(r),
		response:  NewResponse(w),
		ctx:       r.Context(),
		params:    make(map[string]string),
		handlers:  nil,
		index:     -1,
		aborted:   false,
		store:     make(map[string]interface{}),
		validator: validate,
	}
}

// Request trả về đối tượng Request hiện tại.
//
// Returns:
//   - Request: Đối tượng Request trừu tượng hóa http.Request
func (c *forkContext) Request() Request {
	return c.request
}

// Response trả về đối tượng Response hiện tại.
//
// Returns:
//   - Response: Đối tượng Response trừu tượng hóa http.ResponseWriter
func (c *forkContext) Response() Response {
	return c.response
}

// Context trả về context.Context gốc từ request.
//
// Returns:
//   - context.Context: Context gốc của request
func (c *forkContext) Context() context.Context {
	return c.ctx
}

// WithContext thay thế context.Context hiện tại và trả về context đã cập nhật.
//
// Params:
//   - ctx: context.Context mới
//
// Returns:
//   - Context: Context đã cập nhật context.Context
func (c *forkContext) WithContext(ctx context.Context) Context {
	c.ctx = ctx
	return c
}

// Next thực thi handler tiếp theo trong chuỗi middleware.
//
// Được sử dụng để chuyển điều khiển đến middleware tiếp theo trong pipeline.
// Nếu đã gọi Abort thì sẽ dừng lại.
func (c *forkContext) Next() {
	// Tăng index để trỏ đến handler tiếp theo
	c.index++
	// Thực thi tất cả handlers còn lại cho đến khi kết thúc hoặc bị abort
	for c.index < len(c.handlers) && !c.aborted {
		c.handlers[c.index](c)
		c.index++
	}
}

// Abort đánh dấu context là đã bị abort để dừng thực thi các handlers tiếp theo.
func (c *forkContext) Abort() {
	c.aborted = true
}

// IsAborted kiểm tra context có đã bị abort hay không.
//
// Returns:
//   - bool: true nếu đã bị abort, ngược lại là false
func (c *forkContext) IsAborted() bool {
	return c.aborted
}

// Set lưu trữ một giá trị vào context với key được chỉ định.
//
// Params:
//   - key: Tên key
//   - value: Giá trị lưu trữ (interface{})
func (c *forkContext) Set(key string, value interface{}) {
	c.store[key] = value
}

// Get lấy giá trị từ context dựa theo key.
//
// Params:
//   - key: Tên key
//
// Returns:
//   - interface{}: Giá trị lưu trữ
//   - bool: true nếu tồn tại, false nếu không
func (c *forkContext) Get(key string) (interface{}, bool) {
	value, exists := c.store[key]
	return value, exists
}

// GetString lấy giá trị string từ context dựa theo key.
//
// Params:
//   - key: Tên key
//
// Returns:
//   - string: Giá trị string, trả về "" nếu không tồn tại hoặc không phải kiểu string
func (c *forkContext) GetString(key string) string {
	// Lấy giá trị từ context
	if val, ok := c.Get(key); ok {
		// Kiểm tra và ép kiểu về string
		if str, ok := val.(string); ok {
			return str
		}
	}
	// Trả về chuỗi rỗng nếu không tìm thấy hoặc không phải kiểu string
	return ""
}

// GetBool lấy giá trị boolean từ context dựa theo key.
//
// Params:
//   - key: Tên key
//
// Returns:
//   - bool: Giá trị boolean, trả về false nếu không tồn tại hoặc không phải kiểu boolean
func (c *forkContext) GetBool(key string) bool {
	// Lấy giá trị từ context
	if val, ok := c.Get(key); ok {
		// Kiểm tra và ép kiểu về boolean
		if b, ok := val.(bool); ok {
			return b
		}
	}
	// Trả về false nếu không tìm thấy hoặc không phải kiểu boolean
	return false
}

// GetInt lấy giá trị số nguyên (int) từ context dựa theo key.
//
// Params:
//   - key: Tên key
//
// Returns:
//   - int: Giá trị int, trả về 0 nếu không tồn tại hoặc không phải kiểu int
func (c *forkContext) GetInt(key string) int {
	// Lấy giá trị từ context
	if val, ok := c.Get(key); ok {
		// Kiểm tra và ép kiểu về int
		if i, ok := val.(int); ok {
			return i
		}
	}
	// Trả về 0 nếu không tìm thấy hoặc không phải kiểu int
	return 0
}

// GetInt64 lấy giá trị số nguyên 64-bit từ context dựa theo key.
//
// Params:
//   - key: Tên key
//
// Returns:
//   - int64: Giá trị int64, trả về 0 nếu không tồn tại hoặc không phải kiểu int64
func (c *forkContext) GetInt64(key string) int64 {
	// Lấy giá trị từ context
	if val, ok := c.Get(key); ok {
		// Kiểm tra và ép kiểu về int64
		if i, ok := val.(int64); ok {
			return i
		}
	}
	return 0
}

// GetFloat64 lấy giá trị số thực 64-bit từ context dựa theo key.
//
// Params:
//   - key: Tên key
//
// Returns:
//   - float64: Giá trị float64, trả về 0 nếu không tồn tại hoặc không phải kiểu float64
func (c *forkContext) GetFloat64(key string) float64 {
	// Lấy giá trị từ context
	if val, ok := c.Get(key); ok {
		// Kiểm tra và ép kiểu về float64
		if f, ok := val.(float64); ok {
			return f
		}
	}
	// Trả về 0 nếu không tìm thấy hoặc không phải kiểu float64
	return 0
}

// GetTime trả về giá trị thời gian được lưu trữ trong context với key cho trước.
//
// Params:
//   - key: Tên key
//
// Returns:
//   - time.Time: Giá trị thời gian, trả về time.Time{} nếu không tồn tại hoặc không đúng kiểu
func (c *forkContext) GetTime(key string) time.Time {
	// Lấy giá trị từ store
	if val, ok := c.Get(key); ok {
		// Kiểm tra kiểu dữ liệu và thực hiện type assertion
		if t, ok := val.(time.Time); ok {
			return t
		}
	}
	// Trả về giá trị mặc định nếu không tìm thấy hoặc không thể chuyển đổi
	return time.Time{}
}

// GetDuration trả về giá trị thời gian (duration) được lưu trữ trong context với key cho trước.
//
// Params:
//   - key: Tên key
//
// Returns:
//   - time.Duration: Giá trị duration, trả về 0 nếu không tồn tại hoặc không đúng kiểu
func (c *forkContext) GetDuration(key string) time.Duration {
	// Lấy giá trị từ store
	if val, ok := c.Get(key); ok {
		// Kiểm tra kiểu dữ liệu và thực hiện type assertion
		if d, ok := val.(time.Duration); ok {
			return d
		}
	}
	// Trả về giá trị mặc định nếu không tìm thấy hoặc không thể chuyển đổi
	return 0
}

// GetStringSlice trả về mảng chuỗi được lưu trữ trong context với key cho trước.
//
// Params:
//   - key: Tên key
//
// Returns:
//   - []string: Mảng chuỗi, trả về nil nếu không tồn tại hoặc không đúng kiểu
func (c *forkContext) GetStringSlice(key string) []string {
	// Lấy giá trị từ store
	if val, ok := c.Get(key); ok {
		// Kiểm tra kiểu dữ liệu và thực hiện type assertion
		if ss, ok := val.([]string); ok {
			return ss
		}
	}
	// Trả về giá trị mặc định nếu không tìm thấy hoặc không thể chuyển đổi
	return nil
}

// GetStringMap trả về map[string]interface{} được lưu trữ trong context với key cho trước.
//
// Params:
//   - key: Tên key
//
// Returns:
//   - map[string]interface{}: Map, trả về nil nếu không tồn tại hoặc không đúng kiểu
func (c *forkContext) GetStringMap(key string) map[string]interface{} {
	// Lấy giá trị từ store
	if val, ok := c.Get(key); ok {
		// Kiểm tra kiểu dữ liệu và thực hiện type assertion
		if sm, ok := val.(map[string]interface{}); ok {
			return sm
		}
	}
	// Trả về giá trị mặc định nếu không tìm thấy hoặc không thể chuyển đổi
	return nil
}

// GetStringMapString trả về map[string]string được lưu trữ trong context với key cho trước.
//
// Params:
//   - key: Tên key
//
// Returns:
//   - map[string]string: Map, trả về nil nếu không tồn tại hoặc không đúng kiểu
func (c *forkContext) GetStringMapString(key string) map[string]string {
	// Lấy giá trị từ store
	if val, ok := c.Get(key); ok {
		// Kiểm tra kiểu dữ liệu và thực hiện type assertion
		if sms, ok := val.(map[string]string); ok {
			return sms
		}
	}
	// Trả về giá trị mặc định nếu không tìm thấy hoặc không thể chuyển đổi
	return nil
}

// Method trả về HTTP method của request hiện tại.
//
// Returns:
//   - string: HTTP method (GET, POST, ...)
func (c *forkContext) Method() string {
	return c.request.Method()
}

// Path trả về đường dẫn URL của request hiện tại.
//
// Returns:
//   - string: Đường dẫn URL
func (c *forkContext) Path() string {
	return c.request.URL().Path
}

// RawPath trả về đường dẫn gốc (chưa được decode) của request hiện tại.
//
// Returns:
//   - string: Đường dẫn gốc
func (c *forkContext) RawPath() string {
	return c.request.URL().RawPath
}

// Param trả về giá trị của tham số route theo tên.
//
// Params:
//   - name: Tên tham số
//
// Returns:
//   - string: Giá trị tham số, trả về "" nếu không tồn tại
func (c *forkContext) Param(name string) string {
	return c.GetString("param:" + name)
}

// ParamMap trả về tất cả các tham số route dưới dạng map[string]string.
//
// Returns:
//   - map[string]string: Map các tham số route
func (c *forkContext) ParamMap() map[string]string {
	params := make(map[string]string)
	for key, value := range c.store {
		if len(key) > 6 && key[:6] == "param:" {
			paramName := key[6:]
			if strValue, ok := value.(string); ok {
				params[paramName] = strValue
			}
		}
	}
	return params
}

// ParamArray trả về mảng giá trị của tham số route theo tên (hiện chỉ hỗ trợ 1 giá trị).
//
// Params:
//   - name: Tên tham số
//
// Returns:
//   - []string: Mảng giá trị, trả về nil nếu không tồn tại
func (c *forkContext) ParamArray(name string) []string {
	// Current implementation only supports single values per parameter
	// This method is provided for API consistency and future extensibility
	value := c.Param(name)
	if value == "" {
		return nil
	}
	return []string{value}
}

// Query trả về giá trị query string theo tên.
//
// Params:
//   - name: Tên tham số query
//
// Returns:
//   - string: Giá trị query, trả về "" nếu không tồn tại
func (c *forkContext) Query(name string) string {
	return c.request.URL().Query().Get(name)
}

// DefaultQuery trả về giá trị query string theo tên, hoặc defaultValue nếu không tồn tại.
//
// Params:
//   - name: Tên tham số query
//   - defaultValue: Giá trị mặc định
//
// Returns:
//   - string: Giá trị query hoặc defaultValue
func (c *forkContext) DefaultQuery(name, defaultValue string) string {
	if value := c.Query(name); value != "" {
		return value
	}
	return defaultValue
}

// QueryArray trả về mảng giá trị query string theo tên.
//
// Params:
//   - name: Tên tham số query
//
// Returns:
//   - []string: Mảng giá trị query
func (c *forkContext) QueryArray(name string) []string {
	return c.request.URL().Query()[name]
}

// QueryMap trả về map các query string có prefix chỉ định.
//
// Params:
//   - prefix: Tiền tố filter các key
//
// Returns:
//   - map[string]string: Map các query string
func (c *forkContext) QueryMap(prefix string) map[string]string {
	values := c.request.URL().Query()
	result := make(map[string]string)
	for k, v := range values {
		if len(v) > 0 {
			result[k] = v[0]
		}
	}
	return result
}

// Form trả về giá trị form field theo tên.
//
// Params:
//   - name: Tên field
//
// Returns:
//   - string: Giá trị field, trả về "" nếu không tồn tại
func (c *forkContext) Form(name string) string {
	return c.request.FormValue(name)
}

// DefaultForm trả về giá trị form field theo tên, hoặc defaultValue nếu không tồn tại.
//
// Params:
//   - name: Tên field
//   - defaultValue: Giá trị mặc định
//
// Returns:
//   - string: Giá trị field hoặc defaultValue
func (c *forkContext) DefaultForm(name, defaultValue string) string {
	if value := c.Form(name); value != "" {
		return value
	}
	return defaultValue
}

// FormArray trả về mảng giá trị form field theo tên.
//
// Params:
//   - name: Tên field
//
// Returns:
//   - []string: Mảng giá trị field
func (c *forkContext) FormArray(name string) []string {
	form, err := c.request.MultipartForm()
	if err != nil {
		return nil
	}
	return form.Value[name]
}

// FormMap trả về map các form field có prefix chỉ định.
//
// Params:
//   - prefix: Tiền tố filter các key
//
// Returns:
//   - map[string]string: Map các form field
func (c *forkContext) FormMap(prefix string) map[string]string {
	form, err := c.request.MultipartForm()
	if err != nil {
		return nil
	}
	result := make(map[string]string)
	for k, v := range form.Value {
		if len(v) > 0 {
			result[k] = v[0]
		}
	}
	return result
}

// MultipartForm trả về multipart.Form của request hiện tại.
//
// Returns:
//   - *multipart.Form: Đối tượng multipart form
//   - error: Lỗi nếu không thể parse multipart form
func (c *forkContext) MultipartForm() (*multipart.Form, error) {
	return c.request.MultipartForm()
}

// FormFile trả về file upload từ form theo tên field.
//
// Params:
//   - name: Tên field
//
// Returns:
//   - *multipart.FileHeader: Thông tin file upload
//   - error: Lỗi nếu không tìm thấy hoặc không hợp lệ
func (c *forkContext) FormFile(name string) (*multipart.FileHeader, error) {
	return c.request.FormFile(name)
}

// SaveUploadedFile lưu file upload vào đường dẫn chỉ định.
//
// Params:
//   - file: *multipart.FileHeader file upload
//   - dst: Đường dẫn lưu file
//
// Returns:
//   - error: Lỗi nếu không thể lưu file
func (c *forkContext) SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

// BindJSON đọc request body và chuyển đổi thành struct sử dụng JSON unmarshaling.
//
// Params:
//   - obj: Con trỏ struct nhận dữ liệu
//
// Returns:
//   - error: Lỗi nếu không thể đọc hoặc unmarshal JSON
func (c *forkContext) BindJSON(obj interface{}) error {
	body, err := c.GetRawData()
	if err != nil {
		return err
	}
	return json.Unmarshal(body, obj)
}

// BindXML đọc request body và chuyển đổi thành struct sử dụng XML unmarshaling.
//
// Params:
//   - obj: Con trỏ struct nhận dữ liệu
//
// Returns:
//   - error: Lỗi nếu không thể đọc hoặc unmarshal XML
func (c *forkContext) BindXML(obj interface{}) error {
	body, err := c.GetRawData()
	if err != nil {
		return err
	}
	return xml.Unmarshal(body, obj)
}

// BindQuery liên kết các tham số truy vấn URL vào một struct sử dụng function bind.
//
// Params:
//   - obj: Con trỏ struct nhận dữ liệu
//
// Returns:
//   - error: Lỗi nếu không thể bind
func (c *forkContext) BindQuery(obj interface{}) error {
	values := c.request.URL().Query()
	return bind(values, obj)
}

// BindForm phân tích form trong request và liên kết các giá trị form vào một struct.
//
// Params:
//   - obj: Con trỏ struct nhận dữ liệu
//
// Returns:
//   - error: Lỗi nếu không thể bind
func (c *forkContext) BindForm(obj interface{}) error {
	err := c.request.ParseForm()
	if err != nil {
		return err
	}
	return bind(c.request.Form(), obj)
}

// Bind tự động chọn phương thức binding dựa trên Content-Type của request.
//
// Params:
//   - obj: Con trỏ struct nhận dữ liệu
//
// Returns:
//   - error: Lỗi nếu không hỗ trợ Content-Type hoặc bind thất bại
//
// Exceptions:
//   - ErrUnsupportedBinding: Nếu Content-Type không được hỗ trợ
func (c *forkContext) Bind(obj interface{}) error {
	// Lấy Content-Type của request
	contentType := c.ContentType()
	// Chọn phương thức binding phù hợp dựa vào Content-Type
	switch contentType {
	case "application/json":
		return c.BindJSON(obj)
	case "application/xml", "text/xml":
		return c.BindXML(obj)
	case "application/x-www-form-urlencoded", "multipart/form-data":
		return c.BindForm(obj)
	}
	// Trả về lỗi nếu Content-Type không được hỗ trợ
	return ErrUnsupportedBinding
}

// ShouldBind là wrapper cho Bind, dùng trong handler.
//
// Params:
//   - obj: Con trỏ struct nhận dữ liệu
//
// Returns:
//   - error: Lỗi nếu bind thất bại
func (c *forkContext) ShouldBind(obj interface{}) error {
	return c.Bind(obj)
}

// Status thiết lập HTTP status code cho response.
//
// Params:
//   - code: HTTP status code
func (c *forkContext) Status(code int) {
	c.response.WriteHeader(code)
}

// Header thiết lập header cho HTTP response.
//
// Params:
//   - key: Tên header
//   - value: Giá trị header
func (c *forkContext) Header(key, value string) {
	c.response.Header().Set(key, value)
}

// GetHeader trả về giá trị của header request theo tên.
//
// Params:
//   - key: Tên header
//
// Returns:
//   - string: Giá trị header, trả về "" nếu không tìm thấy
func (c *forkContext) GetHeader(key string) string {
	return c.request.Header().Get(key)
}

// Cookie trả về giá trị của cookie từ request dựa theo tên.
//
// Params:
//   - name: Tên cookie
//
// Returns:
//   - string: Giá trị cookie nếu tìm thấy
//   - error: Lỗi nếu không tìm thấy hoặc có vấn đề khi đọc cookie
//
// Errors:
//   - http.ErrNoCookie: Khi không tìm thấy cookie với tên chỉ định
func (c *forkContext) Cookie(name string) (string, error) {
	// Lấy cookie từ request
	cookie, err := c.request.Cookie(name)
	if err != nil {
		return "", err
	}
	// Trả về giá trị của cookie
	return cookie.Value, nil
}

// SetCookie thiết lập cookie trong HTTP response.
//
// Params:
//   - name: Tên cookie
//   - value: Giá trị cookie
//   - maxAge: Thời gian sống tối đa (giây), 0 cho session cookie, âm để xóa cookie
//   - path: Đường dẫn cookie
//   - domain: Domain cookie
//   - secure: Chỉ gửi cookie qua HTTPS nếu true
//   - httpOnly: Ngăn JavaScript truy cập cookie nếu true
func (c *forkContext) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	// Tạo đối tượng cookie với các tham số đã cung cấp
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		Secure:   secure,
		HttpOnly: httpOnly,
	}
	// Thêm cookie vào header response
	c.response.Header().Add("Set-Cookie", cookie.String())
}

// Cookies trả về tất cả cookies từ request hiện tại.
//
// Returns:
//   - []*http.Cookie: Mảng cookies, có thể rỗng nếu không có
func (c *forkContext) Cookies() []*http.Cookie {
	return c.request.Cookies()
}

// Render render một template với dữ liệu và thiết lập HTTP status code.
//
// Params:
//   - code: HTTP status code
//   - name: Tên template
//   - data: Dữ liệu truyền vào template
//
// Note: Requires templates middleware to be registered
func (c *forkContext) Render(code int, name string, data interface{}) {
	// Try to get template registry first for multi-engine support
	if registry, exists := c.Get("template_registry"); exists {
		if templateRegistry, ok := registry.(interface {
			GetEngineByTemplate(string) (interface{}, interface{}, error)
		}); ok {
			// Get the appropriate engine based on template name/extension
			engine, engineType, err := templateRegistry.GetEngineByTemplate(name)
			if err == nil {
				if templateEngine, ok := engine.(interface {
					Render(string, interface{}) ([]byte, error)
				}); ok {
					// Render template using the detected engine
					result, renderErr := templateEngine.Render(name, data)
					if renderErr != nil {
						// If template rendering fails, return error as JSON
						c.JSON(http.StatusInternalServerError, map[string]string{
							"error":       fmt.Sprintf("Template rendering failed: %v", renderErr),
							"template":    name,
							"engine_type": fmt.Sprintf("%v", engineType),
						})
						return
					}

					// Determine content type based on engine type
					contentType := "text/html; charset=utf-8"
					if engineTypeStr, ok := engineType.(fmt.Stringer); ok {
						switch engineTypeStr.String() {
						case "text":
							contentType = "text/plain; charset=utf-8"
						case "pug", "handlebars", "mustache", "django", "ace", "amber", "slim", "html":
							contentType = "text/html; charset=utf-8"
						case "jet":
							contentType = "text/html; charset=utf-8"
						default:
							contentType = "text/html; charset=utf-8"
						}
					}

					// Set content type and status, then write rendered content
					c.Header("Content-Type", contentType)
					c.Status(code)
					c.response.Write(result)
					return
				}
			}
		}
	}

	// Fallback: Try to get legacy single template engine from context
	if engine, exists := c.Get("template_engine"); exists {
		if templateEngine, ok := engine.(interface {
			Render(string, interface{}) ([]byte, error)
		}); ok {
			// Render template using the legacy engine
			result, err := templateEngine.Render(name, data)
			if err != nil {
				// If template rendering fails, return error as JSON
				c.JSON(http.StatusInternalServerError, map[string]string{
					"error":    fmt.Sprintf("Template rendering failed: %v", err),
					"template": name,
				})
				return
			}

			// Set content type and status, then write rendered content
			c.Header("Content-Type", "text/html; charset=utf-8")
			c.Status(code)
			c.response.Write(result)
			return
		}
	}

	// Fallback: if no template engine available, treat as HTML string
	c.HTML(code, fmt.Sprintf("<!-- Template '%s' not found or template engine not available -->", name))
}

// HTML render nội dung HTML với status code đã cho.
//
// Params:
//   - code: HTTP status code
//   - html: Chuỗi HTML để render
func (c *forkContext) HTML(code int, html string) {
	// Thiết lập Content-Type header cho HTML
	c.Header("Content-Type", "text/html; charset=utf-8")
	// Thiết lập HTTP status code
	c.Status(code)
	// Ghi dữ liệu HTML vào response
	c.response.Write([]byte(html))
}

// String render nội dung text với status code cho trước.
//
// Params:
//   - code: HTTP status code
//   - format: Chuỗi định dạng (fmt.Sprintf)
//   - values: Tham số định dạng
func (c *forkContext) String(code int, format string, values ...interface{}) {
	// Thiết lập Content-Type header cho plain text
	c.Header("Content-Type", "text/plain; charset=utf-8")
	// Thiết lập HTTP status code
	c.Status(code)
	// Kiểm tra xem có tham số values không
	if len(values) > 0 {
		// Nếu có, sử dụng fmt.Fprintf để format chuỗi
		fmt.Fprintf(c.response, format, values...)
	} else {
		// Nếu không, ghi chuỗi format trực tiếp
		c.response.Write([]byte(format))
	}
}

// JSON chuyển đổi một đối tượng thành JSON và ghi vào response.
//
// Params:
//   - code: HTTP status code
//   - obj: Đối tượng cần encode
//
// Nếu encode lỗi sẽ trả về error dạng JSON qua Error()
func (c *forkContext) JSON(code int, obj interface{}) {
	// Thiết lập Content-Type header cho JSON
	c.Header("Content-Type", "application/json; charset=utf-8")
	// Thiết lập HTTP status code
	c.Status(code)
	// Tạo JSON encoder sử dụng response writer
	encoder := json.NewEncoder(c.response)
	// Encode đối tượng thành JSON và ghi vào response
	if err := encoder.Encode(obj); err != nil {
		// Xử lý lỗi nếu quá trình encode thất bại
		c.Error(err)
	}
}

// JSONP chuyển đổi một đối tượng thành JSONP và ghi vào response.
//
// Params:
//   - code: HTTP status code
//   - callback: Tên hàm JavaScript callback
//   - obj: Đối tượng encode
//
// Nếu encode lỗi sẽ trả về error dạng JSON qua Error()
func (c *forkContext) JSONP(code int, callback string, obj interface{}) {
	// Thiết lập Content-Type header cho JavaScript
	c.Header("Content-Type", "application/javascript; charset=utf-8")
	// Thiết lập HTTP status code
	c.Status(code)
	// Ghi phần đầu của callback function
	c.response.Write([]byte(callback + "("))
	// Tạo JSON encoder sử dụng response writer
	encoder := json.NewEncoder(c.response)
	// Encode đối tượng thành JSON và ghi vào response
	if err := encoder.Encode(obj); err != nil {
		c.Error(err)
		return
	}
	// Ghi phần kết thúc của callback function
	c.response.Write([]byte(");"))
}

// XML chuyển đổi đối tượng thành XML và ghi vào response.
//
// Params:
//   - code: HTTP status code
//   - obj: Đối tượng encode
//
// Nếu encode lỗi sẽ trả về error dạng XML qua Error()
func (c *forkContext) XML(code int, obj interface{}) {
	// Thiết lập Content-Type header cho XML
	c.Header("Content-Type", "application/xml; charset=utf-8")
	// Thiết lập HTTP status code
	c.Status(code)
	// Tạo XML encoder sử dụng response writer
	encoder := xml.NewEncoder(c.response)
	// Encode đối tượng thành XML và ghi vào response
	if err := encoder.Encode(obj); err != nil {
		// Xử lý lỗi nếu quá trình encode thất bại
		c.Error(err)
	}
}

// File phục vụ một file từ hệ thống tệp với đường dẫn được chỉ định.
//
// Params:
//   - filepath: Đường dẫn file
func (c *forkContext) File(filepath string) {
	// Sử dụng http.ServeFile để phục vụ file
	http.ServeFile(c.response, c.request.Request(), filepath)
}

// FileAttachment phục vụ một file từ hệ thống tệp với tên file được chỉ định cho việc tải xuống.
//
// Params:
//   - filepath: Đường dẫn file
//   - filename: Tên file tải xuống
func (c *forkContext) FileAttachment(filepath, filename string) {
	// Thiết lập header Content-Disposition để hướng dẫn trình duyệt tải xuống file
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	// Phục vụ file sử dụng phương thức File
	c.File(filepath)
}

// Blob phục vụ dữ liệu nhị phân với kiểu nội dung được chỉ định.
//
// Params:
//   - code: HTTP status code
//   - contentType: Kiểu nội dung
//   - data: Dữ liệu nhị phân
func (c *forkContext) Blob(code int, contentType string, data []byte) {
	// Thiết lập Content-Type header
	c.Header("Content-Type", contentType)
	// Thiết lập HTTP status code
	c.Status(code)
	// Ghi dữ liệu nhị phân vào response
	c.response.Write(data)
}

// Stream phục vụ dữ liệu từ một io.Reader với kiểu nội dung được chỉ định.
//
// Params:
//   - code: HTTP status code
//   - contentType: Kiểu nội dung
//   - r: io.Reader nguồn dữ liệu
func (c *forkContext) Stream(code int, contentType string, r io.Reader) {
	// Thiết lập Content-Type header
	c.Header("Content-Type", contentType)
	// Thiết lập HTTP status code
	c.Status(code)
	// Sao chép dữ liệu từ reader vào response
	io.Copy(c.response, r)
}

// Redirect thực hiện chuyển hướng HTTP đến địa chỉ được chỉ định.
//
// Params:
//   - code: HTTP status code (301, 302, ...)
//   - location: Địa chỉ chuyển hướng
func (c *forkContext) Redirect(code int, location string) {
	// Thiết lập Location header cho redirect
	c.Header("Location", location)
	// Thiết lập HTTP status code (thường là 301, 302, 307, 308)
	c.Status(code)
}

// Error trả về HTTP error với status code và thông báo từ error.
//
// Params:
//   - err: error trả về
//
// Sử dụng http.Error với status code 500 (Internal Server Error)
func (c *forkContext) Error(err error) {
	// Sử dụng http.Error để trả về lỗi với status code 500 (Internal Server Error)
	http.Error(c.response, err.Error(), http.StatusInternalServerError)
}

// ClientIP xác định và trả về địa chỉ IP của client từ các header và thông tin kết nối.
//
// Returns:
//   - string: Địa chỉ IP client
func (c *forkContext) ClientIP() string {
	// Kiểm tra X-Forwarded-For header (thường được đặt bởi proxy servers)
	if ip := c.GetHeader("X-Forwarded-For"); ip != "" {
		return ip
	}
	// Kiểm tra X-Real-IP header (thường được đặt bởi reverse proxies)
	if ip := c.GetHeader("X-Real-IP"); ip != "" {
		return ip
	}
	// Sử dụng địa chỉ remote từ kết nối trực tiếp
	return c.request.RemoteAddr()
}

// ContentType trả về giá trị của Content-Type header trong request.
//
// Returns:
//   - string: Content-Type
func (c *forkContext) ContentType() string {
	return c.GetHeader("Content-Type")
}

// IsWebsocket kiểm tra request hiện tại có phải là websocket connection request không.
//
// Returns:
//   - bool: true nếu là websocket request
func (c *forkContext) IsWebsocket() bool {
	// Kiểm tra các header cần thiết cho websocket protocol (RFC 6455)
	if c.GetHeader("Upgrade") == "websocket" && c.GetHeader("Connection") == "Upgrade" {
		return true
	}
	return false
}

// GetRawData đọc và trả về toàn bộ nội dung của request body.
//
// Returns:
//   - []byte: Dữ liệu body
//   - error: Lỗi nếu không thể đọc body
func (c *forkContext) GetRawData() ([]byte, error) {
	// Sử dụng io.ReadAll để đọc toàn bộ body của request
	return io.ReadAll(c.request.Body())
}

// Handlers trả về danh sách các handlers đã đăng ký cho context này.
//
// Returns:
//   - []func(Context): Danh sách middleware/handlers
func (c *forkContext) Handlers() []func(Context) {
	return c.handlers
}

// SetHandlers thiết lập danh sách mới các handlers cho context này.
//
// Params:
//   - handlers: Danh sách middleware/handlers
func (c *forkContext) SetHandlers(handlers []func(Context)) {
	c.handlers = handlers
}

// ValidateStruct kiểm tra tính hợp lệ của một struct dựa trên validation tags.
//
// Params:
//   - obj: Struct cần validate
//
// Returns:
//   - error: Lỗi nếu không hợp lệ
func (c *forkContext) ValidateStruct(obj interface{}) error {
	// Khởi tạo validator nếu chưa có
	if c.validator == nil {
		c.validator = validator.New()
	}
	// Thực hiện validate struct
	return c.validator.Struct(obj)
}

// ShouldBindAndValidate bind request data vào struct và validate nó, trả về lỗi nếu có.
//
// Params:
//   - obj: Struct nhận dữ liệu
//
// Returns:
//   - error: Lỗi bind hoặc validate
func (c *forkContext) ShouldBindAndValidate(obj interface{}) error {
	// Thực hiện binding trước
	if err := c.ShouldBind(obj); err != nil {
		return err
	}
	return c.ValidateStruct(obj)
}

// BindAndValidate bind request data vào struct, validate, và tự động trả về HTTP error nếu thất bại.
//
// Params:
//   - obj: Struct nhận dữ liệu
//
// Returns:
//   - error: Lỗi bind hoặc validate, đồng thời trả về JSON error response
func (c *forkContext) BindAndValidate(obj interface{}) error {
	// Thực hiện binding trước
	if err := c.Bind(obj); err != nil {
		// Trả về lỗi binding sử dụng fork/errors
		details := map[string]interface{}{
			"error": err.Error(),
		}
		// Tạo HTTP error với status code 400 Bad Request
		httpError := forkerrors.NewBadRequest("Failed to bind request data", details, err)
		// Tự động trả về response JSON với thông tin lỗi
		c.JSON(httpError.StatusCode, httpError)
		return httpError
	}

	// Thực hiện validate sau khi binding thành công
	if err := c.ValidateStruct(obj); err != nil {
		// Kiểm tra xem lỗi có phải là ValidationErrors không
		validationErrors, ok := err.(validator.ValidationErrors)
		if ok {
			// Chuyển đổi validation errors thành cấu trúc chi tiết
			fields := make(map[string]interface{})

			// Xử lý từng lỗi validation
			for _, fieldErr := range validationErrors {
				fieldName := fieldErr.Field()

				// Tạo thông tin lỗi chi tiết cho trường này
				fields[fieldName] = map[string]interface{}{
					"field":      fieldName,
					"tag":        fieldErr.Tag(),
					"value":      fmt.Sprintf("%v", fieldErr.Value()),
					"param":      fieldErr.Param(),
					"namespace":  fieldErr.Namespace(),
					"structName": fieldErr.StructNamespace(),
					"actual":     fieldErr.ActualTag(),
				}
			}

			// Sử dụng fork/errors thay vì ValidationError nội bộ
			httpError := forkerrors.NewUnprocessableEntity("Validation failed", fields, err)
			c.JSON(httpError.StatusCode, httpError)
			return httpError
		}

		// Xử lý trường hợp lỗi validation khác
		details := map[string]interface{}{
			"error": err.Error(),
		}
		httpError := forkerrors.NewUnprocessableEntity("Validation failed", details, err)
		c.JSON(httpError.StatusCode, httpError)
		return httpError
	}

	return nil
}

// RegisterValidation đăng ký một hàm validation tùy chỉnh vào validator.
//
// Params:
//   - tag: Tên tag validation
//   - fn: Hàm validation
//
// Returns:
//   - error: Lỗi nếu đăng ký thất bại
func (c *forkContext) RegisterValidation(tag string, fn validator.Func) error {
	// Khởi tạo validator nếu chưa có
	if c.validator == nil {
		c.validator = validator.New()
	}
	// Đăng ký hàm validation với tag được chỉ định
	return c.validator.RegisterValidation(tag, fn)
}

// GetValidator trả về instance của validator để cho phép cấu hình nâng cao.
//
// Returns:
//   - *validator.Validate: Instance validator
func (c *forkContext) GetValidator() *validator.Validate {
	return c.validator
}
