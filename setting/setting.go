package setting

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/spf13/viper"
)

func GetSetting() (*viper.Viper, error) {
	confPath := ""
	flag.StringVar(&confPath, "c", "", "Configuration file path.")
	if confPath != "" {
		content, err := ioutil.ReadFile(confPath)
		if err != nil {
			return nil, fmt.Errorf("error at setuping viper: %s", err)
		}
		viper.ReadConfig(bytes.NewBuffer(content))
	} else {
		viper.SetConfigName("app")
		viper.AddConfigPath(".")
		if err := viper.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("error at setuping viper: %s", err)
		}
	}

	return viper.GetViper(), nil
}
