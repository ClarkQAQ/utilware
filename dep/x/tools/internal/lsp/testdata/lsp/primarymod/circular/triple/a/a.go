package a

import (
	_ "utilware/dep/x/tools/internal/lsp/circular/triple/b" //@diag("_ \"utilware/dep/x/tools/internal/lsp/circular/triple/b\"", "compiler", "import cycle not allowed", "error")
)
