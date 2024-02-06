package service

import (
	"github.com/yitter/idgenerator-go/idgen"
	"strings"
	"time"
	"wusthelper-manager-go/app/model"
)

func (s *Service) GetPlatformList() (*[]string, error) {
	platforms, err := s.dao.GetPlatformList()
	if err != nil {
		return nil, err
	}

	return platforms, nil
}

func (s *Service) GetConfigList(platform string) (*[]model.Config, int64, error) {
	configList, total, err := s.dao.GetConfigList(platform)
	if err != nil {
		return nil, 0, err
	}

	return configList, total, nil
}

type ConfigAddParam struct {
	Name           string
	Value          string
	Type           int8
	Describe       string
	Platform       []string
	PossibleValues []string
}

func (s *Service) AddConfig(param *ConfigAddParam) error {
	configList := make([]model.Config, len(param.Platform))
	now := time.Now()
	status := model.NormalStatus
	for i, platform := range param.Platform {
		p := strings.Clone(platform)
		configList[i] = model.Config{
			ID:             idgen.NextId(),
			Name:           &param.Name,
			Value:          &param.Value,
			Type:           &param.Type,
			Describe:       &param.Describe,
			PossibleValues: &param.PossibleValues,
			Platform:       &p,
			CreateTime:     &now,
			UpdateTime:     &now,
			Status:         &status,
		}
	}

	_, err := s.dao.AddConfigBatch(&configList)
	if err != nil {
		return err
	}

	return nil
}

type ConfigModifyParam struct {
	Id             int64
	Name           *string
	Value          *string
	Describe       *string
	PossibleValues *[]string
}

func (s *Service) ModifyConfig(param *ConfigModifyParam) error {
	now := time.Now()
	config := model.Config{
		ID:             param.Id,
		Name:           param.Name,
		Value:          param.Value,
		Describe:       param.Describe,
		PossibleValues: param.PossibleValues,
		UpdateTime:     &now,
	}

	_, err := s.dao.UpdateConfig(&config)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteConfig(id int64) error {
	err := s.dao.DeleteConfig(id)
	if err != nil {
		return err
	}

	return nil
}
