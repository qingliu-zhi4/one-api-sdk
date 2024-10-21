package OneApiSdk

const contentTypeJson = "application/json"

type CommonAPIRes struct {
	sessionID string
	res       *oneAPIRes
}

type oneAPIRes struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
