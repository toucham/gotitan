package config

// Proxy configurations that is parsed from files
type StaticProxyConfig struct {
}

// Proxy configurations that is retrieved from user input and can be changed during run time
type DynamicProxyConfig struct {
}

type ProxyConfig struct {
	configPath string
	StaticProxyConfig
	DynamicProxyConfig
}

func NewProxyConfig(configPath string, cliInput string) *ProxyConfig {
	p := &ProxyConfig{
		configPath: configPath,
	}
	p.getStaticConfig()
	p.applyDynamicConfig(cliInput)

	return p
}

func (p *ProxyConfig) getStaticConfig() {
	panic("not implemented yet")
}

func (p *ProxyConfig) applyDynamicConfig(input string) DynamicProxyConfig {
	panic("not implemented yet")
}
