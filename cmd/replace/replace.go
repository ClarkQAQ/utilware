package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"utilware/logger"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("lack of config file, eg: ${path_of_file} ${old_string} ${new_string}")
		os.Exit(1)
	}

	oldBytes, newBytes := []byte(os.Args[2]), []byte(os.Args[3])

	filepath.Walk(os.Args[1], func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		b, e := os.ReadFile(path)
		if e != nil {
			return e
		}

		if bytes.Contains(b, oldBytes) {
			if e := os.WriteFile(path,
				bytes.ReplaceAll(b, oldBytes, newBytes), // 替换目标字符串
				info.Mode()); e != nil {
				return e
			}

			logger.Info("SUCCESS: ", path)
		}

		return nil
	})
}
