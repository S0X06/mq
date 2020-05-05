package conf

import (
	"flag"
	"log"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	confPath string
	Conf     = &Config{}
)

type Config struct {
	RunMode  string
	Name     string
	Addr     string
	RabbitMq *RabbitMq
	Grpc     *Grpc
	MongoDB  *MongoDB
	Cron     *Cron
}

type RabbitMq struct {
	Addr     string
	Port     string
	UserName string
	PassWord string
	ConnNum  string
}

type MongoDB struct {
	Auth      string
	Addr      string
	Port      string
	UserName  string
	PassWord  string
	DdataBase string
}

type Grpc struct {
	Port string
}

type Cron struct {
	NotifySpec string
	SendSpec   string
	LockSpec   string
}

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
}

func Init() (err error) {

	if confPath != "" {
		viper.SetConfigFile(confPath) // 如果指定了配置文件，则解析指定的配置文件
	} else {
		viper.AddConfigPath("./conf") // 如果没有指定配置文件，则解析默认的配置文件
		viper.SetConfigName("conf")
		confPath = "./conf/conf"
	}
	viper.SetConfigType("yaml") // 设置配置文件格式为YAML
	viper.AutomaticEnv()        // 读取匹配的环境变量
	// viper.SetEnvPrefix("SERVER") // 读取环境变量的前缀为SERVER
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	if err := viper.ReadInConfig(); err != nil { // viper解析配置文件
		return err
	}

	//基础配置
	Conf.Name = viper.GetString("app.name")
	Conf.RunMode = viper.GetString("app.runmode")
	Conf.Addr = viper.GetString("app.addr")

	//Mq
	Conf.RabbitMq = &RabbitMq{
		Addr:     viper.GetString("rabbitmq.addr"),
		Port:     viper.GetString("rabbitmq.port"),
		UserName: viper.GetString("rabbitmq.username"),
		PassWord: viper.GetString("rabbitmq.password"),
		ConnNum:  viper.GetString("rabbitmq.conn_num"),
	}

	//GRPC
	Conf.Grpc = &Grpc{
		Port: viper.GetString("grpc.port"),
	}

	Conf.MongoDB = &MongoDB{
		Addr:      viper.GetString("mongodb.addr"),
		Port:      viper.GetString("mongodb.port"),
		UserName:  viper.GetString("mongodb.username"),
		PassWord:  viper.GetString("mongodb.password"),
		Auth:      viper.GetString("mongodb.auth"),
		DdataBase: viper.GetString("mongodb.database"),
	}

	Conf.Cron = &Cron{
		NotifySpec: viper.GetString("cron.notifySpec"),
		SendSpec:   viper.GetString("cron.sendSpec"),
		LockSpec:   viper.GetString("cron.lockSpec"),
	}

	return

}

// 监控配置文件变化并热加载程序
func (c *Config) watch() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Config file changed: %s", e.Name)
		// log.Infof("Config file changed: %s", e.Name)
	})
}
