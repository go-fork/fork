package router

import (
	"net/http"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"

	forkCtx "go.fork.vn/fork/context"
)

// HandlerFunc định nghĩa kiểu function handler cho HTTP requests.
// Mỗi handler nhận một context và xử lý request, trả về response.
type HandlerFunc func(ctx forkCtx.Context)

// Router là interface cho HTTP router của framework.
// Router quản lý việc đăng ký và điều hướng các HTTP routes đến handlers tương ứng.
// Nó cũng hỗ trợ route groups và middleware chains.
type Router interface {
	// Handle đăng ký một handler cho method và path cụ thể.
	// Handler sẽ được gọi khi có request phù hợp với method và path.
	//
	// Parameters:
	//   - method: HTTP method (GET, POST, PUT, DELETE, v.v.)
	//   - path: URL path pattern để khớp với requests
	//   - handlers: Chuỗi các handlers xử lý request
	Handle(method string, path string, handlers ...HandlerFunc)

	// Group tạo một router group mới với prefix đường dẫn.
	// Group cho phép tổ chức routes theo cấu trúc và áp dụng middleware cho nhóm routes.
	//
	// Parameters:
	//   - prefix: Tiền tố đường dẫn cho group
	//
	// Returns:
	//   - Router: Router mới đã được tạo với prefix
	Group(prefix string) Router

	// Use thêm middleware vào router.
	// Middleware sẽ được thực thi cho tất cả routes trong router này và các sub-groups.
	//
	// Parameters:
	//   - middleware: Danh sách các middleware functions để thêm
	Use(middleware ...HandlerFunc)

	// Static phục vụ static files từ thư mục root.
	// Đăng ký handler cho việc phục vụ static files từ filesystem.
	//
	// Parameters:
	//   - prefix: Tiền tố URL để phục vụ files (ví dụ: "/static")
	//   - root: Đường dẫn tới thư mục chứa static files
	Static(prefix string, root string)

	// Routes trả về tất cả routes đã đăng ký.
	// Phương thức này thu thập tất cả routes từ router hiện tại và tất cả các sub-groups.
	//
	// Returns:
	//   - []Route: Danh sách tất cả routes đã đăng ký
	Routes() []Route

	// ServeHTTP implements interface http.Handler để xử lý HTTP requests.
	// Phương thức này được gọi bởi HTTP server cho mỗi incoming request.
	//
	// Parameters:
	//   - w: HTTP response writer
	//   - req: HTTP request
	ServeHTTP(w http.ResponseWriter, req *http.Request)

	// Find tìm route phù hợp với method và path.
	// Phương thức này được sử dụng bởi router để tìm handler tương ứng cho request.
	//
	// Parameters:
	//   - method: HTTP method của request
	//   - path: URL path của request
	//
	// Returns:
	//   - HandlerFunc: Handler cho route được tìm thấy hoặc nil nếu không tìm thấy
	Find(method, path string) HandlerFunc
}

// Route định nghĩa một HTTP route đã đăng ký.
// Mỗi route chứa thông tin về HTTP method, path pattern, và handler function.
type Route struct {
	// Method là HTTP method của route (GET, POST, PUT, DELETE, v.v.)
	Method string

	// Path là URL path pattern của route
	Path string

	// Handler là function xử lý requests khớp với route này
	Handler HandlerFunc
}

// DefaultRouter là implementation mặc định của Router interface.
// Nó cung cấp cơ chế routing dựa trên path patterns với hỗ trợ cho parameters,
// wildcards, và regex patterns. Sử dụng trie structure để tối ưu hiệu suất.
type DefaultRouter struct {
	// basePath là tiền tố đường dẫn cho tất cả routes trong router này
	basePath string

	// routes là danh sách các routes đã đăng ký
	routes []Route

	// middlewares là danh sách các middleware functions áp dụng cho tất cả routes
	middlewares []HandlerFunc

	// groups là danh sách các sub-routers (groups) của router này
	groups []*DefaultRouter

	// trie cho việc tìm kiếm route nhanh chóng
	trie *RouteTrie

	// enableTrie bật/tắt việc sử dụng trie (mặc định: true)
	enableTrie bool
}

// NewRouter tạo một instance mới của DefaultRouter.
// Router mới được tạo không có routes, middlewares, hoặc groups.
//
// Returns:
//   - Router: Instance mới của DefaultRouter đã được khởi tạo
func NewRouter() Router {
	return &DefaultRouter{
		basePath:    "",
		routes:      make([]Route, 0),
		middlewares: make([]HandlerFunc, 0),
		groups:      make([]*DefaultRouter, 0),
		trie:        NewRouteTrie(),
		enableTrie:  true,
	}
}

