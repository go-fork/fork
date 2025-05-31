package fork

import (
	"errors"
	"net/http"
)

// HTTP Methods định nghĩa các HTTP methods chuẩn theo RFC7231 và RFC5789.
// Các constants này được sử dụng trong toàn bộ framework để định danh các HTTP methods.
const (
	// MethodGet là HTTP GET method, được sử dụng để truy xuất dữ liệu.
	MethodGet = "GET"

	// MethodPost là HTTP POST method, được sử dụng để tạo dữ liệu mới.
	MethodPost = "POST"

	// MethodPut là HTTP PUT method, được sử dụng để cập nhật toàn bộ dữ liệu hiện có.
	MethodPut = "PUT"

	// MethodDelete là HTTP DELETE method, được sử dụng để xóa dữ liệu.
	MethodDelete = "DELETE"

	// MethodPatch là HTTP PATCH method, được sử dụng để cập nhật một phần dữ liệu.
	MethodPatch = "PATCH"

	// MethodHead là HTTP HEAD method, tương tự GET nhưng chỉ trả về headers.
	MethodHead = "HEAD"

	// MethodOptions là HTTP OPTIONS method, được sử dụng để truy vấn về các options giao tiếp.
	MethodOptions = "OPTIONS"

	// MethodTrace là HTTP TRACE method, được sử dụng để debugging.
	MethodTrace = "TRACE"

	// MethodConnect là HTTP CONNECT method, được sử dụng để thiết lập tunnel đến server.
	MethodConnect = "CONNECT"
)

// Content types định nghĩa các MIME types phổ biến cho HTTP content.
// Các constants này được sử dụng để thiết lập và kiểm tra Content-Type header.
const (
	// ContentTypeJSON là MIME type cho dữ liệu JSON.
	ContentTypeJSON = "application/json"

	// ContentTypeXML là MIME type cho dữ liệu XML.
	ContentTypeXML = "application/xml"

	// ContentTypeForm là MIME type cho form data được mã hóa URL.
	ContentTypeForm = "application/x-www-form-urlencoded"

	// ContentTypeFormMultipart là MIME type cho multipart form data.
	ContentTypeFormMultipart = "multipart/form-data"

	// ContentTypePlain là MIME type cho plain text.
	ContentTypePlain = "text/plain"

	// ContentTypeHTML là MIME type cho HTML content.
	ContentTypeHTML = "text/html"

	// ContentTypeEventStream là MIME type cho server-sent events.
	ContentTypeEventStream = "text/event-stream"

	// ContentTypeProtobuf là MIME type cho Protocol Buffers.
	ContentTypeProtobuf = "application/protobuf"

	// ContentTypeOctetStream là MIME type cho binary data.
	ContentTypeOctetStream = "application/octet-stream"

	// ContentTypePDF là MIME type cho PDF documents.
	ContentTypePDF = "application/pdf"
)

