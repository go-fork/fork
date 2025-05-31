package errors

import (
	"fmt"
	"net/http"
)

// HttpError là kiểu dữ liệu đại diện cho lỗi HTTP.
// Nó bao gồm mã trạng thái HTTP, thông báo lỗi, và các thông tin chi tiết tùy chọn.
// HttpError tuân theo interface error của Go và cung cấp thêm thông tin về lỗi HTTP.
type HttpError struct {
	// StatusCode là mã trạng thái HTTP liên quan đến lỗi.
	// Ví dụ: 400 cho Bad Request, 404 cho Not Found, 500 cho Internal Server Error.
	StatusCode int `json:"status_code"`

	// Message là thông báo mô tả lỗi.
	// Đây là thông báo được hiển thị cho người dùng hoặc client.
	Message string `json:"message"`

	// Details chứa thông tin chi tiết tùy chỉnh về lỗi.
	// Map này có thể chứa các thông tin bổ sung để giúp debugging hoặc cung cấp
	// thêm thông tin cho client.
	Details map[string]interface{} `json:"details,omitempty"`

	// Err là lỗi gốc gây ra HttpError này, nếu có.
	// Trường này không được serialize trong JSON để tránh rò rỉ thông tin nhạy cảm.
	Err error `json:"-"`
}

// Error triển khai interface error của Go.
// Phương thức này trả về một chuỗi đại diện cho lỗi HTTP, bao gồm
// mã trạng thái, thông báo, và lỗi gốc (nếu có).
//
// Returns:
//   - string: Chuỗi đại diện cho lỗi HTTP
func (e *HttpError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("HTTP Error %d: %s - %s", e.StatusCode, e.Message, e.Err.Error())
	}
	return fmt.Sprintf("HTTP Error %d: %s", e.StatusCode, e.Message)
}

// Unwrap trả về lỗi gốc nằm bên trong HttpError.
// Phương thức này hỗ trợ cơ chế errors.Unwrap của Go để kiểm tra chuỗi lỗi.
//
// Returns:
//   - error: Lỗi gốc hoặc nil nếu không có lỗi gốc
func (e *HttpError) Unwrap() error {
	return e.Err
}

// NewHttpError tạo một HttpError với các thông số tùy chỉnh.
// Phương thức này cho phép xác định đầy đủ các thuộc tính của HttpError.
//
// Parameters:
//   - statusCode: Mã trạng thái HTTP
//   - message: Thông báo mô tả lỗi
//   - details: Map chứa thông tin chi tiết về lỗi, có thể là nil
//   - err: Lỗi gốc gây ra HttpError, có thể là nil
//
// Returns:
//   - *HttpError: Một instance mới của HttpError
func NewHttpError(statusCode int, message string, details map[string]interface{}, err error) *HttpError {
	return &HttpError{
		StatusCode: statusCode,
		Message:    message,
		Details:    details,
		Err:        err,
	}
}

// SimpleHttpError tạo một HttpError đơn giản chỉ với mã trạng thái và thông báo.
// Phương thức này là cách nhanh để tạo HttpError mà không cần chi tiết và lỗi gốc.
//
// Parameters:
//   - statusCode: Mã trạng thái HTTP
//   - message: Thông báo mô tả lỗi
//
// Returns:
//   - *HttpError: Một instance mới của HttpError với details và err là nil
func SimpleHttpError(statusCode int, message string) *HttpError {
	return NewHttpError(statusCode, message, nil, nil)
}

// NewBadRequest tạo một HttpError với mã trạng thái 400 Bad Request.
// Phương thức này được sử dụng khi client gửi request không hợp lệ.
//
// Parameters:
//   - message: Thông báo mô tả lỗi, nếu rỗng sẽ sử dụng "Bad Request"
//   - details: Map chứa thông tin chi tiết về lỗi, có thể là nil
//   - err: Lỗi gốc gây ra HttpError, có thể là nil
//
// Returns:
//   - *HttpError: Một instance mới của HttpError với StatusCode là 400
func NewBadRequest(message string, details map[string]interface{}, err error) *HttpError {
	if message == "" {
		message = "Bad Request"
	}
	return NewHttpError(http.StatusBadRequest, message, details, err)
}

