package service

import (
	"github.com/yitter/idgenerator-go/idgen"
	"time"
	"wusthelper-manager-go/app/model"
)

func (s *Service) GetTermList() (*[]model.Term, error) {
	termList, err := s.dao.GetTermList()
	if err != nil {
		return nil, err
	}

	return termList, nil
}

type TermAddParam struct {
	Term  string
	Start time.Time
}

func (s *Service) AddTerm(param *TermAddParam) error {
	now := time.Now()
	status := model.NormalStatus
	term := model.Term{
		ID:         idgen.NextId(),
		Term:       &param.Term,
		Start:      &param.Start,
		CreateTime: &now,
		UpdateTime: &now,
		Status:     &status,
	}

	_, err := s.dao.AddTerm(&term)
	if err != nil {
		return err
	}

	return nil
}

type TermModifyParam struct {
	ID    int64
	Term  *string
	Start *time.Time
}

func (s *Service) ModifyTerm(param *TermModifyParam) error {
	now := time.Now()
	term := model.Term{
		ID:         param.ID,
		Term:       param.Term,
		Start:      param.Start,
		UpdateTime: &now,
	}

	_, err := s.dao.UpdateTerm(&term)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteTerm(id int64) error {
	err := s.dao.DeleteTerm(id)
	if err != nil {
		return err
	}

	return nil
}
