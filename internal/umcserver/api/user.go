package api

import (
	"errors"
	"github.com/gangdoufu/umc/internal/umcserver/global"
	"github.com/gangdoufu/umc/internal/umcserver/middleware"
	"github.com/gangdoufu/umc/internal/umcserver/model"
	"github.com/gangdoufu/umc/internal/umcserver/redis"
	"github.com/gangdoufu/umc/internal/umcserver/service"
	"github.com/gangdoufu/umc/internal/umcserver/service/vo"
	"github.com/gangdoufu/umc/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"time"
)

// 注册
func Register(c *gin.Context) {
	userVo := &vo.UserVo{}
	var err error
	err = c.ShouldBind(userVo)
	if err != nil {
		response.Error(c, err)
		return
	}
	us := service.NewUserService(middleware.SetGinDB(true, c))
	if err = us.Register(c.Request.Context(), userVo.GetModel()); err != nil {
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
	us := service.NewUserService(middleware.SetGinDB(true, c))
	if err = us.CreateUser(c.Request.Context(), user); err != nil {
		response.Error(c, err)
	} else {
		response.Success(c)
	}
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
	if userInfos, err := us.Login(c.Request.Context(), vo); err != nil {
		response.Error(c, err)
	} else {
		response.SuccessWithData(c, userInfos)
	}
}

// 注销 就是将用户的token放入黑名单
func LoginOut(c *gin.Context) {
	token := middleware.GetToken(c)
	parseToken, err := global.Jwt.ParseToken(token)
	if err != nil {
		response.Error(c, err)
		return
	}
	expiration := parseToken.ExpiresAt.Sub(time.Now())
	err = redis.AddTokenTOBlacklist(c.Request.Context(), token, expiration)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c)
	return
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
	db := middleware.SetGinDB(true, c)
	us := service.NewUserService(db)
	if err = us.AddUserToGroup(c.Request.Context(), userVO); err != nil {
		response.Error(c, err)
	} else {
		response.Success(c)
	}
}

func UserRemoveGroup(c *gin.Context) {
	userVO := &vo.UserGroupVo{}
	var err error
	err = c.ShouldBind(userVO)
	if err != nil {
		response.Error(c, err)
		return
	}
	db := middleware.SetGinDB(true, c)
	us := service.NewUserService(db)
	if err = us.RemoveUserGroup(c.Request.Context(), userVO); err != nil {
		response.Error(c, err)
	} else {
		response.Success(c)
	}
}

func UserRemoveTenant(c *gin.Context) {
	userVO := &vo.UserGroupVo{}
	var err error
	err = c.ShouldBind(userVO)
	if err != nil {
		response.Error(c, err)
		return
	}
	db := middleware.SetGinDB(true, c)
	us := service.NewUserService(db)
	if err = us.RemoveUserTenant(c.Request.Context(), userVO); err != nil {
		response.Error(c, err)
	} else {
		response.Success(c)
	}
}

func UserForgetPassword(c *gin.Context) {
	account := c.Param("user_account")
	if account == "" {
		response.Error(c, errors.New("需要账号"))
		return
	}
	db := global.DB
	us := service.NewUserService(db)
	err := us.ForgetPassword(c.Request.Context(), account)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c)
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
	db := middleware.SetGinDB(true, c)
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
		return
	}
	userId := cast.ToUint(param)
	if userId == 0 {
		response.Error(c, errors.New("need user id"))
		return
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
		return
	}
	userId := cast.ToUint(param)
	if userId == 0 {
		response.Error(c, errors.New("need user id"))
		return
	}
	us := service.NewUserService(global.DB)
	tenants, err := us.ListUserTenants(c.Request.Context(), userId)
	if err != nil {
		response.Error(c, err)
		return
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

func CheckJwt(c *gin.Context) {
	token := middleware.GetToken(c)
	parseToken, err := global.Jwt.ParseToken(token)
	if err != nil {
		response.Error(c, err)
		return
	}
	claims := parseToken.BaseClaims
	response.SuccessWithData(c, &claims)
	return
}
