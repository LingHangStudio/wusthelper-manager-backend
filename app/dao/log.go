package dao

import (
	"go.uber.org/zap"
	"wusthelper-manager-go/app/model"
	"wusthelper-manager-go/common"
	"wusthelper-manager-go/library/ecode"
	"wusthelper-manager-go/library/log"
)

func (d *Dao) GetPublishedLog(platform ...string) (*[]model.Log, error) {
	result := make([]model.Log, 0)
	session := d.db.Where("status = ?", model.VersionPublishedStatus)
	if platform != nil && len(platform) > 0 {
		session.In("platform", platform)
	}

	err := session.Desc("id").Find(&result)
	if err != nil {
		log.Error("获取日志时出现错误", zap.Any("platform", platform), zap.String("err", err.Error()))
		return nil, ecode.InternalError
	}

	return &result, nil
}

func (d *Dao) GetLog(id int64) (*model.Log, error) {
	result := new(model.Log)
	exists, err := d.db.
		Where("id = ?", id).And("status != ?", model.DeletedStatus).
		Get(result)
	if err != nil {
		log.Error("获取日志时出现错误", zap.Int64("id", id), zap.String("err", err.Error()))
		return nil, ecode.InternalError
	} else if !exists {
		return nil, nil
	}

	return result, nil
}

func (d *Dao) GetLogList(paging common.Pagination, platform string) (*[]model.Log, int64, error) {
	countSession := d.db.Where("status != ?", model.DeletedStatus)
	if platform != "" {
		countSession.And("platform = ?", platform)
	}

	total, err := countSession.Count(&model.Log{})
	if err != nil {
		log.Error("获取日志数量时出现错误", zap.String("err", err.Error()))
		return nil, 0, ecode.InternalError
	}

	result := make([]model.Log, 0)
	querySession := d.db.Where("status != ?", model.DeletedStatus)
	if platform != "" {
		querySession.And("platform = ?", platform)
	}
	err = querySession.Desc("status", "id").
		Limit(paging.PageSize, paging.PageSize*(paging.Page-1)).Find(&result)
	if err != nil {
		log.Error("获取日志列表时出现错误", zap.String("err", err.Error()))
		return nil, 0, ecode.InternalError
	}

	return &result, total, nil
}

func (d *Dao) AddLog(logEntity ...model.Log) (int64, error) {
	count, err := d.db.Insert(logEntity)
	if err != nil {
		log.Error("添加日志时出现错误", zap.Any("entity", logEntity), zap.String("err", err.Error()))
		return 0, ecode.InternalError
	}

	return count, nil
}

func (d *Dao) UpdateLog(logEntity *model.Log) (int64, error) {
	count, err := d.db.Omit("id").
		Where("id = ?", logEntity.ID).
		And("status != ?", model.DeletedStatus).
		Update(logEntity)
	if err != nil {
		log.Error("修改日志时出现错误", zap.Any("entity", logEntity), zap.String("err", err.Error()))
		return 0, err
	}

	return count, nil
}

func (d *Dao) PublishLog(id int64) (int64, error) {
	status := model.LogPublishedStatus
	count, err := d.db.Omit("id").
		Where("id = ?", id).And("status != ?", model.DeletedStatus).
		Update(&model.Log{Status: &status})
	if err != nil {
		log.Error("发布日志时出现错误", zap.String("err", err.Error()))
		return 0, ecode.InternalError
	}

	return count, nil
}

func (d *Dao) UpdateLogStatus(status int8, id int64) (int64, error) {
	count, err := d.db.Omit("id").
		Where("id = ?", id).
		And("status != ?", model.DeletedStatus).
		Update(&model.Log{Status: &status})
	if err != nil {
		log.Error("修改日志时出现错误", zap.String("err", err.Error()))
		return 0, err
	}

	return count, nil
}

func (d *Dao) UpdateLogStatusBatch(status int8, id ...int64) (int64, error) {
	count, err := d.db.Omit("id").
		In("id", id).
		And("status != ?", model.DeletedStatus).
		Update(&model.Log{Status: &status})
	if err != nil {
		log.Error("修改日志时出现错误",
			zap.Any("id", id),
			zap.Int8("status", status),
			zap.String("err", err.Error()),
		)
		return 0, err
	}

	return count, nil
}

func (d *Dao) DeleteLog(id int64) error {
	status := model.DeletedStatus
	_, err := d.db.Omit("id").
		Where("id = ?", id).
		Update(&model.Log{Status: &status})
	if err != nil {
		log.Error("删除日志时出现错误", zap.Any("id", id), zap.String("err", err.Error()))
		return ecode.InternalError
	}

	return nil
}
