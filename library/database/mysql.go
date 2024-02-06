package database

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"wusthelper-manager-go/library/log"
	"xorm.io/xorm"
)

const _dsnTemplate = "%s:%s@tcp(%s)/%s"

func NewMysql(c *Config) *xorm.Engine {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", c.Username, c.Password, c.Host, c.Port, c.DbName)
	engine, err := xorm.NewEngineWithParams("mysql", dsn, c.ExtraParams)
	if err != nil {
		log.Panic("[Database-Mysql]: 数据库连接失败", zap.String("errMsg", err.Error()))
	}
	engine.ShowSQL(c.PrintSql)

	return engine
}
