// Package mocks provides reusable mock implementations for http/context interfaces.
// This package is used primarily for testing middleware and context-dependent code.
package mocks

import (
	"bufio"
	"net"
	"net/http"
)

// MockResponse is a reusable mock implementation of httpContext.Response for testing.
// It implements all methods required by the httpContext.Response interface and
// provides additional functionality for testing assertions about response data.
// This mock tracks write operations and status code changes for verification in tests.
type MockResponse struct {
	W   http.ResponseWriter
	Ctx *MockContext

	statusCode   int
	bytesWritten int
	isWritten    bool
	Headers      http.Header
}

// Header returns the header of the response
func (r *MockResponse) Header() http.Header {
	if r.Headers == nil {
		r.Headers = make(http.Header)
	}
	return r.Headers
}

// Write writes data to the response
func (r *MockResponse) Write(data []byte) (int, error) {
	n, err := r.W.Write(data)
	r.isWritten = true
	r.Ctx.BodyWritten = append(r.Ctx.BodyWritten, data...)
	r.bytesWritten += n
	return n, err
}

// WriteHeader sets the HTTP status code for the response
func (r *MockResponse) WriteHeader(code int) {
	if r.isWritten {
		return
	}
	r.statusCode = code
	r.Ctx.StatusCode = code
	r.isWritten = true
	r.W.WriteHeader(code) // propagate to underlying ResponseWriter
}

// Flush writes any buffered data to the client
func (r *MockResponse) Flush() {
	if flusher, ok := r.W.(http.Flusher); ok {
		flusher.Flush()
	}
}

// Status returns the status code of the response
func (r *MockResponse) Status() int {
	return r.statusCode
}

// Size returns the size of the response
func (r *MockResponse) Size() int {
	return r.bytesWritten
}

// Written checks if the response has been written
func (r *MockResponse) Written() bool {
	return r.isWritten
}

// Hijack allows taking control of the connection
func (r *MockResponse) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := r.W.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, http.ErrNotSupported
}

// ResponseWriter returns the original http.ResponseWriter
func (r *MockResponse) ResponseWriter() http.ResponseWriter {
	return r.W
}

// Reset resets the response writer to its initial state
func (r *MockResponse) Reset(w http.ResponseWriter) {
	r.W = w
	r.statusCode = http.StatusOK
	r.bytesWritten = 0
	r.isWritten = false
	r.Headers = make(http.Header)
}

// Pusher returns an http.Pusher if HTTP/2 server push is supported
func (r *MockResponse) Pusher() (http.Pusher, bool) {
	pusher, ok := r.W.(http.Pusher)
	return pusher, ok
}
