package service

import (
	"github.com/yitter/idgenerator-go/idgen"
	"strings"
	"time"
	"wusthelper-manager-go/app/model"
	"wusthelper-manager-go/common"
)

func (s *Service) GetPublishedLog(platform ...string) (*[]model.Log, error) {
	latestLogList, err := s.dao.GetPublishedLog(platform...)
	if err != nil {
		return nil, err
	}

	return latestLogList, nil
}

func (s *Service) GetLogList(pagination common.Pagination, platform string) (*[]model.Log, int64, error) {
	logList, total, err := s.dao.GetLogList(pagination, platform)
	if err != nil {
		return nil, 0, err
	}

	return logList, total, nil
}

type LogAddParam struct {
	Title       string
	Content     string
	VersionText string
	Platform    []string
}

func (s *Service) AddLog(param *LogAddParam) error {
	now := time.Now()
	logs := make([]model.Log, len(param.Platform))
	for i, platform := range param.Platform {
		status := model.NormalStatus
		p := strings.Clone(platform)
		logs[i] = model.Log{
			ID:          idgen.NextId(),
			Title:       &param.Title,
			Content:     &param.Content,
			VersionText: &param.VersionText,
			Platform:    &p,
			CreateTime:  &now,
			UpdateTime:  &now,
			Status:      &status,
		}
	}

	_, err := s.dao.AddLog(logs...)
	if err != nil {
		return err
	}

	return nil
}

type LogModifyParam struct {
	Id          int64
	Title       *string
	Content     *string
	VersionText *string
	Platform    *[]string
	Status      *int8
}

func (s *Service) ModifyLog(param *LogModifyParam) error {
	now := time.Now()
	logEntity := model.Log{
		ID:          param.Id,
		Title:       param.Title,
		Content:     param.Content,
		VersionText: param.VersionText,
		UpdateTime:  &now,
		Status:      param.Status,
	}

	*logEntity.UpdateTime = time.Now()
	_, err := s.dao.UpdateLog(&logEntity)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteLog(id int64) error {
	err := s.dao.DeleteLog(id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) PublishLog(id int64) error {
	_, err := s.dao.UpdateLogStatus(model.LogPublishedStatus, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) PublishLogBatch(id ...int64) error {
	_, err := s.dao.UpdateLogStatusBatch(model.LogPublishedStatus, id...)
	if err != nil {
		return err
	}

	return nil
}
