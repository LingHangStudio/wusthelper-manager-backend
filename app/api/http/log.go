package http

import (
	"github.com/gin-gonic/gin"
	"wusthelper-manager-go/app/model"
	"wusthelper-manager-go/app/service"
	"wusthelper-manager-go/common"
	"wusthelper-manager-go/library/ecode"
)

type LogResp struct {
	Logid      int64  `json:"logid"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Version    string `json:"version"`
	Status     int8   `json:"status"`
	Platform   string `json:"platform"`
	UpdateTime string `json:"updateTime"`
}

func getLogList(c *gin.Context) {
	req := new(PlatformPaginationReq)
	if err := c.ShouldBind(req); err != nil {
		responseEcode(c, ecode.ParamWrong)
		return
	}

	logInfoList, total, err := srv.GetLogList(common.Pagination{Page: req.Page, PageSize: req.Size}, req.Platform)

	if err != nil {
		responseEcode(c, err)
		return
	}

	resultList := make([]LogResp, len(*logInfoList))
	for i, logInfo := range *logInfoList {
		resultList[i] = LogResp{
			Logid:      logInfo.ID,
			Title:      *logInfo.Title,
			Content:    *logInfo.Content,
			Version:    *logInfo.VersionText,
			Status:     _internalLogStatus2ApiDefineStatus(*logInfo.Status),
			Platform:   *logInfo.Platform,
			UpdateTime: logInfo.UpdateTime.Format(_defaultDateFormat),
		}
	}

	responseData(c, map[string]any{
		"logs": resultList,
		"num":  total,
	})
}

type PublishedLogResp struct {
	Logid      int64  `json:"logid"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Version    string `json:"version"`
	CreateTime string `json:"createTime"`
}

func getPublishedLogList(c *gin.Context) {
	platform := getPlatform(c)
	var resultList *[]model.Log
	var err error
	if platform == "" {
		resultList, err = srv.GetPublishedLog()
	} else {
		resultList, err = srv.GetPublishedLog(platform)
	}

	if err != nil {
		responseEcode(c, err)
		return
	}

	publishedLogList := make([]PublishedLogResp, len(*resultList))
	for i, logInfo := range *resultList {
		publishedLogList[i] = PublishedLogResp{
			Logid:   logInfo.ID,
			Title:   *logInfo.Title,
			Content: *logInfo.Content,
			Version: *logInfo.VersionText,
		}
	}

	responseData(c, publishedLogList)
}

func _internalLogStatus2ApiDefineStatus(internalStatus int8) int8 {
	switch internalStatus {
	case model.NormalStatus:
		return 0
	case model.LogPublishedStatus:
		return 1
	default:
		return 0
	}
}

type LogPublishReq struct {
	Logid []int64 `json:"logid" binging:"required"`
}

func publishLog(c *gin.Context) {
	req := new(LogPublishReq)
	if err := c.ShouldBind(req); err != nil {
		responseEcode(c, ecode.ParamWrong)
		return
	}

	err := srv.PublishLogBatch(req.Logid...)
	if err != nil {
		responseEcode(c, err)
		return
	}

	responseData(c, nil)
}

type LogAddReq struct {
	Title    string   `json:"title" form:"title"`
	Content  string   `json:"content" form:"content"`
	Version  string   `json:"version" form:"version" binding:"required"`
	Platform []string `json:"platform" form:"platform" binding:"required"`
}

func addLog(c *gin.Context) {
	req := new(LogAddReq)
	if err := c.ShouldBind(req); err != nil {
		responseEcode(c, ecode.ParamWrong)
		return
	}

	logInfo := service.LogAddParam{
		Title:       req.Title,
		Content:     req.Content,
		VersionText: req.Version,
		Platform:    req.Platform,
	}

	err := srv.AddLog(&logInfo)
	if err != nil {
		responseEcode(c, err)
		return
	}

	responseData(c, nil)
}

type LogModifyReq struct {
	Logid    int64     `json:"logid" form:"logid" binding:"required"`
	Title    *string   `json:"title" form:"title"`
	Content  *string   `json:"content" form:"content"`
	Version  *string   `json:"version" form:"version"`
	Platform *[]string `json:"platform" form:"platform"`
	Status   *int8     `json:"status" form:"status"`
}

func modifyLog(c *gin.Context) {
	req := new(LogModifyReq)
	if err := c.ShouldBind(req); err != nil {
		responseEcode(c, ecode.ParamWrong)
		return
	}

	logInfo := service.LogModifyParam{
		Id:          req.Logid,
		Title:       req.Title,
		Content:     req.Content,
		VersionText: req.Version,
		Platform:    req.Platform,
		Status:      req.Status,
	}

	err := srv.ModifyLog(&logInfo)
	if err != nil {
		responseEcode(c, err)
		return
	}

	responseData(c, nil)
}

type LogDeleteReq struct {
	Logid int64 `json:"logid" form:"logid" binding:"required"`
}

func deleteLog(c *gin.Context) {
	req := new(LogDeleteReq)
	if err := c.ShouldBind(req); err != nil {
		responseEcode(c, ecode.ParamWrong)
		return
	}

	err := srv.DeleteLog(req.Logid)
	if err != nil {
		responseEcode(c, err)
		return
	}

	responseData(c, nil)
}
