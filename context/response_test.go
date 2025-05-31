package context

import (
	"bufio"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewResponse(t *testing.T) {
	w := httptest.NewRecorder()
	response := NewResponse(w)

	if response == nil {
		t.Fatal("Expected response to be created, got nil")
	}
}

func TestResponseHeader(t *testing.T) {
	w := httptest.NewRecorder()
	response := NewResponse(w)

	// Set a header via response
	response.Header().Set("X-Test", "value")

	// Check the header was set
	if w.Header().Get("X-Test") != "value" {
		t.Errorf("Expected header X-Test to be set to 'value', got %q", w.Header().Get("X-Test"))
	}
}

func TestResponseWrite(t *testing.T) {
	w := httptest.NewRecorder()
	response := NewResponse(w)

	data := []byte("test data")
	n, err := response.Write(data)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if n != len(data) {
		t.Errorf("Expected to write %d bytes, wrote %d", len(data), n)
	}

	if w.Body.String() != string(data) {
		t.Errorf("Expected body %q, got %q", string(data), w.Body.String())
	}

	if !response.Written() {
		t.Error("Expected response to be marked as written")
	}

	if response.Size() != len(data) {
		t.Errorf("Expected response size to be %d, got %d", len(data), response.Size())
	}
}

func TestResponseWriteHeader(t *testing.T) {
	w := httptest.NewRecorder()
	response := NewResponse(w)

	response.WriteHeader(http.StatusCreated)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	if response.Status() != http.StatusCreated {
		t.Errorf("Expected Status() to return %d, got %d", http.StatusCreated, response.Status())
	}

	if !response.Written() {
		t.Error("Expected response to be marked as written")
	}

	// Writing again should not change the status
	response.WriteHeader(http.StatusOK)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code to remain %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestResponseFlush(t *testing.T) {
	// Create a custom response writer that implements Flusher
	customWriter := &customResponseWriter{
		ResponseRecorder: httptest.NewRecorder(),
		flushed:          false,
	}

	response := NewResponse(customWriter)
	response.Flush()

	if !customWriter.flushed {
		t.Error("Expected Flush to be called")
	}
}

// Mock Flusher implementation
type customResponseWriter struct {
	*httptest.ResponseRecorder
	flushed bool
}

func (crw *customResponseWriter) Flush() {
	crw.flushed = true
}

func TestResponseReset(t *testing.T) {
	w1 := httptest.NewRecorder()
	response := NewResponse(w1)

	// Write some data and set status
	response.WriteHeader(http.StatusCreated)
	response.Write([]byte("test"))

	// Verify initial state
	if response.Status() != http.StatusCreated {
		t.Errorf("Initial status should be %d, got %d", http.StatusCreated, response.Status())
	}

	if response.Size() != 4 {
		t.Errorf("Initial size should be 4, got %d", response.Size())
	}

	if !response.Written() {
		t.Error("Initial written flag should be true")
	}

	// Reset with a new writer
	w2 := httptest.NewRecorder()
	response.Reset(w2)

	// Verify reset state
	if response.Status() != http.StatusOK {
		t.Errorf("After reset, status should be %d, got %d", http.StatusOK, response.Status())
	}

	if response.Size() != 0 {
		t.Errorf("After reset, size should be 0, got %d", response.Size())
	}

	if response.Written() {
		t.Error("After reset, written flag should be false")
	}

	// Verify that we're writing to the new writer
	response.Write([]byte("new"))
	if w2.Body.String() != "new" {
		t.Errorf("Expected to write to new writer, got %q", w2.Body.String())
	}
}

func TestResponseResponseWriter(t *testing.T) {
	w := httptest.NewRecorder()
	response := NewResponse(w)

	if response.ResponseWriter() != w {
		t.Error("Expected ResponseWriter() to return the original writer")
	}
}

func TestResponseHijack(t *testing.T) {
	// Create a mock conn for testing Hijack
	mockConn := &mockConn{}
	mockHijacker := &mockResponseWriter{
		conn: mockConn,
		bufw: bufio.NewReadWriter(bufio.NewReader(nil), bufio.NewWriter(nil)),
	}

	response := NewResponse(mockHijacker)
	conn, bufrw, err := response.Hijack()

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if conn != mockConn {
		t.Error("Expected Hijack to return the mock connection")
	}

	if bufrw == nil {
		t.Error("Expected Hijack to return a non-nil bufio.ReadWriter")
	}
}

// Mock objects for testing Hijack
type mockResponseWriter struct {
	http.ResponseWriter
	conn net.Conn
	bufw *bufio.ReadWriter
}

func (m *mockResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return m.conn, m.bufw, nil
}

type mockConn struct {
	net.Conn
}

func (m *mockConn) Read(b []byte) (n int, err error) {
	return 0, nil
}

func (m *mockConn) Write(b []byte) (n int, err error) {
	return len(b), nil
}

func (m *mockConn) Close() error {
	return nil
}

func (m *mockConn) LocalAddr() net.Addr {
	return &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 8080}
}

func (m *mockConn) RemoteAddr() net.Addr {
	return &net.TCPAddr{IP: net.IPv4(192, 168, 0, 1), Port: 12345}
}

func TestResponsePusher(t *testing.T) {
	// Regular response writer (not a pusher)
	w := httptest.NewRecorder()
	response := NewResponse(w)

	pusher, ok := response.Pusher()
	if ok || pusher != nil {
		t.Error("Expected Pusher to return nil and false for non-pusher writer")
	}

	// Mock pusher for HTTP/2
	mockPusher := &mockPusher{}
	response = NewResponse(mockPusher)

	pusher, ok = response.Pusher()
	if !ok || pusher == nil {
		t.Error("Expected Pusher to return non-nil and true for pusher writer")
	}
}

// Mock object for testing Pusher
type mockPusher struct {
	http.ResponseWriter
	pushed bool
}

func (m *mockPusher) Push(target string, opts *http.PushOptions) error {
	m.pushed = true
	return nil
}

func (m *mockPusher) Header() http.Header {
	return make(http.Header)
}

func (m *mockPusher) Write(b []byte) (int, error) {
	return len(b), nil
}

func (m *mockPusher) WriteHeader(statusCode int) {
	// Do nothing
}
