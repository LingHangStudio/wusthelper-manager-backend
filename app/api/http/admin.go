package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"wusthelper-manager-go/app/model"
	"wusthelper-manager-go/app/service"
	"wusthelper-manager-go/library/ecode"
	"wusthelper-manager-go/library/token"
)

type AdminUserLoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AdminUserLoginResp struct {
	Id      int64  `json:"id"`
	Token   string `json:"token"`
	Groupid int8   `json:"groupid"`
}

type AdminUserResp struct {
	Id         int64  `json:"id"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	AddTime    string `json:"addTime"`
	Groupid    int8   `json:"groupid"`
	UpdateTime string `json:"updateTime"`
}

type AdminUserDeleteReq struct {
	Id uint64 `json:"id" binding:"required"`
}

type AdminUserAddReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AdminUserModifyReq struct {
	Id       int64   `json:"id" binding:"required"`
	Username *string `json:"username"`
	Password *string `json:"password"`
	Groupid  *int8   `json:"groupid"`
}

func adminUserLogin(c *gin.Context) {
	req := new(AdminUserLoginReq)
	if err := c.ShouldBindJSON(req); err != nil {
		responseEcode(c, ecode.ParamWrong)
		return
	}

	user, err := srv.AdminUserLogin(req.Username, req.Password)
	if err != nil {
		responseEcode(c, err)
	}

	tokenPayload := token.SignPayload{
		Uid:      user.ID,
		Username: *user.Username,
	}

	t := jwt.Sign(tokenPayload)

	resp := AdminUserLoginResp{
		Id:      user.ID,
		Token:   t,
		Groupid: *user.Group,
	}

	responseData(c, resp)
}

func getAdminUserList(c *gin.Context) {
	userList, err := srv.GetAllAdminUserList()
	if err != nil {
		responseEcode(c, err)
		return
	}

	resp := make([]AdminUserResp, len(*userList))
	for i, user := range *userList {
		createTime := user.CreateTime.Format(_defaultDateTimeFormat)
		updateTime := user.UpdateTime.Format(_defaultDateTimeFormat)
		resp[i] = AdminUserResp{
			Id:         user.ID,
			Username:   *user.Username,
			Password:   fmt.Sprintf("Bcrypt hash: %s", *user.Password),
			AddTime:    createTime,
			Groupid:    *user.Group,
			UpdateTime: updateTime,
		}
	}

	responseData(c, resp)
}

func deleteAdminUser(c *gin.Context) {
	req := new(AdminUserDeleteReq)
	if err := c.ShouldBindJSON(req); err != nil {
		responseEcode(c, ecode.ParamWrong)
		return
	}

	err := srv.DeleteAdminUserById(req.Id)
	if err != nil {
		responseEcode(c, err)
		return
	}

	responseData(c, nil)
}

func addAdminUser(c *gin.Context) {
	req := new(AdminUserAddReq)
	if err := c.ShouldBindJSON(req); err != nil {
		responseEcode(c, ecode.ParamWrong)
		return
	}

	data := service.AdminUserAddParam{
		Username: req.Username,
		Password: req.Password,
	}

	err := srv.AddAdminUser(&data)
	if err != nil {
		responseEcode(c, err)
		return
	}

	responseData(c, nil)
}

func modifyAdminUser(c *gin.Context) {
	uid, err := getUid(c)
	if err != nil {
		responseEcode(c, ecode.ParamWrong)
		return
	}

	// 权限校验，非su不能更改管理员用户
	userGroup, err := srv.GetUserGroup(uid)
	if err != nil || userGroup != model.SuperAdminGroup {
		responseEcode(c, ecode.PermissionDenied)
		return
	}

	req := new(AdminUserModifyReq)
	if err := c.ShouldBindJSON(req); err != nil {
		responseEcode(c, ecode.ParamWrong)
		return
	}

	// 如果用户名有修改（非nil），进行查重
	if req.Username != nil {
		exists, err := srv.CheckUsernameExists(*req.Username)
		if err != nil {
			responseEcode(c, ecode.ParamWrong)
			return
		} else if exists {
			responseEcode(c, ecode.UsernameExists)
			return
		}
	}

	data := service.AdminUserModifyParam{
		Id:       req.Id,
		Username: req.Username,
		Password: req.Password,
		Groupid:  req.Groupid,
	}

	err = srv.ModifyAdminUser(&data)
	if err != nil {
		responseEcode(c, err)
		return
	}

	responseData(c, nil)
}
