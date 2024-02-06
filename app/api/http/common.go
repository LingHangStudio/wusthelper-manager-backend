package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/url"
	"wusthelper-manager-go/library/ecode"
	"wusthelper-manager-go/library/log"
)

const (
	_defaultDateTimeFormat = "2006-01-02 15-04-05"
	_defaultDateFormat     = "2006-01-02"
	_defaultTimeFormat     = "15-04-05"
)

const (
	_platformAndroid = "android"
	_platformIos     = "ios"
	_platformMp      = "mp"
)

type PlatformPaginationReq struct {
	Page     int    `json:"page,default=1" form:"page,default=1" query:"page,default=1"`
	Size     int    `json:"size,default=10" form:"size,default=10" query:"size,default=10"`
	Platform string `json:"platform" form:"platform" query:"platform"`
}

func _convertBoolStr2Bool(str string) bool {
	switch str {
	case "true":
		return true
	default:
		return false
	}
}

func _convertBoolStr2Int(str string) int {
	switch str {
	case "true":
		return 1
	default:
		return 0
	}
}

func _getPicUrl(imgId string) string {
	if imgId == "" {
		return ""
	}

	fileStorageOption := config.Server.FileStorageOption
	basePath := fileStorageOption.ResourceStorageOption.PicStorageBasePath
	domain := fileStorageOption.AliyunOssOption.BucketBindDomain
	u := url.URL{
		Scheme: "https",
		Host:   domain,
		Path:   basePath,
	}

	return u.JoinPath(fmt.Sprintf("%s.jpg", imgId)).String()
}

func _getFileUrl(fileKey string) string {
	if fileKey == "" {
		return ""
	}

	fileStorageOption := config.Server.FileStorageOption
	basePath := fileStorageOption.ResourceStorageOption.VersionFileStorageBasePath
	domain := fileStorageOption.AliyunOssOption.BucketBindDomain
	u := url.URL{
		Scheme: "https",
		Host:   domain,
		Path:   basePath,
	}

	return u.JoinPath(fileKey).String()
}

func getPlatform(c *gin.Context) string {
	return c.GetHeader("Platform")
}

func getUid(c *gin.Context) (uint64, error) {
	// 这里的值是在token校验（auth.UserTokenCheck）的时候设置的
	_oid, ok := c.Get("uid")
	if !ok {
		log.Error("获取uid参数失败, uid为空")
		return 0, ecode.ParamWrong
	}

	oid, ok := _oid.(uint64)
	if !ok {
		log.Error("获取uid参数失败, uid转换失败")
		return 0, ecode.ParamWrong
	}

	return oid, nil
}

func responseEcode(c *gin.Context, code error) {
	switch errCode := code.(type) {
	case ecode.Codes:
		errors.As(code, &errCode)
		respCode, msg := toResponseCode(errCode)
		responseWithCode(c, respCode, msg, nil)
	default:
		responseWithCode(c, ecode.InternalError.Code(), code.Error(), nil)
	}
}

func responseWithCode(c *gin.Context, code int, msg string, data any) {
	resp := gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	}

	c.JSON(200, resp)
	c.Abort()
}

func responseData(c *gin.Context, data any) {
	resp := gin.H{
		"code": ecode.OK.Code(),
		"msg":  "ok",
		"data": data,
	}

	c.JSON(200, resp)
}

func toResponseCode(code ecode.Codes) (respCode int, msg string) {
	return code.Code(), code.Message()
}
