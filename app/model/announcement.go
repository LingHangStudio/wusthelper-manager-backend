package model

import "time"

const (
	AnnouncementNotPublishedStatus      = NormalStatus
	AnnouncementDeletedStatus           = DeletedStatus
	AnnouncementPublishedStatus    int8 = 2
)

type Announcement struct {
	Id         int64      `xorm:"id" db:"id" json:"id" form:"id"`                     //  公告id
	Title      *string    `xorm:"title" db:"title" json:"title" form:"title"`         //  公告标题
	Content    *string    `xorm:"content" db:"content" json:"content" form:"content"` //  公告内容
	Target     *string    `xorm:"target" db:"target" json:"target" form:"target"`     //  发布对象学院
	Platform   *string    `xorm:"platform" db:"platform" json:"platform" form:"platform"`
	CreateTime *time.Time `xorm:"create_time" db:"create_time" json:"create_time" form:"create_time"` //  发布时间
	UpdateTime *time.Time `xorm:"update_time" db:"update_time" json:"update_time" form:"update_time"` //  更新时间
	Status     *int8      `xorm:"status" db:"status" json:"status" form:"status"`                     //  发布状态 0是未发布，1是发布
}

func (Announcement) TableName() string {
	return "announcement"
}