// Header names định nghĩa các HTTP header phổ biến theo RFC7231.
// Các constants này được sử dụng để truy cập và thiết lập HTTP headers.
const (
	// HeaderAccept chỉ định media types được client chấp nhận.
	HeaderAccept = "Accept"

	// HeaderAcceptEncoding chỉ định encoding algorithms được client chấp nhận.
	HeaderAcceptEncoding = "Accept-Encoding"

	// HeaderAcceptLanguage chỉ định ngôn ngữ được client ưa thích.
	HeaderAcceptLanguage = "Accept-Language"

	// HeaderAllow liệt kê các HTTP methods được phép cho resource.
	HeaderAllow = "Allow"

	// HeaderAuthorization chứa thông tin xác thực của client.
	HeaderAuthorization = "Authorization"

	// HeaderContentDisposition chỉ định cách xử lý content (inline hoặc attachment).
	HeaderContentDisposition = "Content-Disposition"

	// HeaderContentEncoding chỉ định encoding đã áp dụng cho body.
	HeaderContentEncoding = "Content-Encoding"

	// HeaderContentLength chỉ định kích thước của body tính bằng bytes.
	HeaderContentLength = "Content-Length"

	// HeaderContentType chỉ định media type của body.
	HeaderContentType = "Content-Type"

	// HeaderCookie chứa HTTP cookies gửi từ client.
	HeaderCookie = "Cookie"

	// HeaderSetCookie được server sử dụng để gửi cookies đến client.
	HeaderSetCookie = "Set-Cookie"

	// HeaderIfModifiedSince được sử dụng với GET để request chỉ khi resource đã thay đổi.
	HeaderIfModifiedSince = "If-Modified-Since"

	// HeaderLastModified chỉ định thời điểm resource được sửa đổi lần cuối.
	HeaderLastModified = "Last-Modified"

	// HeaderLocation chỉ định URL mà client nên chuyển hướng đến.
	HeaderLocation = "Location"

	// HeaderUpgrade yêu cầu server chuyển đổi sang protocol khác.
	HeaderUpgrade = "Upgrade"

	// HeaderVary chỉ định các headers ảnh hưởng đến cách server chọn representation.
	HeaderVary = "Vary"

	// HeaderWWWAuthenticate chỉ định phương thức xác thực cho resource.
	HeaderWWWAuthenticate = "WWW-Authenticate"

	// HeaderXForwardedFor chứa địa chỉ IP của client khi qua proxy.
	HeaderXForwardedFor = "X-Forwarded-For"

	// HeaderXForwardedProto chứa protocol ban đầu của request.
	HeaderXForwardedProto = "X-Forwarded-Proto"

	// HeaderXForwardedProtocol tương tự HeaderXForwardedProto.
	HeaderXForwardedProtocol = "X-Forwarded-Protocol"

	// HeaderXForwardedHost chứa host ban đầu của request.
	HeaderXForwardedHost = "X-Forwarded-Host"

	// HeaderXRealIP chứa địa chỉ IP thật của client khi qua proxy.
	HeaderXRealIP = "X-Real-IP"

	// HeaderXRequestID chứa ID duy nhất của request để tracking.
	HeaderXRequestID = "X-Request-ID"

	// HeaderXRequestedWith chứa thông tin về loại request (AJAX, v.v.).
	HeaderXRequestedWith = "X-Requested-With"

	// HeaderServer chứa thông tin về server phục vụ request.
	HeaderServer = "Server"

	// HeaderOrigin chỉ định origin của request (CORS).
	HeaderOrigin = "Origin"

	// HeaderAccessControlRequestMethod chỉ định method trong preflight request (CORS).
	HeaderAccessControlRequestMethod = "Access-Control-Request-Method"

	// HeaderAccessControlRequestHeaders chỉ định headers trong preflight request (CORS).
	HeaderAccessControlRequestHeaders = "Access-Control-Request-Headers"

	// HeaderAccessControlAllowOrigin chỉ định origins được phép truy cập (CORS).
	HeaderAccessControlAllowOrigin = "Access-Control-Allow-Origin"

	// HeaderAccessControlAllowMethods chỉ định methods được phép (CORS).
	HeaderAccessControlAllowMethods = "Access-Control-Allow-Methods"

	// HeaderAccessControlAllowHeaders chỉ định headers được phép (CORS).
	HeaderAccessControlAllowHeaders = "Access-Control-Allow-Headers"

	// HeaderAccessControlAllowCredentials cho phép gửi credentials (CORS).
	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"

	// HeaderAccessControlExposeHeaders chỉ định headers mà browser có thể truy cập (CORS).
	HeaderAccessControlExposeHeaders = "Access-Control-Expose-Headers"

	// HeaderAccessControlMaxAge chỉ định thời gian cache preflight request (CORS).
	HeaderAccessControlMaxAge = "Access-Control-Max-Age"

	// HeaderStrictTransportSecurity yêu cầu browser chỉ kết nối qua HTTPS (HSTS).
	HeaderStrictTransportSecurity = "Strict-Transport-Security"

	// HeaderXContentTypeOptions ngăn chặn MIME-type sniffing (bảo mật).
	HeaderXContentTypeOptions = "X-Content-Type-Options"

	// HeaderXXSSProtection kích hoạt bảo vệ XSS của browser.
	HeaderXXSSProtection = "X-XSS-Protection"

	// HeaderXFrameOptions kiểm soát việc nhúng trang trong frames (clickjacking).
	HeaderXFrameOptions = "X-Frame-Options"

	// HeaderContentSecurityPolicy định nghĩa chính sách bảo mật content (CSP).
	HeaderContentSecurityPolicy = "Content-Security-Policy"

	// HeaderXCSRFToken chứa token CSRF để bảo vệ chống lại tấn công CSRF.
	HeaderXCSRFToken = "X-CSRF-Token"
)

