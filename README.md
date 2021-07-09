## Golang Utilware

#### 名称灵感来着Aimware, Aimware永远的空枪FW.....

---

> 坚持纯Go拒绝CGO!

##### Utilware 无任何外部依赖, 真正的ALL IN ONE
##### 收集开源项目并去除CGO, 减少依赖后打包....依赖项目统一在 [dep] 

---

#### TODO:

- [ ] 打包工具链
- [ ] 简化打包流程
- [ ] 包的版本管理以及更新订阅

#### Task:
- [ ] Gorm || 7 Day Orm
- [ ] Google QUIC

#### Go Mod ( 使用方式) :

```mod
module ${name: goweb}

go ${go: 1.16}

replace utilware => github.com/ClarkQAQ/utilware ${tag: v0.0.1}

require utilware v0.0.0-00010101000000-000000000000

```

```go
package main

import "utilware/gow"

func main() {
	gow.New().Run(":8080")
}

```



#### 索引:

1. `utilware/bbolt` bolt 纯go高性能kv数据库
2. `utilware/bolthold` bolt 的sql增强....没storm强大
3. `utilware/storm` 你要说这个我就不困了啊,golang 目前最强关系型数据库替代!
4. `utilware/util/safe` 类似java 的 try catch
5. `utilware/util/eodec` 快速编码解码
6. `utilware/util/crypc` 常用加密
7. `utilware/util/sn` 单链表,用于完成一些`[][]byte`完成不了的事情
8. `utilware/util/tsort` 排序拓展
9. `utilware/util/value` 类型转换以及字符串操作
10. `utilware/util/elog` 带Event返回的log库
11. `utilware/util/others` 垃圾桶 (~~不是~~
12. `utilware/util/package` 打包用的工具
13. `utilware/util/weightedrand` 加权随机库, 不过是interface的...爽是肯定的
14. `utilware/decimal` 高精度数字操作
15. `utilware/goquery` golang html 解析,用于爬虫!
16. `utilware/clockwork` 定时任务
17. `utilware/resize` 图片压缩
18. `utilware/request` http网络请求,可搭配sock5或者http代理
19. `utilware/tgbotapi` Telegram机器人API (可在NAT后面使用)
20. `utilware/afero` 文件系统限制和虚拟的库, 封装得挺好....
21. `utilware/goja` JavaScript进入了Golang的身体.......(hso
22. `utilware/gow` 看着`Golang 7 天Web框架` 写的垃圾框架, 类似Gin
23. `utilware/sse` HTTP ServerSendEvents 服务器单方面推送协议 Doc:https://developer.mozilla.org/zh-TW/docs/Web/API/Server-sent_events
24. `grpool` Golang 协程限制和调度
25. `jwt` Golang Json Web Token 封装