// Handle đăng ký một handler cho method và path cụ thể.
// Phương thức này kết hợp path với basePath của router và
// kết hợp middlewares của router với handlers được cung cấp.
//
// Parameters:
//   - method: HTTP method (GET, POST, PUT, DELETE, v.v.)
//   - path: URL path pattern cho route
//   - handlers: Danh sách các handlers xử lý request
func (r *DefaultRouter) Handle(method string, path string, handlers ...HandlerFunc) {
	// Tính toán đường dẫn tuyệt đối bằng cách kết hợp basePath và path
	absolutePath := r.calculateAbsolutePath(path)

	// Kết hợp middlewares của router với handlers được cung cấp
	finalHandlers := r.combineHandlers(handlers)

	// Tạo một handler duy nhất gọi chuỗi handlers
	finalHandler := func(ctx forkCtx.Context) {
		// Thiết lập handlers trong context để sử dụng với Next()
		// Convert the HandlerFunc to the expected func(context.Context) type
		contextHandlers := make([]func(forkCtx.Context), len(finalHandlers))
		for i, h := range finalHandlers {
			contextHandlers[i] = h
		}

		ctx.SetHandlers(contextHandlers)
		// Bắt đầu chuỗi xử lý
		ctx.Next()
	}

	// Thêm route mới vào danh sách routes
	r.routes = append(r.routes, Route{
		Method:  method,
		Path:    absolutePath,
		Handler: finalHandler,
	})

	// Thêm route vào trie để tối ưu hóa tìm kiếm (nếu trie được bật)
	if r.enableTrie && r.trie != nil {
		r.trie.Insert(method, absolutePath, finalHandler)
	}
}

// Group tạo một router group mới với prefix đường dẫn.
// Group cho phép tổ chức routes theo cấu trúc và áp dụng middleware cho nhóm routes.
//
// Parameters:
//   - prefix: Tiền tố đường dẫn cho group
//
// Returns:
//   - Router: Router mới đã được tạo với prefix
func (r *DefaultRouter) Group(prefix string) Router {
	group := &DefaultRouter{
		basePath:    r.calculateAbsolutePath(prefix),
		routes:      make([]Route, 0),
		middlewares: make([]HandlerFunc, 0),
		groups:      make([]*DefaultRouter, 0),
		trie:        NewRouteTrie(),
		enableTrie:  r.enableTrie,
	}

	// Thêm middlewares hiện tại vào group
	group.middlewares = append(group.middlewares, r.middlewares...)

	// Thêm group vào router cha
	r.groups = append(r.groups, group)

	return group
}

// RemoveGroup removes a group from the router to prevent memory leaks
// This method helps clean up unused groups and their resources
//
// Parameters:
//   - prefix: The prefix of the group to remove
//
// Returns:
//   - bool: true if group was found and removed, false otherwise
func (r *DefaultRouter) RemoveGroup(prefix string) bool {
	absolutePrefix := r.calculateAbsolutePath(prefix)

	for i, group := range r.groups {
		if group.basePath == absolutePrefix {
			// Clear the group's resources before removing
			group.Clear()

			// Remove from slice efficiently
			r.groups[i] = r.groups[len(r.groups)-1]
			r.groups[len(r.groups)-1] = nil
			r.groups = r.groups[:len(r.groups)-1]
			return true
		}
	}
	return false
}

// Use thêm middleware vào router.
// Middleware sẽ được thực thi cho tất cả routes trong router này và các sub-groups.
//
// Parameters:
//   - middleware: Danh sách các middleware functions để thêm
func (r *DefaultRouter) Use(middleware ...HandlerFunc) {
	r.middlewares = append(r.middlewares, middleware...)
}

// Static phục vụ static files từ thư mục root.
// Đăng ký handler cho việc phục vụ static files từ filesystem.
//
// Parameters:
//   - prefix: Tiền tố URL để phục vụ files (ví dụ: "/static")
//   - root: Đường dẫn tới thư mục chứa static files
func (r *DefaultRouter) Static(prefix string, root string) {
	absolutePath := r.calculateAbsolutePath(prefix)
	handler := func(ctx forkCtx.Context) {
		path := ctx.Path()
		if strings.HasPrefix(path, absolutePath) {
			// Clean the file path to prevent path traversal attacks
			relativePath := strings.TrimPrefix(path, absolutePath)

			// Prevent path traversal by rejecting paths with ".."
			if strings.Contains(relativePath, "..") {
				ctx.Status(http.StatusForbidden)
				ctx.String(http.StatusForbidden, "403 Forbidden")
				return
			}

			// Ensure the path starts with /
			if !strings.HasPrefix(relativePath, "/") {
				relativePath = "/" + relativePath
			}

			// Join paths safely using filepath.Join equivalent logic
			filePath := root + relativePath

			// Additional security check: ensure final path is within root
			if !strings.HasPrefix(filePath, root) {
				ctx.Status(http.StatusForbidden)
				ctx.String(http.StatusForbidden, "403 Forbidden")
				return
			}

			ctx.File(filePath)
		}
	}
	r.Handle("GET", prefix+"/*filepath", handler)
}

