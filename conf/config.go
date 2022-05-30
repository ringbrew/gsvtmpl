package conf

import (
	"log"
	"net"
	"strconv"
	"github.com/ringbrew/gsvcore/config"
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
	//todo

	return result, nil
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func shouldInt(val string) int {
	result, _ := strconv.Atoi(val)
	return result
}
