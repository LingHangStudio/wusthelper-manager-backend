package dao

import (
	goredis "github.com/redis/go-redis/v9"
	"wusthelper-manager-go/app/conf"
	"wusthelper-manager-go/library/cache/redis"
	"wusthelper-manager-go/library/database"
	"wusthelper-manager-go/library/log"
	"xorm.io/xorm"
)

type Dao struct {
	db    *xorm.Engine
	redis *goredis.Client
}

func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		db:    database.NewMysql(&c.Database),
		redis: redis.NewRedisClient(&c.Redis),
	}

	return
}

func (d *Dao) Close() {
	dbErr := d.db.Close()
	if dbErr != nil {
		log.Warn("[dao]关闭数据库连接出错")
	}

	redisErr := d.redis.Close()
	if redisErr != nil {
		log.Warn("[dao]关闭redis连接出错")
	}
}
