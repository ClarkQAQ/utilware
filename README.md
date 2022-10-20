## Golang Utilware 离线开发工具包

#### 名称灵感来着 Aimware, Aimware 永远的空枪废物.....

---

> 坚持纯 Go 拒绝 CGO!

##### Utilware 无任何外部依赖, 真正的 ALL IN ONE, 用于离线开发以及闭关开发, 基本的工具都包含了.....

##### 收集开源项目并去除 CGO 以及自己实现, 减少依赖后打包....依赖项目统一在 [dep]

---

#### Go Mod ( 使用方式) :

引入:

```bash

# 1.通常方式
# 这种方式后续更新只需要 git pull 就行了

git clone https://github.com/ClarkQAQ/utilware
go mod edit replace utilware => /xxx/utilware

# 2.使用伪版本

replace utilware => github.com/ClarkQAQ/utilware v0.0.0-[时间]-[commit]

```

项目使用:

```go
package main

import "utilware/tig"

func main() {
    logger.Info("hello world!")
	tig.New().Run(":8080")
}

```

#### 索引:

| 名称 | 描述 | 版本 | 作者 | 源链接 |
| ---- | ---- | ---- | ---- | ------ |
| afero | A FileSystem Abstraction System for Go. Golang 文件系统抽象系统, 文件系统限制和虚拟的库, 封装得挺好.... | v0.0.0 | spf13 | https://github.com/spf13/afero |
| arpc | More Effective Network Communication. 更高效的网络通信通知和广播, 多协议自动重连支持加密的双向rpc库 | v0.0.0 | lesismal | https://github.com/lesismal/arpc |
| bbolt | An embedded key/value database for Go. 嵌入式纯go高性能kv数据库 | v0.0.0 | benbjohnson,etcd-io | https://github.com/etcd-io/bbolt |
| bun | SQL-first Golang ORM. SQL 优先的 Golang ORM, 对 Postgres 有很好的支持. | v0.0.0 | uptrace,vmihailenco | https://github.com/uptrace/bun |
| cron | a cron library for go. Golang Cron 定时任务 | v0.0.0 | robfig | https://github.com/robfig/cron |
| csvutil | csvutil provides fast and idiomatic mapping between CSV and Go (golang) values. CSV 和 Go (golang) 之间数据映射, 一直都在用的CSV处理库, 好用! | v0.0.0 | jszwec | https://github.com/jszwec/csvutil |
| gimg | Golang 轻量级图片处理库 | v0.0.0 | ClarkQAQ | https://github.com/ClarkQAQ/utilware |
| gorm | The fantastic ORM library for Golang, aims to be developer friendly. Golang 出色的对开发人员友好的 ORM, 说它是 Golang 最好的 ORM 不为过吧 | v0.0.0 | jinzhu | https://github.com/go-gorm/gorm |
| gtime | GoFrame 拆分出来的时间处理, 具有标准占位符支持. | v0.0.0 | gqcn | https://github.com/gogf/gf/tree/master/os/gtime |
| ini | ini provides INI file read and write functionality in Go. Golang ini 配置文件读写以及映射 | v0.0.0 | unknwon | https://github.com/go-ini/ini |
| logger | 自定义格式以及writer并且支持同步异步的日志库, 附带进度条以及耗时计算, 颜色仅支持 Posix 终端 | v0.0.0 | ClarkQAQ | https://github.com/ClarkQAQ/utilware |
| lua | VM and compiler for Lua in Go. 基于 gopher-lua 魔改后的 lua 脚本引擎/虚拟机, 支持多线程, 但不保证 vm 内线程安全, 有插件可以方便的传入 golang userdata | v0.0.0 | yuin | https://github.com/yuin/gopher-lua |
| msgpack | msgpack.org[Go] MessagePack encoding for Golang. Golang 的 MessagePack 编码实现 | v0.0.0 | vmihailenco | https://github.com/vmihailenco/msgpack |
| request | A concise HTTP request client for Go. 一个超轻的 http 客户端, 追加了 multifrom 支持, 并且支持 socks5 或者 http 代理 URL | v0.0.0 | DavidCai1111 | https://github.com/DavidCai1111/request |
| sqlite | sqlite is a CGo-free port of SQLite/SQLite3. 使用 plan9 asm 实现的 sqlite3 数据库 | v0.0.0 | cznic | https://gitlab.com/cznic/sqlite |
| sqlx | 对 sqlx 拙劣的模仿, 特性是日志功能以及事务计数器 | v0.0.0 | ClarkQAQ | https://github.com/ClarkQAQ/utilware |
| tig | gin 风格的 web 框架, 使用前缀树路由附带 swagger ui, logger, timer 以及 content type 自动纠正 | v0.0.0 | ClarkQAQ | https://github.com/ClarkQAQ/utilware |
| acm | 字符匹配自动机, 用于敏感词识别以及返回错误判断 | v0.0.0 | ClarkQAQ | https://github.com/ClarkQAQ/utilware |
| codec | 简单的编码解码封装 | v0.0.0 | ClarkQAQ | https://github.com/ClarkQAQ/utilware |
| conv | 常用类型转换 | v0.0.0 | ClarkQAQ | https://github.com/ClarkQAQ/utilware |
| crypc | 常用加密/取模封装 | v0.0.0 | ClarkQAQ | https://github.com/ClarkQAQ/utilware |
| mlock | 业务用键锁, 千万锁7秒加载完毕, 2秒全部上锁, 单个占用230Byte, 具有 context 的加锁可防止任务阻塞 | v0.0.0 | ClarkQAQ | https://github.com/ClarkQAQ/utilware |
| module | 简单的依赖注入/启动任务管理 | v0.0.0 | ClarkQAQ | https://github.com/ClarkQAQ/utilware |
| sanitize | sanitize provides functions for sanitizing text in golang strings. 去除 text 里面不安全的 html 内容 | v0.0.0 | kennygrant | https://github.com/kennygrant/sanitize |
| smap | 多功能的泛型 map 实现, 线程安全, 支持多种类型的 key 和 value | v0.0.0 | ClarkQAQ | https://github.com/ClarkQAQ/utilware |
| ug | web 开发常用返回结构 | v0.0.0 | ClarkQAQ | https://github.com/ClarkQAQ/utilware |

