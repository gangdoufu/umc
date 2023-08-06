package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/gangdoufu/umc/internal/umcserver/global"
	"github.com/gangdoufu/umc/internal/umcserver/message/producer"
	"github.com/gangdoufu/umc/internal/umcserver/model"
	"github.com/gangdoufu/umc/internal/umcserver/redis"
	"github.com/gangdoufu/umc/internal/umcserver/sender"
	"github.com/gangdoufu/umc/internal/umcserver/service/vo"
	"github.com/gangdoufu/umc/internal/umcserver/store"
	"github.com/gangdoufu/umc/pkg/common"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"math/rand"
	"net/url"
	"path"
	"time"
)

type IUserService interface {
	Login(ctx context.Context, loginVo *vo.LoginVo) (*model.UserLoginVo, error)
	Register(ctx context.Context, user *model.User) error
	CreateUser(ctx context.Context, user *model.User) error
	RemoveUserGroup(ctx context.Context, groupVo *vo.UserGroupVo) error
	AddUserToGroup(ctx context.Context, groupVo *vo.UserGroupVo) error
	RemoveUserTenant(ctx context.Context, groupVo *vo.UserGroupVo) error
	AddUserToTenant(ctx context.Context, groupVo *vo.UserGroupVo) error
	UpdateUserInfo(ctx context.Context, user *model.User) error
	QueryUserInfoById(ctx context.Context, userId uint) (*model.User, error)
	ListUserTenants(ctx context.Context, userId uint) ([]*model.TenantShowVo, error)
	ListUserGroups(ctx context.Context, groupVo *vo.UserGroupVo) ([]*model.Group, error)
	ForgetPassword(ctx context.Context, account string) error
}
type userService struct {
	factory store.IFactory
}

var loginFunMap = map[vo.LoginType]loginFun{
	vo.AccountPassword: loginByAccountAndPwd,
	vo.TelCode:         loginByTelVerificationCode,
	vo.EmailCode:       loginByEmailVerificationCode,
}

func NewUserService(db *gorm.DB) IUserService {
	return &userService{factory: store.NewFactory(db)}
}

func (s *userService) ForgetPassword(ctx context.Context, account string) error {
	byAccount, err := s.factory.Users().QueryByAccount(account)
	if err != nil {
		return err
	}
	if byAccount == nil {
		return errors.New("account is not exist")
	}
	code := common.RandToken(16)
	err = redis.SetUserResetPasswordCode(ctx, &redis.UserVerificationCodeVo{
		UserId:   byAccount.ID,
		Source:   byAccount.Email,
		Code:     code,
		CodeType: "user_forget_password",
	})
	if err != nil {
		return err
	}
	u, _ := url.Parse(global.Config.Server.GetHostUrl())
	info := sender.NewResetPasswordInfo(byAccount.Email, path.Join(u.Path, fmt.Sprintf("/manger/user/reset_password?code=%s", code)))
	return producer.MessageProducer.SendVerificationCode(ctx, info)
}

// 登录
func (s *userService) Login(ctx context.Context, loginVo *vo.LoginVo) (*model.UserLoginVo, error) {
	users := s.factory.Users()
	var curLoginFun loginFun = loginByAccountAndPwd
	if login, ok := loginFunMap[loginVo.Type]; ok {
		curLoginFun = login
	}
	user, err := curLoginFun(ctx, loginVo.Key, loginVo.Value, users)
	if err != nil {
		return nil, err
	}
	claims := global.Jwt.CreateClaims(common.BaseClaims{
		UserID:   user.ID,
		Account:  user.Account,
		ClientId: uuid.New().String(),
	})
	token, err := global.Jwt.CreateToken(&claims)
	if err != nil {
		return nil, err
	}
	return &model.UserLoginVo{
		ID:       user.ID,
		Account:  user.Account,
		ClientId: claims.ClientId,
		Token:    token,
	}, nil

}

type loginFun func(context.Context, string, string, store.IUserStore) (*model.User, error)

// 使用账号密码登录
func loginByAccountAndPwd(ctx context.Context, account string, password string, users store.IUserStore) (*model.User, error) {
	one, err := users.QueryUserPasswordByAccount(account)
	if one == nil {
		return nil, errors.New("账号不存在")
	}
	user := &model.User{}
	user.Password = one.Password
	if err != nil {
		return nil, err
	}

	if one != nil {
		if !user.CheckPassword(password) {
			return nil, errors.New(" 用户名或密码不正确")
		}
		user.ID = one.ID
		user.Account = one.Account
		return user, nil
	} else {
		return nil, errors.New(" 用户名不存在")
	}

}

