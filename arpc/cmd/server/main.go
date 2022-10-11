package main

import (
	"net/http"
	"utilware/logger"

	"utilware/arpc"
	"utilware/arpc/extension/middleware/coder/gob4aes"
	"utilware/arpc/extension/middleware/router"
	"utilware/arpc/extension/protocol/websocket"
	"utilware/arpc/log"
)

func main() {
	server := arpc.NewServer()

	log.SetLevel(log.LevelInfo | log.LevelWarn | log.LevelError)

	server.Handler.UseCoder(gob4aes.New([]byte("1234567890123456")))

	server.Handler.Handle("/echo", func(c *arpc.Context) {
		rsp := ""
		if e := c.Bind(&rsp); e != nil {
			c.Error(e)
		}

		c.Write(rsp)
	})

	server.Handler.Use(router.Logger())

	ln, _ := websocket.Listen("0.0.0.0:0", nil)
	http.HandleFunc("/rpc", ln.(*websocket.Listener).Handler)

	go func() {
		if e := http.ListenAndServe("localhost:9001", nil); e != nil {
			logger.Fatal("listen failed: %s", e.Error())
		}
	}()

	if e := server.Serve(ln); e != nil {
		logger.Fatal("serve failed: %s", e.Error())
	}
}
