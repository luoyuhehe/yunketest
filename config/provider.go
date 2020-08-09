package config

type Provider interface {
	GetType() string
}

type FileProvider struct {
	File string
}

type EtcdProvider struct {
	EndPoint   string
	Path       string
	ConfigType string
	SecretKey  string
}

type ConsulProvider struct {
	Addr       string
	Key        string
	SecretKey  string
	ConfigType string
}

func (f *FileProvider) GetType() string {
	return "file"
}

func (e *EtcdProvider) GetType() string {
	return "etcd"
}

func (c *ConsulProvider) GetType() string {
	return "consul"
}
