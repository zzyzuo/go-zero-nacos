package config

import (
	"github.com/zeromicro/go-zero/rest"
	"go-zero-nacos/nocosx"
)

type Config struct {
	rest.RestConf
}
type BootstrapConfig struct {
	NacosConfig nocosx.NacosConfig
}
