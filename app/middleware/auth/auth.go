package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wusthelper-manager-go/app/conf"
	"wusthelper-manager-go/library/ecode"
	_token "wusthelper-manager-go/library/token"
)

var (
	jwt *_token.Token
	dev bool
)

func Init(c *conf.Config) {
	jwt = _token.New(c.Server.TokenSecret, c.Server.TokenTimeout)
	dev = c.Server.Env == conf.DevEnv
}

func AdminUserTokenCheck(c *gin.Context) {
	token := c.GetHeader("Token")
	if token == "" {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"code": ecode.TokenInvalid.Code(),
			"msg":  ecode.TokenInvalid.Message(),
		})
		return
	}

	claims, valid := jwt.GetClaimVerify(token)
	if (!dev && !valid) || claims == nil {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"code": ecode.TokenInvalid.Code(),
			"msg":  ecode.TokenInvalid.Message(),
		})
		return
	}

	oid := (*claims)["uid"]
	c.Set("uid", oid)
	c.Next()
}
