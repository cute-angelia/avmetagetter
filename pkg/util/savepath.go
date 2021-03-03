package util

import "path"

type save struct {
}

func NewSave() *save {
	return &save{}
}

// 获取保存路径
// 1. 非转移状态，保存当前文件夹
// 2. 转移状态，保存在 success 目录
func (self *save) GetSavePathInfo(cfg *ConfigStruct, savepath string) string {
	if !cfg.Path.IsTransfer {
		if len(cfg.Path.PathIn) > 0 {
			return cfg.Path.PathIn + "/" + path.Base(savepath)
		} else {
			return path.Base(savepath)
		}
	}
	return savepath
}
