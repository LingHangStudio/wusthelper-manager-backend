package service

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/dustin/go-humanize"
	"github.com/sunshineplan/imgconv"
	"github.com/yitter/idgenerator-go/idgen"
	"go.uber.org/zap"
	"math/rand"
	"os"
	"strings"
	"time"
	"wusthelper-manager-go/app/model"
	"wusthelper-manager-go/common"
	"wusthelper-manager-go/library/ecode"
	"wusthelper-manager-go/library/log"
)

func (s *Service) GetPublishedBanner(platform ...string) (*[]model.Banner, error) {
	latestBannerList, err := s.dao.GetPublishedBanner(platform...)
	if err != nil {
		return nil, err
	}

	return latestBannerList, nil
}

func (s *Service) GetBannerList(pagination common.Pagination, platform string) (*[]model.Banner, int64, error) {
	bannerList, total, err := s.dao.GetBannerList(pagination, platform)
	if err != nil {
		return nil, 0, err
	}

	return bannerList, total, nil
}

type BannerAddParam struct {
	Title    string
	Link     string
	Img      *File
	Platform []string
}

func (s *Service) AddBanner(param *BannerAddParam) error {
	now := time.Now()
	banners := make([]model.Banner, len(param.Platform))
	for i, platform := range param.Platform {
		bannerId := idgen.NextId()
		imgId := ""
		if param.Img != nil {
			// 先存到本地，再去传oss
			localFileLoc, err := s.processBannerImg(param.Img)
			if err != nil {
				return err
			}

			imgId = fmt.Sprintf("%d/v1.%d.%d", bannerId, idgen.NextId(), time.Now().UnixMilli())

			// 上传新文件
			go s.uploadBannerPic(imgId, localFileLoc)
		}

		status := model.NormalStatus
		p := strings.Clone(platform)
		banners[i] = model.Banner{
			ID:         bannerId,
			Title:      &param.Title,
			Link:       &param.Link,
			Img:        &imgId,
			Platform:   &p,
			CreateTime: &now,
			UpdateTime: &now,
			Status:     &status,
		}
	}

	_, err := s.dao.AddBanner(banners...)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) uploadBannerPic(imgId, localFileLoc string) {
	log.Info("banner图片后台上传任务开始")
	resourceStorageOption := s.config.Server.FileStorageOption.ResourceStorageOption
	ossObjectKey := fmt.Sprintf("%s/%s.jpg", resourceStorageOption.PicStorageBasePath, imgId)
	err := s.ossBucket.PutObjectFromFile(ossObjectKey, localFileLoc, oss.Meta("Content-Type", "image/jpeg"))
	if err != nil {
		log.Warn("banner图片上传oss时出现错误",
			zap.String("oss_key", ossObjectKey),
			zap.String("local_source", localFileLoc),
			zap.Error(err),
		)
	} else {
		err = os.Remove(localFileLoc)
		if err != nil {
			log.Warn("移除本地banner图片时出现异常", zap.String("file", localFileLoc), zap.Error(err))
		}
		log.Info("banner图片上传oss成功",
			zap.String("oss_key", ossObjectKey),
			zap.String("local_source", localFileLoc),
		)
	}
	log.Info("banner图片后台上传任务完成")
}

func (s *Service) processBannerImg(imgFile *File) (string, error) {
	storageOption := s.config.Server.FileStorageOption
	tmpId := time.Now().UnixMilli() + rand.Int63()
	tmpImgFile := fmt.Sprintf("%s/tmp-%d-%s", storageOption.UploadFileLocalTmpPath, tmpId, imgFile.FileName)
	processedFileLoc := fmt.Sprintf("%s/%d-%s", storageOption.UploadFileLocalTmpPath, tmpId, imgFile.FileName)

	err := os.WriteFile(tmpImgFile, *imgFile.Data, 0664)
	if err != nil {
		log.Error("写临时图片到本地临时目录时出现错误", zap.String("file", tmpImgFile), zap.Error(err))
		return "", ecode.InternalError
	}

	img, err := imgconv.Open(tmpImgFile)
	if err != nil {
		log.Warn("图片数据读取错误", zap.Error(err))
		return "", ecode.ParamWrong
	}

	option := imgconv.FormatOption{
		Format: imgconv.JPEG,
	}
	if len(*imgFile.Data) > 1*humanize.MByte {
		option.EncodeOption = append(make([]imgconv.EncodeOption, 0, 1), imgconv.Quality(75))
	}

	err = imgconv.Save(processedFileLoc, img, &option)
	if err != nil {
		log.Error("图片处理错误", zap.Error(err))
		return "", ecode.ParamWrong
	}

	err = os.Remove(tmpImgFile)
	if err != nil {
		log.Warn("删除临时存储图片错误", zap.Error(err))
	}

	return processedFileLoc, nil
}

