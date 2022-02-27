package startup

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"log"
)

// LoadConfig 加载配置对象映射
func LoadConfig() error {
	// a more graceful way to create naming client
	client1, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig: &constant.ClientConfig{
				NamespaceId:         "public", //namespace id
				TimeoutMs:           5000,
				NotLoadCacheAtStart: true,
				LogDir:              "/tmp/nacos/log",
				CacheDir:            "/tmp/nacos/cache",
				LogLevel:            "debug",
			},
			ServerConfigs: []constant.ServerConfig{
				*constant.NewServerConfig("120.24.6.140", 8848),
			},
		},
	)
	if err != nil {
		log.Panic(err)
	}

	content, err := client1.GetConfig(vo.ConfigParam{DataId: "faner-user", Group: "DEFAULT_GROUP"})
	if err != nil {
		panic(fmt.Errorf("nacos读取配置错误: %s \n", err))
	}
	fmt.Println(content)

	return nil
}
