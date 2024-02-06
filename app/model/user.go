package model

import "time"

const (
	SuperAdminGroup  = 1
	NormalAdminGroup = 2
)

type AdminUser struct {
	ID         int64      `xorm:"id" db:"id" json:"id" form:"id"`
	Username   *string    `xorm:"username" db:"username" json:"username" form:"username"`
	Password   *string    `xorm:"password" db:"password" json:"password" form:"password"`
	Group      *int8      `xorm:"group" db:"group" json:"group" form:"group"`
	CreateTime *time.Time `xorm:"create_time" db:"create_time" json:"create_time" form:"create_time"`
	UpdateTime *time.Time `xorm:"update_time" db:"update_time" json:"update_time" form:"update_time"`
	Status     *int8      `xorm:"status" db:"status" json:"status" form:"status"`
}

func (AdminUser) TableName() string {
	return "admin_user"
}
