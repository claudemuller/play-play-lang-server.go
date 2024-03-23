# Play-Play Language Server

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

## Attaching the LSP to Neovim

Create a file in `<nvim_config_dir>/after/plugins/<something>`.
```lua
local client = vim.lsp.start_client {
  name = 'langserver',
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