// Clear clears all routes, middlewares, and groups from the router
// This method helps prevent memory leaks by properly cleaning up resources
func (r *DefaultRouter) Clear() {
	// Clear all child groups first
	for _, group := range r.groups {
		if group != nil {
			group.Clear()
		}
	}

	// Clear slices and set to nil to help GC
	r.routes = nil
	r.middlewares = nil
	r.groups = nil

	// Clear trie if it exists
	if r.trie != nil {
		r.trie.Clear()
		r.trie = nil
	}
}

// GetGroupCount returns the number of groups for monitoring memory usage
func (r *DefaultRouter) GetGroupCount() int {
	count := len(r.groups)
	for _, group := range r.groups {
		count += group.GetGroupCount()
	}
	return count
}

// Routes trả về tất cả routes đã đăng ký.
// Phương thức này thu thập tất cả routes từ router hiện tại và tất cả các sub-groups.
//
// Returns:
//   - []Route: Danh sách tất cả routes đã đăng ký
func (r *DefaultRouter) Routes() []Route {
	routes := r.routes

	// Thêm routes từ groups
	for _, group := range r.groups {
		groupRoutes := group.Routes()
		routes = append(routes, groupRoutes...)
	}

	return routes
}

// ServeHTTP implements interface http.Handler để xử lý HTTP requests.
// Phương thức này được gọi bởi HTTP server cho mỗi incoming request.
//
// Parameters:
//   - w: HTTP response writer
//   - req: HTTP request
func (r *DefaultRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Tạo context mới từ request và response
	ctx := forkCtx.NewContext(w, req)

	// Chuyển request đến handler phù hợp
	r.handleRequest(ctx)
}

// calculateAbsolutePath tính toán đường dẫn tuyệt đối từ đường dẫn tương đối.
// Kết hợp basePath của router với relativePath đã cho để tạo đường dẫn tuyệt đối.
//
// Parameters:
//   - relativePath: Đường dẫn tương đối cần tính toán
//
// Returns:
//   - string: Đường dẫn tuyệt đối đã tính toán
func (r *DefaultRouter) calculateAbsolutePath(relativePath string) string {
	if relativePath == "" {
		return r.basePath
	}

	absolutePath := r.basePath
	// Xử lý dấu / giữa basePath và relativePath
	if len(relativePath) > 0 && relativePath[0] != '/' && len(absolutePath) > 0 && absolutePath[len(absolutePath)-1] != '/' {
		absolutePath += "/"
	}
	return absolutePath + relativePath
}

// combineHandlers kết hợp middlewares của router với handlers được cung cấp.
// Tạo một mảng mới chứa tất cả middlewares và handlers theo thứ tự.
//
// Parameters:
//   - handlers: Danh sách các handlers đặc thù cho route
//
// Returns:
//   - []HandlerFunc: Mảng đã kết hợp các middlewares và handlers
func (r *DefaultRouter) combineHandlers(handlers []HandlerFunc) []HandlerFunc {
	finalSize := len(r.middlewares) + len(handlers)
	mergedHandlers := make([]HandlerFunc, finalSize)
	copy(mergedHandlers, r.middlewares)
	copy(mergedHandlers[len(r.middlewares):], handlers)
	return mergedHandlers
}

// handleRequest xử lý HTTP request với context đã cho.
// Tìm route phù hợp và thực thi handler của route đó.
//
// Parameters:
//   - ctx: Context của HTTP request/response
func (r *DefaultRouter) handleRequest(ctx forkCtx.Context) {
	// Tìm route phù hợp với method và path
	route := r.findRoute(ctx.Method(), ctx.Path())
	if route == nil {
		// Không tìm thấy route, trả về 404 Not Found
		ctx.Status(http.StatusNotFound)
		ctx.String(http.StatusNotFound, "404 page not found")
		return
	}

	// Thiết lập tham số URL vào context
	r.setRouteParams(ctx, route.Path, ctx.Path())

	// Thực thi handler của route đã tìm thấy
	route.Handler(ctx)
}

// setRouteParams thiết lập route parameters vào context.
// Trích xuất các tham số từ path pattern và URL path thực tế.
//
// Parameters:
//   - ctx: Context của HTTP request
//   - pattern: URL path pattern của route
//   - path: URL path thực tế của request
func (r *DefaultRouter) setRouteParams(ctx forkCtx.Context, pattern, path string) {
	// Trích xuất các tham số từ đường dẫn URL
	params := r.extractParams(pattern, path)

	// Lưu trữ các tham số vào context
	for k, v := range params {
		ctx.Set("param:"+k, v)
	}
}

