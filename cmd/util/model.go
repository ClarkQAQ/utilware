package util

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"utilware/dep/shortuuid"
)

const (
	ModelFileName = "mod.json"
)

type Model struct {
	Id          string   `json:"id"`          // 编号
	Name        string   `json:"name"`        // 名称
	URL         string   `json:"url"`         // 地址
	Type        string   `json:"type"`        // 类型 git, local
	Author      string   `json:"author"`      // 作者
	Tag         []string `json:"tag"`         // 标签
	Description string   `json:"description"` // 描述
	Version     string   `json:"version"`     // 版本号
	Updated     bool     `json:"updated"`     // 是否更新
	CreateTime  int64    `json:"create_time"` // 创建时间
	UpdateTime  int64    `json:"update_time"` // 更新时间
}

func ScanModels(p string, f func(m *Model) error) error {
	items, e := os.ReadDir(p)
	if e != nil {
		return fmt.Errorf("read dir failed: %s", e.Error())
	}

	for _, item := range items {
		if !item.IsDir() ||
			strings.HasPrefix(item.Name(), ".") {
			continue
		}

		switch item.Name() {
		case "dep", "cmd", "test":
			continue
		case "util":
			if e := ScanModels(filepath.Join(p, item.Name()), f); e != nil {
				return fmt.Errorf("read models %s failed: %s",
					item.Name(), e.Error())
			}

			continue
		}

		fileName := filepath.Join(p, item.Name(), ModelFileName)

		fileBytes, e := os.ReadFile(fileName)
		if e != nil && !os.IsNotExist(e) {
			return fmt.Errorf("read file failed: %s", e.Error())
		}

		model := &Model{}

		if fileBytes != nil {
			if e := json.Unmarshal(fileBytes, model); e != nil {
				return fmt.Errorf("unmarshal failed: %s %s", fileName, e.Error())
			}
		}

		if strings.TrimSpace(model.Id) == "" {
			model = &Model{
				Id:         shortuuid.New(),
				Name:       item.Name(),
				URL:        "",
				Author:     "anonymous",
				Version:    "0.0.0",
				CreateTime: time.Now().Unix(),
				UpdateTime: time.Now().Unix(),
			}
		}

		model.Name = item.Name()

		if e := f(model); e != nil {
			return e
		}

		fileBytes, e = json.MarshalIndent(model, "", "    ")
		if e != nil {
			return fmt.Errorf("marshal %s failed: %s", item.Name(), e.Error())
		}

		if e := os.WriteFile(fileName, fileBytes, os.ModePerm); e != nil {
			return fmt.Errorf("write %s file failed: %s", item.Name(), e.Error())
		}

	}
	return nil
}

func GetModels(p string) ([]*Model, error) {
	models := []*Model{}

	if e := ScanModels(p, func(m *Model) error {
		models = append(models, m)
		return nil
	}); e != nil {
		return nil, e
	}

	return models, nil
}
