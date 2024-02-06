package ecode

func InitEcodeText() {
	texts := map[Code]string{}

	texts[InternalError] = "内部错误"

	texts[TokenInvalid] = "token无效"
	texts[IdNotNullable] = "id不能为空"
	texts[DeleteFailed] = "删除失败"
	texts[QueryFailed] = "查询失败"
	texts[InvalidId] = "id错误"
	texts[PlatformMissing] = "header中平台参数不正确"

	texts[GettingNumOfUserFailed] = "用户人数获取失败"
	texts[GettingNumOfCollegeUserFailed] = "各学院的用户总数获取出错"

	texts[NotAuthorized] = "登陆数据缺失"
	texts[AuthFailed] = "用户名或密码错误"
	texts[UsernameNotNullable] = "账号不能为空"
	texts[PasswordNotNullable] = "密码不能为空"
	texts[UsernameExists] = "用户名已存在"
	texts[PermissionDenied] = "没有权限"
	texts[UserNotExists] = "没有找到用户"

	texts[ContentCannotBeEmpty] = "内容不能都为空"
	texts[LogNotFound] = "找不到此日志"
	texts[AnnouncementPublishFailed] = "发布公告失败"

	texts[AddAdminLogFailed] = "管理端日志添加失败"

	texts[VersionOperationFailed] = "版本信息操作失败"
	texts[ParamWrong] = "参数错误"

	Register(texts)
}