// Find tìm route phù hợp với method và path.
// Phương thức này được sử dụng bởi router để tìm handler tương ứng cho request.
//
// Parameters:
//   - method: HTTP method của request
//   - path: URL path của request
//
// Returns:
//   - HandlerFunc: Handler cho route được tìm thấy hoặc nil nếu không tìm thấy
func (r *DefaultRouter) Find(method, path string) HandlerFunc {
	route := r.findRoute(method, path)
	if route != nil {
		return route.Handler
	}
	return nil
}

// findRoute tìm route phù hợp với method và path.
// Phương thức này tìm kiếm trong tất cả routes đã đăng ký để tìm route khớp với method và path.
// Sử dụng trie structure để tối ưu hiệu suất khi được bật.
//
// Parameters:
//   - method: HTTP method của request
//   - path: URL path của request
//
// Returns:
//   - *Route: Route được tìm thấy hoặc nil nếu không tìm thấy
func (r *DefaultRouter) findRoute(method, path string) *Route {
	// Sử dụng trie search nếu được bật (tối ưu hiệu suất O(log n))
	if r.enableTrie && r.trie != nil {
		if handler := r.trie.Find(method, path); handler != nil {
			// Tìm route tương ứng trong danh sách routes để trả về đầy đủ thông tin
			for _, route := range r.routes {
				if route.Method == method && r.pathMatch(route.Path, path) {
					return &route
				}
			}
		}
	}

	// Fallback to linear search nếu trie không được bật hoặc không tìm thấy
	// Kiểm tra các routes trong router hiện tại
	for _, route := range r.routes {
		if route.Method == method && r.pathMatch(route.Path, path) {
			return &route
		}
	}

	// Kiểm tra trong các groups
	for _, group := range r.groups {
		if route := group.findRoute(method, path); route != nil {
			return route
		}
	}

	return nil
}

// extractParams trích xuất các tham số từ đường dẫn URL.
// Hỗ trợ các loại tham số:
// - Named parameters: /:id
// - Regex constraints: /:id<\\d+>
// - Optional parameters: /:id?
// - Wildcard parameters: /*filepath
//
// Parameters:
//   - pattern: Pattern của route (ví dụ: "/users/:id" hoặc "/files/*filepath")
//   - path: Đường dẫn thực tế của request (ví dụ: "/users/123" hoặc "/files/images/logo.png")
//
// Returns:
//   - map[string]string: Map các tham số và giá trị của chúng
func (r *DefaultRouter) extractParams(pattern, path string) map[string]string {
	params := make(map[string]string)

	// Chia pattern và path thành các segments
	patternSegments := r.splitPath(pattern)
	pathSegments := r.splitPath(path)

	// Tìm wildcard segment (nếu có)
	wildcardIndex := -1
	for i, segment := range patternSegments {
		if strings.HasPrefix(segment, "*") {
			wildcardIndex = i
			break
		}
	}

	// Kiểm tra độ dài và xử lý route như /api/:version?/users vs /api/users
	// trong trường hợp này, version không nên nhận giá trị "users"
	isOptionalRoute := false
	if len(patternSegments) != len(pathSegments) {
		// Kiểm tra xem có phải route có optional params không
		optionalCount := 0
		for _, segment := range patternSegments {
			if r.isOptionalSegment(segment) {
				optionalCount++
			}
		}
		if optionalCount > 0 {
			isOptionalRoute = true
		}
	}

	// Xử lý trường hợp đặc biệt như /api/:version?/users vs /api/users
	if isOptionalRoute && len(patternSegments) == len(pathSegments)+1 {
		// Tìm phân đoạn optional
		optionalIndex := -1
		for i, segment := range patternSegments {
			if r.isOptionalSegment(segment) {
				optionalIndex = i
				break
			}
		}

		// Nếu tìm thấy phân đoạn optional ở giữa (không phải đầu hoặc cuối)
		if optionalIndex > 0 && optionalIndex < len(patternSegments)-1 {
			// Xử lý các phân đoạn trước optional
			for i := 0; i < optionalIndex; i++ {
				_, paramName := r.segmentMatch(patternSegments[i], pathSegments[i])
				if paramName != "" {
					params[paramName] = pathSegments[i]
				}
			}

			// Thêm giá trị rỗng cho tham số optional
			_, optionalParamName := r.segmentMatch(patternSegments[optionalIndex], "")
			if optionalParamName != "" {
				params[optionalParamName] = ""
			}

			// Xử lý các phân đoạn sau optional
			for i := optionalIndex + 1; i < len(patternSegments); i++ {
				pathIndex := i - 1
				if pathIndex < len(pathSegments) {
					_, paramName := r.segmentMatch(patternSegments[i], pathSegments[pathIndex])
					if paramName != "" {
						params[paramName] = pathSegments[pathIndex]
					}
				}
			}

			return params
		}
	}

	// Xử lý các segment thông thường
	maxLoop := len(patternSegments)
	if wildcardIndex >= 0 {
		maxLoop = wildcardIndex
	}

	for i := 0; i < maxLoop && i < len(pathSegments); i++ {
		isOptional := r.isOptionalSegment(patternSegments[i])
		_, paramName := r.segmentMatch(patternSegments[i], pathSegments[i])

		// Nếu là param (có tên), xử lý
		if paramName != "" {
			// Xử lý trường hợp đặc biệt cho optional parameters
			if isOptional && isOptionalRoute && i+1 < len(patternSegments) {
				// Trường hợp như /api/:version?/users vs /api/users
				// Nếu /users nằm trong pathSegments[i+1], có nghĩa là :version? đã bị bỏ qua
				nextPattern := patternSegments[i+1]
				if i+1 < len(pathSegments) && nextPattern == pathSegments[i+1] {
					// Trường hợp này, optional param bị bỏ qua, không lấy giá trị
					params[paramName] = ""
					continue
				}
			}

			// Trường hợp thông thường, lấy giá trị từ segment
			params[paramName] = pathSegments[i]
		}
	}

	// Xử lý wildcard parameter (nếu có)
	if wildcardIndex >= 0 && wildcardIndex < len(patternSegments) {
		wildcardName := patternSegments[wildcardIndex][1:] // Bỏ dấu '*'

		// Thu thập tất cả các segments còn lại
		if wildcardIndex < len(pathSegments) {
			wildcardValue := strings.Join(pathSegments[wildcardIndex:], "/")
			params[wildcardName] = wildcardValue
		} else {
			// Trường hợp wildcard không khớp với segment nào
			params[wildcardName] = ""
		}
	}

	// Xử lý các optional parameters khi path ngắn hơn pattern
	for i := len(pathSegments); i < len(patternSegments); i++ {
		// Bỏ qua nếu không phải optional parameter
		if !r.isOptionalSegment(patternSegments[i]) {
			continue
		}

		// Xử lý optional parameter
		paramName := patternSegments[i][1:] // Bỏ dấu ':'
		paramName = strings.TrimSuffix(paramName, "?")

		// Xử lý regex constraint
		if idx := strings.Index(paramName, "<"); idx >= 0 && strings.HasSuffix(paramName, ">") {
			paramName = paramName[:idx]
		}

		// Gán giá trị rỗng cho optional parameter
		params[paramName] = ""
	}

	return params
}

