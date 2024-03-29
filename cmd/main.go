package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/claudemuller/play-play-lang-server/analysis"
	"github.com/claudemuller/play-play-lang-server/lsp"
	"github.com/claudemuller/play-play-lang-server/rpc"
)

func main() {
	logger := getLogger("/tmp/play-play-langserver.log")
	logger.Println("Started...")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	state := analysis.NewState()

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, content, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("Error: %s", err)
			continue
		}
		handleMsg(os.Stdout, state, method, content, logger)
	}
}

func handleMsg(writer io.Writer, state analysis.State, method string, content []byte, logger *log.Logger) {
	logger.Printf("Received msg with method: %s", method)

	switch method {
	case "initialize":
		var req lsp.InitialiseRequest
		if err := json.Unmarshal(content, &req); err != nil {
			logger.Printf("Error parsing initialize request: %s", err)
		}

		logger.Printf("Connected to: %s %s", req.Params.ClientInfo.Name, req.Params.ClientInfo.Version)

		msg := lsp.NewInitialiseResponse(req.ID)
		writeResponse(writer, msg)

	case "textDocument/didOpen":
		var req lsp.DidOpenTextDocumentNotification
		if err := json.Unmarshal(content, &req); err != nil {
			logger.Printf("Error parsing didOpen request: %s", err)
		}

		logger.Printf("Opended: %s", req.Params.TextDocument.URI)

		state.OpenDocument(req.Params.TextDocument.URI, req.Params.TextDocument.Text)

	case "textDocument/didChange":
		var req lsp.DidChangeTextDocumentNotification
		if err := json.Unmarshal(content, &req); err != nil {
			logger.Printf("Error parsing didChange request: %s", err)
		}

		logger.Printf("Changed: %s", req.Params.TextDocument.URI)

		for _, change := range req.Params.ContentChanges {
			state.UpdateDocument(req.Params.TextDocument.URI, change.Text)
		}

	case "textDocument/hover":
		var req lsp.HoverRequest
		if err := json.Unmarshal(content, &req); err != nil {
			logger.Printf("Error parsing hover request: %s", err)
		}

		res := state.Hover(req.ID, req.Params.TextDocument.URI, req.Params.Position)
		writeResponse(writer, res)

	case "textDocument/definition":
		var req lsp.DefinitionRequest
		if err := json.Unmarshal(content, &req); err != nil {
			logger.Printf("Error parsing definition request: %s", err)
		}

		res := state.Definition(req.ID, req.Params.TextDocument.URI, req.Params.Position)
		writeResponse(writer, res)

	case "textDocument/codeAction":
		var req lsp.CodeActionRequest
		if err := json.Unmarshal(content, &req); err != nil {
			logger.Printf("Error parsing code action request: %s", err)
		}

		res := state.TextDocumentCodeAction(req.ID, req.Params.TextDocument.URI)
		writeResponse(writer, res)

	case "textDocument/completion":
		var req lsp.CompletionRequest
		if err := json.Unmarshal(content, &req); err != nil {
			logger.Printf("Error parsing completion request: %s", err)
		}

		res := state.TextDocumentCompletion(req.ID, req.Params.TextDocument.URI)
		writeResponse(writer, res)
	}
}

func writeResponse(writer io.Writer, msg interface{}) {
	reply := rpc.EncodeMessage(msg)
	_, _ = writer.Write([]byte(reply))
}

func getLogger(filename string) *log.Logger {
	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("File no good :(")
	}

	return log.New(logFile, "[play-play-langserver] ", log.Ldate|log.Ltime|log.Lshortfile)
}
