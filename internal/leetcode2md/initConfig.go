package leetcode2md

import (
	"log/slog"

	"github.com/spf13/viper"
)

const (
	// defaultConfigName 指定了 leetcode2md 服务的默认配置文件名.
	defaultConfigName = "leetcode2md"
)

func initConfig() {
	if configFile != "" {
		// 从命令行选项指定的配置文件中读取
		viper.SetConfigFile(configFile)
	} else {
		// 把当前目录加入到配置文件的搜索路径中
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config")

		// 配置文件名称（没有文件扩展名）
		viper.SetConfigName(defaultConfigName)
	}

	// 读取配置文件。如果指定了配置文件名，则使用指定的配置文件，否则在注册的搜索路径中搜索
	if err := viper.ReadInConfig(); err != nil {
		slog.Warn("使用默认配置, " + err.Error())
		viper.SetDefault("dir.md", "./")
		viper.SetDefault("dir.img", "./assets/")
	}

	slog.Debug("Using config file", "filename", viper.ConfigFileUsed())
}
