// Package mocks provides reusable mock implementations for http/context interfaces.
// This package is used primarily for testing middleware and context-dependent code.
package mocks

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	httpContext "go.fork.vn/fork/context"
)

// MockContext is a reusable mock implementation of httpContext.Context for testing.
// It implements all methods required by the httpContext.Context interface and tracks
// information about method calls for testing assertions.
type MockContext struct {
	RequestObj        *http.Request
	ResponseWriterObj http.ResponseWriter
	Headers           map[string]string
	StatusCode        int
	Aborted           bool
	NextCalled        bool
	BodyWritten       []byte
	RouteParams       map[string]string
	Store             map[string]interface{}
}

// NewMockContext creates a new MockContext instance with default values.
// This is the recommended way to instantiate a MockContext for testing.
func NewMockContext(w http.ResponseWriter, r *http.Request) *MockContext {
	return &MockContext{
		RequestObj:        r,
		ResponseWriterObj: w,
		Headers:           make(map[string]string),
		StatusCode:        http.StatusOK,
		Aborted:           false,
		NextCalled:        false,
		RouteParams:       make(map[string]string),
		Store:             make(map[string]interface{}),
	}
}

// --- Implementation of httpContext.Context interface ---

func (c *MockContext) Request() httpContext.Request {
	return &MockRequest{R: c.RequestObj}
}

func (c *MockContext) Response() httpContext.Response {
	return &MockResponse{W: c.ResponseWriterObj, Ctx: c}
}

func (c *MockContext) Context() context.Context {
	return c.RequestObj.Context()
}

func (c *MockContext) WithContext(ctx context.Context) httpContext.Context {
	c.RequestObj = c.RequestObj.WithContext(ctx)
	return c
}

