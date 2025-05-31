package context

import (
	"bytes"
	gocontext "context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewContext(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	ctx := NewContext(w, req)

	if ctx == nil {
		t.Fatalf("Expected context to be created, got nil")
	}

	if ctx.Request() == nil {
		t.Error("Expected Request() to return non-nil value")
	}

	if ctx.Response() == nil {
		t.Error("Expected Response() to return non-nil value")
	}
}

func TestContextWithContext(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	ctx := NewContext(w, req)

	// Create a new Go context with a value
	type key string
	testKey := key("testKey")
	testValue := "testValue"

	// Get the underlying Go context
	goCtx := ctx.Context()

	// Create a new one with our value
	newGoCtx := gocontext.WithValue(goCtx, testKey, testValue)

	// Update the context
	updatedCtx := ctx.WithContext(newGoCtx)

	// Check if the value is in the updated context
	if updatedCtx.Context().Value(testKey) != testValue {
		t.Errorf("Expected %s, got %v", testValue, updatedCtx.Context().Value(testKey))
	}
}

func TestContextNextAbort(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	ctx := NewContext(w, req).(*forkContext)

	// We'll use a simple test to verify Abort works
	executed := false

	// First handler that will Abort
	ctx.handlers = []func(Context){
		func(c Context) {
			c.Set("step1", true)
			c.Abort()
		},
	}

	// This will execute the first handler
	ctx.Next()

	// The first handler should have set step1
	val1, exists1 := ctx.Get("step1")
	if !exists1 || val1 != true {
		t.Error("Expected step1 to be set to true")
	}

	// Context should be aborted
	if !ctx.IsAborted() {
		t.Error("Expected context to be aborted")
	}

	// Now we'll add a second handler
	ctx.handlers = append(ctx.handlers, func(c Context) {
		executed = true
	})

	// Reset the index but not the aborted state
	ctx.index = 0

	// Call Next again - since we're aborted, the second handler shouldn't execute
	ctx.Next()

	// Verify the second handler wasn't executed
	if executed {
		t.Error("Expected second handler not to be executed after abort")
	}
}

func TestContextSetGet(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	ctx := NewContext(w, req)

	// Test basic Set/Get
	ctx.Set("string", "test")
	val, ok := ctx.Get("string")
	if !ok || val != "test" {
		t.Errorf("Expected 'test', got %v", val)
	}

	// Test GetString
	ctx.Set("string", "test")
	if ctx.GetString("string") != "test" {
		t.Errorf("Expected 'test', got %v", ctx.GetString("string"))
	}

	// Test GetBool
	ctx.Set("bool", true)
	if !ctx.GetBool("bool") {
		t.Error("Expected true, got false")
	}

	// Test GetInt
	ctx.Set("int", 42)
	if ctx.GetInt("int") != 42 {
		t.Errorf("Expected 42, got %v", ctx.GetInt("int"))
	}

	// Test GetInt64
	ctx.Set("int64", int64(42))
	if ctx.GetInt64("int64") != 42 {
		t.Errorf("Expected 42, got %v", ctx.GetInt64("int64"))
	}

	// Test GetFloat64
	ctx.Set("float64", 3.14)
	if ctx.GetFloat64("float64") != 3.14 {
		t.Errorf("Expected 3.14, got %v", ctx.GetFloat64("float64"))
	}

	// Test GetTime
	now := time.Now()
	ctx.Set("time", now)
	if !ctx.GetTime("time").Equal(now) {
		t.Errorf("Expected %v, got %v", now, ctx.GetTime("time"))
	}

	// Test GetDuration
	dur := 5 * time.Second
	ctx.Set("duration", dur)
	if ctx.GetDuration("duration") != dur {
		t.Errorf("Expected %v, got %v", dur, ctx.GetDuration("duration"))
	}

	// Test GetStringSlice
	slice := []string{"a", "b", "c"}
	ctx.Set("slice", slice)
	result := ctx.GetStringSlice("slice")
	if len(result) != len(slice) {
		t.Errorf("Expected %v, got %v", slice, result)
	}
	for i, v := range slice {
		if result[i] != v {
			t.Errorf("Expected %s at index %d, got %s", v, i, result[i])
		}
	}

	// Test GetStringMap
	m := map[string]interface{}{"key": "value"}
	ctx.Set("map", m)
	mapResult := ctx.GetStringMap("map")
	if mapResult["key"] != "value" {
		t.Errorf("Expected %v, got %v", m, mapResult)
	}

	// Test GetStringMapString
	ms := map[string]string{"key": "value"}
	ctx.Set("mapString", ms)
	msResult := ctx.GetStringMapString("mapString")
	if msResult["key"] != "value" {
		t.Errorf("Expected %v, got %v", ms, msResult)
	}
}