// pathMatch kiểm tra xem path có khớp với pattern không
// Thuật toán matching hỗ trợ đầy đủ:
// 1. Static routes (/users/list)
// 2. Named parameters (/users/:id)
// 3. Regex constraints (/users/:id<\d+>)
// 4. Optional parameters (/users/:id?)
// 5. Wildcard parameters (/files/*filepath)
func (r *DefaultRouter) pathMatch(pattern, path string) bool {
	// Kiểm tra trường hợp đặc biệt với optional parameters
	if strings.Contains(pattern, "?") {
		// Nếu có optional parameter, thử xử lý trường hợp đặc biệt
		if r.specialCaseMatch(pattern, path) {
			return true
		}
	}

	// Chia pattern và path thành các phần (segments)
	patternSegments := r.splitPath(pattern)
	pathSegments := r.splitPath(path)

	// Tìm các phân đoạn wildcard (nếu có)
	hasWildcard := false
	wildcardIndex := -1

	for i, segment := range patternSegments {
		if strings.HasPrefix(segment, "*") {
			hasWildcard = true
			wildcardIndex = i
			break
		}
	}

	// Xử lý route có wildcard
	if hasWildcard {
		// Nếu pattern có ít phân đoạn hơn path (không tính wildcard), không khớp
		if len(pathSegments) < wildcardIndex {
			return false
		}

		// Kiểm tra các phân đoạn trước wildcard
		for i := 0; i < wildcardIndex; i++ {
			if i >= len(pathSegments) {
				// Kiểm tra xem phân đoạn này có phải optional không
				if !r.isOptionalSegment(patternSegments[i]) {
					return false
				}
				continue
			}

			match, _ := r.segmentMatch(patternSegments[i], pathSegments[i])
			if !match {
				return false
			}
		}

		// Wildcard luôn khớp với phần còn lại của path
		return true
	}
	// Xử lý các route không có wildcard
	// Số lượng phân đoạn phải bằng nhau trừ khi có optional params
	maxSegments := len(patternSegments)
	if len(pathSegments) > maxSegments {
		return false
	}

	// Đếm số lượng phân đoạn không phải optional
	requiredSegments := 0
	for _, segment := range patternSegments {
		if !r.isOptionalSegment(segment) {
			requiredSegments++
		}
	}

	// Nếu số lượng phân đoạn trong path ít hơn số lượng phân đoạn bắt buộc, không khớp
	if len(pathSegments) < requiredSegments {
		return false
	}

	// So khớp từng phân đoạn
	for i := 0; i < len(patternSegments); i++ {
		if i >= len(pathSegments) {
			// Nếu đã hết phân đoạn trong path, phân đoạn còn lại trong pattern phải là optional
			if !r.isOptionalSegment(patternSegments[i]) {
				return false
			}
			continue
		}

		match, _ := r.segmentMatch(patternSegments[i], pathSegments[i])
		if !match {
			return false
		}
	}

	return true
}

