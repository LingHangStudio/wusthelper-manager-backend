package dao

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
	"wusthelper-manager-go/library/ecode"
	"wusthelper-manager-go/library/log"
)

const (
	_tokenCacheKey = "wusthelper-mp:token:%s"

	_oidSidCacheKey = "wusthelper-mp:oid:%s:sid"
	_sidOidCacheKey = "wusthelper-mp:sid:%s:oid"

	_totalUserCacheKey = "wusthelper-mp:user:total"

	_adminConfigCacheKey = "wusthelper-mp:admin:config"
)

func (d *Dao) StoreWusthelperTokenCache(c *context.Context, token, oid string, ex time.Duration) error {
	key := fmt.Sprintf(_tokenCacheKey, oid)
	err := d.redis.Set(*c, key, token, ex).Err()
	if err != nil {
		log.Error("缓存助手token出现错误", zap.String("err", err.Error()))
		return ecode.InternalError
	}

	return nil
}

func (d *Dao) StoreOidSidCache(c *context.Context, oid, sid string, ex time.Duration) error {
	oidSidCacheKey := fmt.Sprintf(_oidSidCacheKey, oid)
	sidOidCacheKey := fmt.Sprintf(_sidOidCacheKey, sid)

	err := d.redis.Set(*c, oidSidCacheKey, sid, ex).Err()
	err = d.redis.Set(*c, sidOidCacheKey, oid, ex).Err()
	if err != nil {
		log.Error("缓存Oid-Sid出现错误", zap.String("err", err.Error()))
		return ecode.InternalError
	}

	return nil
}

func (d *Dao) GetToken(c *context.Context, oid string) (token string, err error) {
	key := fmt.Sprintf(_tokenCacheKey, oid)
	token, err = d.redis.Get(*c, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		log.Error("获取tokenCache出现错误", zap.String("err", err.Error()))
		return "", ecode.InternalError
	}

	return
}

func (d *Dao) GetSidForOid(c *context.Context, oid string) (sid string, err error) {
	oidSidCacheKey := fmt.Sprintf(_oidSidCacheKey, oid)
	sid, err = d.redis.Get(*c, oidSidCacheKey).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		log.Error("获取oid-sid出现错误", zap.String("err", err.Error()))
		return "", ecode.InternalError
	}

	return sid, nil
}

func (d *Dao) GetTotalUserCountCache(ctx *context.Context) (total int64, err error) {
	total, err = d.redis.Get(*ctx, _totalUserCacheKey).Int64()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		log.Error("获取用户总数缓存出现错误", zap.String("err", err.Error()))
		return 0, ecode.InternalError
	}

	return total, nil
}

func (d *Dao) StoreTotalUserCountCache(ctx *context.Context, count int64, ex time.Duration) (err error) {
	err = d.redis.Set(*ctx, _totalUserCacheKey, count, ex).Err()
	if err != nil {
		log.Error("存储用户总数出现错误", zap.String("err", err.Error()))
		return ecode.InternalError
	}

	return nil
}

func (d *Dao) IncreaseTotalUserCount(ctx *context.Context) (err error) {
	err = d.redis.Incr(*ctx, _totalUserCacheKey).Err()
	if err != nil {
		log.Error("用户总数+1出现错误", zap.String("err", err.Error()))
		return ecode.InternalError
	}

	return nil
}