func TestContextRequestMethods(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/test?q=search&page=1", bytes.NewBufferString("name=test&age=25"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx := NewContext(w, req)

	// Test Method
	if ctx.Method() != "POST" {
		t.Errorf("Expected POST, got %s", ctx.Method())
	}

	// Test Path
	if ctx.Path() != "/test" {
		t.Errorf("Expected /test, got %s", ctx.Path())
	}

	// Test RawPath
	if ctx.RawPath() != "" {
		t.Errorf("Expected empty RawPath, got %s", ctx.RawPath())
	}

	// Test Query
	if ctx.Query("q") != "search" {
		t.Errorf("Expected 'search', got %s", ctx.Query("q"))
	}

	// Test DefaultQuery
	if ctx.DefaultQuery("q", "default") != "search" {
		t.Errorf("Expected 'search', got %s", ctx.DefaultQuery("q", "default"))
	}

	if ctx.DefaultQuery("missing", "default") != "default" {
		t.Errorf("Expected 'default', got %s", ctx.DefaultQuery("missing", "default"))
	}

	// Test Form
	if ctx.Form("name") != "test" {
		t.Errorf("Expected 'test', got %s", ctx.Form("name"))
	}

	// Test DefaultForm
	if ctx.DefaultForm("name", "default") != "test" {
		t.Errorf("Expected 'test', got %s", ctx.DefaultForm("name", "default"))
	}

	if ctx.DefaultForm("missing", "default") != "default" {
		t.Errorf("Expected 'default', got %s", ctx.DefaultForm("missing", "default"))
	}
}

func TestContextResponseMethods(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	ctx := NewContext(w, req)

	// Test Status
	ctx.Status(http.StatusCreated)
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	// Test Header
	ctx.Header("X-Test", "value")
	if w.Header().Get("X-Test") != "value" {
		t.Errorf("Expected header 'value', got %s", w.Header().Get("X-Test"))
	}

	// Test GetHeader
	req.Header.Set("X-Request-Test", "test-value")
	if ctx.GetHeader("X-Request-Test") != "test-value" {
		t.Errorf("Expected header 'test-value', got %s", ctx.GetHeader("X-Request-Test"))
	}
}

func TestContextContentNegotiation(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(`{"name":"test","age":25}`))
	req.Header.Set("Content-Type", "application/json")
	ctx := NewContext(w, req)

	type TestUser struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	var user TestUser
	err := ctx.BindJSON(&user)
	if err != nil {
		t.Errorf("Failed to bind JSON: %v", err)
	}

	if user.Name != "test" || user.Age != 25 {
		t.Errorf("Expected name='test' and age=25, got name='%s' and age=%d", user.Name, user.Age)
	}
}

func TestContextResponding(t *testing.T) {
	tests := []struct {
		name        string
		fn          func(ctx Context)
		contentType string
		statusCode  int
	}{
		{
			name: "JSON",
			fn: func(ctx Context) {
				ctx.JSON(http.StatusOK, map[string]interface{}{"message": "ok"})
			},
			contentType: "application/json",
			statusCode:  http.StatusOK,
		},
		{
			name: "XML",
			fn: func(ctx Context) {
				type xmlTest struct {
					Message string `xml:"message"`
				}
				ctx.XML(http.StatusAccepted, xmlTest{Message: "ok"})
			},
			contentType: "application/xml",
			statusCode:  http.StatusAccepted,
		},
		{
			name: "String",
			fn: func(ctx Context) {
				ctx.String(http.StatusOK, "Hello %s", "World")
			},
			contentType: "text/plain",
			statusCode:  http.StatusOK,
		},
		{
			name: "HTML",
			fn: func(ctx Context) {
				ctx.HTML(http.StatusOK, "<h1>Hello World</h1>")
			},
			contentType: "text/html",
			statusCode:  http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/test", nil)
			ctx := NewContext(w, req)

			tt.fn(ctx)

			if w.Code != tt.statusCode {
				t.Errorf("Expected status %d, got %d", tt.statusCode, w.Code)
			}

			contentType := w.Header().Get("Content-Type")
			if contentType == "" || !bytes.Contains([]byte(contentType), []byte(tt.contentType)) {
				t.Errorf("Expected content type to contain %s, got %s", tt.contentType, contentType)
			}
		})
	}
}

