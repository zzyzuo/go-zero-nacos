package nacosx

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/zeromicro/go-zero/core/logx"
)

type (
	ConfigCenter interface {
		InitConfig(listenConfigCallback ListenConfig) string
	}

	ListenConfig func(data string)

	NacosServerConfig struct {
		IpAddr string
		Port   uint64
	}
	NacosClientConfig struct {
		NamespaceId         string
		TimeoutMs           uint64
		NotLoadCacheAtStart bool
		LogDir              string
		CacheDir            string
		LogLevel            string
	}
	NacosConfig struct {
		ServerConfigs []NacosServerConfig
		ClientConfig  NacosClientConfig
		DataId        string
		Group         string
	}
	defaultNacos struct {
		cfg NacosConfig
	}
)

func NewNacos(cfg NacosConfig) ConfigCenter {
	return &defaultNacos{cfg: cfg}
}
func (n *defaultNacos) InitConfig(listenConfigCallback ListenConfig) string {
	var sc []constant.ServerConfig
	if len(n.cfg.ServerConfigs) == 0 {
		panic("nacos 配置错误")
	}
	for _, serveConfig := range n.cfg.ServerConfigs {
		sc = append(sc, constant.ServerConfig{
			Port:   serveConfig.Port,
			IpAddr: serveConfig.IpAddr,
		},
		)
	}

	cc := constant.ClientConfig{
		NamespaceId:         n.cfg.ClientConfig.NamespaceId,
		TimeoutMs:           n.cfg.ClientConfig.TimeoutMs,
		NotLoadCacheAtStart: n.cfg.ClientConfig.NotLoadCacheAtStart,
		LogDir:              n.cfg.ClientConfig.LogDir,
		CacheDir:            n.cfg.ClientConfig.CacheDir,
		LogLevel:            n.cfg.ClientConfig.LogLevel,
	}

	pa := vo.NacosClientParam{
		ClientConfig:  &cc,
		ServerConfigs: sc,
	}
	configClient, err := clients.NewConfigClient(pa)
	if err != nil {
		panic(err)
	}

	//获取配置中心内容
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: n.cfg.DataId,
		Group:  n.cfg.Group,
	})
	if err != nil {
		panic(err)
	}
	if len(content) == 0 {
		panic("配置文件读取为空")
	}
	//设置配置监听
	if err = configClient.ListenConfig(vo.ConfigParam{
		DataId: n.cfg.DataId,
		Group:  n.cfg.Group,
		OnChange: func(namespace, group, dataId, data string) {
			//配置文件产生变化就会触发
			logx.Infof("配置文件发生了变化:%s", data)
			if len(data) == 0 {
				logx.Error("配置文件读取为空")
				return
			}
			listenConfigCallback(data)
		},
	}); err != nil {
		panic(err)
	}

	return content
}
