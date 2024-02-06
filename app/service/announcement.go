package service

import (
	"github.com/yitter/idgenerator-go/idgen"
	"time"
	"wusthelper-manager-go/app/model"
	"wusthelper-manager-go/common"
)

type AnnouncementAddParam struct {
	Title    *string
	Content  *string
	Target   *string
	Platform *[]string
}

type AnnouncementModifyParam struct {
	Id       int64
	Title    *string
	Content  *string
	Target   *string
	Platform *string
	Status   *int8
}

func (s *Service) GetPublishedAnnouncement(platform string) (*[]model.Announcement, error) {
	announcements, err := s.dao.GetPublishedAnnouncement(platform)
	if err != nil {
		return nil, err
	}

	return announcements, nil
}

func (s *Service) GetAllAnnouncement(paging common.Pagination, platform string) (*[]model.Announcement, int64, error) {
	announcements, total, err := s.dao.GetAnnouncementList(paging, platform)
	if err != nil {
		return nil, 0, err
	}

	return announcements, total, nil
}

func (s *Service) PublishAnnouncement(id int64) error {
	announcement := model.Announcement{
		Id:         id,
		Status:     new(int8),
		UpdateTime: new(time.Time),
	}

	*announcement.Status = model.AnnouncementPublishedStatus
	*announcement.UpdateTime = time.Now()

	_, err := s.dao.UpdateAnnouncement(&announcement)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) PublishAnnouncementBatch(ids []int64) error {
	_, err := s.dao.UpdateAnnouncementStatusBatch(ids, model.AnnouncementPublishedStatus)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) AddAnnouncement(param *AnnouncementAddParam) error {
	for _, platform := range *param.Platform {
		announcement := model.Announcement{
			Id:         idgen.NextId(),
			Title:      param.Title,
			Content:    param.Content,
			Target:     param.Target,
			Platform:   &platform,
			CreateTime: new(time.Time),
			UpdateTime: new(time.Time),
			Status:     new(int8),
		}

		*announcement.CreateTime = time.Now()
		*announcement.UpdateTime = time.Now()
		*announcement.Status = model.AnnouncementNotPublishedStatus

		_, err := s.dao.AddAnnouncement(&announcement)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) DeleteAnnouncement(id int64) error {
	err := s.dao.DeleteAnnouncement(id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) ModifyAnnouncement(param *AnnouncementModifyParam) error {
	announcement := model.Announcement{
		Id:         param.Id,
		Title:      param.Title,
		Content:    param.Content,
		Target:     param.Target,
		Platform:   param.Platform,
		Status:     param.Status,
		UpdateTime: new(time.Time),
	}

	*announcement.UpdateTime = time.Now()

	_, err := s.dao.UpdateAnnouncement(&announcement)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) ModifyAnnouncementBatch(params ...AnnouncementModifyParam) error {
	for _, param := range params {
		now := time.Now()
		announcement := model.Announcement{
			Id:         param.Id,
			Title:      param.Title,
			Content:    param.Content,
			Target:     param.Target,
			Platform:   param.Platform,
			UpdateTime: &now,
		}

		_, err := s.dao.UpdateAnnouncement(&announcement)
		if err != nil {
			return err
		}
	}

	return nil
}
