package model

import "time"

const (
	VersionPublishedStatus int8 = 2
)

type Version struct {
	ID          int64      `xorm:"id"`
	VersionText *string    `xorm:"version_text"`
	Summary     *string    `xorm:"summary"`
	File        *string    `xorm:"file"`
	Platform    *string    `xorm:"platform"`
	CreateTime  *time.Time `xorm:"create_time"`
	UpdateTime  *time.Time `xorm:"update_time"`
	Status      *int8      `xorm:"status"`
}

func (Version) TableName() string {
	return "version"
}
