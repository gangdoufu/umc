package api

import (
	"errors"
	"github.com/gangdoufu/umc/internal/umcserver/global"
	"github.com/gangdoufu/umc/internal/umcserver/middleware"
	"github.com/gangdoufu/umc/internal/umcserver/model"
	"github.com/gangdoufu/umc/internal/umcserver/service"
	"github.com/gangdoufu/umc/internal/umcserver/service/vo"
	"github.com/gangdoufu/umc/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// 注册
func Register(c *gin.Context) {
	user := &model.User{}
	var err error
	err = c.ShouldBind(user)
	if err != nil {
		response.Error(c, err)
		return
	}
	db := global.DB.Begin()
	defer func() {
		if err != nil {
			db.Rollback()
		} else {
			db.Commit()
		}
	}()
	middleware.SetTransaction(c, db)
	us := service.NewUserService(db)
	if err = us.Register(c.Request.Context(), user); err != nil {
		response.Error(c, err)
	} else {
		// todo 注册成功跳转到登录页面,这里需要提供登录页面url
		response.SuccessWithData(c, "login ref url")
	}
}

// 管理员创建用户
func CreateUser(c *gin.Context) {
	user := &model.User{}
	var err error
	err = c.ShouldBind(user)
	if err != nil {
		response.Error(c, err)
		return
	}
	db := global.DB.Begin()
	defer func() {
		if err != nil {
			db.Rollback()
		} else {
			db.Commit()
		}
	}()
	us := service.NewUserService(db)
	if err = us.CreateUser(c.Request.Context(), user); err != nil {
		response.Error(c, err)
	} else {
		response.Success(c)
	}
}

// 激活
func AccountActive(c *gin.Context) {

}

// 登录
func Login(c *gin.Context) {
	vo := &vo.LoginVo{}
	var err error
	err = c.ShouldBind(vo)
	if err != nil {
		response.Error(c, err)
		return
	}
	db := global.DB
	us := service.NewUserService(db)
	if err = us.Login(c.Request.Context(), vo); err != nil {
		response.Error(c, err)
	} else {
		response.Success(c)
	}
}

// 注销

func LoginOut(c *gin.Context) {

}

// 更新用户信息
func UpdateUserInfo(c *gin.Context) {

}

// 修改密码
func ChangePassword(c *gin.Context) {

}

// 忘记密码重置密码
func ResetPassword(c *gin.Context) {

}

// 用户加入用户组
func UserAddInGroup(c *gin.Context) {
	userVO := &vo.UserGroupVo{}
	var err error
	err = c.ShouldBind(userVO)
	if err != nil {
		response.Error(c, err)
		return
	}
	db := global.DB.Begin()
	defer func() {
		if err != nil {
			db.Rollback()
		} else {
			db.Commit()
		}
	}()
	us := service.NewUserService(db)
	if err = us.AddUserToGroup(c.Request.Context(), userVO); err != nil {
		response.Error(c, err)
	} else {
		response.Success(c)
	}
}

// 用户加入租户
func UserAddInTenant(c *gin.Context) {
	userVO := &vo.UserGroupVo{}
	var err error
	err = c.ShouldBind(userVO)
	if err != nil {
		response.Error(c, err)
		return
	}
	db := global.DB.Begin()
	defer func() {
		if err != nil {
			db.Rollback()
		} else {
			db.Commit()
		}
	}()
	us := service.NewUserService(db)
	if err = us.AddUserToTenant(c.Request.Context(), userVO); err != nil {
		response.Error(c, err)
	} else {
		response.Success(c)
	}
}

// 获取验证码

func UserGetVerificationCode(c *gin.Context) {

}

// 锁定用户
func LockUser(c *gin.Context) {
	userVO := &vo.UserGroupVo{}
	var err error
	err = c.ShouldBind(userVO)
	if err != nil {
		response.Error(c, err)
		return
	}
	userInfo := model.User{}
	userInfo.ID = userVO.UserId
	userInfo.SetUserStatusLocking()
	db := global.DB.Begin()
	defer func() {
		if err != nil {
			db.Rollback()
		} else {
			db.Commit()
		}
	}()
	us := service.NewUserService(db)
	if err = us.UpdateUserInfo(c.Request.Context(), &userInfo); err != nil {
		response.Error(c, err)
	} else {
		response.Success(c)
	}
}

// 解锁用户

// 查询用户信息
func GetUserBaseInfo(c *gin.Context) {
	param := c.Param("userid")
	if param == "" {
		response.Error(c, errors.New("need user id"))
	}
	userId := cast.ToUint(param)
	if userId == 0 {
		response.Error(c, errors.New("need user id"))
	}
	us := service.NewUserService(global.DB)
	userInfo, err := us.QueryUserInfoById(c.Request.Context(), userId)
	if err != nil {
		response.Error(c, err)
	} else {
		response.SuccessWithData(c, userInfo)
	}
}

// 查询用户所属于的租户
func GetUserTenants(c *gin.Context) {
	param := c.Param("userid")
	if param == "" {
		response.Error(c, errors.New("need user id"))
	}
	userId := cast.ToUint(param)
	if userId == 0 {
		response.Error(c, errors.New("need user id"))
	}
	us := service.NewUserService(global.DB)
	tenants, err := us.ListUserTenants(c.Request.Context(), userId)
	if err != nil {
		response.Error(c, err)
	} else {
		response.SuccessWithData(c, tenants)
	}
}

func GetUserGroups(c *gin.Context) {
	userVO := &vo.UserGroupVo{}
	var err error
	err = c.ShouldBind(userVO)
	if err != nil {
		response.Error(c, err)
		return
	}
	us := service.NewUserService(global.DB)
	groups, err := us.ListUserGroups(c.Request.Context(), userVO)
	if err != nil {
		response.Error(c, err)
	} else {
		response.SuccessWithData(c, groups)
	}

}
