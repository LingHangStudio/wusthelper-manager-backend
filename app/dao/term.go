package dao

import (
	"go.uber.org/zap"
	"wusthelper-manager-go/app/model"
	"wusthelper-manager-go/library/ecode"
	"wusthelper-manager-go/library/log"
)

func (d *Dao) GetTermList() (*[]model.Term, error) {
	result := make([]model.Term, 0)
	err := d.db.Where("status != ?", model.DeletedStatus).Find(&result)
	if err != nil {
		log.Error("获取学期条目列表时出现错误", zap.String("err", err.Error()))
		return nil, ecode.InternalError
	}

	return &result, nil
}

func (d *Dao) AddTerm(term *model.Term) (int64, error) {
	count, err := d.db.InsertOne(term)
	if err != nil {
		log.Error("添加学期条目时出现错误", zap.String("err", err.Error()))
		return 0, ecode.InternalError
	}

	return count, nil
}

func (d *Dao) UpdateTerm(term *model.Term) (int64, error) {
	count, err := d.db.Omit("id").
		Where("id = ?", term.ID).And("status != ?", model.DeletedStatus).
		Update(term)
	if err != nil {
		log.Error("修改学期条目时出现错误", zap.String("err", err.Error()))
		return 0, err
	}

	return count, nil
}

func (d *Dao) DeleteTerm(id int64) error {
	status := model.DeletedStatus
	_, err := d.db.Where("id = ?", id).Update(&model.Term{Status: &status})
	if err != nil {
		log.Error("删除学期条目时出现错误", zap.String("err", err.Error()))
		return ecode.InternalError
	}

	return nil
}
