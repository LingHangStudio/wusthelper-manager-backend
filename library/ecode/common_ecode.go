package ecode

// All common ecode
var (
	OK            = add(0)
	InternalError = add(1)

	TokenInvalid    = add(10001) // token不存在
	IdNotNullable   = add(10002) // id不能为空
	DeleteFailed    = add(10003) // 删除失败
	QueryFailed     = add(10004) // 查询失败
	InvalidId       = add(10005) // id错误
	PlatformMissing = add(10006) // header中平台参数不正确

	GettingNumOfUserFailed        = add(20102) // 用户人数获取失败
	GettingNumOfCollegeUserFailed = add(20103) // 各学院的用户总数获取出错

	NotAuthorized       = add(20200) // 登陆数据缺失
	AuthFailed          = add(20201) // 用户名或密码错误
	UsernameNotNullable = add(20202) // 账号不能为空
	PasswordNotNullable = add(20203) // 密码不能为空
	UsernameExists      = add(20204) // 用户名已存在
	PermissionDenied    = add(20205) // 没有权限
	UserNotExists       = add(20206) // 没有找到用户

	ContentCannotBeEmpty      = add(20300) // 内容不能都为空
	LogNotFound               = add(20301) // 找不到此日志
	AnnouncementPublishFailed = add(20302) // 发布公告失败

	AddAdminLogFailed = add(40102) // 管理端日志添加失败

	VersionOperationFailed = add(50100) // 版本信息操作失败
	ParamWrong             = add(50101) // 请求的参数不正确
)
