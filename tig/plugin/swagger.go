package plugin

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"path"
	"strings"
	"utilware/tig"
	"utilware/tig/middleware"
)

//go:embed public
var swagfs embed.FS

type SwaggerBasicAuth struct {
	Name     string
	Password string
}

func Swagger(relativePath string, swaggerJson []byte, auth *SwaggerBasicAuth) tig.PluginHandler {
	static, _ := fs.Sub(swagfs, "public/swagger")
	fileServer := http.StripPrefix(relativePath, http.FileServer(http.FS(static)))

	return func(g *tig.RouterGroup) error {
		r := g.NewGroup(relativePath)

		if auth != nil {
			r.Use(func(c *tig.Context) {
				u, p, ok := c.Req.BasicAuth()
				if !ok || u != auth.Name || p != auth.Password {
					c.Writer.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
					// 401 状态码
					c.Writer.WriteHeader(http.StatusUnauthorized)
					c.End()
				}
				c.Next()
			})
		}

		r.ANY("/*filepath", func(c *tig.Context) {
			mime, e := middleware.QueryMimeType(strings.TrimPrefix(path.Ext(c.Param("filepath")), "."))
			if e != nil {
				c.SetHeader(tig.HeaderContentType, "text/plain")
			} else {
				c.SetHeader(tig.HeaderContentType, fmt.Sprintf("%s; charset=utf-8", mime))
			}

			fileServer.ServeHTTP(c.Writer, c.Req)
		})
		r.ANY("/", func(c *tig.Context) {
			c.SetHeader(tig.HeaderContentType, "text/html; charset=UTF-8")
			f, e := static.Open("index.html")
			if e != nil {
				c.String(http.StatusNotFound, "404 not found")
				return
			}

			b, e := io.ReadAll(f)
			if e != nil {
				c.String(http.StatusBadRequest, "400 bad request")
				return
			}

			c.String(200, strings.ReplaceAll(string(b), "[path]", relativePath))
		})
		r.ANY("/swagger.json", func(c *tig.Context) {
			c.SetHeader(tig.HeaderContentType, "application/json; charset=utf-8")
			c.Data(200, swaggerJson)
		})

		return nil
	}
}
