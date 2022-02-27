package startup

import (
	"context"
	"github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	jwt2 "github.com/golang-jwt/jwt/v4"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"os"
	"project/vars"
)

func Init() (*kratos.App, error) {
	httpSrv := http.NewServer(
		http.Address(":8000"),
		http.Middleware(
			recovery.Recovery(),
			selector.Server(
				jwt.Server(func(token *jwt2.Token) (interface{}, error) {
					return []byte(vars.JwtSetting.SecretKey), nil
				}, jwt.WithSigningMethod(jwt2.SigningMethodHS256)),
			).Match(newWhiteListMatcher()).Build(),
		),
	)
	grpcSrv := grpc.NewServer(
		grpc.Address(":9000"),
		grpc.Middleware(
			recovery.Recovery(),
			jwt.Server(func(token *jwt2.Token) (interface{}, error) {
				return []byte(vars.JwtSetting.SecretKey), nil
			}),
		),
	)

	err := RegisterHTTP(httpSrv)
	if err != nil {
		return nil, err
	}

	err = RegisterGRPC(grpcSrv)
	if err != nil {
		return nil, err
	}

	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ServerConfigs: []constant.ServerConfig{
				*constant.NewServerConfig("120.24.6.140", 8848),
			},
		},
	)
	if err != nil {
		return nil, err
	}

	// 引入配置
	err = LoadConfig()
	if err != nil {
		return nil, err
	}

	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", "2022_02_28_002849",
		"service.name", "faner-user",
		"service.version", "1.0",
		"trace_id", tracing.TraceID(),
		"span_id", tracing.SpanID(),
	)
	app := kratos.New(
		kratos.ID("2022_02_28_002849"),
		kratos.Name("faner-user"),
		kratos.Version("1.0"),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			httpSrv,
			grpcSrv,
		),
		kratos.Registrar(nacos.New(client)),
	)

	return app, nil
}

func newWhiteListMatcher() selector.MatchFunc {
	whiteList := make(map[string]struct{})
	whiteList["/api.users.v1.Auth/Login"] = struct{}{}
	return func(ctx context.Context, operation string) bool {
		if _, ok := whiteList[operation]; ok {
			return false
		}
		return true
	}
}