func TestContextFile(t *testing.T) {
	// Create a temporary file for testing
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	// Write some content to the file
	content := []byte("test file content")
	if _, err := tmpfile.Write(content); err != nil {
		tmpfile.Close()
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	ctx := NewContext(w, req)

	ctx.File(tmpfile.Name())

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	if !bytes.Equal(w.Body.Bytes(), content) {
		t.Errorf("Expected body %s, got %s", content, w.Body.Bytes())
	}
}

func TestContextMultipartForm(t *testing.T) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Add form field
	if err := w.WriteField("name", "test"); err != nil {
		t.Fatalf("Failed to write field: %v", err)
	}

	// Add file
	fileContent := []byte("file content")
	fileWriter, err := w.CreateFormFile("file", "test.txt")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	if _, err := fileWriter.Write(fileContent); err != nil {
		t.Fatalf("Failed to write file content: %v", err)
	}

	// Close writer
	if err := w.Close(); err != nil {
		t.Fatalf("Failed to close writer: %v", err)
	}

	// Create request
	req := httptest.NewRequest("POST", "/test", &b)
	req.Header.Set("Content-Type", w.FormDataContentType())

	rec := httptest.NewRecorder()
	ctx := NewContext(rec, req)

	// Test FormFile
	file, err := ctx.FormFile("file")
	if err != nil {
		t.Fatalf("Failed to get form file: %v", err)
	}

	if file.Filename != "test.txt" {
		t.Errorf("Expected filename 'test.txt', got '%s'", file.Filename)
	}

	// Test SaveUploadedFile
	tmpDir, err := os.MkdirTemp("", "upload")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	savePath := filepath.Join(tmpDir, "saved.txt")
	if err := ctx.SaveUploadedFile(file, savePath); err != nil {
		t.Fatalf("Failed to save uploaded file: %v", err)
	}

	// Verify saved file
	savedContent, err := os.ReadFile(savePath)
	if err != nil {
		t.Fatalf("Failed to read saved file: %v", err)
	}

	if !bytes.Equal(savedContent, fileContent) {
		t.Errorf("Expected file content %s, got %s", fileContent, savedContent)
	}
}

// TestParamArrayAndMap checks the ParamArray and ParamMap methods
func TestParamArrayAndMap(t *testing.T) {
	// Create a request with no params
	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	ctx := NewContext(w, req)

	// Test with no params
	if len(ctx.ParamMap()) != 0 {
		t.Error("Expected empty ParamMap for new context")
	}
	if ctx.ParamArray("id") != nil {
		t.Error("Expected nil ParamArray for non-existent param")
	}

	// Set params
	ctx.Set("param:id", "123")
	ctx.Set("param:name", "test")

	// Test ParamMap
	params := ctx.ParamMap()
	if len(params) != 2 {
		t.Errorf("Expected 2 params, got %d", len(params))
	}
	if params["id"] != "123" {
		t.Errorf("Expected id=123, got %s", params["id"])
	}
	if params["name"] != "test" {
		t.Errorf("Expected name=test, got %s", params["name"])
	}

	// Test ParamArray
	idArray := ctx.ParamArray("id")
	if len(idArray) != 1 {
		t.Errorf("Expected 1 item in idArray, got %d", len(idArray))
	}
	if idArray[0] != "123" {
		t.Errorf("Expected idArray[0]=123, got %s", idArray[0])
	}

	// Test ParamArray with non-existent param
	emptyArray := ctx.ParamArray("non-existent")
	if emptyArray != nil {
		t.Error("Expected nil ParamArray for non-existent param")
	}
}