// BadRequest tạo một HttpError 400 đơn giản chỉ với thông báo.
// Phương thức này là cách nhanh để tạo lỗi Bad Request khi không cần chi tiết và lỗi gốc.
//
// Parameters:
//   - message: Thông báo mô tả lỗi, nếu rỗng sẽ sử dụng "Bad Request"
//
// Returns:
//   - *HttpError: Một instance mới của HttpError với StatusCode là 400
func BadRequest(message string) *HttpError {
	if message == "" {
		message = "Bad Request"
	}
	return SimpleHttpError(http.StatusBadRequest, message)
}

// NewUnauthorized tạo một HttpError với mã trạng thái 401 Unauthorized.
// Phương thức này được sử dụng khi client chưa được xác thực.
//
// Parameters:
//   - message: Thông báo mô tả lỗi, nếu rỗng sẽ sử dụng "Unauthorized"
//   - details: Map chứa thông tin chi tiết về lỗi, có thể là nil
//   - err: Lỗi gốc gây ra HttpError, có thể là nil
//
// Returns:
//   - *HttpError: Một instance mới của HttpError với StatusCode là 401
func NewUnauthorized(message string, details map[string]interface{}, err error) *HttpError {
	if message == "" {
		message = "Unauthorized"
	}
	return NewHttpError(http.StatusUnauthorized, message, details, err)
}

// Unauthorized tạo một HttpError 401 đơn giản chỉ với thông báo.
// Phương thức này là cách nhanh để tạo lỗi Unauthorized khi không cần chi tiết và lỗi gốc.
//
// Parameters:
//   - message: Thông báo mô tả lỗi, nếu rỗng sẽ sử dụng "Unauthorized"
//
// Returns:
//   - *HttpError: Một instance mới của HttpError với StatusCode là 401
func Unauthorized(message string) *HttpError {
	if message == "" {
		message = "Unauthorized"
	}
	return SimpleHttpError(http.StatusUnauthorized, message)
}

// NewForbidden tạo một HttpError với mã trạng thái 403 Forbidden.
// Phương thức này được sử dụng khi client không có quyền truy cập resource.
//
// Parameters:
//   - message: Thông báo mô tả lỗi, nếu rỗng sẽ sử dụng "Forbidden"
//   - details: Map chứa thông tin chi tiết về lỗi, có thể là nil
//   - err: Lỗi gốc gây ra HttpError, có thể là nil
//
// Returns:
//   - *HttpError: Một instance mới của HttpError với StatusCode là 403
func NewForbidden(message string, details map[string]interface{}, err error) *HttpError {
	if message == "" {
		message = "Forbidden"
	}
	return NewHttpError(http.StatusForbidden, message, details, err)
}

// Forbidden tạo một HttpError 403 đơn giản chỉ với thông báo.
// Phương thức này là cách nhanh để tạo lỗi Forbidden khi không cần chi tiết và lỗi gốc.
//
// Parameters:
//   - message: Thông báo mô tả lỗi, nếu rỗng sẽ sử dụng "Forbidden"
//
// Returns:
//   - *HttpError: Một instance mới của HttpError với StatusCode là 403
func Forbidden(message string) *HttpError {
	if message == "" {
		message = "Forbidden"
	}
	return SimpleHttpError(http.StatusForbidden, message)
}

