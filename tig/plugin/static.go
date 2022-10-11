package plugin

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"utilware/tig"
	"utilware/tig/middleware"
)

func staticFile(c *tig.Context, vfs fs.FS, filePath string) bool {
	file, e := vfs.Open(filePath)
	if e != nil {
		c.String(http.StatusNotFound, "404 page not found")
		return false
	}

	defer file.Close()

	filestat, e := file.Stat()
	if e != nil {
		c.String(http.StatusBadRequest, "400 bad request")
		return false
	}

	if filestat.IsDir() {
		c.String(http.StatusForbidden, "403 forbidden")
		return false
	}

	b, e := io.ReadAll(file)
	if e != nil {
		c.String(http.StatusBadRequest, "400 bad request")
		return false
	}

	c.Status(http.StatusOK)
	c.Writer.Write(b)
	return true
}

// 静态文件
// @param rootPrefixPath: 静态文件根目录 (如果就是根目录，传入空字符串即可)
// @param relativePath: 路由路径
// @param vfs: 静态文件系统
// @return: tig.PluginHandler
func Static(rootPrefixPath, relativePath string, vfs fs.FS) tig.PluginHandler {
	urlPattern := path.Join(relativePath, "/*filepath")

	return func(g *tig.RouterGroup) error {
		g.ANY(urlPattern, func(c *tig.Context) {
			filePath := path.Join(rootPrefixPath, c.Param("filepath"))

			if !staticFile(c, vfs, filePath) {
				c.End()
			}

			mime, e := middleware.QueryMimeType(strings.TrimPrefix(path.Ext(c.Param("filepath")), "."))
			if e != nil {
				c.SetHeader(tig.HeaderContentType, "text/plain")
			} else {
				c.SetHeader(tig.HeaderContentType, fmt.Sprintf("%s; charset=utf-8", mime))
			}
		})

		return nil
	}
}

func StaticWithRouterMap(rootPrefixPath string, vfs fs.FS, fileMap map[string]string) tig.PluginHandler {
	return func(g *tig.RouterGroup) error {
		for k, v := range fileMap {
			filePath := filepath.Join(rootPrefixPath, v)
			mime, e := middleware.QueryMimeType(strings.TrimPrefix(path.Ext(v), "."))
			if e != nil {
				mime = "text/plain"
			}

			g.ANY(k, func(c *tig.Context) {
				if !staticFile(c, vfs, filePath) {
					c.End()
				}

				c.SetHeader(tig.HeaderContentType, fmt.Sprintf("%s; charset=utf-8", mime))
			})
		}

		return nil
	}
}
