package middles

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rbee3u/gohelp/conflog"
)

func LoggingRequest(l *conflog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		entry := l.WithContext(c.Request.Context())
		entry = entry.WithField("method", c.Request.Method)
		entry = entry.WithField("uri", c.Request.URL.RequestURI())
		entry = entry.WithField("headers", composeHeaders(c.Request.Header))
		entry = entry.WithField("body", requestBody(c))
		entry.Infof("incoming http request")
	}
}

func LoggingResponse(l *conflog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		rw := &responseWriter{Body: new(bytes.Buffer), ResponseWriter: c.Writer}
		c.Writer = rw
		c.Next()
		statusCode := c.Writer.Status()
		statusText := http.StatusText(statusCode)
		entry := l.WithContext(c.Request.Context())
		entry = entry.WithField("status", fmt.Sprintf("%v %s", statusCode, statusText))
		entry = entry.WithField("headers", composeHeaders(c.Writer.Header()))
		entry = entry.WithField("body", rw.Body.String())
		entry.Infof("outgoing http response")
	}
}

func composeHeaders(headers http.Header) string {
	pairs := make([]string, 0, len(headers))
	for key, values := range headers {
		pairs = append(pairs, fmt.Sprintf("%s: %s", key, strings.Join(values, ", ")))
	}

	sort.Strings(pairs)

	return strings.Join(pairs, "; ")
}

func requestBody(c *gin.Context) string {
	if c.Request.Body == nil || c.Request.Body == http.NoBody {
		return ""
	}

	body, _ := ioutil.ReadAll(c.Request.Body)
	_ = c.Request.Body.Close()
	c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))

	return string(body)
}

type responseWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func (w responseWriter) Write(body []byte) (int, error) {
	w.Body.Write(body)

	n, err := w.ResponseWriter.Write(body)
	if err != nil {
		return 0, fmt.Errorf("failed to write: %w", err)
	}

	return n, nil
}