func (c *MockContext) Next()                              { c.NextCalled = true }
func (c *MockContext) Abort()                             { c.Aborted = true }
func (c *MockContext) IsAborted() bool                    { return c.Aborted }
func (c *MockContext) Set(key string, value interface{})  { c.Store[key] = value }
func (c *MockContext) Get(key string) (interface{}, bool) { v, ok := c.Store[key]; return v, ok }
func (c *MockContext) GetString(key string) string {
	v := c.Store[key]
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
func (c *MockContext) GetBool(key string) bool {
	v := c.Store[key]
	if b, ok := v.(bool); ok {
		return b
	}
	return false
}
func (c *MockContext) GetInt(key string) int {
	v := c.Store[key]
	if i, ok := v.(int); ok {
		return i
	}
	return 0
}
func (c *MockContext) GetInt64(key string) int64 {
	v := c.Store[key]
	if i, ok := v.(int64); ok {
		return i
	}
	return 0
}
func (c *MockContext) GetFloat64(key string) float64 {
	v := c.Store[key]
	if f, ok := v.(float64); ok {
		return f
	}
	return 0
}
func (c *MockContext) GetTime(key string) time.Time {
	v := c.Store[key]
	if t, ok := v.(time.Time); ok {
		return t
	}
	return time.Time{}
}
func (c *MockContext) GetDuration(key string) time.Duration {
	v := c.Store[key]
	if d, ok := v.(time.Duration); ok {
		return d
	}
	return 0
}
func (c *MockContext) GetStringSlice(key string) []string {
	v := c.Store[key]
	if s, ok := v.([]string); ok {
		return s
	}
	return nil
}
func (c *MockContext) GetStringMap(key string) map[string]interface{} {
	v := c.Store[key]
	if m, ok := v.(map[string]interface{}); ok {
		return m
	}
	return nil
}
func (c *MockContext) GetStringMapString(key string) map[string]string {
	v := c.Store[key]
	if m, ok := v.(map[string]string); ok {
		return m
	}
	return nil
}
func (c *MockContext) Method() string              { return c.RequestObj.Method }
func (c *MockContext) Path() string                { return c.RequestObj.URL.Path }
func (c *MockContext) RawPath() string             { return c.RequestObj.URL.RawPath }
func (c *MockContext) Param(name string) string    { return c.RouteParams[name] }
func (c *MockContext) ParamMap() map[string]string { return c.RouteParams }
func (c *MockContext) ParamArray(name string) []string {
	value := c.Param(name)
	if value == "" {
		return nil
	}
	return []string{value}
}
func (c *MockContext) Query(name string) string { return c.RequestObj.URL.Query().Get(name) }
func (c *MockContext) DefaultQuery(name, defaultValue string) string {
	v := c.Query(name)
	if v == "" {
		return defaultValue
	}
	return v
}
func (c *MockContext) QueryArray(name string) []string          { return c.RequestObj.URL.Query()[name] }
func (c *MockContext) QueryMap(prefix string) map[string]string { return map[string]string{} }
func (c *MockContext) Form(name string) string                  { return c.RequestObj.FormValue(name) }
func (c *MockContext) DefaultForm(name, defaultValue string) string {
	v := c.Form(name)
	if v == "" {
		return defaultValue
	}
	return v
}
func (c *MockContext) FormArray(name string) []string          { return c.RequestObj.Form[name] }
func (c *MockContext) FormMap(prefix string) map[string]string { return map[string]string{} }
func (c *MockContext) MultipartForm() (*multipart.Form, error) {
	return c.RequestObj.MultipartForm, nil
}
func (c *MockContext) FormFile(name string) (*multipart.FileHeader, error) {
	_, fh, err := c.RequestObj.FormFile(name)
	return fh, err
}
func (c *MockContext) SaveUploadedFile(file *multipart.FileHeader, dst string) error { return nil }
func (c *MockContext) BindJSON(obj interface{}) error                                { return nil }
func (c *MockContext) BindXML(obj interface{}) error                                 { return nil }
func (c *MockContext) BindQuery(obj interface{}) error                               { return nil }
func (c *MockContext) BindForm(obj interface{}) error                                { return nil }
func (c *MockContext) Bind(obj interface{}) error                                    { return nil }
func (c *MockContext) ShouldBind(obj interface{}) error                              { return nil }
func (c *MockContext) Status(code int) {
	c.StatusCode = code
	c.Response().WriteHeader(code)
}
func (c *MockContext) Header(key, value string)                         { c.Headers[key] = value }
func (c *MockContext) GetHeader(key string) string                      { return c.RequestObj.Header.Get(key) }
func (c *MockContext) ClientIP() string                                 { return "127.0.0.1" }
func (c *MockContext) ContentType() string                              { return c.Headers["Content-Type"] }
func (c *MockContext) IsWebsocket() bool                                { return false }
func (c *MockContext) GetRawData() ([]byte, error)                      { return nil, nil }
func (c *MockContext) Handlers() []func(httpContext.Context)            { return nil }
func (c *MockContext) SetHandlers(handlers []func(httpContext.Context)) {}
func (c *MockContext) ValidateStruct(obj interface{}) error             { return nil }
func (c *MockContext) ShouldBindAndValidate(obj interface{}) error {
	if err := c.ShouldBind(obj); err != nil {
		return err
	}
	return nil
}
func (c *MockContext) BindAndValidate(obj interface{}) error {
	if err := c.Bind(obj); err != nil {
		return err
	}
	return nil
}
func (c *MockContext) RegisterValidation(tag string, fn validator.Func) error         { return nil }
func (c *MockContext) RegisterCustomTranslation(tag string, translation string) error { return nil }
func (c *MockContext) GetValidator() *validator.Validate                              { return validator.New() }
func (c *MockContext) Cookie(name string) (string, error) {
	cookie, err := c.RequestObj.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
func (c *MockContext) Cookies() []*http.Cookie                  { return c.RequestObj.Cookies() }
func (c *MockContext) Error(err error)                          { c.StatusCode = http.StatusInternalServerError }
func (c *MockContext) File(filepath string)                     {}
func (c *MockContext) FileAttachment(filepath, filename string) {}
func (c *MockContext) Blob(code int, contentType string, data []byte) {
	c.StatusCode = code
	c.Headers["Content-Type"] = contentType
	c.BodyWritten = append(c.BodyWritten, data...)
}
func (c *MockContext) HTML(code int, html string) {
	c.StatusCode = code
	c.Headers["Content-Type"] = "text/html; charset=utf-8"
	c.BodyWritten = append(c.BodyWritten, []byte(html)...)
}

func (c *MockContext) JSON(code int, obj interface{}) {
	c.StatusCode = code
	c.Headers["Content-Type"] = "application/json; charset=utf-8"
	// Simplistic conversion without error handling for mocking purposes
	data, _ := json.Marshal(obj)
	c.BodyWritten = append(c.BodyWritten, data...)
}

func (c *MockContext) JSONP(code int, callback string, obj interface{}) {
	c.StatusCode = code
	c.Headers["Content-Type"] = "application/javascript; charset=utf-8"
	data, _ := json.Marshal(obj)
	result := []byte(callback + "(")
	result = append(result, data...)
	result = append(result, []byte(")")...)
	c.BodyWritten = append(c.BodyWritten, result...)
}

func (c *MockContext) XML(code int, obj interface{}) {
	c.StatusCode = code
	c.Headers["Content-Type"] = "application/xml; charset=utf-8"
	data, _ := xml.Marshal(obj)
	c.BodyWritten = append(c.BodyWritten, data...)
}

func (c *MockContext) Redirect(code int, location string) {
	c.StatusCode = code
	c.Headers["Location"] = location
}

func (c *MockContext) Render(code int, name string, data interface{}) {
	c.StatusCode = code
	// Mock implementation just stores the name and data without actual rendering
	c.Set("render_template", name)
	c.Set("render_data", data)
}

func (c *MockContext) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     path,
		Domain:   domain,
		MaxAge:   maxAge,
		Secure:   secure,
		HttpOnly: httpOnly,
	}
	http.SetCookie(c.ResponseWriterObj, cookie)
}

func (c *MockContext) Stream(code int, contentType string, r io.Reader) {
	c.StatusCode = code
	c.Headers["Content-Type"] = contentType
	// Read data from reader and store in bodyWritten for testing
	data, _ := io.ReadAll(r)
	c.BodyWritten = append(c.BodyWritten, data...)
}

func (c *MockContext) String(code int, format string, values ...interface{}) {
	c.StatusCode = code
	c.Headers["Content-Type"] = "text/plain; charset=utf-8"
	var data []byte
	if len(values) > 0 {
		data = []byte(c.FormatWithValues(format, values...))
	} else {
		data = []byte(format)
	}
	c.BodyWritten = append(c.BodyWritten, data...)
}

// Helper method to format a string with values
func (c *MockContext) FormatWithValues(format string, values ...interface{}) string {
	var result string
	if len(values) > 0 {
		result = c.sprintf(format, values...)
	} else {
		result = format
	}
	return result
}

// Helper method for sprintf (simplified)
func (c *MockContext) sprintf(format string, values ...interface{}) string {
	// This is a simplified implementation for mock purposes
	return format
}
