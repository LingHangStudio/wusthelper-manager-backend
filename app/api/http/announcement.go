package http

import (
	"github.com/gin-gonic/gin"
	"wusthelper-manager-go/app/model"
	"wusthelper-manager-go/app/service"
	"wusthelper-manager-go/common"
	"wusthelper-manager-go/library/ecode"
)

type PublishedAnnouncementResp struct {
	Id         int64   `json:"newsid"`
	Title      *string `json:"title"`
	Content    *string `json:"content"`
	Obj        *string `json:"obj"`
	UpdateTime *string `json:"updateTime"`
}

func getPublishedAnnouncement(c *gin.Context) {
	platform := c.GetHeader("Platform")
	announcements, err := srv.GetPublishedAnnouncement(platform)
	if err != nil {
		responseEcode(c, err)
		return
	}

	resp := make([]PublishedAnnouncementResp, len(*announcements))
	for i, announcement := range *announcements {
		updateTime := announcement.UpdateTime.Format(_defaultDateTimeFormat)
		resp[i] = PublishedAnnouncementResp{
			Id:         announcement.Id,
			Title:      announcement.Title,
			Content:    announcement.Content,
			Obj:        announcement.Target,
			UpdateTime: &updateTime,
		}
	}

	responseData(c, resp)
}

type AnnouncementAddReq struct {
	Title    string    `json:"title" binging:"required"`
	Content  string    `json:"content" binging:"required"`
	Obj      string    `json:"obj" binging:"required"`
	Platform *[]string `json:"platform"`
}

func addAnnouncement(c *gin.Context) {
	req := new(AnnouncementAddReq)
	if err := c.ShouldBindJSON(req); err != nil {
		responseEcode(c, ecode.ParamWrong)
		return
	}

	announcement := service.AnnouncementAddParam{
		Title:   &req.Title,
		Content: &req.Content,
		Target:  &req.Obj,
	}

	// 平台参数为空，则默认全部平台
	if req.Platform == nil || len(*req.Platform) == 0 {
		*announcement.Platform = make([]string, 3)
		(*announcement.Platform)[1] = "app"
		(*announcement.Platform)[2] = "iod"
		(*announcement.Platform)[3] = "mp"
	} else {
		announcement.Platform = req.Platform
	}

	err := srv.AddAnnouncement(&announcement)
	if err != nil {
		responseEcode(c, err)
		return
	}

	responseData(c, nil)
}

type AnnouncementAdminResp struct {
	Id         int64  `json:"newsid"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Obj        string `json:"obj"`
	Status     int8   `json:"status"`
	Platform   string `json:"platform"`
	UpdateTime string `json:"updateTime"`
}

func getAnnouncement(c *gin.Context) {
	query := new(PlatformPaginationReq)
	if err := c.ShouldBindQuery(query); err != nil {
		responseEcode(c, ecode.ParamWrong)
		return
	}

	pageParam := common.Pagination{
		Page:     query.Page,
		PageSize: query.Size,
	}

	announcements, total, err := srv.GetAllAnnouncement(pageParam, query.Platform)
	if err != nil {
		responseEcode(c, err)
		return
	}

	resp := make([]AnnouncementAdminResp, len(*announcements))
	for i, announcement := range *announcements {
		resp[i] = AnnouncementAdminResp{
			Id:         announcement.Id,
			Title:      *announcement.Title,
			Content:    *announcement.Content,
			Obj:        *announcement.Target,
			Status:     _internalAnnouncementStatus2ApiDefineStatus(*announcement.Status),
			Platform:   *announcement.Platform,
			UpdateTime: announcement.UpdateTime.Format(_defaultDateTimeFormat),
		}
	}

	responseData(c, map[string]interface{}{
		"notices": resp,
		"num":     total,
	})
}

func _internalAnnouncementStatus2ApiDefineStatus(internalStatus int8) int8 {
	switch internalStatus {
	case model.AnnouncementNotPublishedStatus:
		return 0
	case model.AnnouncementPublishedStatus:
		return 1
	default:
		return 0
	}
}

type AnnouncementModifyReq struct {
	Id       int64   `json:"newsid" binging:"required"`
	Title    *string `json:"title"`
	Content  *string `json:"content"`
	Obj      *string `json:"obj"`
	Platform *string `json:"platform"`
	Status   *int8   `json:"status"`
}

func modifyAnnouncement(c *gin.Context) {
	req := new(AnnouncementModifyReq)
	if err := c.ShouldBindJSON(req); err != nil {
		responseEcode(c, ecode.ParamWrong)
		return
	}

	announcement := service.AnnouncementModifyParam{
		Id:       req.Id,
		Title:    req.Title,
		Content:  req.Content,
		Target:   req.Obj,
		Platform: req.Platform,
		Status:   _announcementApiDefineStatus2InternalStatus(req.Status),
	}

	err := srv.ModifyAnnouncement(&announcement)
	if err != nil {
		responseEcode(c, err)
		return
	}

	responseData(c, nil)
}

func _announcementApiDefineStatus2InternalStatus(apiDefStatus *int8) *int8 {
	if apiDefStatus == nil {
		return nil
	}

	var status = model.NormalStatus
	// status字段值转换
	switch *apiDefStatus {
	case 0:
		status = model.NormalStatus
	case 1:
		status = model.AnnouncementPublishedStatus
	}

	return &status
}

type AnnouncementDeleteReq struct {
	Id int64 `json:"newsid"`
}

func deleteAnnouncement(c *gin.Context) {
	req := new(AnnouncementDeleteReq)
	if err := c.ShouldBindJSON(req); err != nil {
		responseEcode(c, ecode.ParamWrong)
		return
	}

	err := srv.DeleteAnnouncement(req.Id)
	if err != nil {
		responseEcode(c, err)
		return
	}

	responseData(c, nil)
}

type AnnouncementPublishReq struct {
	Ids []int64 `json:"newsid"`
}

func publishAnnouncement(c *gin.Context) {
	req := new(AnnouncementPublishReq)
	if err := c.ShouldBindJSON(req); err != nil {
		responseEcode(c, ecode.ParamWrong)
		return
	}

	err := srv.PublishAnnouncementBatch(req.Ids)
	if err != nil {
		responseEcode(c, err)
		return
	}

	responseData(c, nil)
}