// MIME types định nghĩa các media type đầy đủ với charset.
// Các constants này được sử dụng để thiết lập Content-Type header.
const (
	// MIMEWebAppJSON là MIME type cho JSON.
	MIMEWebAppJSON = "application/json"

	// MIMEWebAppJSONCharsetUTF8 là MIME type cho JSON với charset UTF-8.
	MIMEWebAppJSONCharsetUTF8 = "application/json; charset=utf-8"

	// MIMEWebAppXML là MIME type cho XML.
	MIMEWebAppXML = "application/xml"

	// MIMEWebAppXMLCharsetUTF8 là MIME type cho XML với charset UTF-8.
	MIMEWebAppXMLCharsetUTF8 = "application/xml; charset=utf-8"

	// MIMETextXML là MIME type cho XML dạng text.
	MIMETextXML = "text/xml"

	// MIMETextXMLCharsetUTF8 là MIME type cho XML dạng text với charset UTF-8.
	MIMETextXMLCharsetUTF8 = "text/xml; charset=utf-8"

	// MIMEWebAppForm là MIME type cho form data được mã hóa URL.
	MIMEWebAppForm = "application/x-www-form-urlencoded"

	// MIMEWebAppProtobuf là MIME type cho Protocol Buffers.
	MIMEWebAppProtobuf = "application/protobuf"

	// MIMEWebAppMsgpack là MIME type cho MessagePack.
	MIMEWebAppMsgpack = "application/msgpack"

	// MIMETextHTML là MIME type cho HTML.
	MIMETextHTML = "text/html"

	// MIMETextHTMLCharsetUTF8 là MIME type cho HTML với charset UTF-8.
	MIMETextHTMLCharsetUTF8 = "text/html; charset=utf-8"

	// MIMETextPlain là MIME type cho plain text.
	MIMETextPlain = "text/plain"

	// MIMETextPlainCharsetUTF8 là MIME type cho plain text với charset UTF-8.
	MIMETextPlainCharsetUTF8 = "text/plain; charset=utf-8"

	// MIMEMultipartForm là MIME type cho multipart form data.
	MIMEMultipartForm = "multipart/form-data"

	// MIMEOctetStream là MIME type cho binary data.
	MIMEOctetStream = "application/octet-stream"
)

