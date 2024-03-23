package lsp

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

type ServerCapabilities struct{}

func NewInitialiseResponse(id int) InitialiseResponse {
	return InitialiseResponse{
		Response: Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: InitialiseResult{
			Capabilities: ServerCapabilities{},
			ServerInfo: ServerInfo{
				Name:    "langserver",
				Version: "0.0.1",
			},
		},
	}
}
