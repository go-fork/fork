package context

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewRequest(t *testing.T) {
	req := httptest.NewRequest("GET", "/test?q=search", nil)
	request := NewRequest(req)

	if request == nil {
		t.Fatal("Expected request to be created, got nil")
	}
}

func TestRequestMethod(t *testing.T) {
	req := httptest.NewRequest("POST", "/test", nil)
	request := NewRequest(req)

	if request.Method() != "POST" {
		t.Errorf("Expected POST, got %s", request.Method())
	}
}

func TestRequestURL(t *testing.T) {
	req := httptest.NewRequest("GET", "/test?q=search&page=1", nil)
	request := NewRequest(req)

	url := request.URL()
	if url.Path != "/test" {
		t.Errorf("Expected path /test, got %s", url.Path)
	}

	if url.RawQuery != "q=search&page=1" {
		t.Errorf("Expected query q=search&page=1, got %s", url.RawQuery)
	}
}

func TestRequestHost(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/test", nil)
	request := NewRequest(req)

	if request.Host() != "example.com" {
		t.Errorf("Expected host example.com, got %s", request.Host())
	}
}

func TestRequestProtocol(t *testing.T) {
	testCases := []struct {
		name           string
		protoMajor     int
		proto          string
		headers        map[string]string
		expectedResult string
	}{
		{
			name:           "HTTP/1.1 Protocol",
			protoMajor:     1,
			proto:          "HTTP/1.1",
			headers:        map[string]string{},
			expectedResult: "HTTP/1.1",
		},
		{
			name:           "HTTP/2 Protocol via ProtoMajor",
			protoMajor:     2,
			proto:          "HTTP/2.0",
			headers:        map[string]string{},
			expectedResult: "HTTP/2",
		},
		{
			name:           "HTTP/2 Protocol via Proto string",
			protoMajor:     1, // Wrong ProtoMajor but Proto string should overrule
			proto:          "HTTP/2.0",
			headers:        map[string]string{},
			expectedResult: "HTTP/2",
		},
		{
			name:           "HTTP/3 Protocol via ProtoMajor",
			protoMajor:     3,
			proto:          "HTTP/1.1", // Wrong Proto but ProtoMajor should overrule
			headers:        map[string]string{},
			expectedResult: "HTTP/3",
		},
		{
			name:           "HTTP/3 Protocol via Proto string",
			protoMajor:     1,
			proto:          "HTTP/3.0",
			headers:        map[string]string{},
			expectedResult: "HTTP/3",
		},
		{
			name:           "HTTP/3 Protocol via Alt-Used header",
			protoMajor:     1,
			proto:          "HTTP/1.1",
			headers:        map[string]string{"Alt-Used": "example.com:443"},
			expectedResult: "HTTP/3",
		},
		{
			name:           "HTTP/2 Cleartext Upgrade",
			protoMajor:     1,
			proto:          "HTTP/1.1",
			headers:        map[string]string{"Upgrade": "h2c"},
			expectedResult: "h2c-upgrade",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://example.com/test", nil)
			req.ProtoMajor = tc.protoMajor
			req.Proto = tc.proto

			// Add headers
			for key, value := range tc.headers {
				req.Header.Set(key, value)
			}

			request := NewRequest(req)
			result := request.Protocol()

			if result != tc.expectedResult {
				t.Errorf("Expected protocol %s, got %s", tc.expectedResult, result)
			}
		})
	}
}

func TestRequestRemoteAddr(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:1234"
	request := NewRequest(req)

	if request.RemoteAddr() != "192.168.1.1:1234" {
		t.Errorf("Expected remote addr 192.168.1.1:1234, got %s", request.RemoteAddr())
	}
}

func TestRequestHeader(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Test", "value")
	request := NewRequest(req)

	header := request.Header()
	if header.Get("X-Test") != "value" {
		t.Errorf("Expected header X-Test=value, got %s", header.Get("X-Test"))
	}
}

