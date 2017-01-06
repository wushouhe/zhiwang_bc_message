package config

type Config struct {
	Mysql        map[string]string `yaml:"mysql"`
	RPC          map[string]string `yaml:"rpc"`
	BlocChanSize int `yaml:"blocChanSize"`
	ThreadSize   int `yaml:"threadSize"`
}

func init() {

}
