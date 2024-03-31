package logging

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func HandleRequestBody(req *http.Request) string {
	var requestBodyBytes []byte
	if req.Body == nil {
		return ""
	}

	requestBodyBytes, _ = io.ReadAll(req.Body)
	req.Body = io.NopCloser(bytes.NewBuffer(requestBodyBytes))
	return string(requestBodyBytes)
}

func HandleResponseBody(rw gin.ResponseWriter) *BodyLogWriter {
	return &BodyLogWriter{Body: bytes.NewBufferString(""), ResponseWriter: rw}
}

func FormatRequestAndResponse(rw gin.ResponseWriter, req *http.Request, responseBody string, requestId string, requestBody string) string {
	if req.URL.String() == "/metrics" {
		return ""
	}

	return fmt.Sprintf("[Request ID: %s], Status: [%d], Method: [%s], Url: %s Request Body: %s, Response Body: %s",
		requestId, rw.Status(), req.Method, req.URL.String(), requestBody, responseBody)
}