// NewNotFound tạo một HttpError với mã trạng thái 404 Not Found.
// Phương thức này được sử dụng khi resource không tồn tại.
//
// Parameters:
//   - message: Thông báo mô tả lỗi, nếu rỗng sẽ sử dụng "Not Found"
//   - details: Map chứa thông tin chi tiết về lỗi, có thể là nil
//   - err: Lỗi gốc gây ra HttpError, có thể là nil
//
// Returns:
//   - *HttpError: Một instance mới của HttpError với StatusCode là 404
func NewNotFound(message string, details map[string]interface{}, err error) *HttpError {
	if message == "" {
		message = "Not Found"
	}
	return NewHttpError(http.StatusNotFound, message, details, err)
}

// NotFound tạo một HttpError 404 đơn giản chỉ với thông báo.
// Phương thức này là cách nhanh để tạo lỗi Not Found khi không cần chi tiết và lỗi gốc.
//
// Parameters:
//   - message: Thông báo mô tả lỗi, nếu rỗng sẽ sử dụng "Not Found"
//
// Returns:
//   - *HttpError: Một instance mới của HttpError với StatusCode là 404
func NotFound(message string) *HttpError {
	if message == "" {
		message = "Not Found"
	}
	return SimpleHttpError(http.StatusNotFound, message)
}

// NewMethodNotAllowed tạo một HttpError với mã trạng thái 405 Method Not Allowed.
// Phương thức này được sử dụng khi HTTP method không được hỗ trợ cho resource.
//
// Parameters:
//   - message: Thông báo mô tả lỗi, nếu rỗng sẽ sử dụng "Method Not Allowed"
//   - details: Map chứa thông tin chi tiết về lỗi, có thể là nil
//   - err: Lỗi gốc gây ra HttpError, có thể là nil
//
// Returns:
//   - *HttpError: Một instance mới của HttpError với StatusCode là 405
func NewMethodNotAllowed(message string, details map[string]interface{}, err error) *HttpError {
	if message == "" {
		message = "Method Not Allowed"
	}
	return NewHttpError(http.StatusMethodNotAllowed, message, details, err)
}

// MethodNotAllowed tạo một HttpError 405 đơn giản chỉ với thông báo.
// Phương thức này là cách nhanh để tạo lỗi Method Not Allowed khi không cần chi tiết và lỗi gốc.
//
// Parameters:
//   - message: Thông báo mô tả lỗi, nếu rỗng sẽ sử dụng "Method Not Allowed"
//
// Returns:
//   - *HttpError: Một instance mới của HttpError với StatusCode là 405
func MethodNotAllowed(message string) *HttpError {
	if message == "" {
		message = "Method Not Allowed"
	}
	return SimpleHttpError(http.StatusMethodNotAllowed, message)
}

// NewNotAcceptable tạo một HttpError với mã trạng thái 406 Not Acceptable.
// Phương thức này được sử dụng khi server không thể tạo response phù hợp với danh sách Accept headers.
//
// Parameters:
//   - message: Thông báo mô tả lỗi, nếu rỗng sẽ sử dụng "Not Acceptable"
//   - details: Map chứa thông tin chi tiết về lỗi, có thể là nil
//   - err: Lỗi gốc gây ra HttpError, có thể là nil
//
// Returns:
//   - *HttpError: Một instance mới của HttpError với StatusCode là 406
func NewNotAcceptable(message string, details map[string]interface{}, err error) *HttpError {
	if message == "" {
		message = "Not Acceptable"
	}
	return NewHttpError(http.StatusNotAcceptable, message, details, err)
}

// NotAcceptable tạo một HttpError 406 đơn giản chỉ với thông báo.
// Phương thức này là cách nhanh để tạo lỗi Not Acceptable khi không cần chi tiết và lỗi gốc.
//
// Parameters:
//   - message: Thông báo mô tả lỗi, nếu rỗng sẽ sử dụng "Not Acceptable"
//
// Returns:
//   - *HttpError: Một instance mới của HttpError với StatusCode là 406
func NotAcceptable(message string) *HttpError {
	if message == "" {
		message = "Not Acceptable"
	}
	return SimpleHttpError(http.StatusNotAcceptable, message)
}

