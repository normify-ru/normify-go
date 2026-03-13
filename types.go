package normify

type ProcessRequest struct {
	Entity string     `json:"entity"`
	Data   []DataItem `json:"data"`
}

type DataItem struct {
	ID    string      `json:"id"`
	Value interface{} `json:"value"`
}

type ProcessResponse struct {
	Success bool        `json:"success"`
	Data    ProcessData `json:"data"`
}

type ProcessData struct {
	Entity string        `json:"entity"`
	Result ProcessResult `json:"result"`
}

type ProcessResult struct {
	Output []ProcessOutputItem `json:"output"`
	Errors []interface{}       `json:"errors"`
}

type ProcessOutputItem struct {
	ID       string                 `json:"id"`
	Value    interface{}            `json:"value"`
	Metadata map[string]interface{} `json:"metadata"`
}
