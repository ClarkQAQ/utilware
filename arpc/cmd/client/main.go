package main

import (
	"fmt"
	"net"
	"net/url"
	"sync"
	"time"

	"utilware/logger"

	"utilware/dep/x/net/websocket"

	"utilware/arpc"
	"utilware/arpc/extension/middleware/coder/gob4aes"
	"utilware/arpc/extension/middleware/router"
	"utilware/arpc/log"
)

func main() {
	log.SetLevel(log.LevelInfo | log.LevelWarn | log.LevelError)
	urlString := "ws://localhost:9001/rpc"

	client, e := arpc.NewClient(func() (net.Conn, error) {
		u, e := url.Parse(urlString)
		if e != nil {
			return nil, e
		}

		ws, e := websocket.Dial(urlString, "", "http://"+u.Host)
		if e != nil {
			return nil, e
		}
		ws.PayloadType = websocket.BinaryFrame

		return ws, nil
	})
	if e != nil {
		logger.Fatal("connect to server failed: %s", e.Error())
	}

	client.Handler.UseCoder(gob4aes.New([]byte("1234567890123456")))

	client.Handler.Use(router.Logger())

	defer client.Stop()

	wg := sync.WaitGroup{}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			req := fmt.Sprintf("114514 %d", i)
			rsp := ""

			if e = client.Call("/echo", &req, &rsp, time.Second*5); e != nil {
				logger.Fatal("call failed: %s", e.Error())
			} else {
				logger.Info("call success: %s", rsp)
			}
		}(i)
	}

	wg.Wait()

}
