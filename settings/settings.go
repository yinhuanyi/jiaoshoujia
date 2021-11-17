/**
 * @Author: Robby
 * @File name: settings.go
 * @Create date: 2021-05-18
 * @Function:
 **/

package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/viper"
)

// 这里是用来保存所有程序的配置信息, 这里使用new函数，相当于 var Conf &AppConfig
var Conf = new(AppConfig)

// 由于viper使用的是mapstructure来解析配置文件，因此，这里的结构体tag使用mapstructure
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type MysqlConfig struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DbName       string `mapstructure:"dbname"`
	Port         int    `mapstructure:"port"`
	MaxOpenConns int    `mapstructure:"max_open_connection"`
	MaxIdleConns int    `mapstructure:"max_idle_connection"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
	Db       int    `mapstructure:"dbname"`
	PoolSize int    `mapstructure:"port"`
}

// 第一层结构体
type AppConfig struct {
	Name      string `mapstructure:"name"`
	Mode      string `mapstructure:"mode"`
	Version   string `mapstructure:"version"`
	Port      int    `mapstructure:"port"`
	StartTime string `mapstructure:"start_time"`
	MachineId int64  `mapstructure:"machine_id"`

	// 这里是三个匿名字段，访问这个字段就是类型的名字
	*LogConfig   `mapstructure:"log"`
	*MysqlConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
}

func Init(filePath string) (err error) {
	// 这个filename从命令行传递进来
	viper.SetConfigFile(filePath) // 这里就是直接指定了文件，所以用它比较靠谱
	//viper.SetConfigName("config")
	//viper.SetConfigType("yaml") // 这个只是在远程获取配置文件有用，其他情况下没有用, 如果本地有config.json和config.yaml，那么就会报错
	//viper.AddConfigPath(".")
	if err = viper.ReadInConfig(); err != nil {
		fmt.Printf("配置文件读取失败: %v\n", err)
		return
	}
	// 将viper的配置信息，映射到Conf这个结构体中, 这是一个反序列化的过程
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("配置文件映射到结构体失败")
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Printf("配置文件修改了")
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("配置文件映射到结构体失败")
		}
	})
	return

}
