package http

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"wusthelper-manager-go/app/model"
	"wusthelper-manager-go/app/service"
	"wusthelper-manager-go/library/ecode"
)

func getConfigListPublic(c *gin.Context) {
	platform := getPlatform(c)
	if platform == "" {
		responseEcode(c, ecode.ParamWrong)
		return
	}

	// 获取配置
	configList, _, err := srv.GetConfigList(platform)
	if err != nil {
		responseEcode(c, err)
		return
	}

	// 转换以前的鬼格式（真是麻了）
	resp := map[string]any{}
	switch platform {
	case _platformMp:
		resp["menuList"] = map[string]string{}
		resp["schedule"] = map[string]string{}
		for _, conf := range *configList {
			if *conf.Type == model.ConfigValueTypeBool {
				resp[*conf.Name] = _convertBoolStr2Bool(*conf.Value)
			} else {
				switch *conf.Name {
				// 小程序的一些配置项（都动态配置了还搞这种玩意放到menulist和schedule这些干嘛）
				case "news":
				case "volunteer":
					resp["menuList"].(map[string]string)[*conf.Name] = *conf.Value
				case "refreshSchedule":
				case "scheduleVersion":
					resp["schedule"].(map[string]string)[*conf.Name] = *conf.Value
				default:
					resp[*conf.Name] = *conf.Value
				}
			}
		}
	default:
		for _, conf := range *configList {
			if *conf.Type == model.ConfigValueTypeBool {
				// android和ios端需要的是数字型的bool
				resp[*conf.Name] = _convertBoolStr2Int(*conf.Value)
			} else {
				resp[*conf.Name] = *conf.Value
			}
		}
	}

	// 获取最新版本
	latestVersion, err := srv.GetLatestVersion(platform)
	if err != nil {
		responseEcode(c, err)
		return
	}

	if latestVersion != nil {
		resp["updateContent"] = latestVersion.Summary
		resp["version"] = latestVersion.VersionText
		resp["apkUrl"] = _getFileUrl(*latestVersion.File)
	}

	// 获取学期
	terms, err := srv.GetTermList()
	if err != nil {
		responseEcode(c, err)
		return
	}

	if terms != nil {
		termResp := make([]TermListResp, len(*terms))
		for i, term := range *terms {
			switch platform {
			default:
				termResp[i] = TermListResp{
					Id:        term.ID,
					Term:      *term.Term,
					StartDate: term.Start.Format(_defaultDateFormat),
				}
			case _platformAndroid:
				// 安卓端需要转换成时间戳
				termResp[i] = TermListResp{
					Id:        term.ID,
					Term:      *term.Term,
					StartDate: strconv.FormatInt(term.Start.UnixMilli(), 10),
				}
			}
		}

		// 令人疑惑的termList
		switch platform {
		case _platformMp:
			resp["termList"] = termResp
		default:
			resp["termSetting"] = termResp
		}
	}

	// 可算是结束了
	responseData(c, resp)
}

type ConfigItemResp struct {
	Id             int64          `json:"id"`
	SettingName    string         `json:"settingName"`
	CurrentSetting string         `json:"currentSetting"`
	Type           int8           `json:"type"`
	Content        string         `json:"content"`
	Platform       string         `json:"platform"`
	UpdateTime     string         `json:"updateTime"`
	OptionsList    []ConfigOption `json:"optionsList"`
}

type ConfigOption struct {
	OptionName string `json:"optionName"`
}

func getConfigList(c *gin.Context) {
	req := new(PlatformPaginationReq)
	if err := c.ShouldBind(req); err != nil {
		responseEcode(c, ecode.ParamWrong)
		return
	}

	resultList, total, err := srv.GetConfigList(req.Platform)
	if err != nil {
		responseEcode(c, err)
		return
	}

	configRespList := make([]ConfigItemResp, len(*resultList))
	for i, result := range *resultList {
		options := make([]ConfigOption, len(*result.PossibleValues))
		for i2, s := range *result.PossibleValues {
			options[i2] = ConfigOption{OptionName: s}
		}

		configRespList[i] = ConfigItemResp{
			Id:             result.ID,
			SettingName:    *result.Name,
			CurrentSetting: *result.Value,
			Type:           *result.Type,
			Content:        *result.Describe,
			Platform:       *result.Platform,
			UpdateTime:     result.UpdateTime.Format(_defaultDateTimeFormat),
			OptionsList:    options,
		}
	}

	responseData(c, map[string]any{
		"configs": configRespList,
		"num":     total,
	})
}

type PlatformResp struct {
	Platform []string `json:"platform"`
}

func getPlatformList(c *gin.Context) {
	result, err := srv.GetPlatformList()
	if err != nil {
		responseEcode(c, err)
	}

	responseData(c, PlatformResp{Platform: *result})
}

type ConfigAddReq struct {
	SettingName    string   `json:"settingName" binding:"required"`
	CurrentSetting string   `json:"currentSetting" binding:"required"`
	Type           int8     `json:"type" binding:"required"`
	Content        string   `json:"content"`
	Platform       []string `json:"platform" binding:"required"`
	OptionList     []string `json:"optionList" binding:"required"`
}

func addConfig(c *gin.Context) {
	req := new(ConfigAddReq)
	if err := c.ShouldBind(req); err != nil {
		responseEcode(c, ecode.ParamWrong)
		return
	}

	conf := service.ConfigAddParam{
		Name:           req.SettingName,
		Value:          req.CurrentSetting,
		Type:           req.Type,
		Describe:       req.Content,
		PossibleValues: req.OptionList,
		Platform:       req.Platform,
	}

	if req.Type == model.ConfigValueTypeBool {
		conf.PossibleValues = []string{"true", "false"}
	}

	err := srv.AddConfig(&conf)

	if err != nil {
		responseEcode(c, err)
		return
	}

	responseData(c, nil)
}

type ConfigModifyReq struct {
	Id             int64     `json:"configId" binding:"required"`
	ConfigName     *string   `json:"configName"`
	CurrentSetting *string   `json:"currentSetting"`
	Content        *string   `json:"content"`
	OptionList     *[]string `json:"optionList"`
}

func modifyConfig(c *gin.Context) {
	req := new(ConfigModifyReq)
	if err := c.ShouldBind(req); err != nil {
		responseEcode(c, ecode.ParamWrong)
		return
	}

	err := srv.ModifyConfig(&service.ConfigModifyParam{
		Id:             req.Id,
		Name:           req.ConfigName,
		Value:          req.CurrentSetting,
		Describe:       req.Content,
		PossibleValues: req.OptionList,
	})

	if err != nil {
		responseEcode(c, err)
		return
	}

	responseData(c, nil)
}

type ConfigDeleteReq struct {
	ConfigId int64 `json:"configId" form:"configId" query:"configId" binding:"required"`
}

func deleteConfig(c *gin.Context) {
	req := new(ConfigDeleteReq)
	if err := c.ShouldBind(req); err != nil {
		responseEcode(c, ecode.ParamWrong)
		return
	}

	err := srv.DeleteConfig(req.ConfigId)
	if err != nil {
		responseEcode(c, err)
		return
	}

	responseData(c, nil)
}