// NewConflict tạo một HttpError với mã trạng thái 409 Conflict.
// Phương thức này được sử dụng khi request xung đột với trạng thái hiện tại của server.
//
// Parameters:
//   - message: Thông báo mô tả lỗi, nếu rỗng sẽ sử dụng "Conflict"
//   - details: Map chứa thông tin chi tiết về lỗi, có thể là nil
//   - err: Lỗi gốc gây ra HttpError, có thể là nil
//
// Returns:
//   - *HttpError: Một instance mới của HttpError với StatusCode là 409
func NewConflict(message string, details map[string]interface{}, err error) *HttpError {
	if message == "" {
		message = "Conflict"
	}
	return NewHttpError(http.StatusConflict, message, details, err)
}

// Conflict tạo một HttpError 409 đơn giản chỉ với thông báo.
// Phương thức này là cách nhanh để tạo lỗi Conflict khi không cần chi tiết và lỗi gốc.
//
// Parameters:
//   - message: Thông báo mô tả lỗi, nếu rỗng sẽ sử dụng "Conflict"
//
// Returns:
//   - *HttpError: Một instance mới của HttpError với StatusCode là 409
func Conflict(message string) *HttpError {
	if message == "" {
		message = "Conflict"
	}
	return SimpleHttpError(http.StatusConflict, message)
}

// NewGone tạo một HttpError với mã trạng thái 410 Gone.
// Phương thức này được sử dụng khi resource không còn tồn tại vĩnh viễn.
//
// Parameters:
//   - message: Thông báo mô tả lỗi, nếu rỗng sẽ sử dụng "Gone"
//   - details: Map chứa thông tin chi tiết về lỗi, có thể là nil
//   - err: Lỗi gốc gây ra HttpError, có thể là nil
//
// Returns:
//   - *HttpError: Một instance mới của HttpError với StatusCode là 410
func NewGone(message string, details map[string]interface{}, err error) *HttpError {
	if message == "" {
		message = "Gone"
	}
	return NewHttpError(http.StatusGone, message, details, err)
}

// Gone tạo một HttpError 410 đơn giản chỉ với thông báo.
// Phương thức này là cách nhanh để tạo lỗi Gone khi không cần chi tiết và lỗi gốc.
//
// Parameters:
//   - message: Thông báo mô tả lỗi, nếu rỗng sẽ sử dụng "Gone"
//
// Returns:
//   - *HttpError: Một instance mới của HttpError với StatusCode là 410
func Gone(message string) *HttpError {
	if message == "" {
		message = "Gone"
	}
	return SimpleHttpError(http.StatusGone, message)
}

// NewUnsupportedMediaType tạo một HttpError với mã trạng thái 415 Unsupported Media Type.
// Phương thức này được sử dụng khi server không hỗ trợ định dạng media yêu cầu.
//
// Parameters:
//   - message: Thông báo mô tả lỗi, nếu rỗng sẽ sử dụng "Unsupported Media Type"
//   - details: Map chứa thông tin chi tiết về lỗi, có thể là nil
//   - err: Lỗi gốc gây ra HttpError, có thể là nil
//
// Returns:
//   - *HttpError: Một instance mới của HttpError với StatusCode là 415
func NewUnsupportedMediaType(message string, details map[string]interface{}, err error) *HttpError {
	if message == "" {
		message = "Unsupported Media Type"
	}
	return NewHttpError(http.StatusUnsupportedMediaType, message, details, err)
}

// UnsupportedMediaType tạo một HttpError 415 đơn giản chỉ với thông báo.
// Phương thức này là cách nhanh để tạo lỗi Unsupported Media Type khi không cần chi tiết và lỗi gốc.
//
// Parameters:
//   - message: Thông báo mô tả lỗi, nếu rỗng sẽ sử dụng "Unsupported Media Type"
//
// Returns:
//   - *HttpError: Một instance mới của HttpError với StatusCode là 415
func UnsupportedMediaType(message string) *HttpError {
	if message == "" {
		message = "Unsupported Media Type"
	}
	return SimpleHttpError(http.StatusUnsupportedMediaType, message)
}

