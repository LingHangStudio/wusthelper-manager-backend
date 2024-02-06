package model

import "time"

type Term struct {
	ID         int64      `xorm:"id" db:"id" json:"id" form:"id"`
	Term       *string    `xorm:"term" db:"term" json:"term" form:"term"`
	Start      *time.Time `xorm:"start" db:"start" json:"start" form:"start"`
	CreateTime *time.Time `xorm:"create_time" db:"create_time" json:"create_time" form:"create_time"`
	UpdateTime *time.Time `xorm:"update_time" db:"update_time" json:"update_time" form:"update_time"`
	Status     *int8      `xorm:"status" db:"status" json:"status" form:"status"`
}

func (Term) TableName() string {
	return "terms"
}
