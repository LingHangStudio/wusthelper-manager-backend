package service

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/yitter/idgenerator-go/idgen"
	"go.uber.org/zap"
	"os"
	"time"
	"wusthelper-manager-go/app/model"
	"wusthelper-manager-go/common"
	"wusthelper-manager-go/library/ecode"
	"wusthelper-manager-go/library/log"
)

func (s *Service) GetLatestVersion(platform ...string) (*model.Version, error) {
	latestVersionList, err := s.dao.GetLatestVersion(platform...)
	if err != nil {
		return nil, err
	}

	return latestVersionList, nil
}

func (s *Service) GetVersionList(pagination common.Pagination, platform string) (*[]model.Version, int64, error) {
	versionList, total, err := s.dao.GetVersionList(pagination, platform)
	if err != nil {
		return nil, 0, err
	}

	return versionList, total, nil
}

type VersionAddParam struct {
	Version    string
	Summary    string
	Platform   string
	UploadFile *File
}

func (s *Service) AddVersion(param *VersionAddParam) error {
	// 先存到本地，再去传oss
	storageOption := s.config.Server.FileStorageOption
	localFileLoc := fmt.Sprintf("%s/%s", storageOption.UploadFileLocalTmpPath, param.UploadFile.FileName)

	fileKey := ""
	if param.UploadFile != nil {
		err := os.WriteFile(localFileLoc, *param.UploadFile.Data, 0664)
		if err != nil {
			log.Error("写版本文件到本地临时目录时出现错误", zap.String("file", localFileLoc), zap.Error(err))
			return ecode.InternalError
		}

		fileKey = fmt.Sprintf("%d/%s", idgen.NextId(), param.UploadFile.FileName)

		// 上传新文件
		go func() {
			log.Info("新版本文件处理后台任务开始")
			resourceStorageOption := storageOption.ResourceStorageOption
			ossObjectKey := fmt.Sprintf("%s/%s", resourceStorageOption.VersionFileStorageBasePath, fileKey)
			err = s.ossBucket.PutObjectFromFile(ossObjectKey, localFileLoc)
			if err != nil {
				log.Warn("版本文件上传oss时出现错误",
					zap.String("oss_key", ossObjectKey),
					zap.String("local_source", localFileLoc),
					zap.Error(err),
				)
			} else {
				err = os.Remove(localFileLoc)
				if err != nil {
					log.Warn("移除本地版本文件时出现异常", zap.String("file", localFileLoc), zap.Error(err))
				}
				log.Info("版本文件上传oss成功",
					zap.String("oss_key", ossObjectKey),
					zap.String("local_source", localFileLoc),
				)
			}
			log.Info("新版本文件处理后台任务完成")
		}()
	}

	now := time.Now()
	version := model.Version{
		ID:          idgen.NextId(),
		VersionText: &param.Version,
		Summary:     &param.Summary,
		Platform:    &param.Platform,
		File:        &fileKey,
		Status:      new(int8),
		CreateTime:  &now,
		UpdateTime:  &now,
	}

	*version.Status = model.NormalStatus

	_, err := s.dao.AddVersion(&version)
	if err != nil {
		return err
	}

	return nil
}

type VersionModifyParam struct {
	Id         int64
	Version    *string
	Summary    *string
	Platform   *string
	UploadFile *File
}

