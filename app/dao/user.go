package dao

import (
	"go.uber.org/zap"
	"wusthelper-manager-go/app/model"
	"wusthelper-manager-go/library/ecode"
	"wusthelper-manager-go/library/log"
)

const (
	_hasAdminUserSql = "select 1 from `admin_user` where `username` = ? and `status` = 0"

	_getAdminUserByUsernameSql = "select * from `admin_user` where `username` = ? and `status` = 0"
	_getAdminUserByIdSql       = "select * from `admin_user` where `id` = ? and `status` = 0"
	_getAllAdminUserSql        = "select * from `admin_user` where `status` = 0 limit 2048"

	_deleteAdminUserByUsernameSql = "update `admin_user` set `status` = ? where `username` = ?"
	_deleteAdminUserByIdSql       = "update `admin_user` set `status` = ? where `id` = ?"
)

func (d *Dao) HasAdminUserName(username string) (bool, error) {
	result, err := d.db.SQL(_hasAdminUserSql, username).Exist()

	if err != nil {
		log.Error("查询管理员用户是否存在时出现错误", zap.String("err", err.Error()))
		return false, ecode.InternalError
	}

	return result, nil
}

func (d *Dao) GetAdminUserByUsername(username string) (*model.AdminUser, error) {
	user := new(model.AdminUser)
	has, err := d.db.Where("username = ?", username).And("status != ?", model.DeletedStatus).Get(user)
	if err != nil {
		log.Error("从用户名获取管理员用户出现错误", zap.Error(err))
		return nil, ecode.InternalError
	}

	if !has {
		return nil, nil
	} else {
		return user, nil
	}
}

func (d *Dao) GetAdminUserById(id uint64) (*model.AdminUser, error) {
	user := new(model.AdminUser)
	has, err := d.db.SQL(_getAdminUserByIdSql, id).Get(user)
	if err != nil {
		log.Error("从id获取管理员用户出现错误", zap.Error(err))
		return nil, ecode.InternalError
	}

	if !has {
		return nil, nil
	} else {
		return user, nil
	}
}

func (d *Dao) GetAllAdminUser() (*[]model.AdminUser, error) {
	userList := make([]model.AdminUser, 0)
	err := d.db.SQL(_getAllAdminUserSql).Find(&userList)
	if err != nil {
		log.Error("获取所有管理员用户时出现错误", zap.Error(err))
		return nil, ecode.InternalError
	}

	return &userList, nil
}

func (d *Dao) DeleteAdminUserByUsername(username string) error {
	_, err := d.db.Exec(_deleteAdminUserByUsernameSql, model.DeletedStatus, username)
	if err != nil {
		log.Error("删除管理员用户时出现错误", zap.Error(err))
		return ecode.InternalError
	}

	return nil
}

func (d *Dao) DeleteAdminUserById(id uint64) error {
	_, err := d.db.Exec(_deleteAdminUserByIdSql, model.DeletedStatus, id)
	if err != nil {
		log.Error("删除管理员用户时出现错误", zap.Error(err))
		return ecode.InternalError
	}

	return nil
}

func (d *Dao) AddAdminUser(user *model.AdminUser) (int64, error) {
	result, err := d.db.InsertOne(user)
	if err != nil {
		log.Error("添加管理员用户时出现错误", zap.Error(err))
		return 0, ecode.InternalError
	}

	return result, nil
}

func (d *Dao) UpdateAdminUser(user *model.AdminUser) (int64, error) {
	result, err := d.db.
		Where("id = ?", user.ID).
		And("status != ?", model.DeletedStatus).
		Update(user)
	if err != nil {
		log.Error("更新管理员用户时出现错误", zap.Error(err))
		return 0, err
	}

	return result, nil
}
