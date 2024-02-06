package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"runtime/debug"
	"wusthelper-manager-go/library/ecode"
	"wusthelper-manager-go/library/log"
)

func GlobalPanicRecover(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("捕获到的错误",
				zap.String("req_url", c.Request.URL.String()),
				zap.Any("error", err),
				zap.ByteString("stack", debug.Stack()),
			)

			debug.PrintStack()
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": ecode.InternalError,
				"msg":  "未知错误",
			})
			c.Abort()
		}
	}()

	c.Next()
}