// HTTP status định nghĩa các HTTP status codes theo RFC7231, RFC6585 và RFC4918.
// Các constants này giúp xác định trạng thái của HTTP response.
const (
	// StatusContinue (100) - Server đã nhận được request headers và client nên tiếp tục gửi request body.
	StatusContinue = http.StatusContinue // 100

	// StatusSwitchingProtocols (101) - Server đang chuyển đổi sang protocol khác theo yêu cầu của client.
	StatusSwitchingProtocols = http.StatusSwitchingProtocols // 101

	// StatusProcessing (102) - Server đang xử lý request nhưng chưa có response sẵn có.
	StatusProcessing = http.StatusProcessing // 102

	// StatusOK (200) - Request đã thành công.
	StatusOK = http.StatusOK // 200

	// StatusCreated (201) - Request đã thành công và một resource mới đã được tạo.
	StatusCreated = http.StatusCreated // 201

	// StatusAccepted (202) - Request đã được chấp nhận để xử lý nhưng chưa hoàn thành.
	StatusAccepted = http.StatusAccepted // 202

	// StatusNonAuthoritativeInfo (203) - Response chứa thông tin không chính thức.
	StatusNonAuthoritativeInfo = http.StatusNonAuthoritativeInfo // 203

	// StatusNoContent (204) - Server đã xử lý thành công request nhưng không trả về body.
	StatusNoContent = http.StatusNoContent // 204

	// StatusResetContent (205) - Server đã xử lý request và client nên reset document view.
	StatusResetContent = http.StatusResetContent // 205

	// StatusPartialContent (206) - Server đã xử lý một phần của request (range request).
	StatusPartialContent = http.StatusPartialContent // 206

	// StatusMultiStatus (207) - Response chứa nhiều status codes cho nhiều sub-requests.
	StatusMultiStatus = http.StatusMultiStatus // 207

	// StatusAlreadyReported (208) - Members đã được báo cáo trước đó.
	StatusAlreadyReported = http.StatusAlreadyReported // 208

	// StatusIMUsed (226) - Server đã thực hiện request cho resource và response là biểu diễn của tập hợp các manipulations.
	StatusIMUsed = http.StatusIMUsed // 226

	// StatusMultipleChoices (300) - Request có nhiều hơn một khả năng phản hồi.
	StatusMultipleChoices = http.StatusMultipleChoices // 300

	// StatusMovedPermanently (301) - Resource đã được di chuyển vĩnh viễn sang URL mới.
	StatusMovedPermanently = http.StatusMovedPermanently // 301

	// StatusFound (302) - Resource tạm thời được tìm thấy ở một URL khác.
	StatusFound = http.StatusFound // 302

	// StatusSeeOther (303) - Client nên lấy resource từ URL khác bằng GET method.
	StatusSeeOther = http.StatusSeeOther // 303

	// StatusNotModified (304) - Resource không bị sửa đổi kể từ phiên bản đã cached.
	StatusNotModified = http.StatusNotModified // 304

	// StatusUseProxy (305) - Resource phải được truy cập thông qua proxy.
	StatusUseProxy = http.StatusUseProxy // 305

	// StatusTemporaryRedirect (307) - Resource tạm thời được chuyển hướng sang một URL khác.
	StatusTemporaryRedirect = http.StatusTemporaryRedirect // 307

	// StatusPermanentRedirect (308) - Resource đã được di chuyển vĩnh viễn sang URL mới.
	StatusPermanentRedirect = http.StatusPermanentRedirect // 308

	// StatusBadRequest (400) - Server không thể xử lý request do cú pháp không hợp lệ.
	StatusBadRequest = http.StatusBadRequest // 400

	// StatusUnauthorized (401) - Request yêu cầu xác thực người dùng.
	StatusUnauthorized = http.StatusUnauthorized // 401

	// StatusPaymentRequired (402) - Yêu cầu thanh toán để tiếp tục xử lý.
	StatusPaymentRequired = http.StatusPaymentRequired // 402

	// StatusForbidden (403) - Server hiểu request nhưng từ chối cấp quyền.
	StatusForbidden = http.StatusForbidden // 403

	// StatusNotFound (404) - Server không tìm thấy resource được yêu cầu.
	StatusNotFound = http.StatusNotFound // 404

	// StatusMethodNotAllowed (405) - Method không được cho phép cho resource.
	StatusMethodNotAllowed = http.StatusMethodNotAllowed // 405

	// StatusNotAcceptable (406) - Không thể tạo response thỏa mãn header Accept.
	StatusNotAcceptable = http.StatusNotAcceptable // 406

	// StatusProxyAuthRequired (407) - Yêu cầu xác thực proxy.
	StatusProxyAuthRequired = http.StatusProxyAuthRequired // 407

	// StatusRequestTimeout (408) - Server đã hết thời gian chờ đợi request.
	StatusRequestTimeout = http.StatusRequestTimeout // 408

	// StatusConflict (409) - Request xung đột với trạng thái hiện tại của resource.
	StatusConflict = http.StatusConflict // 409

	// StatusGone (410) - Resource đã không còn tồn tại vĩnh viễn.
	StatusGone = http.StatusGone // 410

	// StatusLengthRequired (411) - Content-Length header là bắt buộc.
	StatusLengthRequired = http.StatusLengthRequired // 411

	// StatusPreconditionFailed (412) - Điều kiện tiên quyết không được đáp ứng.
	StatusPreconditionFailed = http.StatusPreconditionFailed // 412

	// StatusRequestEntityTooLarge (413) - Request entity quá lớn.
	StatusRequestEntityTooLarge = http.StatusRequestEntityTooLarge // 413

	// StatusRequestURITooLong (414) - Request URI quá dài.
	StatusRequestURITooLong = http.StatusRequestURITooLong // 414

	// StatusUnsupportedMediaType (415) - Media type không được hỗ trợ.
	StatusUnsupportedMediaType = http.StatusUnsupportedMediaType // 415

	// StatusRequestedRangeNotSatisfiable (416) - Range request không thể được đáp ứng.
	StatusRequestedRangeNotSatisfiable = http.StatusRequestedRangeNotSatisfiable // 416

	// StatusExpectationFailed (417) - Expect header không thể được đáp ứng.
	StatusExpectationFailed = http.StatusExpectationFailed // 417

	// StatusTeapot (418) - "I'm a teapot" - RFC 2324 (April Fools' joke).
	StatusTeapot = http.StatusTeapot // 418

	// StatusMisdirectedRequest (421) - Request đã được chuyển hướng đến server không thể tạo response.
	StatusMisdirectedRequest = http.StatusMisdirectedRequest // 421

	// StatusUnprocessableEntity (422) - Request hợp lệ nhưng không thể xử lý.
	StatusUnprocessableEntity = http.StatusUnprocessableEntity // 422

	// StatusLocked (423) - Resource bị khóa.
	StatusLocked = http.StatusLocked // 423

	// StatusFailedDependency (424) - Request thất bại do phụ thuộc vào request khác đã thất bại.
	StatusFailedDependency = http.StatusFailedDependency // 424

	// StatusTooEarly (425) - Server không sẵn sàng xử lý request vì có thể bị tấn công replay.
	StatusTooEarly = http.StatusTooEarly // 425

	// StatusUpgradeRequired (426) - Client nên chuyển sang protocol khác.
	StatusUpgradeRequired = http.StatusUpgradeRequired // 426

	// StatusPreconditionRequired (428) - Origin server yêu cầu request phải có điều kiện.
	StatusPreconditionRequired = http.StatusPreconditionRequired // 428

	// StatusTooManyRequests (429) - Client đã gửi quá nhiều requests trong một khoảng thời gian.
	StatusTooManyRequests = http.StatusTooManyRequests // 429

	// StatusRequestHeaderFieldsTooLarge (431) - Header fields quá lớn.
	StatusRequestHeaderFieldsTooLarge = http.StatusRequestHeaderFieldsTooLarge // 431

	// StatusUnavailableForLegalReasons (451) - Resource không khả dụng vì lý do pháp lý.
	StatusUnavailableForLegalReasons = http.StatusUnavailableForLegalReasons // 451

	// StatusInternalServerError (500) - Lỗi nội bộ server không xác định.
	StatusInternalServerError = http.StatusInternalServerError // 500

	// StatusNotImplemented (501) - Method không được server hỗ trợ.
	StatusNotImplemented = http.StatusNotImplemented // 501

	// StatusBadGateway (502) - Server, trong vai trò gateway, nhận được response không hợp lệ.
	StatusBadGateway = http.StatusBadGateway // 502

	// StatusServiceUnavailable (503) - Server tạm thời không khả dụng.
	StatusServiceUnavailable = http.StatusServiceUnavailable // 503

	// StatusGatewayTimeout (504) - Server, trong vai trò gateway, không nhận được response kịp thời.
	StatusGatewayTimeout = http.StatusGatewayTimeout // 504

	// StatusHTTPVersionNotSupported (505) - Phiên bản HTTP không được hỗ trợ.
	StatusHTTPVersionNotSupported = http.StatusHTTPVersionNotSupported // 505

	// StatusVariantAlsoNegotiates (506) - Cấu hình nội bộ server không chính xác.
	StatusVariantAlsoNegotiates = http.StatusVariantAlsoNegotiates // 506

	// StatusInsufficientStorage (507) - Server không đủ storage để hoàn thành request.
	StatusInsufficientStorage = http.StatusInsufficientStorage // 507

	// StatusLoopDetected (508) - Server phát hiện vòng lặp vô hạn trong request.
	StatusLoopDetected = http.StatusLoopDetected // 508

	// StatusNotExtended (510) - Yêu cầu thêm extensions để hoàn thành request.
	StatusNotExtended = http.StatusNotExtended // 510

	// StatusNetworkAuthenticationRequired (511) - Client cần xác thực để truy cập network.
	StatusNetworkAuthenticationRequired = http.StatusNetworkAuthenticationRequired // 511
)

