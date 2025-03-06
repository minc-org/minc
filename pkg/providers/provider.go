package providers

type Provider interface {
	Name() string
	Info() (*ProviderInfo, error)
	Create() error
	WaitForAPI() error
	GetKubeConfig() ([]byte, error)
	Delete() error
	List() error
}

type ProviderInfo struct {
	Rootless bool
	CGroupV2 bool
}