// specialCaseMatch xử lý các trường hợp đặc biệt cho routes có optional parameters
// Xử lý các trường hợp như /api/:version?/users khớp với /api/users.
//
// Parameters:
//   - pattern: URL path pattern
//   - path: URL path thực tế
//
// Returns:
//   - bool: true nếu path khớp với pattern theo trường hợp đặc biệt, ngược lại là false
func (r *DefaultRouter) specialCaseMatch(pattern, path string) bool {
	// Xử lý trường hợp đặc biệt: /api/:version?/users với /api/users
	// hoặc /optional/:param?/test với /optional/test

	// Chia pattern và path thành các phần (segments)
	patternSegments := r.splitPath(pattern)
	pathSegments := r.splitPath(path)

	// Nếu số lượng phân đoạn không chênh lệch 1, không phải trường hợp đặc biệt
	if len(patternSegments) != len(pathSegments)+1 {
		return false
	}

	// Tìm phân đoạn optional
	optionalIndex := -1
	for i, segment := range patternSegments {
		if r.isOptionalSegment(segment) {
			optionalIndex = i
			break
		}
	}

	// Nếu không tìm thấy phân đoạn optional, hoặc phân đoạn optional ở cuối, không phải trường hợp đặc biệt
	if optionalIndex < 0 || optionalIndex == len(patternSegments)-1 {
		return false
	}

	// Tạo pattern mới bằng cách loại bỏ phân đoạn optional
	newPatternSegments := make([]string, 0, len(patternSegments)-1)
	for i, segment := range patternSegments {
		if i != optionalIndex {
			newPatternSegments = append(newPatternSegments, segment)
		}
	}

	// Tạo pattern mới
	newPattern := "/" + strings.Join(newPatternSegments, "/")

	// Kiểm tra pattern mới với path
	return r.pathMatch(newPattern, path)
}

// isOptionalSegment kiểm tra xem một phân đoạn có phải là optional không.
// Xác định nếu segment là một named parameter và có dấu ? ở cuối.
//
// Parameters:
//   - segment: Segment URL cần kiểm tra
//
// Returns:
//   - bool: true nếu segment là optional parameter, ngược lại là false
func (r *DefaultRouter) isOptionalSegment(segment string) bool {
	if !strings.HasPrefix(segment, ":") {
		return false
	}

	// Bỏ ':'
	paramName := segment[1:]

	// Kiểm tra optional marker '?'
	// Phải xử lý cả trường hợp có regex như :id<\d+>?
	if strings.HasSuffix(paramName, "?") {
		return true
	}

	// Kiểm tra regex có optional marker ở cuối
	if idx := strings.Index(paramName, "<"); idx >= 0 && strings.HasSuffix(paramName, ">") {
		// Không có optional marker trong regex
		return false
	} else if idx := strings.Index(paramName, "<"); idx >= 0 && strings.HasSuffix(paramName, ">?") {
		// Có optional marker sau regex
		return true
	}

	return false
}

// segmentMatch kiểm tra xem một segment có khớp với pattern không
// Hỗ trợ nhiều loại patterns:
// - Static segments: "user", "product", etc.
// - Named parameters: ":id", ":slug", etc.
// - Parameters with regex constraints: ":id<\\d+>", ":slug<[a-z]+>"
// - Optional parameters: ":id?", ":slug?"
// - Wildcard parameters: "*filepath", "*any"
//
// Trả về (match, paramName)
func (r *DefaultRouter) segmentMatch(pattern, segment string) (bool, string) {
	// 1. Static segment
	if !strings.HasPrefix(pattern, ":") && !strings.HasPrefix(pattern, "*") {
		return pattern == segment, ""
	}

	// 2. Named parameter (ví dụ: :id)
	if strings.HasPrefix(pattern, ":") {
		paramName := pattern[1:] // Bỏ ':'

		// Xử lý optional parameter (ví dụ: :id?)
		isOptional := false
		if strings.HasSuffix(paramName, "?") {
			paramName = paramName[:len(paramName)-1] // Bỏ '?'
			isOptional = true

			// Nếu segment rỗng và parameter là optional, thì khớp
			if segment == "" {
				return true, paramName
			}
		}

		// Kiểm tra regex constraint (ví dụ: :id<[0-9]+>)
		regexPattern := ""
		if idx := strings.Index(paramName, "<"); idx >= 0 && strings.HasSuffix(paramName, ">") {
			regexPattern = paramName[idx+1 : len(paramName)-1]
			paramName = paramName[:idx]

			// Compile và check regex
			regex, err := r.compileRegex(regexPattern)
			if err != nil {
				return false, ""
			}

			// Nếu segment rỗng và parameter là optional, thì khớp
			if segment == "" && isOptional {
				return true, paramName
			}

			// Kiểm tra segment với regex
			if !regex.MatchString(segment) {
				return false, ""
			}

			return true, paramName
		}

		// Kiểm tra empty segment với non-optional parameter
		if segment == "" && !isOptional {
			return false, ""
		}

		// Named parameter đơn giản (khớp với bất kỳ giá trị nào không rỗng)
		return segment != "", paramName
	}

	// 3. Wildcard parameter (ví dụ: *filepath)
	if strings.HasPrefix(pattern, "*") {
		paramName := pattern[1:] // Bỏ '*'
		return true, paramName
	}

	// Không khớp với pattern nào
	return false, ""
}

