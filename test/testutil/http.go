package testutil

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

// TestRequest represents an HTTP test request
type TestRequest struct {
	Method  string
	Path    string
	Body    interface{}
	Headers map[string]string
}

// TestResponse represents an HTTP test response
type TestResponse struct {
	Code int
	Body interface{}
}

// ExecuteRequest executes a test request and returns the response
func ExecuteRequest(router *gin.Engine, req TestRequest) *httptest.ResponseRecorder {
	// Convert body to JSON if it's not nil
	var bodyReader *bytes.Buffer
	if req.Body != nil {
		jsonBody, err := json.Marshal(req.Body)
		if err != nil {
			panic(err)
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	// Create HTTP request
	httpReq := httptest.NewRequest(req.Method, req.Path, bodyReader)
	httpReq.Header.Set("Content-Type", "application/json")

	// Add custom headers
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	// Create response recorder
	w := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(w, httpReq)

	return w
}

// ParseResponse parses the response body into the given interface
func ParseResponse(w *httptest.ResponseRecorder, v interface{}) error {
	return json.NewDecoder(w.Body).Decode(v)
}

// CreateTestRouter creates a new Gin router in test mode
func CreateTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

// AddTestRoute adds a test route to the router
func AddTestRoute(router *gin.Engine, method, path string, handler gin.HandlerFunc) {
	router.Handle(method, path, handler)
}
