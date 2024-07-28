package mapstructure

// 一些功能修改和添加尽可能的迁移到这个文件

func WithSquash(config *DecoderConfig) {
	config.Squash = true
}
