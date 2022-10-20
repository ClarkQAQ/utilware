package main

import (
	"utilware/cmd/util"
	"utilware/logger"
)

func main() {
	if e := util.ScanModels(".", func(m *util.Model) error {
		// 如果不是 git 仓库，跳过
		if m.Type == "local" {
			return nil
		}

		logger.Info("%s %s", m.Id, m.Name)

		return nil
	}); e != nil {
		logger.Fatal("scan models failed: %s", e.Error())
	}
}