// 电话号码+验证码登录
func loginByTelVerificationCode(ctx context.Context, tel, code string, users store.IUserStore) (*model.User, error) {
	one, err := users.FindOne(&model.User{Tel: tel})
	if err != nil {
		return nil, err
	}
	if one != nil {
		uv := &redis.UserVerificationCodeVo{
			UserId:   one.ID,
			Source:   tel,
			Code:     code,
			CodeType: "login_tel_code",
		}
		if !redis.CheckUserVerificationCode(ctx, uv) {
			return nil, errors.New("验证码不正确")
		}
		return one, nil
	} else {
		return nil, errors.New("用户名不存在")
	}
}

// 邮箱验证码登录
func loginByEmailVerificationCode(ctx context.Context, email, code string, users store.IUserStore) (*model.User, error) {
	one, err := users.FindOne(&model.User{Email: email})
	if err != nil {
		return nil, err
	}
	if one != nil {
		uv := &redis.UserVerificationCodeVo{
			UserId:   one.ID,
			Source:   email,
			Code:     code,
			CodeType: "login_email_code",
		}
		if !redis.CheckUserVerificationCode(ctx, uv) {
			return nil, errors.New("验证码不正确")
		}
		return one, nil
	} else {
		return nil, errors.New("用户名不存在")
	}
}

// 获取验证码
func (s userService) GetVerificationCode(ctx context.Context, sourceType, source, codeType string) (string, error) {
	var query = &model.User{}
	var sender sendCodeFun
	switch sourceType {
	case "tel":
		query.Tel = source
		sender = codeTelSender
	case "email":
		query.Email = source
		sender = codeEmailSender
	default:
		return "", errors.New("not support")
	}
	users := s.factory.Users()
	one, err := users.FindOne(query)
	if err != nil {
		return "", err
	}
	if one != nil {
		r := rand.New(rand.NewSource(99))
		code := common.IntSliceToString(r.Perm(6))
		uv := &redis.UserVerificationCodeVo{
			UserId:   one.ID,
			Source:   source,
			Code:     code,
			CodeType: codeType,
		}
		if err = redis.SetUserVerificationCode(ctx, uv); err != nil {
			return "", err
		}
		if err = sender(source, code, codeType); err != nil {
			return "", err
		}
		return code, nil
	}
	return "", errors.New("not found")
}

type sendCodeFun func(source, code, codeType string) error

// 通过邮件发送验证码
func codeEmailSender(emailAddress, code, codeType string) error {

	return nil
}

// 通过短信发送验证码
func codeTelSender(tel, code, codeType string) error {

	return nil
}

// 注销
func (s userService) LoginOut(ctx context.Context, token string, expiration time.Duration) error {
	// 将token放入黑名单
	return redis.AddTokenTOBlacklist(ctx, token, expiration)
}

// 用户在被邀请之后会发送一个邀请邮件或者短信给对方。对方点击之后会确认key并允许对方重新设置密码
func (s userService) createActiveKeys(ctx context.Context, userId uint, account string) (string, error) {
	code := uuid.New().String()
	uv := &redis.UserVerificationCodeVo{
		UserId:   userId,
		Source:   account,
		Code:     code,
		CodeType: "user_active_code",
	}
	return code, redis.SetUserVerificationCode(ctx, uv)
}

// 激活账号
func (s userService) ActiveUser(ctx context.Context, account, password, code string) error {
	users := s.factory.Users()
	user, err := users.QueryByAccount(account)
	if err != nil {
		return err
	}
	uv := &redis.UserVerificationCodeVo{
		UserId:   user.ID,
		Source:   account,
		Code:     code,
		CodeType: "user_active_code",
	}
	// key 校验通过,表示
	if redis.CheckUserVerificationCode(ctx, uv) {
		updateUser := new(model.User)
		updateUser.ID = user.ID
		updateUser.SetPassword(password)
		updateUser.SetUserStatusNormal()
		err := users.Update(updateUser)
		if err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("code is error")
	}

}

