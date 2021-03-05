package other

import (
	references "utilware/dep/x/tools/internal/lsp/references"
)

func _() {
	references.Q = "hello" //@mark(assignExpQ, "Q")
	bob := func(_ string) {}
	bob(references.Q) //@mark(bobExpQ, "Q")
}
