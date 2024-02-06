package model

import "time"

const (
	ConfigValueTypeString = 0
	ConfigValueTypeBool   = 2
)

type Config struct {
	ID             int64      `xorm:"id" db:"id" json:"id" form:"id"`
	Name           *string    `xorm:"name" db:"name" json:"name" form:"name"`
	Value          *string    `xorm:"value" db:"value" json:"value" form:"value"`
	PossibleValues *[]string  `xorm:"possible_values" db:"possible_values" json:"possible_values" form:"possible_values"`
	Type           *int8      `xorm:"type" db:"type" json:"type" form:"type"`                 //  0代表输入框 1代表选择框 2代表switch开关
	Describe       *string    `xorm:"describe" db:"describe" json:"describe" form:"describe"` //  描述
	Platform       *string    `xorm:"platform" db:"platform" json:"platform" form:"platform"`
	CreateTime     *time.Time `xorm:"create_time" db:"create_time" json:"create_time" form:"create_time"`
	UpdateTime     *time.Time `xorm:"update_time" db:"update_time" json:"update_time" form:"update_time"`
	Status         *int8      `xorm:"status" db:"status" json:"status" form:"status"`
}

func (Config) TableName() string {
	return "config"
}
