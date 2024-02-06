package dao

import (
	"go.uber.org/zap"
	"wusthelper-manager-go/app/model"
	"wusthelper-manager-go/common"
	"wusthelper-manager-go/library/ecode"
	"wusthelper-manager-go/library/log"
)

func (d *Dao) GetPublishedBanner(platform ...string) (*[]model.Banner, error) {
	result := make([]model.Banner, 0)
	session := d.db.Where("status = ?", model.VersionPublishedStatus)
	if platform != nil && len(platform) > 0 {
		session.In("platform", platform)
	}

	err := session.Desc("id").Find(&result)
	if err != nil {
		log.Error("获取轮播图时出现错误", zap.Any("platform", platform), zap.String("err", err.Error()))
		return nil, ecode.InternalError
	}

	return &result, nil
}

func (d *Dao) GetBanner(id int64) (*model.Banner, error) {
	result := new(model.Banner)
	exists, err := d.db.
		Where("id = ?", id).And("status != ?", model.DeletedStatus).
		Get(result)
	if err != nil {
		log.Error("获取轮播图时出现错误", zap.Int64("id", id), zap.String("err", err.Error()))
		return nil, ecode.InternalError
	} else if !exists {
		return nil, nil
	}

	return result, nil
}

func (d *Dao) GetBannerList(paging common.Pagination, platform string) (*[]model.Banner, int64, error) {
	countSession := d.db.Where("status != ?", model.DeletedStatus)
	if platform != "" {
		countSession.And("platform = ?", platform)
	}

	total, err := countSession.Count(&model.Banner{})
	if err != nil {
		log.Error("获取轮播图数量时出现错误", zap.String("err", err.Error()))
		return nil, 0, ecode.InternalError
	}

	result := make([]model.Banner, 0)
	querySession := d.db.Where("status != ?", model.DeletedStatus)
	if platform != "" {
		querySession.And("platform = ?", platform)
	}
	err = querySession.Desc("status", "id").
		Limit(paging.PageSize, paging.PageSize*(paging.Page-1)).Find(&result)
	if err != nil {
		log.Error("获取轮播图列表时出现错误", zap.String("err", err.Error()))
		return nil, 0, ecode.InternalError
	}

	return &result, total, nil
}

func (d *Dao) AddBanner(banner ...model.Banner) (int64, error) {
	count, err := d.db.Insert(banner)
	if err != nil {
		log.Error("添加轮播图时出现错误", zap.Any("entity", banner), zap.String("err", err.Error()))
		return 0, ecode.InternalError
	}

	return count, nil
}

func (d *Dao) UpdateBanner(banner *model.Banner) (int64, error) {
	count, err := d.db.Omit("id").
		Where("id = ?", banner.ID).
		And("status != ?", model.DeletedStatus).
		Update(banner)
	if err != nil {
		log.Error("修改轮播图时出现错误", zap.Any("entity", banner), zap.String("err", err.Error()))
		return 0, err
	}

	return count, nil
}

func (d *Dao) PublishBanner(id int64) (int64, error) {
	status := model.BannerPublishedStatus
	count, err := d.db.Omit("id").
		Where("id = ?", id).And("status != ?", model.DeletedStatus).
		Update(&model.Banner{Status: &status})
	if err != nil {
		log.Error("发布轮播图时出现错误", zap.String("err", err.Error()))
		return 0, ecode.InternalError
	}

	return count, nil
}

func (d *Dao) UpdateBannerStatus(status int8, id int64) (int64, error) {
	count, err := d.db.Omit("id").
		Where("id = ?", id).
		And("status != ?", model.DeletedStatus).
		Update(&model.Banner{Status: &status})
	if err != nil {
		log.Error("修改轮播图时出现错误", zap.String("err", err.Error()))
		return 0, err
	}

	return count, nil
}

func (d *Dao) UpdateBannerStatusBatch(status int8, id ...int64) (int64, error) {
	count, err := d.db.Omit("id").
		In("id", id).
		And("status != ?", model.DeletedStatus).
		Update(&model.Banner{Status: &status})
	if err != nil {
		log.Error("修改轮播图时出现错误",
			zap.Any("id", id),
			zap.Int8("status", status),
			zap.String("err", err.Error()),
		)
		return 0, err
	}

	return count, nil
}

func (d *Dao) DeleteBanner(id int64) error {
	status := model.DeletedStatus
	_, err := d.db.Omit("id").
		Where("id = ?", id).
		Update(&model.Banner{Status: &status})
	if err != nil {
		log.Error("删除轮播图时出现错误", zap.Any("id", id), zap.String("err", err.Error()))
		return ecode.InternalError
	}

	return nil
}
