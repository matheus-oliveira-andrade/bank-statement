package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

func InitConfigFile() {
	viper.AddConfigPath("configs")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	environment := viper.GetString("environment")

	viper.SetConfigName(fmt.Sprint("configs", ".", environment))
	viper.SetConfigType("json")
	viper.AddConfigPath("configs")

	err = viper.MergeInConfig()
	if err != nil {
		panic(err)
	}
}