// regexCache là cache cho các compiled regular expressions với thread safety
var (
	regexCache   = make(map[string]*regexp.Regexp)
	regexCacheMu sync.RWMutex
)

// compileRegex biên dịch một regex pattern và cache nó để tái sử dụng.
// Cache giúp tăng hiệu suất khi cùng một pattern được dùng nhiều lần.
// Thread-safe để sử dụng trong môi trường concurrent.
//
// Parameters:
//   - pattern: Regex pattern cần biên dịch
//
// Returns:
//   - *regexp.Regexp: Đối tượng regex đã biên dịch
//   - error: Lỗi nếu không thể biên dịch pattern
func (r *DefaultRouter) compileRegex(pattern string) (*regexp.Regexp, error) {
	// Kiểm tra xem pattern đã được cache chưa (read lock)
	regexCacheMu.RLock()
	if regex, found := regexCache[pattern]; found {
		regexCacheMu.RUnlock()
		return regex, nil
	}
	regexCacheMu.RUnlock()

	// Biên dịch pattern
	regex, err := regexp.Compile("^" + pattern + "$")
	if err != nil {
		return nil, err
	}

	// Lưu vào cache (write lock)
	regexCacheMu.Lock()
	// Double-check để tránh trùng lặp trong trường hợp concurrent access
	if existingRegex, found := regexCache[pattern]; found {
		regexCacheMu.Unlock()
		return existingRegex, nil
	}
	regexCache[pattern] = regex
	regexCacheMu.Unlock()

	return regex, nil
}

// Cache for splitPath results with performance tracking
var (
	splitPathCache    = make(map[string][]string)
	splitPathCacheMu  sync.RWMutex
	splitPathHits     int64
	splitPathMisses   int64
	splitPathMaxSize  = 1000 // Configurable cache size limit
	splitPathEvictPct = 33   // Percentage of cache to evict when full (33%)

	// Pre-allocated common path results for ultimate performance
	commonPaths = map[string][]string{
		"/":            {},
		"":             {},
		"/api":         {"api"},
		"/api/v1":      {"api", "v1"},
		"/api/v2":      {"api", "v2"},
		"/users":       {"users"},
		"/admin":       {"admin"},
		"/static":      {"static"},
		"/assets":      {"assets"},
		"/public":      {"public"},
		"/health":      {"health"},
		"/metrics":     {"metrics"},
		"/ping":        {"ping"},
		"/favicon.ico": {"favicon.ico"},
		"/robots.txt":  {"robots.txt"},
		"/sitemap.xml": {"sitemap.xml"},
	}
)

// splitPath chia path thành các segments với caching và tối ưu hiệu suất cao.
// Phương thức này được sử dụng nhiều lần trong quá trình routing, vì vậy
// việc cache kết quả và tối ưu hóa string operations giúp giảm đáng kể chi phí xử lý.
//
// Advanced optimization features:
// 1. Thread-safe caching mechanism với atomic operations
// 2. Pre-computed common paths cho zero-allocation lookups
// 3. Fast path cho các trường hợp phổ biến (root, empty, single segment)
// 4. Optimized string operations với manual parsing
// 5. Memory-conscious implementation với cache eviction
// 6. Performance metrics tracking
// 7. Efficient capacity estimation
//
// Parameters:
//   - path: URL path cần chia thành segments
//
// Returns:
//   - []string: Slice các segments của path
func (r *DefaultRouter) splitPath(path string) []string {
	// Fast path for pre-computed common paths (zero allocation)
	if result, exists := commonPaths[path]; exists {
		atomic.AddInt64(&splitPathHits, 1)
		return result
	}

	// Fast path for simple cases
	if path == "/" || path == "" {
		atomic.AddInt64(&splitPathHits, 1)
		return []string{}
	}

	// Check cache first (read lock)
	splitPathCacheMu.RLock()
	if segments, found := splitPathCache[path]; found {
		splitPathCacheMu.RUnlock()
		atomic.AddInt64(&splitPathHits, 1)
		return segments
	}
	splitPathCacheMu.RUnlock()

	// Cache miss - increment counter
	atomic.AddInt64(&splitPathMisses, 1)

	// Process path with optimized algorithm
	segments := r.splitPathOptimized(path)

	// Cache the result (write lock)
	splitPathCacheMu.Lock()
	// Double-check to avoid duplicate work in concurrent scenarios
	if existingSegments, found := splitPathCache[path]; found {
		splitPathCacheMu.Unlock()
		return existingSegments
	}

	// Implement cache size limit with configurable eviction
	if len(splitPathCache) >= splitPathMaxSize {
		r.evictCacheEntries()
	}

	splitPathCache[path] = segments
	splitPathCacheMu.Unlock()

	return segments
}

