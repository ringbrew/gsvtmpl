package conf

import (
	"github.com/ringbrew/gsv/config"
)

type Config struct {
	Env   string `yaml:"environment"`
	Debug bool   `yaml:"debug"`
	Port  Port   `yaml:"port"`
	Mysql Mysql  `yaml:"mysql"`
	Redis Redis  `yaml:"redis"`
	Trace Trace  `yaml:"trace"`
}

type Port struct {
	Api   int `yaml:"api"`
	Admin int `yaml:"admin"`
	Rpc   int `yaml:"rpc"`
	Proxy int `yaml:"proxy"`
}

type Mysql struct {
	UserName string `yaml:"user_name"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Database string `yaml:"database"`
}

type Redis struct {
	Host     string `yaml:"host"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type Trace struct {
	Type     string `yaml:"type"`
	Endpoint string `yaml:"endpoint"`
}

func Load(path string) (Config, error) {
	var result Config
	loader := config.NewLoader(config.LoaderTypeYml, path)
	if err := loader.Load(&result); err != nil {
		return Config{}, err
	}

	return result, nil
}
