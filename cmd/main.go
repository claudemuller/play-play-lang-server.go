package main

import (
	"bufio"
	"encoding/json"
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
		handleMsg(state, method, content, logger)
	}
}

func handleMsg(state analysis.State, method string, content []byte, logger *log.Logger) {
	logger.Printf("Received msg with method: %s", method)

	switch method {
	case "initialize":
		var req lsp.InitialiseRequest
		if err := json.Unmarshal(content, &req); err != nil {
			logger.Printf("Error parsing initialize request: %s", err)
		}

		logger.Printf("Connected to: %s %s", req.Params.ClientInfo.Name, req.Params.ClientInfo.Version)

		msg := lsp.NewInitialiseResponse(req.ID)
		reply := rpc.EncodeMessage(msg)
		writer := os.Stdout
		_, _ = writer.Write([]byte(reply))

		logger.Print("Sent reply")

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
	}
}

func getLogger(filename string) *log.Logger {
	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("File no good :(")
	}

	return log.New(logFile, "[play-play-langserver] ", log.Ldate|log.Ltime|log.Lshortfile)
}
