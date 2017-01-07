package config

type Config struct {
	Mysql        struct {
			     Ip       string `yaml:"ip"`
			     Port     string `yaml:"port"`
			     Username string `yaml:"username"`
			     Passwd   string `yaml:"passwd"`
			     BaseName string `yaml:"baseName"`
		     }

	RPC          struct {
			     Protocol string `yaml:"protocol"`
			     Ip       string `yaml:"ip"`
			     Port     string `yaml:"port"`
		     }

	BlocChanSize int `yaml:"blocChanSize"`
	ThreadSize   int `yaml:"threadSize"`
}

var Cfg Config