func (s *Service) ModifyVersion(param *VersionModifyParam) error {
	version := model.Version{
		ID:          param.Id,
		VersionText: param.Version,
		Summary:     param.Summary,
		File:        nil,
		Platform:    param.Platform,
		UpdateTime:  new(time.Time),
	}

	// 新版本文件需要修改
	if param.UploadFile != nil {
		existsVersion, err := s.dao.GetVersion(param.Id)
		if err != nil {
			return err
		}

		// 如果没查到数据
		if existsVersion == nil {
			return ecode.VersionOperationFailed
		}

		// 先保存新文件到本地
		storageOption := s.config.Server.FileStorageOption
		localFileLoc := fmt.Sprintf("%s/%s", storageOption.UploadFileLocalTmpPath, param.UploadFile.FileName)
		err = os.WriteFile(localFileLoc, *param.UploadFile.Data, 0664)
		if err != nil {
			log.Error("写版本文件到本地临时目录时出现错误", zap.String("file", localFileLoc), zap.Error(err))
			return ecode.InternalError
		}

		// 如果原来有文件就删除
		resourceStorageOption := storageOption.ResourceStorageOption
		ossObjectKey := fmt.Sprintf("%s/%s", resourceStorageOption.VersionFileStorageBasePath, *existsVersion.File)
		if version.File != nil && *version.File != "" {
			err = s.ossBucket.SetObjectACL(ossObjectKey, oss.ACLPrivate)
			if err != nil {
				log.Warn("删除已存在的oss文件出现错误", zap.String("oss_key", ossObjectKey), zap.Error(err))
			} else {
				log.Info("删除已存在的oss文件完成", zap.String("oss_key", ossObjectKey))
			}
		}

		fileKey := fmt.Sprintf("%d/%s", idgen.NextId(), param.UploadFile.FileName)
		version.File = &fileKey

		// 上传新文件
		go func() {
			log.Info("新版本文件处理后台任务开始")
			ossObjectKey = fmt.Sprintf("%s/%s", resourceStorageOption.VersionFileStorageBasePath, fileKey)
			err = s.ossBucket.PutObjectFromFile(ossObjectKey, localFileLoc)
			if err != nil {
				log.Warn("版本文件上传oss时出现错误",
					zap.String("oss_key", ossObjectKey),
					zap.String("local_source", localFileLoc),
					zap.Error(err),
				)
				fileKey = ""
				_, _ = s.dao.UpdateVersion(&model.Version{ID: param.Id, File: &fileKey})
			} else {
				err = os.Remove(localFileLoc)
				if err != nil {
					log.Warn("移除本地版本文件时出现异常", zap.String("file", localFileLoc), zap.Error(err))
				}
				log.Info("版本文件上传oss成功",
					zap.String("oss_key", ossObjectKey),
					zap.String("local_source", localFileLoc),
				)
			}
			log.Info("新版本文件处理后台任务完成")
		}()
	}

	*version.UpdateTime = time.Now()
	_, err := s.dao.UpdateVersion(&version)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteVersion(id int64) error {
	existsVersion, err := s.dao.GetVersion(id)
	if err != nil {
		return err
	}

	if existsVersion == nil {
		return ecode.VersionOperationFailed
	}

	err = s.dao.DeleteVersion(id)
	if err != nil {
		return err
	}

	// 如果有文件记录，删除oss文件（设置为私有不可见）
	if existsVersion.File != nil && *existsVersion.File != "" {
		resourceStorageOption := s.config.Server.FileStorageOption.ResourceStorageOption
		ossObjectKey := fmt.Sprintf("%s/%s", resourceStorageOption.VersionFileStorageBasePath, *existsVersion.File)
		err = s.ossBucket.SetObjectACL(ossObjectKey, oss.ACLPrivate)
		if err != nil {
			log.Warn("删除oss文件出现错误", zap.String("oss_key", ossObjectKey), zap.Error(err))
		} else {
			log.Info("删除oss文件完成", zap.String("oss_key", ossObjectKey))
		}
	}

	return nil
}

func (s *Service) PublishVersion(id int64) error {
	_, err := s.dao.PublishVersion(id)
	if err != nil {
		return err
	}

	go func() {
		log.Info("新版本发布后台任务开始")
		version, err := s.dao.GetVersion(id)
		if err != nil {
			return
		}

		if version.File == nil || *version.File == "" {
			log.Info("该平台版本无文件，不需处理", zap.Int64("id", id))
			return
		}

		resourceStorageOption := s.config.Server.FileStorageOption.ResourceStorageOption
		ossObjectKey := fmt.Sprintf("%s/%s", resourceStorageOption.VersionFileStorageBasePath, *version.File)
		err = s.ossBucket.PutSymlink(resourceStorageOption.WusthelperReleaseFileKey, ossObjectKey)
		if err != nil {
			log.Error("创建助手网页发布文件oss软连接时出现错误", zap.Int64("id", id), zap.Error(err))
			return
		}

		log.Info("新版本发布后台任务完毕")
	}()

	return nil
}
