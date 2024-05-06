package http

import (
	"bytes"
	"io"
	"item-service/config"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

type GinMiddleware struct {
}

func (e *GinMiddleware) JSONMiddleware(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Next()
}

func (e *GinMiddleware) CORS(c *gin.Context) {
	cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	})
}

func (e *GinMiddleware) Logger(zapLogger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		fields := []zapcore.Field{}

		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		var body []byte
		var buf bytes.Buffer
		tee := io.TeeReader(c.Request.Body, &buf)
		body, _ = io.ReadAll(tee)

		c.Request.Body = io.NopCloser(&buf)
		fields = append(fields, zap.String("Body", string(body)))
		fields = append(fields, zap.String("Path", path))

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		comment := c.Errors.ByType(gin.ErrorTypePrivate).String()

		fields = append(fields, zap.String("Raw", raw))
		fields = append(fields, zap.String("ClientIP", clientIP))
		fields = append(fields, zap.String("Method", method))
		fields = append(fields, zap.Int("StatusCode", statusCode))
		fields = append(fields, zap.String("Comment", comment))
		fields = append(fields, zap.Duration("Latency", latency))
		fields = append(fields, zap.String("Response", blw.body.String()))

		if comment != "" {
			zapLogger.Error("Request", fields...)

		} else {
			zapLogger.Info("Request", fields...)
		}
	}
}

func (e *GinMiddleware) JWT(config config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
			auth := c.Request.Header.Get("Authorization")
			if auth == "" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "Unauthorized",
				})
				return
			}
			token := strings.Split(auth, " ")[1]
			if token == "" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "Unauthorized",
				})
				return
			}
		c.Next()
	}
}

func (e *GinMiddleware) ErrorHandler(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Next()
		if len(c.Errors) > 0 {
			//logger.Error("Error", zap.String("Error", c.Errors.String()))
			//var message string
			//statusCode := http.StatusInternalServerError
			//err := c.Errors.Last()
			//
			//if err.IsType(gin.ErrorTypePrivate) {
			//	message = utils.ErrInternalServerError.Error()
			//} else {
			//	message = err.Error()
			//	statusCode = err.Meta.(int)
			//}
			//c.JSON(statusCode, dto.ResponseError{
			//	Message: message,
			//	Errors:  c.Errors.Errors(),
			//})
			//return
		}
	}
}

func InitMiddleware() *GinMiddleware {
	return &GinMiddleware{}
}
