package model

import "time"

const (
	BannerPublishedStatus int8 = 2
)

type Banner struct {
	ID         int64      `xorm:"id" db:"id" json:"id" form:"id"`
	Title      *string    `xorm:"title" db:"title" json:"title" form:"title"`
	Link       *string    `xorm:"link" db:"link" json:"link" form:"link"`
	Img        *string    `xorm:"img" db:"img" json:"img" form:"img"`
	Platform   *string    `xorm:"platform" db:"platform" json:"platform" form:"platform"`
	CreateTime *time.Time `xorm:"create_time" db:"create_time" json:"create_time" form:"create_time"`
	UpdateTime *time.Time `xorm:"update_time" db:"update_time" json:"update_time" form:"update_time"`
	Status     *int8      `xorm:"status" db:"status" json:"status" form:"status"`
}

func (Banner) TableName() string {
	return "banner"
}
