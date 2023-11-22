package dinfra

import "github.com/mitchellh/mapstructure"

type (
	// 对应微服务一个实例
	MicroService struct {
		ID       string            // 实例唯一 ID
		Name     string            // 服务名称
		Address  string            // 可访问地址
		Port     int               // 可访问端口
		Env      string            // 归属的环境
		Tags     []string          // 标签
		Metadata map[string]string // 服务注册时附加数据
	}

	// 注册中心
	ServiceRegister interface {
		// 获取 name 服务对应的信息
		Get(name string) (*MicroService, error)

		// 向注册中心注册服务
		Register(service *MicroService) error
	}
)

// 将附加数据绑定到结构体
func (ms *MicroService) MetaBindToStruct(dst any) error {
	return mapstructure.WeakDecode(ms.Metadata, dst)
}