type BannerModifyParam struct {
	Id       int64
	Title    *string
	Link     *string
	Img      *File
	Platform *string
	Status   *int8
}

func (s *Service) ModifyBanner(param *BannerModifyParam) error {
	now := time.Now()
	banner := model.Banner{
		ID:         param.Id,
		Title:      param.Title,
		Link:       param.Link,
		Img:        nil,
		Platform:   param.Platform,
		UpdateTime: &now,
		Status:     param.Status,
	}

	// banner图片需要修改
	if param.Img != nil {
		existsBanner, err := s.dao.GetBanner(param.Id)
		if err != nil {
			return err
		}

		// 如果没查到数据
		if existsBanner == nil {
			return ecode.InvalidId
		}

		// 先保存新文件到本地
		storageOption := s.config.Server.FileStorageOption
		localFileLoc, err := s.processBannerImg(param.Img)
		if err != nil {
			return err
		}

		// 如果原来有文件就删除
		resourceStorageOption := storageOption.ResourceStorageOption
		ossObjectKey := fmt.Sprintf("%s/%s.jpg", resourceStorageOption.PicStorageBasePath, *existsBanner.Img)
		if banner.Img != nil && *banner.Img != "" {
			err = s.ossBucket.SetObjectACL(ossObjectKey, oss.ACLPrivate)
			if err != nil {
				log.Warn("删除已存在的oss文件出现错误", zap.String("oss_key", ossObjectKey), zap.Error(err))
			} else {
				log.Info("删除已存在的oss文件完成", zap.String("oss_key", ossObjectKey))
			}
		}

		imgId := fmt.Sprintf("%d/v1.%d.%d", param.Id, idgen.NextId(), time.Now().UnixMilli())
		banner.Img = &imgId

		// 上传新文件
		go s.uploadBannerPic(imgId, localFileLoc)
	}

	*banner.UpdateTime = time.Now()
	_, err := s.dao.UpdateBanner(&banner)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteBanner(id int64) error {
	existsBanner, err := s.dao.GetBanner(id)
	if err != nil {
		return err
	}

	if existsBanner == nil {
		return ecode.InvalidId
	}

	err = s.dao.DeleteBanner(id)
	if err != nil {
		return err
	}

	// 如果有文件记录，删除oss文件（仅设置不可见）
	if existsBanner.Img != nil && *existsBanner.Img != "" {
		resourceStorageOption := s.config.Server.FileStorageOption.ResourceStorageOption
		ossObjectKey := fmt.Sprintf("%s/%s.jpg", resourceStorageOption.PicStorageBasePath, *existsBanner.Img)
		err = s.ossBucket.SetObjectACL(ossObjectKey, oss.ACLPrivate)
		if err != nil {
			log.Warn("删除oss文件出现错误", zap.String("oss_key", ossObjectKey), zap.Error(err))
		} else {
			log.Info("删除oss文件完成", zap.String("oss_key", ossObjectKey))
		}
	}

	return nil
}

func (s *Service) PublishBanner(id int64) error {
	_, err := s.dao.UpdateBannerStatus(model.BannerPublishedStatus, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) PublishBannerBatch(id ...int64) error {
	_, err := s.dao.UpdateBannerStatusBatch(model.BannerPublishedStatus, id...)
	if err != nil {
		return err
	}

	return nil
}
