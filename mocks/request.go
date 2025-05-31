// Package mocks provides reusable mock implementations for http/context interfaces.
// This package is used primarily for testing middleware and context-dependent code.
package mocks

import (
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
)

// MockRequest is a reusable mock implementation of httpContext.Request for testing.
// It wraps an http.Request and implements all methods required by the
// httpContext.Request interface to facilitate testing of components
// that depend on this interface.
type MockRequest struct {
	R *http.Request
}

// Method returns the HTTP method of the request
func (r *MockRequest) Method() string {
	return r.R.Method
}

// URL returns the full URL of the request
func (r *MockRequest) URL() *url.URL {
	return r.R.URL
}

// Header returns the header of the request
func (r *MockRequest) Header() http.Header {
	return r.R.Header
}

// Cookies returns all cookies of the request
func (r *MockRequest) Cookies() []*http.Cookie {
	return r.R.Cookies()
}

// Cookie returns a cookie by name
func (r *MockRequest) Cookie(name string) (*http.Cookie, error) {
	return r.R.Cookie(name)
}

// Body returns the body of the request
func (r *MockRequest) Body() io.ReadCloser {
	return r.R.Body
}

// Form returns the form values of the request
func (r *MockRequest) Form() url.Values {
	return r.R.Form
}

// PostForm returns the post form values of the request
func (r *MockRequest) PostForm() url.Values {
	return r.R.PostForm
}

// FormValue returns a form value by name
func (r *MockRequest) FormValue(name string) string {
	return r.R.FormValue(name)
}

// PostFormValue returns a post form value by name
func (r *MockRequest) PostFormValue(name string) string {
	return r.R.PostFormValue(name)
}

// MultipartForm returns the multipart form of the request
func (r *MockRequest) MultipartForm() (*multipart.Form, error) {
	if r.R.MultipartForm == nil {
		if err := r.R.ParseMultipartForm(32 << 20); err != nil {
			return nil, err
		}
	}
	return r.R.MultipartForm, nil
}

// FormFile returns an uploaded file by name
func (r *MockRequest) FormFile(name string) (*multipart.FileHeader, error) {
	_, fileHeader, err := r.R.FormFile(name)
	return fileHeader, err
}

// RemoteAddr returns the client's address
func (r *MockRequest) RemoteAddr() string {
	return r.R.RemoteAddr
}

// ParseForm parses the form of the request
func (r *MockRequest) ParseForm() error {
	return r.R.ParseForm()
}

// ParseMultipartForm parses the multipart form of the request
func (r *MockRequest) ParseMultipartForm(maxMemory int64) error {
	return r.R.ParseMultipartForm(maxMemory)
}

// UserAgent returns the user agent of the request
func (r *MockRequest) UserAgent() string {
	return r.R.UserAgent()
}

// Referer returns the referer of the request
func (r *MockRequest) Referer() string {
	return r.R.Referer()
}

// ContentLength returns the content length of the request
func (r *MockRequest) ContentLength() int64 {
	return r.R.ContentLength
}

// Host returns the host of the request
func (r *MockRequest) Host() string {
	return r.R.Host
}

// RequestURI returns the URI of the request
func (r *MockRequest) RequestURI() string {
	return r.R.RequestURI
}

// Scheme returns the scheme of the request (http or https)
func (r *MockRequest) Scheme() string {
	if r.R.TLS != nil {
		return "https"
	}
	return "http"
}

// IsSecure checks if the request is secure (HTTPS)
func (r *MockRequest) IsSecure() bool {
	return r.Scheme() == "https"
}

// Protocol returns the protocol of the request
func (r *MockRequest) Protocol() string {
	return r.R.Proto
}

// Request returns the original http.Request
func (r *MockRequest) Request() *http.Request {
	return r.R
}
