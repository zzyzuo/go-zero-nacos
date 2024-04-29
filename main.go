package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zzyzuo/go-zero-nacos/config"
	"github.com/zzyzuo/go-zero-nacos/nacosx"
	"time"
)

var bootstrapFile = flag.String("f", "nacos.yaml", "the config file")

func main() {
	//解析bootstrap config
	flag.Parse()
	var bootstrapConfig config.BootstrapConfig
	conf.MustLoad(*bootstrapFile, &bootstrapConfig)

	//解析业务配置
	var c config.Config
	nacos := nacosx.NewNacos(bootstrapConfig.NacosConfig)
	serviceConfigContent := nacos.InitConfig(func(data string) {
		err := conf.LoadFromYamlBytes([]byte(data), &c)
		if err != nil {
			panic(err)
		}
	})
	err := conf.LoadFromYamlBytes([]byte(serviceConfigContent), &c)
	if err != nil {
		panic(err)
	}

	for {
		fmt.Printf("c : %+v \n", c)
		time.Sleep(10 * time.Second)
	}
}
