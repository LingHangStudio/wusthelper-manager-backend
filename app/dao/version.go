package dao

import (
	"go.uber.org/zap"
	"wusthelper-manager-go/app/model"
	"wusthelper-manager-go/common"
	"wusthelper-manager-go/library/ecode"
	"wusthelper-manager-go/library/log"
	"xorm.io/xorm"
)

func (d *Dao) GetLatestVersion(platform ...string) (*model.Version, error) {
	result := new(model.Version)
	session := d.db.Where("status = ?", model.VersionPublishedStatus).Desc("version_text")
	if platform != nil && len(platform) > 0 {
		session.In("platform", platform)
	}

	has, err := session.Get(result)
	if err != nil {
		log.Error("获取版本信息时出现错误", zap.Any("platform", platform), zap.String("err", err.Error()))
		return nil, ecode.InternalError
	}

	if !has {
		return nil, nil
	}

	return result, nil
}

func (d *Dao) GetVersion(id int64) (*model.Version, error) {
	result := new(model.Version)
	exists, err := d.db.
		Where("id = ?", id).And("status != ?", model.DeletedStatus).
		Get(result)
	if err != nil {
		log.Error("获取版本信息时出现错误", zap.Int64("id", id), zap.String("err", err.Error()))
		return nil, ecode.InternalError
	} else if !exists {
		return nil, nil
	}

	return result, nil
}

func (d *Dao) GetVersionList(paging common.Pagination, platform string) (*[]model.Version, int64, error) {
	countSession := d.db.Where("status != ?", model.DeletedStatus)
	if platform != "" {
		countSession.And("platform = ?", platform)
	}

	total, err := countSession.Count(model.Version{})
	if err != nil {
		log.Error("获取版本信息数量时出现错误", zap.String("err", err.Error()))
		return nil, 0, ecode.InternalError
	}

	result := make([]model.Version, 0)
	querySession := d.db.Where("status != ?", model.DeletedStatus)
	if platform != "" {
		querySession.And("platform = ?", platform)
	}
	err = querySession.Limit(paging.PageSize, paging.PageSize*(paging.Page-1)).Find(&result)
	if err != nil {
		log.Error("获取版本信息列表时出现错误", zap.String("err", err.Error()))
		return nil, 0, ecode.InternalError
	}

	return &result, total, nil
}

func (d *Dao) AddVersion(version *model.Version) (int64, error) {
	count, err := d.db.InsertOne(version)
	if err != nil {
		log.Error("添加版本信息时出现错误", zap.Any("entity", version), zap.String("err", err.Error()))
		return 0, ecode.InternalError
	}

	return count, nil
}

func (d *Dao) UpdateVersion(version *model.Version) (int64, error) {
	count, err := d.db.Omit("id").NoVersionCheck().
		Where("id = ?", version.ID).And("status != ?", model.DeletedStatus).
		Update(version)
	if err != nil {
		log.Error("修改版本信息时出现错误", zap.Any("entity", version), zap.String("err", err.Error()))
		return 0, err
	}

	return count, nil
}

func (d *Dao) PublishVersion(id int64) (int64, error) {
	transaction := d.db.NewSession()
	defer func(transaction *xorm.Session) {
		err := transaction.Close()
		if err != nil {
			log.Warn("发布版本信息时出现错误，事务session关闭时出现异常", zap.Int64("id", id), zap.Error(err))
		}
	}(transaction)

	if err := transaction.Begin(); err != nil {
		log.Error("发布版本信息时出现错误，事务开启时出现异常", zap.Int64("id", id), zap.Error(err))
		return 0, ecode.InternalError
	}

	// 先获取一遍待发布的版本信息，获取其平台，修改该平台其他版本状态为普通状态，当前版本设置为发布，保证一个平台只有一个已发布版本
	version := new(model.Version)
	has, err := transaction.
		Cols("platform").
		Where("id = ?", id).And("status != ?", model.DeletedStatus).
		Get(version)

	if err != nil {
		log.Error("发布版本信息时出现错误，获取待发布版本信息时出现异常",
			zap.Int64("id", id), zap.String("err", err.Error()),
		)
		return 0, ecode.InternalError
	} else if !has {
		log.Error("发布版本信息时出现错误，id不存在", zap.Int64("id", id), zap.String("err", err.Error()))
		return 0, ecode.QueryFailed
	}

	// 修改当前平台其他版本状态为普通状态
	publishedStatus, normalStatus := model.VersionPublishedStatus, model.NormalStatus
	_, err = transaction.Omit("id").
		In("platform", *version.Platform).And("status = ?", publishedStatus).
		Update(&model.Version{Status: &normalStatus})

	if err != nil {
		log.Error("发布版本信息时出现错误，切换其他版本信息状态时出现异常",
			zap.Int64("id", id),
			zap.String("platform", *version.Platform),
			zap.String("err", err.Error()),
		)
		return 0, ecode.InternalError
	}

	count, err := transaction.Omit("id").
		Where("id = ?", id).And("status != ?", model.DeletedStatus).
		Update(&model.Version{Status: &publishedStatus})
	if err != nil {
		log.Error("发布版本信息时出现错误", zap.String("err", err.Error()))
		return 0, ecode.InternalError
	}

	err = transaction.Commit()
	if err != nil {
		log.Error("发布版本信息时出现错误，提交事务时出现异常", zap.Int64("id", id), zap.Error(err))
		return 0, ecode.InternalError
	}

	return count, nil
}

func (d *Dao) UpdateVersionStatus(status int8, id int64) (int64, error) {
	count, err := d.db.Omit("id").
		Where("id = ?", id).
		And("status != ?", model.DeletedStatus).
		Update(&model.Version{Status: &status})
	if err != nil {
		log.Error("修改版本信息时出现错误", zap.String("err", err.Error()))
		return 0, err
	}

	return count, nil
}

func (d *Dao) UpdateVersionStatusBatch(status int8, id ...int64) (int64, error) {
	count, err := d.db.Omit("id").
		In("id", id).
		And("status != ?", model.DeletedStatus).
		Update(&model.Version{Status: &status})
	if err != nil {
		log.Error("修改版本信息时出现错误",
			zap.Any("id", id),
			zap.Int8("status", status),
			zap.String("err", err.Error()),
		)
		return 0, err
	}

	return count, nil
}

func (d *Dao) DeleteVersion(id int64) error {
	status := model.DeletedStatus
	_, err := d.db.Omit("id").
		Where("id = ?", id).
		Update(&model.Version{Status: &status})
	if err != nil {
		log.Error("删除版本信息时出现错误", zap.Any("id", id), zap.String("err", err.Error()))
		return ecode.InternalError
	}

	return nil
}
