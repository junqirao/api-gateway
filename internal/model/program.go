package model

// DeleteProgramInfoInput delete program
type DeleteProgramInfoInput struct {
	ServiceName string `json:"service_name"`
	Name        string `json:"name"`
}

// SetGlobalVariablesInput set global
type SetGlobalVariablesInput struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