func TestRequestCookie(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{Name: "test", Value: "value"})
	request := NewRequest(req)

	cookie, err := request.Cookie("test")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if cookie.Value != "value" {
		t.Errorf("Expected cookie value='value', got '%s'", cookie.Value)
	}
}

func TestRequestCookies(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{Name: "test1", Value: "value1"})
	req.AddCookie(&http.Cookie{Name: "test2", Value: "value2"})
	request := NewRequest(req)

	cookies := request.Cookies()
	if len(cookies) != 2 {
		t.Fatalf("Expected 2 cookies, got %d", len(cookies))
	}

	if cookies[0].Name != "test1" || cookies[0].Value != "value1" {
		t.Errorf("Expected cookie test1=value1, got %s=%s", cookies[0].Name, cookies[0].Value)
	}

	if cookies[1].Name != "test2" || cookies[1].Value != "value2" {
		t.Errorf("Expected cookie test2=value2, got %s=%s", cookies[1].Name, cookies[1].Value)
	}
}

func TestRequestFormValue(t *testing.T) {
	formData := "name=test&age=25"
	req := httptest.NewRequest("POST", "/test", strings.NewReader(formData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request := NewRequest(req)

	if request.FormValue("name") != "test" {
		t.Errorf("Expected form value name=test, got %s", request.FormValue("name"))
	}

	if request.FormValue("age") != "25" {
		t.Errorf("Expected form value age=25, got %s", request.FormValue("age"))
	}
}

func TestRequestMultipartForm(t *testing.T) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Add form field
	if err := w.WriteField("name", "test"); err != nil {
		t.Fatalf("Failed to write field: %v", err)
	}

	// Close writer
	if err := w.Close(); err != nil {
		t.Fatalf("Failed to close writer: %v", err)
	}

	// Create request
	req := httptest.NewRequest("POST", "/test", &b)
	req.Header.Set("Content-Type", w.FormDataContentType())
	request := NewRequest(req)

	form, err := request.MultipartForm()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if name := form.Value["name"][0]; name != "test" {
		t.Errorf("Expected name=test, got %s", name)
	}
}

func TestRequestFormFile(t *testing.T) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Add file
	fileWriter, err := w.CreateFormFile("file", "test.txt")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	if _, err := fileWriter.Write([]byte("file content")); err != nil {
		t.Fatalf("Failed to write file content: %v", err)
	}

	// Close writer
	if err := w.Close(); err != nil {
		t.Fatalf("Failed to close writer: %v", err)
	}

	// Create request
	req := httptest.NewRequest("POST", "/test", &b)
	req.Header.Set("Content-Type", w.FormDataContentType())
	request := NewRequest(req)

	file, err := request.FormFile("file")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if file.Filename != "test.txt" {
		t.Errorf("Expected filename test.txt, got %s", file.Filename)
	}
}

func TestRequestRequest(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	request := NewRequest(req)

	if request.Request() != req {
		t.Error("Expected Request() to return the original request")
	}
}

func TestRequestBody(t *testing.T) {
	// Read body
	body := "test body"
	req := httptest.NewRequest("POST", "/test", strings.NewReader(body))
	request := NewRequest(req)

	// Get the body reader
	bodyReader := request.Body()
	data, err := io.ReadAll(bodyReader)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if string(data) != body {
		t.Errorf("Expected body %s, got %s", body, string(data))
	}
}

func TestRequestUserAgent(t *testing.T) {
	userAgent := "Mozilla/5.0 Test User Agent"
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("User-Agent", userAgent)
	request := NewRequest(req)

	if request.UserAgent() != userAgent {
		t.Errorf("Expected User-Agent %s, got %s", userAgent, request.UserAgent())
	}
}

func TestRequestReferer(t *testing.T) {
	referer := "http://example.com/referer"
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Referer", referer)
	request := NewRequest(req)

	if request.Referer() != referer {
		t.Errorf("Expected Referer %s, got %s", referer, request.Referer())
	}
}
