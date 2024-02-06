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

type LatestVersionResp struct {
	Version       string `json:"version"`
	UpdateContent string `json:"updateContent"`
	ApkUrl        string `json:"apkUrl"`
}

func getLatestVersion(c *gin.Context) {
	platform := getPlatform(c)
	if platform == "" {
		responseEcode(c, ecode.ParamWrong)
		return
	}

	result, err := srv.GetLatestVersion(platform)
	if err != nil {
		responseEcode(c, err)
		return
	}

	resp := LatestVersionResp{
		Version:       *result.VersionText,
		UpdateContent: *result.Summary,
		ApkUrl:        _getFileUrl(*result.File),
	}

	responseData(c, resp)
}

type VersionInfoResp struct {
	Id            int64  `json:"id"`
	Version       string `json:"version"`
	UpdateContent string `json:"updateContent"`
	ApkUrl        string `json:"apkUrl"`
	Status        int8   `json:"status"`
	Platform      string `json:"platform"`
	CreateTime    string `json:"createTime"`
}

func getVersionList(c *gin.Context) {
	req := new(PlatformPaginationReq)
	if err := c.ShouldBind(req); err != nil {
		responseEcode(c, ecode.ParamWrong)
		return
	}

	versionList, total, err := srv.GetVersionList(common.Pagination{Page: req.Page, PageSize: req.Size}, req.Platform)

	if err != nil {
		responseEcode(c, err)
		return
	}

	resultList := make([]VersionInfoResp, len(*versionList))
	for i, version := range *versionList {
		resultList[i] = VersionInfoResp{
			Id:            version.ID,
			Version:       *version.VersionText,
			UpdateContent: *version.Summary,
			ApkUrl:        _getFileUrl(*version.File),
			Status:        _internalVersionStatus2ApiDefineStatus(*version.Status),
			Platform:      *version.Platform,
			CreateTime:    version.CreateTime.Format(_defaultDateTimeFormat),
		}
	}

	responseData(c, map[string]any{
		"versionControls": resultList,
		"num":             total,
	})
}

func _internalVersionStatus2ApiDefineStatus(internalStatus int8) int8 {
	switch internalStatus {
	case model.NormalStatus:
		return 0
	case model.VersionPublishedStatus:
		return 1
	default:
		return 0
	}
}

type VersionPublishReq struct {
	Id       int64  `json:"id" binging:"required"`
	Platform string `json:"platform" binging:"required"`
}

func publishVersion(c *gin.Context) {
	req := new(VersionPublishReq)
	if err := c.ShouldBind(req); err != nil {
		responseEcode(c, ecode.ParamWrong)
		return
	}

	err := srv.PublishVersion(req.Id)
	if err != nil {
		responseEcode(c, err)
		return
	}

	responseData(c, nil)
}

type VersionAddReq struct {
	Version       string                `form:"version" binding:"required"`
	UpdateContent string                `form:"updateContent" binding:"required"`
	Platform      string                `form:"platform" binding:"required"`
	File          *multipart.FileHeader `form:"file"`
}

func addVersion(c *gin.Context) {
	req := new(VersionAddReq)
	if err := c.ShouldBind(req); err != nil {
		responseEcode(c, ecode.ParamWrong)
		return
	}

	var uploadFile *service.File = nil
	if req.File != nil {
		// 限制文件100mb以内
		if req.File.Size > 100*humanize.MByte {
			responseEcode(c, ecode.VersionOperationFailed)
			return
		}

		fileName := req.File.Filename
		file, err := req.File.Open()
		if err != nil {
			responseEcode(c, ecode.InternalError)
			return
		}

		// 如果文件比较多或者比较大，service.VersionAddParam.File.Data应该用io.Reader，需要使用流式的读写，
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

	version := service.VersionAddParam{
		Version:    req.Version,
		Summary:    req.UpdateContent,
		Platform:   req.Platform,
		UploadFile: uploadFile,
	}

	err := srv.AddVersion(&version)
	if err != nil {
		responseEcode(c, err)
		return
	}

	responseData(c, nil)
}

type VersionModifyReq struct {
	Id            int64                 `form:"id" binding:"required"`
	Version       *string               `form:"version"`
	UpdateContent *string               `form:"updateContent"`
	Platform      *string               `form:"platform"`
	File          *multipart.FileHeader `form:"file"`
}

func modifyVersion(c *gin.Context) {
	req := new(VersionModifyReq)
	if err := c.ShouldBind(req); err != nil {
		responseEcode(c, ecode.ParamWrong)
		return
	}

	var uploadFile *service.File = nil
	if req.File != nil {
		// 限制文件100mb以内
		if req.File.Size > 100*humanize.MByte {
			responseEcode(c, ecode.VersionOperationFailed)
			return
		}

		fileName := req.File.Filename
		file, err := req.File.Open()
		if err != nil {
			responseEcode(c, ecode.InternalError)
			return
		}

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

	version := service.VersionModifyParam{
		Id:         req.Id,
		Version:    req.Version,
		Summary:    req.UpdateContent,
		Platform:   req.Platform,
		UploadFile: uploadFile,
	}

	err := srv.ModifyVersion(&version)
	if err != nil {
		responseEcode(c, err)
		return
	}

	responseData(c, nil)
}

type VersionDeleteReq struct {
	Id int64 `form:"id" binding:"required"`
}

func deleteVersion(c *gin.Context) {
	req := new(VersionDeleteReq)
	if err := c.ShouldBind(req); err != nil {
		responseEcode(c, ecode.ParamWrong)
		return
	}

	err := srv.DeleteVersion(req.Id)
	if err != nil {
		responseEcode(c, err)
		return
	}

	responseData(c, nil)
}