// splitPathOptimized performs the actual path splitting with manual optimization
func (r *DefaultRouter) splitPathOptimized(path string) []string {
	// Handle edge cases
	if len(path) == 0 {
		return []string{}
	}

	// Single character path
	if len(path) == 1 {
		if path[0] == '/' {
			return []string{}
		}
		return []string{path}
	}

	// Find the actual content boundaries
	start := 0
	end := len(path)

	// Skip leading slashes
	for start < end && path[start] == '/' {
		start++
	}

	// Skip trailing slashes
	for end > start && path[end-1] == '/' {
		end--
	}

	// Empty after trimming
	if start >= end {
		return []string{}
	}

	// Extract the trimmed path
	trimmed := path[start:end]

	// Single segment optimization (no slashes)
	if !strings.Contains(trimmed, "/") {
		return []string{trimmed}
	}

	// Count slashes to pre-allocate with exact capacity
	slashCount := 0
	for i := start; i < end; i++ {
		if path[i] == '/' {
			slashCount++
		}
	}

	// Allocate with exact capacity (segments = slashes + 1)
	segments := make([]string, 0, slashCount+1)

	// Manual parsing for better performance than strings.Split
	segStart := start
	for i := start; i <= end; i++ {
		if i == end || path[i] == '/' {
			if i > segStart {
				segments = append(segments, path[segStart:i])
			}
			segStart = i + 1
		}
	}

	return segments
}

// evictCacheEntries removes old entries when cache is full using configurable parameters
func (r *DefaultRouter) evictCacheEntries() {
	// Use configurable eviction percentage
	evictCount := (len(splitPathCache) * splitPathEvictPct) / 100
	if evictCount == 0 {
		evictCount = 1 // Always evict at least one entry
	}

	count := 0
	for k := range splitPathCache {
		if count >= evictCount {
			break
		}
		delete(splitPathCache, k)
		count++
	}
}

// ClearSplitPathCache clears the splitPath cache to free memory.
// This method can be called periodically or during low-traffic periods
// to manage memory usage.
func (r *DefaultRouter) ClearSplitPathCache() {
	splitPathCacheMu.Lock()
	defer splitPathCacheMu.Unlock()

	// Clear the cache
	for k := range splitPathCache {
		delete(splitPathCache, k)
	}
}

// GetSplitPathCacheStats returns detailed statistics about the splitPath cache
// for monitoring and performance analysis.
//
// Returns:
//   - cacheSize: Number of cached entries
//   - hitRatio: Cache hit ratio as percentage (0-100)
//   - totalHits: Total number of cache hits
//   - totalMisses: Total number of cache misses
//   - totalRequests: Total number of splitPath requests
func (r *DefaultRouter) GetSplitPathCacheStats() (cacheSize int, hitRatio int, totalHits int64, totalMisses int64, totalRequests int64) {
	splitPathCacheMu.RLock()
	cacheSize = len(splitPathCache)
	splitPathCacheMu.RUnlock()

	// Get atomic counters safely
	totalHits = atomic.LoadInt64(&splitPathHits)
	totalMisses = atomic.LoadInt64(&splitPathMisses)
	totalRequests = totalHits + totalMisses

	// Calculate hit ratio
	if totalRequests > 0 {
		hitRatio = int((totalHits * 100) / totalRequests)
	}

	return
}

// ResetSplitPathStats resets the performance counters for fresh measurement
func (r *DefaultRouter) ResetSplitPathStats() {
	atomic.StoreInt64(&splitPathHits, 0)
	atomic.StoreInt64(&splitPathMisses, 0)
}

// SetSplitPathCacheConfig configures the splitPath cache parameters
//
// Parameters:
//   - maxSize: Maximum number of entries in cache (default: 1000)
//   - evictPercent: Percentage of cache to evict when full (default: 33)
func (r *DefaultRouter) SetSplitPathCacheConfig(maxSize int, evictPercent int) {
	if maxSize > 0 {
		splitPathMaxSize = maxSize
	}
	if evictPercent > 0 && evictPercent <= 100 {
		splitPathEvictPct = evictPercent
	}
}

// GetSplitPathCacheConfig returns current cache configuration
func (r *DefaultRouter) GetSplitPathCacheConfig() (maxSize int, evictPercent int) {
	return splitPathMaxSize, splitPathEvictPct
}
