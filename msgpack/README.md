# MessagePack encoding for Golang

[![Build Status](https://travis-ci.org/vmihailenco/msgpack.svg)](https://travis-ci.org/vmihailenco/msgpack)
[![PkgGoDev](https://pkg.go.dev/badge/utilware/msgpack)](https://pkg.go.dev/utilware/msgpack)
[![Documentation](https://img.shields.io/badge/msgpack-documentation-informational)](https://msgpack.uptrace.dev/)
[![Chat](https://discordapp.com/api/guilds/752070105847955518/widget.png)](https://discord.gg/rWtp5Aj)

> msgpack is brought to you by :star: [**uptrace/uptrace**](https://github.com/uptrace/uptrace).
> Uptrace is an open source and blazingly fast
> [distributed tracing tool](https://get.uptrace.dev/compare/distributed-tracing-tools.html) powered
> by OpenTelemetry and ClickHouse. Give it a star as well!

## Resources

- [Documentation](https://msgpack.uptrace.dev)
- [Chat](https://discord.gg/rWtp5Aj)
- [Reference](https://pkg.go.dev/utilware/msgpack)
- [Examples](https://pkg.go.dev/utilware/msgpack#pkg-examples)

## Features

- Primitives, arrays, maps, structs, time.Time and interface{}.
- Appengine \*datastore.Key and datastore.Cursor.
- [CustomEncoder]/[CustomDecoder] interfaces for custom encoding.
- [Extensions](https://pkg.go.dev/utilware/msgpack#example-RegisterExt) to encode
  type information.
- Renaming fields via `msgpack:"my_field_name"` and alias via `msgpack:"alias:another_name"`.
- Omitting individual empty fields via `msgpack:",omitempty"` tag or all
  [empty fields in a struct](https://pkg.go.dev/utilware/msgpack#example-Marshal-OmitEmpty).
- [Map keys sorting](https://pkg.go.dev/utilware/msgpack#Encoder.SetSortMapKeys).
- Encoding/decoding all
  [structs as arrays](https://pkg.go.dev/utilware/msgpack#Encoder.UseArrayEncodedStructs)
  or
  [individual structs](https://pkg.go.dev/utilware/msgpack#example-Marshal-AsArray).
- [Encoder.SetCustomStructTag] with [Decoder.SetCustomStructTag] can turn msgpack into drop-in
  replacement for any tag.
- Simple but very fast and efficient
  [queries](https://pkg.go.dev/utilware/msgpack#example-Decoder.Query).

[customencoder]: https://pkg.go.dev/utilware/msgpack#CustomEncoder
[customdecoder]: https://pkg.go.dev/utilware/msgpack#CustomDecoder
[encoder.setcustomstructtag]:
  https://pkg.go.dev/utilware/msgpack#Encoder.SetCustomStructTag
[decoder.setcustomstructtag]:
  https://pkg.go.dev/utilware/msgpack#Decoder.SetCustomStructTag

## Installation

msgpack supports 2 last Go versions and requires support for
[Go modules](https://github.com/golang/go/wiki/Modules). So make sure to initialize a Go module:

```shell
go mod init github.com/my/repo
```

And then install msgpack/v5 (note _v5_ in the import; omitting it is a popular mistake):

```shell
go get utilware/msgpack
```

## Quickstart

```go
import "utilware/msgpack"

func ExampleMarshal() {
    type Item struct {
        Foo string
    }

    b, err := msgpack.Marshal(&Item{Foo: "bar"})
    if err != nil {
        panic(err)
    }

    var item Item
    err = msgpack.Unmarshal(b, &item)
    if err != nil {
        panic(err)
    }
    fmt.Println(item.Foo)
    // Output: bar
}
```

## See also

- [Golang ORM](https://utilware/bun) for PostgreSQL, MySQL, MSSQL, and SQLite
- [Golang PostgreSQL](https://bun.uptrace.dev/postgres/)
- [Golang HTTP router](https://utilware/bunrouter)
- [Golang ClickHouse ORM](https://github.com/uptrace/go-clickhouse)

## Contributors

Thanks to all the people who already contributed!

<a href="https://github.com/vmihailenco/msgpack/graphs/contributors">
  <img src="https://contributors-img.web.app/image?repo=vmihailenco/msgpack" />
</a>
