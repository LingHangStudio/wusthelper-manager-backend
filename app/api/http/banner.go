package http

import (
	"github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"wusthelper-manager-go/app/model"
	"wusthelper-manager-go/app/service"
	"wusthelper-manager-go/common"
	"wusthelper-manager-go/library/ecode"
	"wusthelper-manager-go/library/log"
)

type BannerInfoResp struct {
	Actid      int64  `json:"actid"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	ImgUrl     string `json:"imgUrl"`
	Status     int8   `json:"status"`
	Platform   string `json:"platform"`
	UpdateTime string `json:"updateTime"`
}

func getBannerList(c *gin.Context) {
	req := new(PlatformPaginationReq)
	if err := c.ShouldBind(req); err != nil {
		responseEcode(c, ecode.ParamWrong)
		return
	}

	bannerList, total, err := srv.GetBannerList(common.Pagination{Page: req.Page, PageSize: req.Size}, req.Platform)

	if err != nil {
		responseEcode(c, err)
		return
	}

	resultList := make([]BannerInfoResp, len(*bannerList))
	for i, banner := range *bannerList {
		resultList[i] = BannerInfoResp{
			Actid:      banner.ID,
			Title:      *banner.Title,
			Content:    *banner.Link,
			ImgUrl:     _getPicUrl(*banner.Img),
			Status:     _internalBannerStatus2ApiDefineStatus(*banner.Status),
			Platform:   *banner.Platform,
			UpdateTime: banner.UpdateTime.Format(_defaultDateTimeFormat),
		}
	}

	responseData(c, map[string]any{
		"actList": resultList,
		"num":     total,
	})
}

type PublishedBannerResp struct {
	Actid      int64  `json:"actid"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	ImgUrl     string `json:"imgUrl"`
	UpdateTime string `json:"updateTime"`
}

func getPublishedBannerList(c *gin.Context) {
	platform := getPlatform(c)
	var resultList *[]model.Banner
	var err error
	if platform == "" {
		resultList, err = srv.GetPublishedBanner()
	} else {
		resultList, err = srv.GetPublishedBanner(platform)
	}

	if err != nil {
		responseEcode(c, err)
		return
	}

	publishedBannerList := make([]PublishedBannerResp, len(*resultList))
	for i, banner := range *resultList {
		publishedBannerList[i] = PublishedBannerResp{
			Actid:      int64(i),
			Title:      *banner.Title,
			Content:    *banner.Link,
			ImgUrl:     _getPicUrl(*banner.Img),
			UpdateTime: banner.UpdateTime.Format(_defaultDateTimeFormat),
		}
	}

	responseData(c, publishedBannerList)
}

func _internalBannerStatus2ApiDefineStatus(internalStatus int8) int8 {
	switch internalStatus {
	case model.NormalStatus:
		return 0
	case model.BannerPublishedStatus:
		return 1
	default:
		return 0
	}
}

type BannerPublishReq struct {
	Actid    []int64 `json:"actid" binging:"required"`
	Platform string  `json:"platform" binging:"required"`
}

func publishBanner(c *gin.Context) {
	req := new(BannerPublishReq)
	if err := c.ShouldBind(req); err != nil {
		responseEcode(c, ecode.ParamWrong)
		return
	}

	err := srv.PublishBannerBatch(req.Actid...)
	if err != nil {
		responseEcode(c, err)
		return
	}

	responseData(c, nil)
}

type BannerAddReq struct {
	Title    string                `json:"title" form:"title" binding:"required"`
	Content  string                `json:"content" form:"content" binding:"required"`
	Platform []string              `json:"platform" form:"platform" binding:"required"`
	File     *multipart.FileHeader `json:"file" form:"file" binding:"required"`
}

func addBanner(c *gin.Context) {
	req := new(BannerAddReq)
	if err := c.ShouldBind(req); err != nil {
		responseEcode(c, ecode.ParamWrong)
		return
	}

	var uploadFile *service.File = nil
	if req.File != nil {
		// 限制文件100mb以内
		if req.File.Size > 100*humanize.MByte {
			responseEcode(c, ecode.ParamWrong)
			return
		}

		fileName := req.File.Filename
		file, err := req.File.Open()
		if err != nil {
			responseEcode(c, ecode.InternalError)
			return
		}

		// 如果文件比较多或者比较大，service.BannerAddParam.File.Data应该用io.Reader，需要使用流式的读写，
		// 此处会直接将上传的文件全部读入内存里边，传的文件多了或者大了就会导致内存占用太大，甚至oom
		// 使用io.Reader和io.Writer是为了流式读写，防止内存占用过大，
		// 不过目前的使用情景用不到太多内存（在请求的时候就限制住大小，但是只是限制了单次上传，多人并发上传也还是会有同样的问题）
		fileData, err := io.ReadAll(file)
		if err != nil {
			log.Error("读取上传文件时出现错误", zap.Error(err))
			responseEcode(c, ecode.InternalError)
			return
		}

		uploadFile = &service.File{
			Data:     &fileData,
			FileName: fileName,
		}
	}

	banner := service.BannerAddParam{
		Title:    req.Title,
		Link:     req.Content,
		Img:      uploadFile,
		Platform: req.Platform,
	}

	err := srv.AddBanner(&banner)
	if err != nil {
		responseEcode(c, err)
		return
	}

	responseData(c, nil)
}

type BannerModifyReq struct {
	Actid    int64   `json:"actid" form:"actid" binding:"required"`
	Title    *string `json:"title" form:"title"`
	Content  *string `json:"content" form:"content"`
	Platform *string `json:"platform" form:"platform"`
	Status   *int8   `json:"status" form:"status"`
}

func modifyBanner(c *gin.Context) {
	req := new(BannerModifyReq)
	if err := c.ShouldBind(req); err != nil {
		responseEcode(c, ecode.ParamWrong)
		return
	}

	banner := service.BannerModifyParam{
		Id:       req.Actid,
		Title:    req.Title,
		Link:     req.Content,
		Platform: req.Platform,
		Status:   req.Status,
	}

	err := srv.ModifyBanner(&banner)
	if err != nil {
		responseEcode(c, err)
		return
	}

	responseData(c, nil)
}

type BannerDeleteReq struct {
	Actid int64 `form:"actid" binding:"required"`
}

func deleteBanner(c *gin.Context) {
	req := new(BannerDeleteReq)
	if err := c.ShouldBind(req); err != nil {
		responseEcode(c, ecode.ParamWrong)
		return
	}

	err := srv.DeleteBanner(req.Actid)
	if err != nil {
		responseEcode(c, err)
		return
	}

	responseData(c, nil)
}
