package service

import (
	"github.com/yitter/idgenerator-go/idgen"
	"golang.org/x/crypto/bcrypt"
	"time"
	"wusthelper-manager-go/app/model"
	"wusthelper-manager-go/library/ecode"
)

const (
	DefaultUserGroup = 2
)

type AdminUserAddParam struct {
	Username string
	Password string
}

type AdminUserModifyParam struct {
	Id       int64
	Username *string
	Password *string
	Groupid  *int8
}

func (s *Service) AdminUserLogin(username, password string) (*model.AdminUser, error) {
	user, err := s.dao.GetAdminUserByUsername(username)
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, ecode.UserNotExists
	}

	err = bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(password))
	if err != nil {
		return nil, ecode.AuthFailed
	}

	return user, nil
}

func (s *Service) GetAllAdminUserList() (*[]model.AdminUser, error) {
	adminUserList, err := s.dao.GetAllAdminUser()
	if err != nil {
		return nil, err
	}

	return adminUserList, nil
}

func (s *Service) DeleteAdminUserById(id uint64) error {
	err := s.dao.DeleteAdminUserById(id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) AddAdminUser(param *AdminUserAddParam) error {
	usernameExists, err := s.dao.HasAdminUserName(param.Username)
	if err != nil {
		return err
	}

	if usernameExists {
		return ecode.UsernameExists
	}

	hashedPasswordBytes, _ := bcrypt.GenerateFromPassword([]byte(param.Password), bcrypt.DefaultCost)
	hashedPassword := string(hashedPasswordBytes)

	now := time.Now()
	user := model.AdminUser{
		ID:         idgen.NextId(),
		Username:   &param.Username,
		Password:   &hashedPassword,
		Group:      new(int8),
		Status:     new(int8),
		CreateTime: &now,
		UpdateTime: &now,
	}

	*user.Group = DefaultUserGroup
	*user.Status = model.NormalStatus

	_, err = s.dao.AddAdminUser(&user)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetUserGroup(id uint64) (int8, error) {
	user, err := s.dao.GetAdminUserById(id)
	if err != nil {
		return 0, err
	} else if user == nil {
		return 0, ecode.UserNotExists
	}

	return *user.Group, nil
}

func (s *Service) CheckUsernameExists(username string) (bool, error) {
	exists, err := s.dao.HasAdminUserName(username)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (s *Service) ModifyAdminUser(param *AdminUserModifyParam) error {
	var hashedPassword *string
	if param.Password != nil {
		hashedBytes, _ := bcrypt.GenerateFromPassword([]byte(*param.Password), bcrypt.DefaultCost)
		*hashedPassword = string(hashedBytes)
	}

	newUser := model.AdminUser{
		ID:         param.Id,
		Username:   param.Username,
		Password:   hashedPassword,
		Group:      param.Groupid,
		Status:     new(int8),
		UpdateTime: new(time.Time),
	}

	*newUser.Status = model.NormalStatus
	*newUser.UpdateTime = time.Now()

	_, err := s.dao.UpdateAdminUser(&newUser)
	if err != nil {
		return err
	}

	return nil
}
