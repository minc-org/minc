package providers

type Provider interface {
	Name() string
	Info() (*ProviderInfo, error)
	Create() error
	Delete() error
	List() error
}

type ProviderInfo struct {
	Rootless bool
	CGroupV2 bool
}
