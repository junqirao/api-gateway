package model

type ProgramInfo struct {
	Expr        string `json:"expr"` // base64 encoded
	Name        string `json:"name"`
	ServiceName string `json:"service_name"`
}
