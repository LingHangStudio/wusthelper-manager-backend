package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yitter/idgenerator-go/idgen"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	v1 "wusthelper-manager-go/app/api/http"
	"wusthelper-manager-go/app/conf"
	"wusthelper-manager-go/app/middleware/auth"
	"wusthelper-manager-go/library/ecode"
	"wusthelper-manager-go/library/log"
)

func main() {
	// config加载必须在最先位置
	loadConfig()

	setupIdGenerator()
	setupLogger()

	ecode.InitEcodeText()

	// server启动必须在最后
	startServer()
}

func loadConfig() {
	if confErr := conf.Init(); confErr != nil {
		panic(confErr)
	}
}

func setupIdGenerator() {
	options := idgen.NewIdGeneratorOptions(0)
	idgen.SetIdGenerator(options)
}

func setupLogger() {
	stdoutLogLevel := log.WarnLevel
	if conf.Conf.Server.Env == conf.DevEnv {
		stdoutLogLevel = log.InfoLevel
	}
	tees := []log.TeeOption{
		{
			Out: os.Stdout,
			LevelEnablerFunc: func(level log.Level) bool {
				return level >= stdoutLogLevel
			},
		},
	}

	logFileLocation := conf.Conf.Server.LogLocation
	if logFileLocation != "" {
		tees = append(tees, log.TeeOption{
			Out: &lumberjack.Logger{
				Filename: logFileLocation,
				MaxSize:  128,
				MaxAge:   30,
				Compress: true,
			},
			LevelEnablerFunc: func(level log.Level) bool {
				return level >= log.InfoLevel
			},
		})
	}

	logger := log.NewTee(tees)
	log.SetDefault(logger)
}

func startServer() {
	if conf.Conf.Server.Env == conf.ProdEnv {
		gin.SetMode(gin.ReleaseMode)
	}

	setupMiddleware()

	engine, err := v1.NewEngine(conf.Conf, conf.Conf.Server.BaseUrl)
	if err != nil {
		panic(err)
	}

	addr := fmt.Sprintf("%s:%d", conf.Conf.Server.Address, conf.Conf.Server.Port)
	err = engine.Run(addr)
	if err != nil {
		panic(err)
	}
}

func setupMiddleware() {
	auth.Init(conf.Conf)
}
