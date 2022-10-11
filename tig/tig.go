package tig

import (
	"net"
	"net/http"
	"sync"
	"time"
)

// 主框架结构
type Tig struct {
	*RouterGroup

	pLocker *sync.RWMutex
	p       *Pattern
}

func New() *Tig {
	t := &Tig{}

	t.pLocker = &sync.RWMutex{}
	t.p = newPattern()
	t.RouterGroup = &RouterGroup{t, "/", nil, nil}

	return t
}

func (t *Tig) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := t.newContext(w, req)
	t.p.mainHandle(c)

	if c.index < -1 {
		panic(nil)
	}

	if c.writer.exported {
		return
	}

	for k, v := range c.writer.header {
		w.Header()[k] = v
	}

	w.WriteHeader(c.writer.statusCode)

	if _, e := w.Write(c.writer.bodyBuffer.Bytes()); e != nil {
		panic(e)
	}
}

func (t *Tig) Run(addr string) (*http.Server, error) {
	net, e := net.Listen("tcp", addr)
	if e != nil {
		return nil, e
	}

	return t.RunWithListener(net)
}

func (t *Tig) RunWithListener(l net.Listener) (*http.Server, error) {
	http := &http.Server{Handler: t}
	ch := make(chan error)
	defer close(ch)

	go func() {
		if e := http.Serve(l); e != nil {
			ch <- e
			return
		}
	}()

	select {
	case e := <-ch:
		return http, e
	case <-time.After(time.Millisecond * 100):
		return http, nil
	}
}
