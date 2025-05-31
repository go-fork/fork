// package app cung cấp một hệ thống HTTP linh hoạt và có thể mở rộng cho ứng dụng.
//
// Package này sử dụng thiết kế dựa trên adapter pattern, cho phép sử dụng nhiều HTTP engine
// khác nhau (net/http, fasthttp, http2, quic-h3) với khả năng tùy chỉnh và mở rộng cao.
// Hỗ trợ đầy đủ các thành phần cần thiết cho việc xây dựng RESTful API, web services
// và ứng dụng web hiện đại.
//
// Các tính năng chính:
//   - Hỗ trợ nhiều HTTP engine thông qua adapter pattern
//   - Router mạnh mẽ với hỗ trợ params và wildcard
//   - Hệ thống middleware linh hoạt với middleware groups
//   - Context API mạnh mẽ để quản lý request/response
//   - Binding và validation cho request data (JSON, XML, Form)
//   - Hỗ trợ nhiều định dạng response (JSON, XML, HTML, Blob)
//   - Hệ thống xử lý lỗi HTTP nhất quán và mở rộng
//   - Struct HttpError chung với status_code, message, details và error gốc
//   - Hàm tạo lỗi cho tất cả các mã trạng thái HTTP phổ biến
//   - Tích hợp sẵn với validation errors (422 Unprocessable Entity)
//   - Chi tiết lỗi tùy chỉnh và lỗi gốc để debug
//   - Triển khai error interface và hỗ trợ unwrapping
//   - File uploads và file serving
//   - Quản lý session và cookies
//   - Tích hợp sẵn với Dependency Injection container
//   - Middleware chains thống nhất giữa các adapter khác nhau
//   - Xử lý handlers dạng chuỗi với Context.Next()
//   - Thứ tự thực thi middleware nhất quán giữa các adapters
//   - Xử lý chính xác chuỗi middleware khi abort request
//   - Hỗ trợ đồng thời cho HTTP/1.1, HTTP/2 và HTTP/3 (QUIC)
//
// Cấu trúc package:
//   - adapter/: Các adapter cho các HTTP engine khác nhau
//   - net.go: Adapter cho chuẩn net/http của Go
//   - fasthttp.go: Adapter cho fasthttp - HTTP server hiệu suất cao
//   - xhttp2.go: Adapter hỗ trợ HTTP/2
//   - quic-h3.go: Adapter hỗ trợ HTTP/3 qua QUIC
//   - context/: Quản lý request/response context
//   - context.go: Context API chính
//   - request.go: Request wrapper và utilities
//   - response.go: Response wrapper và utilities
//   - router/: Router và định tuyến đường dẫn
//   - router.go: Router core, hỗ trợ routing và middleware
//   - application.go: Quản lý ứng dụng HTTP, entry point chính
//   - constants.go: Các hằng số và error messages
//   - config.go: Cấu hình HTTP server và behavior
//   - provider.go: ServiceProvider để tích hợp với container DI
//
// Ví dụ sử dụng cơ bản:
//
//	// Khởi tạo container
//	c := container.NewContainer()
//
//	// Khởi tạo application
//	app := http.NewWebApp(c)
//
//	// Đăng ký route
//	app.GET("/users", func(ctx context.Context) {
//		ctx.JSON(200, GetUsers())
//	})
//
//	// Sử dụng middleware
//	app.Use(Logger(), Recovery())
//
//	// Tạo route group với prefix
//	api := app.Group("/api")
//	api.Use(Auth()) // Middleware chỉ cho group này
//
//	// Đăng ký route trong group
//	api.GET("/profile", GetProfile)
//	api.POST("/settings", UpdateSettings)
//
//	// Khởi động server
//	app.Listen(":8080")
//
// Thiết kế bảo mật:
//   - Cung cấp các struct private (httpContext, httpRequest, httpResponse)
//     chỉ có thể truy cập thông qua interface công khai
//   - Kiểm soát chặt chẽ việc tiếp xúc với đối tượng nền tảng
//   - Validation tích hợp để ngăn chặn các lỗ hổng phổ biến
//   - Xử lý middleware chains an toàn, không cho phép truy cập trực tiếp vào handlers
//   - Chuỗi xử lý request rõ ràng và có thể kiểm soát thông qua Context.Next()
//   - Bảo vệ trạng thái Context với cơ chế abort để ngăn chặn việc xử lý không mong muốn
//   - Đảm bảo thứ tự thực thi handlers nhất quán giữa các adapter khác nhau
//   - Cơ chế nhất quán để quản lý middleware chains trong mọi adapter (Net, FastHTTP, HTTP/2, QUIC)
//   - Bảo vệ toàn vẹn của chuỗi middleware thông qua các method Handlers() và SetHandlers()
//   - Ngăn chặn truy cập trực tiếp vào handlers nội bộ của httpContext qua interface
//   - Xử lý graceful shutdown để đảm bảo tất cả request đang xử lý được hoàn thành
//   - Kiểm soát chặt chẽ việc tiếp xúc với các đối tượng http.Request và http.ResponseWriter
//
// Testing:
//   - Test files đầy đủ cho tất cả thành phần: context, request, response
//   - Tests cho tất cả adapter (net/http, fasthttp, http2, quic-h3)
//   - Tests bao phủ router và application
//   - Đảm bảo code coverage cao và chất lượng code
//   - Tests cho middleware chains và thứ tự thực thi
//   - Tests cho các tình huống lỗi và edge cases
//   - Kiểm tra đầy đủ quá trình xử lý chuỗi middleware và Context.Next()
//   - Verify thứ tự thực thi middleware trong các adapter khác nhau
//   - Tests cho các handlers dùng Abort() trong middleware chains
//   - Kiểm tra phương thức Handlers() và SetHandlers() của Context
//   - Nhất quán middleware execution giữa các HTTP engines
//   - Benchmark performace của các adapter trong các tình huống khác nhau
//
// Cách thức xử lý Middleware:
//   - Mỗi middleware là một hàm nhận Context và gọi Next() để truyền request đến middleware tiếp theo
//   - Các middleware được đăng ký thông qua phương thức Use() của WebApp hoặc RouteGroup
//   - Thứ tự đăng ký middleware quyết định thứ tự thực thi
//   - Khái niệm "onion model": middleware pre-process trước khi vào handler chính, post-process sau khi handler thực thi
//   - Middleware có thể dùng Abort() để dừng chuỗi xử lý
//   - Việc triển khai middleware nhất quán giữa các adapters khác nhau
//   - WebApp hỗ trợ global middlewares và group-specific middlewares
//
// Package này được thiết kế theo nguyên tắc SOLID, dễ dàng mở rộng
// và tùy chỉnh, phù hợp cho cả ứng dụng nhỏ và các hệ thống phân tán phức tạp.
//
// Xử lý lỗi HTTP:
//   - PAMM cung cấp struct HttpError để đại diện cho lỗi HTTP với mã trạng thái, thông báo và chi tiết
//   - Hỗ trợ tạo các lỗi HTTP chuẩn (400, 401, 403, 404, 422, 500, v.v.)
//   - Có thể đính kèm thông tin chi tiết tùy chỉnh và lỗi gốc
//   - Các hàm tạo lỗi đơn giản: BadRequest, Unauthorized, Forbidden, NotFound, v.v.
//   - Các hàm tạo lỗi đầy đủ: NewBadRequest, NewUnauthorized, NewForbidden, NewNotFound, v.v.
//   - Tuân thủ interface error và hỗ trợ unwrapping thông qua errors.Unwrap
//   - Có thể sử dụng với errors.Is và errors.As
//   - Giữ cho xử lý lỗi nhất quán trong toàn bộ ứng dụng
//
// Ví dụ sử dụng HttpError:
//
//	// Tạo lỗi 400 đơn giản
//	err := http.BadRequest("Thông tin không hợp lệ")
//	ctx.JSON(err.StatusCode, err)
//
//	// Tạo lỗi 404 với chi tiết và lỗi gốc
//	origErr := sql.ErrNoRows
//	details := map[string]interface{}{
//		"user_id": 123,
//	}
//	err := http.NewNotFound("Người dùng không tồn tại", details, origErr)
//	ctx.JSON(err.StatusCode, err)
//
// Xem thêm ví dụ trong /examples/infra/http/error_handling
package fork
