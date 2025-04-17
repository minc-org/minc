package types

type CreateType struct {
	Provider      string
	UShiftVersion string
	UShiftConfig  string
}

type StatusType struct {
	Container string `json:"container"`
	APIServer string `json:"apiserver"`
	Error     string `json:"error,omitempty"`
}
