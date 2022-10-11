## Golang Utilware 离线开发工具包

#### 名称灵感来着Aimware, Aimware永远的空枪废物.....

---


> 坚持纯Go拒绝CGO!


##### Utilware 无任何外部依赖, 真正的ALL IN ONE, 用于离线开发以及闭关开发, 基本的工具都包含了.....


##### 收集开源项目并去除CGO以及自己实现, 减少依赖后打包....依赖项目统一在 [dep]


---


#### Go Mod ( 使用方式) :



```bash
git clone https://github.com/ClarkQAQ/utilware

go mod edit replace utilware => /xxx/utilware
```



```go
package main

import "utilware/tig"

func main() {
    logger.Info("hello world!")
	tig.New().Run(":8080")
}

```



#### 索引:




- `afero` [v0.0.0] `A FileSystem Abstraction System for Go. Golang 文件系统抽象系统, 文件系统限制和虚拟的库, 封装得挺好....` [源链接](https://github.com/spf13/afero)



- `arpc` [v0.0.0] `More Effective Network Communication. 更高效的网络通信通知和广播, 多协议自动重连支持加密的双向rpc库` [源链接](https://github.com/lesismal/arpc)



- `bbolt` [v0.0.0] `An embedded key/value database for Go. 嵌入式纯go高性能kv数据库` [源链接](https://github.com/etcd-io/bbolt)



- `cron` [v0.0.0] `a cron library for go. Golang Cron 定时任务` [源链接](https://github.com/robfig/cron)



- `csvutil` [v0.0.0] `csvutil provides fast and idiomatic mapping between CSV and Go (golang) values. CSV 和 Go (golang) 之间数据映射, 一直都在用的CSV处理库, 好用!` [源链接](https://github.com/jszwec/csvutil)



- `gimg` [v0.0.0] `Golang 轻量级图片处理库` [源链接](https://github.com/ClarkQAQ/utilware)



- `gorm` [v0.0.0] `The fantastic ORM library for Golang, aims to be developer friendly. Golang 出色的对开发人员友好的 ORM, 说它是 Golang 最好的 ORM 不为过吧` [源链接](https://github.com/go-gorm/gorm)



- `gtime` [v0.0.0] `GoFrame 拆分出来的时间处理, 具有标准占位符支持.` [源链接](https://github.com/gogf/gf/tree/master/os/gtime)



- `ini` [v0.0.0] `ini provides INI file read and write functionality in Go. Golang ini 配置文件读写以及映射` [源链接](https://github.com/go-ini/ini)



- `logger` [v0.0.0] `自定义格式以及writer并且支持同步异步的日志库, 附带进度条以及耗时计算, 颜色仅支持 Posix 终端` [源链接](https://github.com/ClarkQAQ/utilware)



- `lua` [v0.0.0] `VM and compiler for Lua in Go. 基于 gopher-lua 魔改后的 lua 脚本引擎/虚拟机, 支持多线程, 但不保证 vm 内线程安全, 有插件可以方便的传入 golang userdata` [源链接](https://github.com/yuin/gopher-lua)



- `request` [v0.0.0] `A concise HTTP request client for Go. 一个超轻的 http 客户端, 追加了 multifrom 支持, 并且支持 socks5 或者 http 代理 URL` [源链接](https://github.com/DavidCai1111/request)



- `sqlite` [v0.0.0] `sqlite is a CGo-free port of SQLite/SQLite3. 使用 plan9 asm 实现的 sqlite3 数据库` [源链接](https://gitlab.com/cznic/sqlite)



- `sqlx` [v0.0.0] `对 sqlx 拙劣的模仿, 特性是日志功能以及事务计数器` [源链接](https://github.com/ClarkQAQ/utilware)



- `tig` [v0.0.0] `gin 风格的 web 框架, 使用前缀树路由附带 swagger ui, logger, timer 以及 content type 自动纠正` [源链接](https://github.com/ClarkQAQ/utilware)



- `acm` [v0.0.0] `字符匹配自动机, 用于敏感词识别以及返回错误判断` [源链接](https://github.com/ClarkQAQ/utilware)



- `codec` [v0.0.0] `简单的编码解码封装` [源链接](https://github.com/ClarkQAQ/utilware)



- `conv` [v0.0.0] `常用类型转换` [源链接](https://github.com/ClarkQAQ/utilware)



- `crypc` [v0.0.0] `常用加密/取模封装` [源链接](https://github.com/ClarkQAQ/utilware)



- `mlock` [v0.0.0] `业务用键锁, 千万锁7秒加载完毕, 2秒全部上锁, 单个占用230Byte, 具有 context 的加锁可防止任务阻塞` [源链接](https://github.com/ClarkQAQ/utilware)



- `module` [v0.0.0] `简单的依赖注入/启动任务管理` [源链接](https://github.com/ClarkQAQ/utilware)



- `sanitize` [v0.0.0] `sanitize provides functions for sanitizing text in golang strings. 去除 text 里面不安全的 html 内容` [源链接](https://github.com/kennygrant/sanitize)



- `smap` [v0.0.0] `多功能的泛型 map 实现, 线程安全, 支持多种类型的 key 和 value` [源链接](https://github.com/ClarkQAQ/utilware)



- `ug` [v0.0.0] `web 开发常用返回结构` [源链接](https://github.com/ClarkQAQ/utilware)



- `weightedrand` [v0.0.0] `` [源链接]()

