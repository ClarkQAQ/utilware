package circular

import (
	_ "utilware/dep/x/tools/internal/lsp/circular" //@diag("_ \"utilware/dep/x/tools/internal/lsp/circular\"", "compiler", "import cycle not allowed", "error"),diag("\"utilware/dep/x/tools/internal/lsp/circular\"", "compiler", "could not import utilware/dep/x/tools/internal/lsp/circular (no package for import utilware/dep/x/tools/internal/lsp/circular)", "error")
)
