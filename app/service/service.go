package service

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"os"
	"time"
	"wusthelper-manager-go/app/conf"
	"wusthelper-manager-go/app/dao"
	"wusthelper-manager-go/app/rpc/http/wusthelper/v3"
	"wusthelper-manager-go/library/log"
)

const (
	wusthelperTokenExpiration = time.Hour * 24
)

type Service struct {
	config    *conf.Config
	dao       *dao.Dao
	ossBucket *oss.Bucket
	rpc       *v3.WusthelperHttpRpc
}

func New(c *conf.Config) (*Service, error) {
	service := &Service{
		config: c,
		dao:    dao.New(c),
		rpc:    v3.NewRpcClient(&c.Wusthelper),
	}

	uploadFileLocalTmpPath := c.Server.FileStorageOption.UploadFileLocalTmpPath
	if err := os.MkdirAll(uploadFileLocalTmpPath, 0660); err != nil {
		return nil, fmt.Errorf("初始化临时文件目录失败：%s", err.Error())
	}

	aliyunOssOption := c.Server.FileStorageOption.AliyunOssOption
	client, err := oss.New(aliyunOssOption.Endpoint, aliyunOssOption.AccessKeyId, aliyunOssOption.AccessKeySecret)
	if err != nil {
		log.Warn("阿里云oss客户端初始化失败")
		return nil, fmt.Errorf("阿里云oss客户端初始化失败：%s", err.Error())
	}

	service.ossBucket, err = client.Bucket(aliyunOssOption.Bucket)
	if err != nil {
		log.Warn("阿里云oss bucket初始化失败")
		return nil, fmt.Errorf("阿里云oss bucket初始化失败，bucket: %s，err: %s", aliyunOssOption.Bucket, err.Error())
	}

	return service, nil
}
