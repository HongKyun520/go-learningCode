package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

type LogMiddlewareBuilder struct {
	logFn         func(ctx *gin.Context, l AccessLog)
	allowReqBody  bool
	allowRespBody bool
}

type AccessLog struct {
	Path         string        `json:"path"`
	Method       string        `json:"method"`
	RequestBody  string        `json:"request_body"`
	ResponseBody string        `json:"response_body"`
	Duration     time.Duration `json:"duration"`
	Status       int           `json:"status"`
}

type responseWriter struct {
	gin.ResponseWriter
	al *AccessLog
}

func NewLogMiddlewareBuilder(logFn func(ctx *gin.Context, l AccessLog)) *LogMiddlewareBuilder {
	return &LogMiddlewareBuilder{
		logFn: logFn,
	}
}

func (l *LogMiddlewareBuilder) AllowRespBody() *LogMiddlewareBuilder {
	l.allowRespBody = true
	return l
}

func (l *LogMiddlewareBuilder) AllowReqBody() *LogMiddlewareBuilder {
	l.allowReqBody = true
	return l
}

func (r *responseWriter) Write(p []byte) (int, error) {
	r.al.ResponseBody = string(p)
	return r.ResponseWriter.Write(p)
}

func (r *responseWriter) WriteHeader(statusCode int) {
	r.al.Status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *responseWriter) WriteString(s string) (int, error) {
	r.al.ResponseBody = s
	return r.ResponseWriter.WriteString(s)
}

// Build 构建日志中间件
//
// 实现功能:
// 1. 记录请求路径,限制长度为1024字符
// 2. 记录请求方法
// 3. 如果allowReqBody为true,记录请求体内容,限制长度为1024字符
// 4. 记录请求处理时长
// 5. 如果allowRespBody为true,记录响应体内容
//
// 实现步骤:
// 1. 获取并处理请求路径和方法
// 2. 创建AccessLog实例记录请求信息
// 3. 包装ResponseWriter以捕获响应信息
// 4. 如果允许记录请求体,读取并记录请求体内容
// 5. 记录开始时间
// 6. defer函数中记录耗时和响应信息
// 7. 调用后续处理器
//
// 参数:
//   - 无
//
// 返回值:
//   - gin.HandlerFunc: gin中间件处理函数
func (l *LogMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		path := ctx.Request.URL.Path
		if len(path) > 1024 {
			path = path[:1024]
		}

		method := ctx.Request.Method
		al := AccessLog{
			Path:   path,
			Method: method,
		}

		if l.allowRespBody {
			ctx.Writer = &responseWriter{
				ResponseWriter: ctx.Writer,
				al:             &al,
			}
		}

		// 如果允许记录请求体
		if l.allowReqBody {
			// 读取原始请求体数据
			bodyBytes, _ := ctx.GetRawData()
			// 将请求体转换为字符串并记录到AccessLog中
			if len(bodyBytes) > 1024 {
				al.RequestBody = string(bodyBytes[:1024])
			} else {
				al.RequestBody = string(bodyBytes)
			}
			// 由于GetRawData会消耗掉请求体,需要重新设置Request.Body
			// 使用NopCloser包装bytes.Buffer,重新构造可读取的请求体
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		start := time.Now()

		// 记录日志
		defer func() {
			al.Duration = time.Since(start)
			l.logFn(ctx, al)
		}()

		// 如果允许记录响应体
		// 1. 调用Next()执行后续的处理器
		// 2. 从响应写入器中获取响应体内容
		// 3. 将响应体记录到AccessLog中
		// 注意:这里存在问题,因为gin.ResponseWriter没有Body字段
		// TODO: 需要使用ResponseWriter的包装器来捕获响应体
		ctx.Next()
	}
}
