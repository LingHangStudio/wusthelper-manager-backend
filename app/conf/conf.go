package conf

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/viper"
	"time"
	"wusthelper-manager-go/library/cache/redis"
	"wusthelper-manager-go/library/database"
)

const (
	DevEnv  = "dev"
	ProdEnv = "prod"
)

var (
	Conf = &Config{}
)

type Config struct {
	Server     ServerConf
	Wusthelper WusthelperConf
	Database   database.Config
	Redis      redis.Config
}

type ServerConf struct {
	Env          string
	Port         int
	Address      string
	BaseUrl      string
	TokenSecret  string
	TokenTimeout time.Duration
	LogLocation  string

	FileStorageOption FileStorageOption
}

type FileStorageOption struct {
	UploadFileLocalTmpPath string

	ResourceStorageOption ResourceStorageOption
	AliyunOssOption       AliyunOssOption
}

type ResourceStorageOption struct {
	WusthelperReleaseFileKey string

	VersionFileStorageBasePath string
	PicStorageBasePath         string
	DefaultPicUrl              string
}

type AliyunOssOption struct {
	AccessKeyId      string
	AccessKeySecret  string
	Endpoint         string
	Bucket           string
	BucketBindDomain string
}

type WusthelperConf struct {
	Upstream     string
	Timeout      time.Duration
	Proxy        string
	TokenKey     string
	AdminBaseUrl string
}

func Init() (err error) {
	viper.AddConfigPath(".")
	viper.AddConfigPath("./conf")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/wusthelper-manager")
	viper.AddConfigPath("$HOME/.wusthelper-manager")

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(Conf)
	if err != nil {
		return
	}

	if Conf.Server.Env == DevEnv {
		jsonByte, _ := jsoniter.Marshal(Conf)
		fmt.Println(string(jsonByte))
	}

	return
}
