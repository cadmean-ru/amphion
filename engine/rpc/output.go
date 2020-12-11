package rpc

type FunctionOutput struct {
	Error  int                    `json:"error"`
	Result interface{}            `json:"result"`
	Meta   map[string]interface{} `json:"meta"`
}
