package v3

import (
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"time"
	"wusthelper-manager-go/app/conf"
	"wusthelper-manager-go/library/ecode"
	"wusthelper-manager-go/library/log"
)

type WusthelperHttpRpc struct {
	client *resty.Client
}

func NewRpcClient(c *conf.WusthelperConf) (rpc *WusthelperHttpRpc) {
	httpClient := resty.New()
	httpClient.SetBaseURL(c.Upstream)
	httpClient.SetHeader("User-Agent", "wusthelper-manager/0.0.1")
	httpClient.SetHeader("Content-Type", "application/json")
	httpClient.SetTimeout(time.Second * c.Timeout)

	if c.Proxy != "" {
		httpClient.SetProxy(c.Proxy)
	}

	rpc = &WusthelperHttpRpc{
		client: httpClient,
	}

	return
}

const (
	adminSuccess     = 0
	success          = 10000 // 成功
	localSuc         = 11000 // 本地成功
	limitRequest     = 99999 // 操作频繁，请稍后再试
	err              = 11111 // 失败
	baseParamErr     = 10001 // 参数缺失或参数类型错误
	requestDenied    = 10002 // 请求错误
	requestCalledErr = 10003 // 接口调用错误:请求出错，data中为错误信息，请检查并详细参照api文档使用该接口
	clientAbortErr   = 10004 // 连接中断

	authErr                = 20001 // 用户身份验证模块异常
	authErrTokenMiss       = 21001 // 用户身份校验失败:请求缺失令牌token
	authErrTokenInvalid    = 21002 // 用户身份校验失败:token无效或过期
	authDecodeForStuNumSuc = 10000 // token解密成功
	authDecodeForStuNumErr = 21101 // token解密失败

	// 本科生
	jwcPwdNedUpd                       = 30001 // 教务处密码已修改,请重新登录
	jwcLoginSuc                        = 10000 // 教务处登录成功
	jwcLoginErrJwcErr                  = 30104 // 登录失败:教务系统异常，请重试
	jwcLoginErrJwcFinInfoErr           = 30101 // 登陆失败:账号或密码错误
	jwcLoginErrJwcModDefPwd            = 30102 // 登录失败:请前往教务官网修改默认密码
	jwcLoginErrJwcErrInfoErr           = 30103 // 登录失败:账号或密码错误，或教务异常，请重试
	jwcLoginErrJwcRetryTooManyTimesErr = 30101 // 登录失败:密码重试次数已达三次，继续重试可能会导致暂时封禁
	jwcLoginErrJwcUserWasBannedErr     = 30101 // 登录失败:密码总重试次数已达五次，已被暂时封禁
	jwcLoginSucJwcErrLocSuc            = 11000 // 教务异常:缓存登录成功
	jwcGetStuInfoSuc                   = 10000 // 教务处获取用户信息成功
	jwcGetStuInfoSucJwcErrLocSuc       = 11000 // 教务异常:缓存学生信息获取成功
	jwcGetStuInfoErrNoSuchStu          = 30201 // 教务处获取用户信息失败:用户不存在
	jwcGetCurrCourseSuc                = 10000 // 教务处获取课表信息成功
	jwcGetCurrCourseSucJwcErrLocSuc    = 11000 // 教务异常:缓存课表信息获取成功
	jwcGetCurrCourseErrJwcErr          = 30301 // 教务异常:获取课表失败
	jwcGetGradeSuc                     = 10000 // 教务处成绩获取成功
	jwcGetGradeErrNeedEvaluateTeaching = 30601 // 教务处成绩获取失败:请前往教务官网评教
	jwcGetGradeErrJwcErrLocSuc         = 11000 // 教务异常:缓存成绩获取成功
	jwcTrainingSchemeSuc               = 10000 // 教务处培养方案获取成功
	jwcTrainingSchemeErr               = 30701 // 教务异常，培养方案获取失败

	// 研究生
	yjsLoginSuc              = 10000 // 研究生官网登录成功
	yjsLoginErrYjsFinInfoErr = 80002 // 登录失败:账号或密码错误
	yjsLoginError            = 11111 // 登录失败，研究生官网异常，请稍后再试
	yjsGetStuInfoSuc         = 10000 // 研究生官网获取用户信息成功
	yjsGetStuInfoErr         = 11111 // 研究生官网异常，获取用户信息失败，请稍后再试
	yjsGetGradeSuc           = 10000 // 研究生官网成绩获取成功
	yjsGetGradeErr           = 11111 // 研究生官网异常，获取成绩失败，请稍后再试
	yjsGetCourseSuc          = 10000 // 研究生官网课表获取成功
	yjsGetCourseErr          = 11111 // 研究生官网异常，刷新课表失败，请稍后再试
	yjsGetTrainingPlanSuc    = 10000 // 研究生官网培养管理获取成功
	yjsGetTrainingPlanErr    = 11111 // 研究生官网异常，获取培养管理失败，请稍后再试
	yjsPwdNedUpd             = 80001 // 研究生官网密码已修改,请重新登录
)

func toEcode(code int, service string) ecode.Code {
	switch code {
	default:
		log.Warn("未处理的助手上游响应码", zap.Int("code", code), zap.String("service", service))
	}

	return ecode.InternalError
}