// 新用户注册
func (s userService) Register(ctx context.Context, user *model.User) error {
	users := s.factory.Users()
	if s.CheckAccountExist(ctx, user.Account) {
		return errors.New("account is exist")
	}
	err := users.Create(user)
	if err != nil {
		return err
	}
	return nil
}

// 检查账号是否已存在
func (s userService) CheckAccountExist(ctx context.Context, account string) bool {
	users := s.factory.Users()
	user, err := users.QueryByAccount(account)

	if user != nil || err != nil {
		return true
	}
	return false

}

// 管理员创建用户
func (s userService) CreateUser(ctx context.Context, user *model.User) error {
	users := s.factory.Users()
	if s.CheckAccountExist(ctx, user.Account) {
		return errors.New("account is exist")
	}
	user.SetPassword("a1234567")
	err := users.Create(user)
	if err != nil {
		return err
	}
	return nil
}

// 更新用户信息
func (s userService) UpdateUserInfo(ctx context.Context, user *model.User) error {
	users := s.factory.Users()
	err := users.Update(user)
	if err != nil {
		return err
	}
	return nil
}

// 修改密码
func (s userService) ChangeUserPassword(ctx context.Context, userAccount, oldPassword, newPassword string) error {
	users := s.factory.Users()
	userRes, err := users.QueryByAccount(userAccount)
	if err != nil {
		return err
	}
	if userRes == nil {
		return errors.New("user not found")
	}
	if userRes.CheckPassword(oldPassword) {
		userRes.SetPassword(newPassword)
		err = users.Update(userRes)
		if err != nil {
			return err
		}
	}
	return nil
}

// 用户添加到用户组
func (s userService) AddUserToGroup(ctx context.Context, groupVo *vo.UserGroupVo) error {
	var userId, groupId = groupVo.UserId, groupVo.GroupId
	group := s.factory.GroupUsers()
	userGroup := &model.GroupUser{}
	userGroup.GroupId = groupId
	userGroup.UserId = userId
	err := group.Create(userGroup)
	if err != nil {
		return err
	}
	return nil
}

func (s userService) RemoveUserGroup(ctx context.Context, groupVo *vo.UserGroupVo) error {
	var userId, groupId = groupVo.UserId, groupVo.GroupId
	group := s.factory.GroupUsers()
	userGroup := &model.GroupUser{}
	userGroup.GroupId = groupId
	userGroup.UserId = userId
	return group.RemoveUserGroup(userId, groupId)
}

func (s userService) RemoveUserTenant(ctx context.Context, groupVo *vo.UserGroupVo) error {
	var userId, tenantId = groupVo.UserId, groupVo.TenantId
	group := s.factory.GroupUsers()
	userGroup := &model.GroupUser{}
	userGroup.TenantId = tenantId
	userGroup.UserId = userId
	return group.RemoveUserTenant(userId, tenantId)
}

// 用户添加到租户 添加到租户 都是添加到默认的 guest 组
func (s userService) AddUserToTenant(ctx context.Context, vo *vo.UserGroupVo) error {
	userId, tenantId := vo.UserId, vo.TenantId
	groups := s.factory.Groups()
	defaultGroup, err := groups.QueryTenantDefaultGroup(tenantId)
	if err != nil {
		return err
	}
	if defaultGroup == nil {
		defaultGroup, err = groups.CreateGuestGroup(tenantId)
		if err != nil {
			return err
		}
	}
	groupUsers := s.factory.GroupUsers()
	userGroup := &model.GroupUser{}
	userGroup.GroupId = defaultGroup.ID
	userGroup.UserId = userId
	err = groupUsers.Create(userGroup)
	return err
}

func (s userService) ListUserTenants(ctx context.Context, userId uint) ([]*model.TenantShowVo, error) {
	tenants := s.factory.Tenants()
	return tenants.ListUserTenants(userId)
}

func (s userService) ListUserGroups(ctx context.Context, groupVo *vo.UserGroupVo) ([]*model.Group, error) {
	userId, tenantId := groupVo.UserId, groupVo.TenantId
	groups := s.factory.Groups()
	return groups.QueryUserGroupList(tenantId, userId)
}

func (s userService) QueryUserInfoById(ctx context.Context, userId uint) (*model.User, error) {
	users := s.factory.Users()
	return users.FindById(userId)
}
