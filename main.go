package main

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/liangguifeng/kratos-app"
	"github.com/liangguifeng/kratos-app/app"
	"google.golang.org/grpc"
	"project/startup"
)

// 此处定义当前服务的名字
const APP_NAME = "faner-user"

func main() {
	runer, err := app.NewRunner(&kratos.Application{
		Name:       APP_NAME,
		Type:       kratos.APP_TYPE_GRPC,
		LoadConfig: startup.LoadConfig,
		SetupVars:  startup.SetupVars,
		RegisterCallback: map[kratos.CallbackPos]func() error{
			kratos.POS_LOAD_CONFIG: startup.LoadConfigCallback,
			kratos.POS_SETUP_VARS:  startup.SetupVarsCallback,
			kratos.POS_NEW_RUNNER:  startup.RunNewRunnerCallback,
		},
	})
	if err != nil {
		log.Fatalf("app.NewRunner err: %v", err)
	}

	err = runer.ListenGRPCServer(&kratos.GRPCApplication{
		RegisterGRPCServer:      startup.RegisterGRPCServer,
		RegisterHttpRoute:       startup.RegisterHTTPServer,
		UnaryServerInterceptors: []grpc.UnaryServerInterceptor{},
	})
	if err != nil {
		log.Fatalf("runer.ListenGRPCServer err: %v", err)
	}
}
