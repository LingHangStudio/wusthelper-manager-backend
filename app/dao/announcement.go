package dao

import (
	"go.uber.org/zap"
	"wusthelper-manager-go/app/model"
	"wusthelper-manager-go/common"
	"wusthelper-manager-go/library/ecode"
	"wusthelper-manager-go/library/log"
)

const (
	_announcementTableName = "announcement"

	_hasAnnouncementSql = "select 1 from `announcement` where `id` = ? and `status` != 1"

	_getAnnouncementSql    = "select * from `announcement` where `id` = ? and `status` != 1"
	_getAllAnnouncementSql = "select * from `announcement` where `status` != 1 and platform = ? limit 2048"

	_getPublishedAnnouncementSql = "select * from `announcement` where `status` = 2 limit 2048"
	_deleteAnnouncementSql       = "update `announcement` set `status` = 1 where `id` = ?"
)

func (d *Dao) HasAnnouncement(id string) (bool, error) {
	result, err := d.db.Table(_announcementTableName).
		Where("id = ?", id).And("status != ?", model.DeletedStatus).
		Exist()

	if err != nil {
		log.Error("查询公告是否存在时出现错误", zap.String("err", err.Error()))
		return false, ecode.InternalError
	}

	return result, nil
}

func (d *Dao) GetAnnouncement(id int64) (*model.Announcement, error) {
	announcement := new(model.Announcement)
	has, err := d.db.
		Where("id = ?", id).And("status != ?", model.DeletedStatus).
		Get(announcement)
	if err != nil {
		log.Error("从id获取公告出现错误", zap.Error(err))
		return nil, ecode.InternalError
	}

	if !has {
		return nil, nil
	} else {
		return announcement, nil
	}
}

func (d *Dao) GetAnnouncementList(paging common.Pagination, platform string) (*[]model.Announcement, int64, error) {
	countSession := d.db.Table(_announcementTableName).Where("status != ?", model.DeletedStatus)
	if platform != "" {
		countSession.And("platform = ?", platform)
	}

	total, err := countSession.Count()
	if err != nil {
		if err != nil {
			log.Error("获取公告总数出现错误", zap.Error(err))
			return nil, 0, ecode.InternalError
		}
	}

	announcement := make([]model.Announcement, 0)
	querySession := d.db.Where("status != ?", model.DeletedStatus)
	if platform != "" {
		querySession.And("platform = ?", platform)
	}

	err = querySession.Desc("status", "id").Asc("platform").
		Limit(paging.PageSize, paging.PageSize*(paging.Page-1)).
		Find(&announcement)
	if err != nil {
		log.Error("获取公告列表出现错误", zap.Error(err))
		return nil, 0, ecode.InternalError
	}

	return &announcement, total, nil
}

func (d *Dao) GetPublishedAnnouncement(platform string) (*[]model.Announcement, error) {
	announcementList := make([]model.Announcement, 0)
	session := d.db.Where("status = ?", model.AnnouncementPublishedStatus)
	if platform != "" {
		session.And("platform = ?", platform)
	}

	err := session.Asc("id").Find(&announcementList)
	if err != nil {
		log.Error("获取所有公告时出现错误", zap.Error(err))
		return nil, ecode.InternalError
	}

	return &announcementList, nil
}

func (d *Dao) DeleteAnnouncement(id int64) error {
	_, err := d.db.Exec(_deleteAnnouncementSql, id)
	if err != nil {
		log.Error("删除公告时出现错误", zap.Error(err))
		return ecode.InternalError
	}

	return nil
}

func (d *Dao) AddAnnouncement(announcement *model.Announcement) (int64, error) {
	result, err := d.db.InsertOne(announcement)
	if err != nil {
		log.Error("添加公告时出现错误", zap.Error(err))
		return 0, ecode.InternalError
	}

	return result, nil
}

func (d *Dao) UpdateAnnouncement(announcement *model.Announcement) (int64, error) {
	result, err := d.db.
		Where("id = ?", announcement.Id).And("status != ?", model.AnnouncementDeletedStatus).
		Update(announcement)
	if err != nil {
		log.Error("更新公告时出现错误", zap.Error(err))
		return 0, err
	}

	return result, nil
}

func (d *Dao) UpdateAnnouncementStatusBatch(ids []int64, status int8) (int64, error) {
	result, err := d.db.
		In("id", ids).And("status != ?", model.AnnouncementDeletedStatus).
		MustCols("status").
		Update(&model.Announcement{Status: &status})
	if err != nil {
		log.Error("更新公告时出现错误", zap.Error(err))
		return 0, err
	}

	return result, nil
}
