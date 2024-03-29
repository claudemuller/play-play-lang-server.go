package lsp

import (
	"encoding/json"
	"log"
)

type InitialiseRequest struct {
	Request
	Params InitialiseRequestParams `json:"params"`
}

type InitialiseRequestParams struct {
	ClientInfo *ClientInfo `json:"clientInfo"`
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type InitialiseResponse struct {
	Response
	Result InitialiseResult `json:"result"`
}

type InitialiseResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   ServerInfo         `json:"serverInfo"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ServerCapabilities struct {
	TextDocumentSync   int                    `json:"textDocumentSync"`
	HoverProvider      bool                   `json:"hoverProvider"`
	DefinitionProvider bool                   `json:"definitionProvider"`
	CodeActionProvider bool                   `json:"codeActionProvider"`
	CompletionProvider map[string]interface{} `json:"completionProvider"`
}

func HandleInitialise(content []byte, logger *log.Logger) InitialiseResponse {
	var req InitialiseRequest
	if err := json.Unmarshal(content, &req); err != nil {
		logger.Printf("Error parsing initialize request: %s", err)
	}

	logger.Printf("Connected to: %s %s", req.Params.ClientInfo.Name, req.Params.ClientInfo.Version)

	return InitialiseResponse{
		Response: Response{
			RPC: "2.0",
			ID:  &req.ID,
		},
		Result: InitialiseResult{
			Capabilities: ServerCapabilities{
				TextDocumentSync:   1,
				HoverProvider:      true,
				DefinitionProvider: true,
				CodeActionProvider: true,
				CompletionProvider: map[string]interface{}{},
			},
			ServerInfo: ServerInfo{
				Name:    "play-play-langserver",
				Version: "0.0.1",
			},
		},
	}
}
