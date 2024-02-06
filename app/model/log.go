package model

import "time"

const (
	LogPublishedStatus int8 = 2
)

type Log struct {
	ID          int64      `xorm:"id" db:"id" json:"id" form:"id"`
	Title       *string    `xorm:"title" db:"title" json:"title" form:"title"`
	Content     *string    `xorm:"content" db:"content" json:"content" form:"content"`
	VersionText *string    `xorm:"version_text" db:"version_text" json:"version_text" form:"version_text"`
	Platform    *string    `xorm:"platform" db:"platform" json:"platform" form:"platform"`
	CreateTime  *time.Time `xorm:"create_time" db:"create_time" json:"create_time" form:"create_time"`
	UpdateTime  *time.Time `xorm:"update_time" db:"update_time" json:"update_time" form:"update_time"`
	Status      *int8      `xorm:"status" db:"status" json:"status" form:"status"`
}

func (Log) TableName() string {
	return "log"
}
