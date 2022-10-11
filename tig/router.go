package tig

import (
	"net/http"
	"path"
	"reflect"
	"strings"
)

// 路由组
// 群组可以嵌套, 并且可以有自己的中间件
// 中间件将是递归生效
type RouterGroup struct {
	t          *Tig          // 主框架
	prefix     string        // 该组的前缀
	middleware []HandlerFunc // 组内中间件
	field      *RouterGroup  // 上级组
}

// 请求路径
// 每个请求路径只有一个节点
// 顶级路径是/
type Pattern struct {
	field   *Pattern            // 上级节点
	part    string              // 当前节点路径
	method  map[string]*Method  // 请求方法
	pattern map[string]*Pattern // 下级路径
}

// 请求方法
// 上级是路径, 然后同时保留路由组
type Method struct {
	method      string       // 请求方法
	pattern     *Pattern     // 路径
	routerGroup *RouterGroup // 请求组
	handler     HandlerFunc  // 请求处理
}

type PluginHandler func(t *RouterGroup) error

// 新建节点
func newPattern() *Pattern {
	return &Pattern{
		field:   nil,
		part:    "",
		method:  make(map[string]*Method),
		pattern: make(map[string]*Pattern),
	}
}

func parsePattern(pattern string) []string {
	patterns := strings.Split(pattern, "/")
	parts := []string{}

	for i := 0; i < len(patterns); i++ {
		if patterns[i] != "" {
			parts = append(parts, patterns[i])

			// 判断到通配符就停止
			// 谁也不想处理两个通配符吧
			if patterns[i][0] == '*' {
				break
			}
		}
	}

	return parts
}

// 添加路径
func (p *Pattern) addPattern(method string, pattern string, g *RouterGroup, h HandlerFunc) {
	if strings.TrimSpace(method) == "" ||
		strings.TrimSpace(pattern) == "" ||
		g == nil || h == nil {
		return
	}

	parts := parsePattern(path.Join(g.prefix, pattern))

	for i := 0; i < len(parts); i++ {
		t := parts[i]

		// 判断是否通配
		if t[0] == ':' || t[0] == '*' {
			t = "*"
		}

		// 路由结构是否存在
		if _, ok := p.pattern[t]; !ok {
			p.pattern[t] = newPattern()
			p.pattern[t].field = p
		}

		// 指向
		p = p.pattern[t]
		p.part = parts[i]
	}

	// 添加请求
	p.method[method] = &Method{
		method:      method,
		routerGroup: g,
		pattern:     p,
		handler:     h,
	}
}

func (p *Pattern) getRoute(method string, pattern string) (*Method, map[string]string) {
	parts := parsePattern(pattern)
	params := make(map[string]string)

	for i := 0; i < len(parts); i++ {
		t := parts[i]

		// 查找以及通配
		if _, ok := p.pattern[t]; !ok {
			t = "*"
			if _, ok := p.pattern[t]; !ok {
				return nil, nil
			}
		}

		p = p.pattern[t]

		if p.part[0] == ':' {
			params[p.part[1:]] = parts[i]
		}

		if p.part[0] == '*' {
			params[p.part[1:]] = strings.Join(parts[i:], "/")
			i = len(parts) - 1
		}
	}

	if op, ok := p.method[method]; ok {
		return op, params
	}

	if op, ok := p.method["*"]; ok {
		return op, params
	}

	return nil, nil
}

func (g *RouterGroup) NewGroup(prefix string) *RouterGroup {
	return &RouterGroup{
		t:      g.t,
		prefix: g.prefix + prefix,
		field:  g,
	}
}

func (g *RouterGroup) Group(prefix string, f func(g *RouterGroup)) {
	f(g.NewGroup(prefix))
}

func (g *RouterGroup) upMiddlewares() (handles []HandlerFunc) {
	if g.field != nil {
		handles = append(g.field.upMiddlewares(), g.middleware...)
	}

	return handles
}

func (p *Pattern) mainHandle(c *Context) {
	c.handlerList = append(c.handlerList, c.t.middleware...)

	n, params := p.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		c.handlerList = append(c.handlerList, n.routerGroup.upMiddlewares()...)
		c.handlerList = append(c.handlerList, n.handler)
	} else {
		c.handlerList = append(c.handlerList, func(c *Context) {
			c.Status(http.StatusNotFound)
			c.Sprintf(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
			c.End()
		})
	}

	c.Next()
}

func (g *RouterGroup) Use(middleware ...HandlerFunc) {
	g.middleware = append(g.middleware, middleware...)
}

func (g *RouterGroup) Plugin(h PluginHandler) error {
	return h(g)
}

func (g *RouterGroup) Method(method, pattern string, handler HandlerFunc) {
	g.t.p.addPattern(method, pattern, g, handler)
}

func (g *RouterGroup) GET(pattern string, handler HandlerFunc) {
	g.Method("GET", pattern, handler)
}

func (g *RouterGroup) POST(pattern string, handler HandlerFunc) {
	g.Method("POST", pattern, handler)
}

func (g *RouterGroup) PUT(pattern string, handler HandlerFunc) {
	g.Method("PUT", pattern, handler)
}

func (g *RouterGroup) DELETE(pattern string, handler HandlerFunc) {
	g.Method("DELETE", pattern, handler)
}

func (g *RouterGroup) OPTION(pattern string, handler HandlerFunc) {
	g.Method("OPTION", pattern, handler)
}

func (g *RouterGroup) HEAD(pattern string, handler HandlerFunc) {
	g.Method("HEAD", pattern, handler)
}

func (g *RouterGroup) PATCH(pattern string, handler HandlerFunc) {
	g.Method("PATCH", pattern, handler)
}

func (g *RouterGroup) ANY(pattern string, handler HandlerFunc) {
	g.Method("*", pattern, handler)
}

func (g *RouterGroup) Object(pattern string, object interface{}) {
	// 通过反射获取结构体
	v := reflect.ValueOf(object)
	t := v.Type()

	if v.Kind() == reflect.Struct {
		newValue := reflect.New(t)
		newValue.Elem().Set(v)
		v = newValue
		t = v.Type()
	}

	g.Group(pattern, func(r *RouterGroup) {
		for i := 0; i < v.NumMethod(); i++ {
			methodName := strings.ToUpper(t.Method(i).Name)

			switch methodName {
			case "MIDDLEWARE":
				if h, ok := v.Method(i).Interface().(func(*Context)); ok {
					r.Use(h)
				}
			case "PLUGIN":
				if h, ok := v.Method(i).Interface().(func(t *RouterGroup) error); ok {
					r.Plugin(h)
				}
			default:
				if h, ok := v.Method(i).Interface().(func(*Context)); ok {
					r.Method(methodName, "/", h)
				}
			}
		}
	})
}