// NewUnprocessableEntity tạo một HttpError với mã trạng thái 422 Unprocessable Entity.
// Phương thức này được sử dụng khi request có định dạng đúng nhưng không thể xử lý do lỗi ngữ nghĩa.
//
// Parameters:
//   - message: Thông báo mô tả lỗi, nếu rỗng sẽ sử dụng "Unprocessable Entity"
//   - details: Map chứa thông tin chi tiết về lỗi, có thể là nil
//   - err: Lỗi gốc gây ra HttpError, có thể là nil
//
// Returns:
//   - *HttpError: Một instance mới của HttpError với StatusCode là 422
func NewUnprocessableEntity(message string, details map[string]interface{}, err error) *HttpError {
	if message == "" {
		message = "Unprocessable Entity"
	}
	return NewHttpError(http.StatusUnprocessableEntity, message, details, err)
}

// UnprocessableEntity tạo một HttpError 422 đơn giản chỉ với thông báo.
// Phương thức này là cách nhanh để tạo lỗi Unprocessable Entity khi không cần chi tiết và lỗi gốc.
//
// Parameters:
//   - message: Thông báo mô tả lỗi, nếu rỗng sẽ sử dụng "Unprocessable Entity"
//
// Returns:
//   - *HttpError: Một instance mới của HttpError với StatusCode là 422
func UnprocessableEntity(message string) *HttpError {
	if message == "" {
		message = "Unprocessable Entity"
	}
	return SimpleHttpError(http.StatusUnprocessableEntity, message)
}

// NewTooManyRequests tạo một HttpError với mã trạng thái 429 Too Many Requests.
// Phương thức này được sử dụng khi client gửi quá nhiều requests trong một khoảng thời gian (rate limiting).
//
// Parameters:
//   - message: Thông báo mô tả lỗi, nếu rỗng sẽ sử dụng "Too Many Requests"
//   - details: Map chứa thông tin chi tiết về lỗi, có thể là nil
//   - err: Lỗi gốc gây ra HttpError, có thể là nil
//
// Returns:
//   - *HttpError: Một instance mới của HttpError với StatusCode là 429
func NewTooManyRequests(message string, details map[string]interface{}, err error) *HttpError {
	if message == "" {
		message = "Too Many Requests"
	}
	return NewHttpError(http.StatusTooManyRequests, message, details, err)
}

// TooManyRequests tạo một HttpError 429 đơn giản chỉ với thông báo.
// Phương thức này là cách nhanh để tạo lỗi Too Many Requests khi không cần chi tiết và lỗi gốc.
//
// Parameters:
//   - message: Thông báo mô tả lỗi, nếu rỗng sẽ sử dụng "Too Many Requests"
//
// Returns:
//   - *HttpError: Một instance mới của HttpError với StatusCode là 429
func TooManyRequests(message string) *HttpError {
	if message == "" {
		message = "Too Many Requests"
	}
	return SimpleHttpError(http.StatusTooManyRequests, message)
}

// NewInternalServerError tạo một HttpError với mã trạng thái 500 Internal Server Error.
// Phương thức này được sử dụng khi có lỗi không mong muốn xảy ra trên server.
//
// Parameters:
//   - message: Thông báo mô tả lỗi, nếu rỗng sẽ sử dụng "Internal Server Error"
//   - details: Map chứa thông tin chi tiết về lỗi, có thể là nil
//   - err: Lỗi gốc gây ra HttpError, có thể là nil
//
// Returns:
//   - *HttpError: Một instance mới của HttpError với StatusCode là 500
func NewInternalServerError(message string, details map[string]interface{}, err error) *HttpError {
	if message == "" {
		message = "Internal Server Error"
	}
	return NewHttpError(http.StatusInternalServerError, message, details, err)
}

