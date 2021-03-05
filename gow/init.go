package gow

import (
	"embed"
	"errors"
	"html/template"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"time"
)

// HandlerFunc defines the request handler used by gee
type HandlerFunc func(*Context)

// Engine implement the interface of ServeHTTP
type (
	RouterGroup struct {
		prefix      string
		middlewares []HandlerFunc // support middleware
		parent      *RouterGroup  // support nesting
		engine      *Engine       // all groups share a Engine instance
	}

	Engine struct {
		*RouterGroup
		http          *http.Server
		https         *http.Server
		router        *router
		groups        []*RouterGroup     // store all groups
		htmlTemplates *template.Template // for html render
		htmlFiles     map[string][]byte  // for html render
		funcMap       template.FuncMap   // for html render
	}
)

// New is the constructor of gee.Engine
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// 初始化路由
func (r *Engine) NewRouter() {
	r.router = &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// Group is defined to create a new RouterGroup
// remember all groups share the same Engine instance
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// Use is defined to add middleware to the group
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp

	group.engine.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

//add Get and Post request
func (group *RouterGroup) RESTFUL(pattern string, handler HandlerFunc) {
	s := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"}
	for i := 0; i < len(s); i++ {
		group.addRoute(s[i], pattern, handler)
	}
}

//add Get and Post request
func (group *RouterGroup) CustomMethods(methods []string, pattern string, handler HandlerFunc) {
	for i := 0; i < len(methods); i++ {
		group.addRoute(strings.ToUpper(methods[i]), pattern, handler)
	}
}

//Custom request
func (group *RouterGroup) Custom(method, pattern string, handler HandlerFunc) {
	group.addRoute(strings.ToUpper(method), pattern, handler)
}

// create static handler
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))

	return func(c *Context) {
		// Check if file exists and/or if we have permission to access it
		//fmt.Println(c.Param("filepath"))
		file, err := fs.Open(c.Param("filepath"))
		if err != nil {
			c.String(404, "404 page not found")
			return
		}

		filestat, e := file.Stat()
		if e != nil {
			c.Status(http.StatusBadRequest)
			c.String(400, "400 bad request")
			return
		}

		if filestat.IsDir() {
			c.String(403, "403 forbidden")
			return
		}

		c.StatusCode = 200
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// serve static files
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	// Register GET handlers
	group.GET(urlPattern, handler)
}

func (group *RouterGroup) FsStatic(relativePath, headPath string, fs embed.FS) {
	urlPattern := path.Join(relativePath, "/*filepath")
	//absolutePath := path.Join(group.prefix, relativePath)

	// Register GET handlers
	group.GET(urlPattern, func(c *Context) {
		file, e := fs.Open(headPath + c.Param("filepath"))
		if e != nil {
			c.String(404, "404 page not found")
			return
		}

		filestat, e := file.Stat()
		if e != nil {
			c.Status(http.StatusBadRequest)
			c.String(400, "400 bad request")
			return
		}

		if filestat.IsDir() {
			c.String(403, "403 forbidden")
			return
		}

		bt, e := fs.ReadFile(headPath + c.Param("filepath"))
		if e != nil {
			c.String(404, "404 page not found")
			return
		}

		c.Data(200, bt)
		return
	})
}

// for custom render function
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

func (engine *Engine) LoadTemplateGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlFiles = make(map[string][]byte)

	rd, _ := ioutil.ReadDir(pattern)
	for _, fi := range rd {
		if str, err := ioutil.ReadFile(pattern + "/" + fi.Name()); err == nil {
			engine.htmlFiles[fi.Name()] = str
		}
	}
}

func (engine *Engine) HasHTML(name string) bool {
	if _, status := engine.htmlFiles[name]; status {
		return true
	}
	return false
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) RunTls(addr string, crt string, key string) (err error) {
	return http.ListenAndServeTLS(addr, crt, key, engine)
}

func (engine *Engine) RunServe(addr string) error {
	if engine.http != nil {
		return errors.New("http server is start!")
	}
	engine.http = &http.Server{Addr: addr, Handler: engine}
	var err error = nil
	go func() {
		err = engine.http.ListenAndServe()
	}()
	time.Sleep(1 * time.Second) //等待1秒钟看看Goroutine有没有报错!
	return err
}

func (engine *Engine) ShutdownServe() error {
	if engine.http == nil {
		return errors.New("http server not start!")
	}
	engine.http = nil
	return engine.http.Shutdown(nil)
}

func (engine *Engine) RunTlsServe(addr, crt, key string) error {
	if engine.https != nil {
		return errors.New("http server is start!")
	}
	engine.https = &http.Server{Addr: addr, Handler: engine}
	var err error = nil
	go func() {
		err = engine.https.ListenAndServeTLS(crt, key)
	}()
	time.Sleep(1 * time.Second) //等待1秒钟看看Goroutine有没有报错!
	return err
}

func (engine *Engine) ShutdownTlsServe() error {
	if engine.https == nil {
		return errors.New("http server not start!")
	}
	engine.https = nil
	return engine.https.Shutdown(nil)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	c.engine = engine
	engine.router.handle(c)
}

func (engine *Engine) RootHandler(f func(*Context) bool) {
	engine.router.root_handler = f
}
