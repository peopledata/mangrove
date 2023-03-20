package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/viper"
)

func Init() (err error) {
	viper.SetConfigName("config") // 指定配置文件名称（不要带后缀）
	viper.SetConfigType("yaml")   // 指定配置文件类型
	viper.AddConfigPath("./conf")      // 指定查找配置文件的路径（相对路径）
	// 读取配置信息
	if err = viper.ReadInConfig(); err != nil {
		return
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了...")
	})
	return
}
