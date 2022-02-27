package auth_server

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	jwt2 "github.com/golang-jwt/jwt/v4"
	"project/pkg/utils"
	user "project/proto/faner_user_proto"
	"time"
)

type AuthServer struct {
	user.UnimplementedAuthServer
}

// NewAuthService 登录服务
func NewAuthServer() *AuthServer {
	return &AuthServer{}
}

// Login 登录
func (s *AuthServer) Login(ctx context.Context, in *user.LoginRequest) (*user.LoginResponse, error) {
	myClaims := utils.MyClaims{
		"test",
		jwt2.RegisteredClaims{
			ExpiresAt: jwt2.NewNumericDate(time.Now().Add(time.Hour * 2)), //设置JWT过期时间,此处设置为2小时
			Issuer:    "test",                                             //设置签发人
		},
	}
	claims := jwt2.NewWithClaims(jwt2.SigningMethodHS256, myClaims)
	//加盐
	signedString, err := claims.SignedString([]byte("testKey"))
	if err != nil {
		return nil, err
	}

	return &user.LoginResponse{
		Code:    0,
		Message: signedString,
	}, nil
}

// Logout 退出登录
func (s *AuthServer) Logout(ctx context.Context, in *user.LogoutRequest) (*user.LogoutResponse, error) {
	token, ok := jwt.FromContext(ctx)
	if !ok {
		panic("sssss")
	}
	fmt.Println(token)
	return &user.LogoutResponse{
		Code:    0,
		Message: "",
	}, nil
}
