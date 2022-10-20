package mysqldialect

import (
	"reflect"

	"utilware/bun/schema"
)

func scanner(typ reflect.Type) schema.ScannerFunc {
	return schema.Scanner(typ)
}
