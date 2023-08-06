package config

type Server struct {
	Model     string `mapstructure:"model" json:"model" yaml:"model"`
	Health    bool   `mapstructure:"health" json:"health" yaml:"health"`
	Profiling bool   `mapstructure:"profiling" json:"profiling" yaml:"profiling"`
	Host      string `mapstructure:"host" json:"host" yaml:"host"`
	UseSSL    bool   `mapstructure:"use-ssl" json:"use-ssl" yaml:"use-ssl"`
}

func (s *Server) GetHostUrl() string {
	if s.UseSSL {
		return "https://" + s.Host
	} else {
		return "http://" + s.Host
	}
}

type Grpc struct {
	BindInfo
}

type Http struct {
	BindInfo
}

type SSL struct {
	BindInfo
	CertFile       string `mapstructure:"cert-file" json:"cert-file" yaml:"cert-file"`
	PrivateKeyFile string `mapstructure:"private-key-file" json:"private-key-file" yaml:"private-key-file"`
}

type BindInfo struct {
	BindAddress string `mapstructure:"bind-address" json:"bind-address" yaml:"bind-address"`
	BindPort    int    `mapstructure:"bind-port" json:"bind-port" yaml:"bind-port"`
}
