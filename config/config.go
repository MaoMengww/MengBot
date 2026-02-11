package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Master MasterConfig `mapstructure:"master"`
	Bot BotConfig `mapstructure:"bot"`
	Prompt PromptConfig `mapstructure:"prompt"`
}

type MasterConfig struct {
	MasterID string `mapstructure:"master_id"`
	MasterName string `mapstructure:"master_name"`
}

type MemoryConfig struct {
	WindowLength int `mapstructure:"window_length"`
}

type BotConfig struct {
	Name string `mapstructure:"name"`
	Memory MemoryConfig `mapstructure:"memory"`
}

type PromptConfig struct {
	ChatPrompt string `mapstructure:"chat_prompt"`
	RouterPrompt string `mapstructure:"router_prompt"`
}



var Conf Config

func Init() {
	// 获取当前工作目录
	workDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting working directory: %s", err)
	}

	// 设置配置文件名
	viper.SetConfigName("config")
	// 设置配置文件类型
	viper.SetConfigType("yaml")
	// 设置查找路径
	viper.AddConfigPath(filepath.Join(workDir, "config"))
	viper.AddConfigPath(workDir)
	viper.AddConfigPath(filepath.Join(workDir, "../../config"))

	// 读取配置
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	// 映射到结构体
	if err := viper.Unmarshal(&Conf); err != nil {
		log.Fatalf("Unable to decode into struct: %s", err)
	}

	log.Println("Config loaded successfully!")
}
