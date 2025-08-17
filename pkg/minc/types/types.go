package types

type CreateType struct {
	Provider      string
	UShiftVersion string
	UShiftImage   string
	UShiftConfig  string
	HTTPSPort     int
	HTTPPort      int
}

type StatusType struct {
	Container string `json:"container"`
	APIServer string `json:"apiserver"`
	Error     string `json:"error,omitempty"`
}
