package config

import (
	"github.com/zeromicro/go-zero/rest"
	"go-zero-nacos/nacosx"
)

type Config struct {
	rest.RestConf
}
type BootstrapConfig struct {
	NacosConfig nacosx.NacosConfig
}
