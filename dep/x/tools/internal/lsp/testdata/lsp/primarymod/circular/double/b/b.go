package b

import (
	_ "utilware/dep/x/tools/internal/lsp/circular/double/one" //@diag("_ \"utilware/dep/x/tools/internal/lsp/circular/double/one\"", "compiler", "import cycle not allowed", "error"),diag("\"utilware/dep/x/tools/internal/lsp/circular/double/one\"", "compiler", "could not import utilware/dep/x/tools/internal/lsp/circular/double/one (no package for import utilware/dep/x/tools/internal/lsp/circular/double/one)", "error")
)
