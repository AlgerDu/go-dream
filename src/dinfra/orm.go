package dinfra

import "gorm.io/gorm"

type (
	// 数据库
	Database struct {
		*gorm.DB // 简单内嵌 Gorm
	}

	// 数据库 provider
	DatabaseProvider interface {
		// 通过服务 name 获取对应的数据库实例访问对象
		Provide(name string) (*Database, error)
	}
)
