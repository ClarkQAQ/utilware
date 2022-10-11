package main

import (
	"fmt"
	"utilware/cmd/util"
	"utilware/logger"
)

func main() {
	if e := util.ScanModels(".", func(m *util.Model) error {
		logger.Info("%s %s", m.Id, m.Name)

		v := "Y"
		fmt.Print("是否Git [Y/n]:")
		fmt.Scanln(&v)

		if v == "Y" || v == "y" || v == "" {
			m.Type = "git"
		} else {
			m.Type = "local"
		}

		return nil
	}); e != nil {
		logger.Fatal("scan models failed: %s", e.Error())
	}
}
