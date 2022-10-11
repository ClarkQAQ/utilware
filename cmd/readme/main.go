package main

import (
	"bytes"
	_ "embed"
	"os"
	"text/template"
	"utilware/cmd/util"
	"utilware/logger"
)

var (
	//go:embed readme.tmpl
	readmeTemplate string
)

func main() {
	tmpl, e := template.New("readme").Parse(readmeTemplate)
	if e != nil {
		logger.Fatal("parse template failed: %s", e.Error())
	}

	models, e := util.GetModels(".")
	if e != nil {
		logger.Fatal("scan models failed: %s", e.Error())
	}

	buf := bytes.NewBuffer(nil)

	if e := tmpl.Execute(buf, models); e != nil {
		logger.Fatal("execute template failed: %s", e.Error())
	}

	if e := os.WriteFile("README.md", buf.Bytes(), os.ModePerm); e != nil {
		logger.Fatal("write file failed: %s", e.Error())
	}
}
