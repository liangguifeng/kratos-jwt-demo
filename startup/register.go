package startup

import (
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	user "project/proto/faner_user_proto"
	"project/server/auth_server"
)

// RegisterHTTPServer 此处注册http接口
func RegisterHTTPServer(httpSrv *http.Server) error {
	user.RegisterAuthHTTPServer(httpSrv, auth_server.NewAuthServer())

	return nil
}

// RegisterGRPCServer 此处注册grpc接口
func RegisterGRPCServer(grpcSrv *grpc.Server) error {
	user.RegisterAuthServer(grpcSrv, auth_server.NewAuthServer())

	return nil
}
