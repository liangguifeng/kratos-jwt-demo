package main

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	jwt2 "github.com/golang-jwt/jwt/v4"
	"log"
	v1 "project/proto/faner_user_proto"
	"time"
)

func NewWhiteListMatcher() selector.MatchFunc {
	whiteList := make(map[string]struct{})
	whiteList["/api.users.v1.Auth/Login"] = struct{}{}
	return func(ctx context.Context, operation string) bool {
		if _, ok := whiteList[operation]; ok {
			return false
		}
		return true
	}
}

type AuthServer struct {
	v1.UnimplementedAuthServer

	hc v1.AuthClient
}

func NewAuthService() *AuthServer {
	return &AuthServer{}
}

type MyClaims struct {
	UserName string `json:"username"`
	jwt2.RegisteredClaims
}

func (s *AuthServer) Login(ctx context.Context, in *v1.LoginRequest) (*v1.LoginResponse, error) {
	myClaims := MyClaims{
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

	return &v1.LoginResponse{
		Code:    0,
		Message: signedString,
	}, nil
}

func (s *AuthServer) Logout(ctx context.Context, in *v1.LogoutRequest) (*v1.LogoutResponse, error) {
	token, ok := jwt.FromContext(ctx)
	if !ok {
		panic("sssss")
	}
	fmt.Println(token)
	return &v1.LogoutResponse{
		Code:    0,
		Message: "",
	}, nil
}

func main() {
	testKey := "testKey"
	httpSrv := http.NewServer(
		http.Address(":8000"),
		http.Middleware(
			recovery.Recovery(),
			selector.Server(
				jwt.Server(func(token *jwt2.Token) (interface{}, error) {
					return []byte(testKey), nil
				}, jwt.WithSigningMethod(jwt2.SigningMethodHS256)),
			).Match(NewWhiteListMatcher()).Build(),
		),
	)
	grpcSrv := grpc.NewServer(
		grpc.Address(":9000"),
		grpc.Middleware(
			recovery.Recovery(),
			jwt.Server(func(token *jwt2.Token) (interface{}, error) {
				return []byte(testKey), nil
			}),
		),
	)
	//serviceTestKey := "serviceTestKey"
	//con, _ := grpc.DialInsecure(
	//	context.Background(),
	//	grpc.WithEndpoint("dns:///127.0.0.1:9001"),
	//	grpc.WithMiddleware(
	//		jwt.Client(func(token *jwt2.Token) (interface{}, error) {
	//			return []byte(serviceTestKey), nil
	//		}),
	//	),
	//)
	//s := &server{
	//	hc: v1.NewAuthClient(con),
	//}
	v1.RegisterAuthServer(grpcSrv, NewAuthService())
	v1.RegisterAuthHTTPServer(httpSrv, NewAuthService())
	app := kratos.New(
		kratos.Name("v1"),
		kratos.Server(
			httpSrv,
			grpcSrv,
		),
	)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
