# Recommended IDEs for Go

Below are popular editors and IDEs with strong Go support.

## [GoLand](https://www.jetbrains.com/go/)
JetBrains' dedicated Go IDE. Provides deep code intelligence, integrated debugger, test runner, and refactoring tools out of the box. Best choice if you want a fully-featured IDE with minimal setup.

## [Visual Studio Code](https://code.visualstudio.com/)
Free, lightweight editor from Microsoft. Install the official **[Go extension](https://marketplace.visualstudio.com/items?itemName=golang.Go)** (powered by `gopls`) to get autocompletion, go-to-definition, linting, debugging, and test support.

## [Neovim](https://neovim.io/)
Highly extensible terminal-based editor. Configure with [nvim-lspconfig](https://github.com/neovim/nvim-lspconfig) + `gopls` for a full Go language server experience. Popular plugin managers: [lazy.nvim](https://github.com/folke/lazy.nvim) or [packer.nvim](https://github.com/wbthomason/packer.nvim).

## [Vim](https://www.vim.org/)
Classic terminal editor. The [vim-go](https://github.com/fatih/vim-go) plugin adds syntax highlighting, formatting, imports, linting, and debugging support.

## [Emacs](https://www.gnu.org/software/emacs/)
Powerful, extensible editor. Use [go-mode](https://github.com/dominikh/go-mode.el) along with [lsp-mode](https://emacs-lsp.github.io/lsp-mode/) and `gopls` for a complete Go development experience.

## [Zed](https://zed.dev/)
Fast, modern editor with built-in Go support via its extensions marketplace and native LSP integration.

---

## Shared Tooling

Regardless of editor, these tools integrate with all of the above:

| Tool | Purpose | Install |
|------|---------|---------|
| [`gopls`](https://pkg.go.dev/golang.org/x/tools/gopls) | Official Go language server | `go install golang.org/x/tools/gopls@latest` |
| [`golangci-lint`](https://golangci-lint.run/) | Aggregated linter | See [install docs](https://golangci-lint.run/usage/install/) |
| [`dlv` (Delve)](https://github.com/go-delve/delve) | Go debugger | `go install github.com/go-delve/delve/cmd/dlv@latest` |
| [`goimports`](https://pkg.go.dev/golang.org/x/tools/cmd/goimports) | Format + manage imports | `go install golang.org/x/tools/cmd/goimports@latest` |