// Errors định nghĩa các lỗi tiêu chuẩn của framework.
// Các error constants này được sử dụng để thông báo lỗi trong toàn bộ framework.
var (
	// ErrServerNotRunning được trả về khi cố gắng thao tác với server chưa được khởi động.
	// Điều này xảy ra khi cố gắng thực hiện các hoạt động yêu cầu server đang chạy.
	ErrServerNotRunning = errors.New("http: server not running")

	// ErrInvalidMethod được trả về khi HTTP method không hợp lệ hoặc không được hỗ trợ.
	// Điều này xảy ra khi client sử dụng một HTTP method không chuẩn hoặc không được server hỗ trợ.
	ErrInvalidMethod = errors.New("http: invalid method")

	// ErrUnsupportedMediaType được trả về khi nội dung request có Content-Type không được hỗ trợ.
	// Điều này xảy ra khi server không thể xử lý định dạng dữ liệu mà client gửi.
	ErrUnsupportedMediaType = errors.New("http: unsupported media type")

	// ErrNotFound được trả về khi resource không tồn tại.
	// Điều này xảy ra khi client yêu cầu một resource không có trên server.
	ErrNotFound = errors.New("http: not found")

	// ErrMethodNotAllowed được trả về khi method không được phép cho resource.
	// Điều này xảy ra khi resource tồn tại nhưng không hỗ trợ HTTP method đã yêu cầu.
	ErrMethodNotAllowed = errors.New("http: method not allowed")

	// ErrBadRequest được trả về khi request không hợp lệ.
	// Điều này xảy ra khi định dạng request không đúng hoặc thiếu thông tin bắt buộc.
	ErrBadRequest = errors.New("http: bad request")

	// ErrInternalServerError được trả về khi có lỗi không xác định từ server.
	// Điều này xảy ra khi server gặp lỗi mà không thể xử lý một cách bình thường.
	ErrInternalServerError = errors.New("http: internal server error")

	// ErrUnauthorized được trả về khi client chưa được xác thực.
	// Điều này xảy ra khi client cố gắng truy cập resource yêu cầu xác thực nhưng chưa đăng nhập.
	ErrUnauthorized = errors.New("http: unauthorized")

	// ErrForbidden được trả về khi client bị từ chối quyền truy cập.
	// Điều này xảy ra khi client đã xác thực nhưng không có đủ quyền truy cập resource.
	ErrForbidden = errors.New("http: forbidden")

	// ErrRequestTimeout được trả về khi xử lý request vượt quá thời gian cho phép.
	// Điều này xảy ra khi server không thể hoàn thành xử lý request trong thời gian giới hạn.
	ErrRequestTimeout = errors.New("http: request timeout")

	// ErrRouteNotFound được trả về khi không tìm thấy route cho request.
	// Điều này xảy ra khi URL path được yêu cầu không khớp với bất kỳ route nào đã đăng ký.
	ErrRouteNotFound = errors.New("http: route not found")

	// ErrUnsupportedBinding được trả về khi kiểu binding không được hỗ trợ.
	// Điều này xảy ra khi cố gắng bind request body vào một kiểu dữ liệu không hỗ trợ.
	ErrUnsupportedBinding = errors.New("http: unsupported binding")

	// ErrUnsupportedRendering được trả về khi kiểu rendering không được hỗ trợ.
	// Điều này xảy ra khi cố gắng render response trong một định dạng không hỗ trợ.
	ErrUnsupportedRendering = errors.New("http: unsupported rendering")

	// ErrValidationFailed được trả về khi validation của request body thất bại.
	// Điều này xảy ra khi dữ liệu đầu vào không đáp ứng các quy tắc validation.
	ErrValidationFailed = errors.New("http: validation failed")

	// ErrInvalidBinding được trả về khi không thể bind request body vào struct.
	// Điều này xảy ra khi định dạng của request body không khớp với cấu trúc đích.
	ErrInvalidBinding = errors.New("http: invalid binding")

	// ErrInvalidCertificate được trả về khi tệp chứng chỉ SSL/TLS không hợp lệ hoặc không tồn tại.
	// Điều này xảy ra khi cố gắng khởi động HTTPS server với chứng chỉ không hợp lệ.
	ErrInvalidCertificate = errors.New("http: invalid certificate")

	// ErrInvalidRequestBody được trả về khi request body không hợp lệ hoặc không đọc được.
	// Điều này xảy ra khi body không thể được parse hoặc có định dạng không hợp lệ.
	ErrInvalidRequestBody = errors.New("http: invalid request body")

	// ErrAdapterNotFound được trả về khi không tìm thấy adapter được yêu cầu.
	// Điều này xảy ra khi cố gắng sử dụng một adapter không tồn tại hoặc không được đăng ký.
	ErrAdapterNotFound = errors.New("http: adapter not found")

	// ErrAdapterNotSet được trả về khi cố gắng chạy server mà không thiết lập adapter.
	// Điều này xảy ra khi gọi Run() hoặc RunTLS() trước khi gọi SetAdapter().
	ErrAdapterNotSet = errors.New("http: adapter not set")

	// ErrInvalidContext được trả về khi context không hợp lệ hoặc đã hết hạn.
	// Điều này xảy ra khi cố gắng truy cập context đã hết hạn hoặc bị hủy.
	ErrInvalidContext = errors.New("invalid or expired context")

	// ErrServerAlreadyRunning được trả về khi cố gắng khởi động server đã chạy.
	// Điều này xảy ra khi gọi Run() hoặc RunTLS() nhiều lần.
	ErrServerAlreadyRunning = errors.New("server is already running")

	// ErrInvalidConfiguration được trả về khi cấu hình không hợp lệ.
	// Điều này xảy ra khi các giá trị cấu hình không đúng định dạng hoặc nằm ngoài phạm vi cho phép.
	ErrInvalidConfiguration = errors.New("invalid configuration")
)

// Adapter engine types định nghĩa các loại HTTP server engine hỗ trợ.
// Các constants này được sử dụng để xác định loại adapter đang được sử dụng.
const (
	// AdapterTypeNetHTTP là adapter cho Go's standard net/http package.
	AdapterTypeNetHTTP = "net/http"

	// AdapterTypeFastHTTP là adapter cho thư viện fasthttp với hiệu suất cao hơn.
	AdapterTypeFastHTTP = "fasthttp"

	// AdapterTypeHTTP2 là adapter cho HTTP/2 protocol.
	AdapterTypeHTTP2 = "http2"

	// AdapterTypeQuicH3 là adapter cho QUIC và HTTP/3 protocol.
	AdapterTypeQuicH3 = "quic-h3"
)
