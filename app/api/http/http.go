package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
	"wusthelper-manager-go/app/conf"
	"wusthelper-manager-go/app/middleware"
	"wusthelper-manager-go/app/middleware/auth"
	"wusthelper-manager-go/app/service"
	"wusthelper-manager-go/library/token"
)

var (
	config *conf.Config
	srv    *service.Service
	jwt    *token.Token
)

func NewEngine(c *conf.Config, baseUrl string) (*gin.Engine, error) {
	config = c
	engine := gin.Default()
	//中间件在路由配置开始前才生效
	engine.Use(middleware.GlobalPanicRecover)
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"*"}
	corsConfig.AllowMethods = []string{"*"}
	corsConfig.AllowPrivateNetwork = true
	corsConfig.MaxAge = time.Second
	engine.Use(cors.New(corsConfig))

	rootRouter := engine.RouterGroup.Group(baseUrl)
	//rootRouter.Use(gin.LoggerWithWriter(*log.DefaultWriter().))

	setupRouter(rootRouter)

	var err error
	srv, err = service.New(c)
	if err != nil {
		return nil, err
	}

	jwt = token.New(c.Server.TokenSecret, c.Server.TokenTimeout)

	return engine, nil
}

// setupRouter api路由设置，详见api文档
func setupRouter(rootRouter *gin.RouterGroup) {
	setupAdminRouter(rootRouter)
	setupPublicApiRouter(rootRouter)
}

func setupAdminRouter(rootRouter *gin.RouterGroup) {
	admin := rootRouter.Group("/admin")
	{
		// 活动（轮播图）管理端相关路由
		banner := admin.Group("/act", auth.AdminUserTokenCheck)
		{
			banner.PUT("/addActAndFile", addBanner)   // 添加活动
			banner.DELETE("/deleteAct", deleteBanner) // 删除活动
			banner.PATCH("/chAct", modifyBanner)      // 修改活动
			banner.GET("/getActs", getBannerList)     // 查询活动
			banner.POST("/publishAct", publishBanner) // 发布活动
		}

		// 管理员日志相关
		adminLog := admin.Group("/admin_log", auth.AdminUserTokenCheck)
		{
			adminLog.GET("/get")
			adminLog.PUT("/add")
			adminLog.DELETE("/delete")
		}

		// 日志相关
		mainLog := admin.Group("/log", auth.AdminUserTokenCheck)
		{
			mainLog.PUT("/addLog", addLog)
			mainLog.GET("/getLog", getLogList)
			mainLog.PATCH("/chLog", modifyLog)
			mainLog.DELETE("/deleteLog", deleteLog)
			mainLog.GET("/getVersion", func(context *gin.Context) {
				responseData(context, map[string]string{"version": "1.0.0", "time": "2021-02-19 17:25:30"})
			})
			mainLog.POST("/publishLog", publishLog)
		}

		configRouter := admin.Group("/config", auth.AdminUserTokenCheck)
		{
			configRouter.GET("/getAllConfig", getConfigList)
			configRouter.PUT("/addConfig", addConfig)
			configRouter.PATCH("/chConfig", modifyConfig)
			configRouter.DELETE("/deleteConfig", deleteConfig)
			configRouter.GET("/getAllPlatform", getPlatformList)
		}

		user := admin.Group("/data", auth.AdminUserTokenCheck)
		{
			user.GET("/getAllUser")    // 通过学院和专业来查询学生
			user.GET("/getOne")        // 通过姓名和学号来查询学生
			user.GET("/getCollegeNum") // 查询学院及其专业信息
			user.GET("/getAddUser")    // 查询增长的用户数
			user.GET("/getNum")        // 查询用户总数
		}

		adminUser := admin.Group("/user")
		{
			adminUser.POST("/login", adminUserLogin)
			adminUser.GET("/getAllAdmin", auth.AdminUserTokenCheck, getAdminUserList)
			adminUser.DELETE("/deleteAdmin", auth.AdminUserTokenCheck, deleteAdminUser)
			adminUser.PUT("/addAdmin", auth.AdminUserTokenCheck, addAdminUser)
			adminUser.POST("/chAdmin", auth.AdminUserTokenCheck, modifyAdminUser)
		}

		announcement := admin.Group("/notice", auth.AdminUserTokenCheck)
		{
			announcement.PUT("/addNotice", addAnnouncement)
			announcement.GET("/getNotice", getAnnouncement)
			announcement.PATCH("/chNotice", modifyAnnouncement)
			announcement.DELETE("/deleteNotice", deleteAnnouncement)
			announcement.POST("/publishNotice", publishAnnouncement)
		}

		operationRecord := admin.Group("/operationRecord", auth.AdminUserTokenCheck)
		{
			operationRecord.GET("/operationRecord")
			operationRecord.GET("/operationRecordCounts")
		}

		termConfigure := admin.Group("/term")
		{
			termConfigure.GET("/getAllTerm", getTermList)
			termConfigure.PUT("/addTerm", auth.AdminUserTokenCheck, addTerm)
			termConfigure.PATCH("/chTerm", auth.AdminUserTokenCheck, modifyTerm)
			termConfigure.DELETE("/deleteTerm", auth.AdminUserTokenCheck, deleteTerm)
		}

		versionConfigure := admin.Group("/version", auth.AdminUserTokenCheck)
		{
			versionConfigure.GET("/getAll", getVersionList)
			versionConfigure.PATCH("/update", modifyVersion)
			versionConfigure.PUT("/add", addVersion)
			versionConfigure.DELETE("/delete", deleteVersion)
			versionConfigure.POST("/publish", publishVersion)
		}

		// 接口已弃用
		webConfigure := admin.Group("/website", auth.AdminUserTokenCheck)
		{
			webConfigure.POST("/add-website-pic")
			webConfigure.GET("/get-pic-location")
			webConfigure.POST("/set-displaying")
			webConfigure.POST("/set-one-displaying")
			webConfigure.GET("/displaying-pic")
			webConfigure.POST("/set-admin-info")
			webConfigure.POST("/update-admin-info")
			webConfigure.GET("/admin-info-list")
			webConfigure.POST("/add-inform")
			webConfigure.POST("/update-inform")
			webConfigure.GET("/inform-list")
			webConfigure.POST("/deleteAdminInfo")
			webConfigure.POST("/deleteInform")
			webConfigure.POST("/addProduction")
			webConfigure.POST("/updateProduction")
			webConfigure.POST("/deleteProduction")
			webConfigure.GET("/productions")
			webConfigure.GET("/production-info")
		}
	}
}

func setupPublicApiRouter(rootRouter *gin.RouterGroup) {
	joinUs := rootRouter.Group("/join-us")
	{
		joinUs.POST("/sendApplication")
		joinUs.GET("/refuse")
		joinUs.GET("/accept")
		joinUs.GET("/queryAll")
		joinUs.GET("/send")
	}

	wusthelper := rootRouter.Group("/wusthelper")
	{
		wusthelper.GET("/notice", getPublishedAnnouncement)
		wusthelper.GET("/act", getPublishedBannerList)
		wusthelper.GET("/config", getConfigListPublic)
		wusthelper.GET("/log", getPublishedLogList)
		wusthelper.GET("/version", getLatestVersion)
	}
}
