package gin

import (
	"bytes"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/clin211/miniblog-v2/internal/pkg/contextx"
	"github.com/clin211/miniblog-v2/internal/pkg/log"
	"github.com/clin211/miniblog-v2/pkg/token"
	"github.com/gin-gonic/gin"
)

const (
	// MultipartFormData 文件上传的 Content-Type
	MultipartFormData = "multipart/form-data"
	// MaxResponseLogSize 响应日志记录的最大大小 (10KB)
	MaxResponseLogSize = 10 * 1024
	// ResponseTooLargeMessage 响应过大时的提示信息
	ResponseTooLargeMessage = "Response data size is too Large to log"
)

// bufferPool 用于复用 bytes.Buffer，减少 GC 压力
var bufferPool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

// getBuffer 从池中获取一个 buffer
func getBuffer() *bytes.Buffer {
	return bufferPool.Get().(*bytes.Buffer)
}

// putBuffer 将 buffer 重置后放回池中
func putBuffer(buf *bytes.Buffer) {
	buf.Reset()
	bufferPool.Put(buf)
}

// bodyLogWriter 包装 gin.ResponseWriter，用于拦截响应数据
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write 实现 io.Writer 接口，拦截写入的响应数据
// 让gin写响应的时候先写到 bodyLogWriter 再写 gin.ResponseWriter，
// 这样利用中间件里输出访问日志时就能拿到响应了
// https://stackoverflow.com/questions/38501325/how-to-log-response-body-in-gin
func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// AccessLogger 返回一个访问日志中间件
func AccessLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//保存body
		var reqBody []byte
		contentType := ctx.GetHeader("Content-Type")
		// multipart/form-data 文件上传请求, 不在日志里记录body
		if !strings.Contains(contentType, MultipartFormData) {
			reqBody, _ = io.ReadAll(ctx.Request.Body)
			ctx.Request.Body = io.NopCloser(bytes.NewReader(reqBody))
		}
		// 从请求头中获取 `x-request-id`
		requestID := contextx.RequestID(ctx.Request.Context())
		start := time.Now()

		// 从池中获取 buffer
		responseBuffer := getBuffer()
		blw := &bodyLogWriter{body: responseBuffer, ResponseWriter: ctx.Writer}
		ctx.Writer = blw

		// 记录访问开始
		accessLog(ctx, "access_start", time.Since(start), requestID, reqBody, nil)

		defer func() {
			var responseLogging string
			if ctx.Writer.Size() > MaxResponseLogSize { // 响应大于10KB 不记录
				responseLogging = ResponseTooLargeMessage
			} else {
				responseLogging = blw.body.String()
			}

			// 记录访问结束
			accessLog(ctx, "access_end", time.Since(start), requestID, reqBody, responseLogging)

			// 将 buffer 放回池中
			putBuffer(responseBuffer)
		}()
		ctx.Next()
	}
}

// accessLog 记录访问日志
func accessLog(c *gin.Context, accessType string, dur time.Duration, requestID string, body []byte, dataOut interface{}) {
	req := c.Request
	bodyStr := string(body)
	query := req.URL.RawQuery
	path := req.URL.Path

	userID, _ := token.ParseRequest(c)
	log.Infow("AccessLog",
		"type", accessType,
		"ip", c.ClientIP(),
		"userID", userID,
		"requestID", requestID,
		"method", req.Method,
		"path", path,
		"query", query,
		"body", bodyStr,
		"output", dataOut,
		"time(ms)", int64(dur/time.Millisecond),
		"status", c.Writer.Status(),
	)
}
