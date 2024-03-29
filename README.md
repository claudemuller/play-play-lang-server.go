# Play-Play Language Server

[![go](https://github.com/claudemuller/play-play-langserver.go/actions/workflows/go.yml/badge.svg)](https://github.com/claudemuller/play-play-langserver.go/actions/workflows/go.yml)

A play-play implementation of the [Language Server Protocol](https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/) in Go.

## Requirements

- [Go](https://go.dev/)

## Install Project Dependencies

```bash
go mod tidy
```

## Build

```bash
go build cmd/main.go
```

## Run  tests

```bash
go test ./...
```

## Attaching the LSP to Neovim

Create a file in `<nvim_config_dir>/after/plugins/<something>`.
```lua
local client = vim.lsp.start_client {
  name = 'play-play-langserver',
  cmd = { '<location_of_binary>' },
  -- on_attach = require("??").on_attach,
}

if not client then
  vim.notify 'Client not found'
  return
end

vim.api.nvim_create_autocmd('FileType', {
  pattern = 'markdown',
  callback = function()
    vim.lsp.buf_attach_client(0, client)
  end,
})
```
