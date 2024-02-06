package dao

import (
	"go.uber.org/zap"
	"wusthelper-manager-go/app/model"
	"wusthelper-manager-go/library/ecode"
	"wusthelper-manager-go/library/log"
)

func (d *Dao) GetPlatformList() (*[]string, error) {
	result := make([]string, 0)
	err := d.db.
		Table("config").
		Where("status != ?", model.DeletedStatus).
		GroupBy("platform").
		Find(&result)

	if err != nil {
		log.Error("获取平台列表时出现错误", zap.Error(err))
		return nil, err
	}

	return &result, nil
}

func (d *Dao) GetConfigList(platform string) (*[]model.Config, int64, error) {
	result := make([]model.Config, 0)
	countSession := d.db.Where("status != ?", model.DeletedStatus)
	if platform != "" {
		countSession.In("platform", platform)
	}

	total, err := countSession.Count(&model.Config{})
	if err != nil {
		log.Error("获取配置条目数量时出现错误", zap.Any("platform", platform), zap.String("err", err.Error()))
		return nil, 0, ecode.InternalError
	}

	querySession := d.db.Where("status != ?", model.DeletedStatus)
	if platform != "" {
		querySession.In("platform", platform)
	}

	err = querySession.Find(&result)
	if err != nil {
		log.Error("获取配置条目列表时出现错误", zap.Any("platform", platform), zap.String("err", err.Error()))
		return nil, 0, ecode.InternalError
	}

	return &result, total, nil
}

func (d *Dao) AddConfig(config *model.Config) (int64, error) {
	count, err := d.db.InsertOne(config)
	if err != nil {
		log.Error("添加配置条目时出现错误", zap.String("err", err.Error()))
		return 0, ecode.InternalError
	}

	return count, nil
}

func (d *Dao) AddConfigBatch(config *[]model.Config) (int64, error) {
	count, err := d.db.Insert(config)
	if err != nil {
		log.Error("批量添加配置条目时出现错误", zap.String("err", err.Error()))
		return 0, ecode.InternalError
	}

	return count, nil
}

func (d *Dao) UpdateConfig(config *model.Config) (int64, error) {
	count, err := d.db.Omit("id").
		Where("id = ?", config.ID).And("status != ?", model.DeletedStatus).
		Update(config)
	if err != nil {
		log.Error("修改配置条目时出现错误", zap.Int64("id", config.ID), zap.String("err", err.Error()))
		return 0, err
	}

	return count, nil
}

func (d *Dao) DeleteConfig(id int64) error {
	status := model.DeletedStatus
	_, err := d.db.Where("id = ?", id).Update(&model.Config{Status: &status})
	if err != nil {
		log.Error("删除配置条目时出现错误", zap.Int64("id", id), zap.String("err", err.Error()))
		return ecode.InternalError
	}

	return nil
}