// InternalServerError tạo một HttpError 500 đơn giản chỉ với thông báo.
// Phương thức này là cách nhanh để tạo lỗi Internal Server Error khi không cần chi tiết và lỗi gốc.
//
// Parameters:
//   - message: Thông báo mô tả lỗi, nếu rỗng sẽ sử dụng "Internal Server Error"
//
// Returns:
//   - *HttpError: Một instance mới của HttpError với StatusCode là 500
func InternalServerError(message string) *HttpError {
	if message == "" {
		message = "Internal Server Error"
	}
	return SimpleHttpError(http.StatusInternalServerError, message)
}

// NewNotImplemented tạo một HttpError với mã trạng thái 501 Not Implemented.
// Phương thức này được sử dụng khi server không hỗ trợ chức năng cần thiết để xử lý request.
//
// Parameters:
//   - message: Thông báo mô tả lỗi, nếu rỗng sẽ sử dụng "Not Implemented"
//   - details: Map chứa thông tin chi tiết về lỗi, có thể là nil
//   - err: Lỗi gốc gây ra HttpError, có thể là nil
//
// Returns:
//   - *HttpError: Một instance mới của HttpError với StatusCode là 501
func NewNotImplemented(message string, details map[string]interface{}, err error) *HttpError {
	if message == "" {
		message = "Not Implemented"
	}
	return NewHttpError(http.StatusNotImplemented, message, details, err)
}

// NotImplemented tạo một HttpError 501 đơn giản chỉ với thông báo.
// Phương thức này là cách nhanh để tạo lỗi Not Implemented khi không cần chi tiết và lỗi gốc.
//
// Parameters:
//   - message: Thông báo mô tả lỗi, nếu rỗng sẽ sử dụng "Not Implemented"
//
// Returns:
//   - *HttpError: Một instance mới của HttpError với StatusCode là 501
func NotImplemented(message string) *HttpError {
	if message == "" {
		message = "Not Implemented"
	}
	return SimpleHttpError(http.StatusNotImplemented, message)
}

// NewBadGateway tạo một HttpError với mã trạng thái 502 Bad Gateway
func NewBadGateway(message string, details map[string]interface{}, err error) *HttpError {
	if message == "" {
		message = "Bad Gateway"
	}
	return NewHttpError(http.StatusBadGateway, message, details, err)
}

// BadGateway tạo một HttpError 502 đơn giản chỉ với thông báo
func BadGateway(message string) *HttpError {
	if message == "" {
		message = "Bad Gateway"
	}
	return SimpleHttpError(http.StatusBadGateway, message)
}

// NewServiceUnavailable tạo một HttpError với mã trạng thái 503 Service Unavailable
func NewServiceUnavailable(message string, details map[string]interface{}, err error) *HttpError {
	if message == "" {
		message = "Service Unavailable"
	}
	return NewHttpError(http.StatusServiceUnavailable, message, details, err)
}

// ServiceUnavailable tạo một HttpError 503 đơn giản chỉ với thông báo
func ServiceUnavailable(message string) *HttpError {
	if message == "" {
		message = "Service Unavailable"
	}
	return SimpleHttpError(http.StatusServiceUnavailable, message)
}

// NewGatewayTimeout tạo một HttpError với mã trạng thái 504 Gateway Timeout
func NewGatewayTimeout(message string, details map[string]interface{}, err error) *HttpError {
	if message == "" {
		message = "Gateway Timeout"
	}
	return NewHttpError(http.StatusGatewayTimeout, message, details, err)
}

// GatewayTimeout tạo một HttpError 504 đơn giản chỉ với thông báo
func GatewayTimeout(message string) *HttpError {
	if message == "" {
		message = "Gateway Timeout"
	}
	return SimpleHttpError(http.StatusGatewayTimeout, message)
}
