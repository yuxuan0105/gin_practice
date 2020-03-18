package setting

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func GetSetting(confPath string) (*viper.Viper, error) {
	if confPath != "" {
		content, err := ioutil.ReadFile(confPath)
		if err != nil {
			return nil, fmt.Errorf("error at setuping viper: %s", err)
		}
		viper.ReadConfig(bytes.NewBuffer(content))
	} else {
		viper.SetConfigName("app")

		if gin.Mode() == gin.TestMode {
			viper.AddConfigPath("../conf")
		} else {
			viper.AddConfigPath("./conf")
		}

		if err := viper.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("error at setuping viper: %s", err)
		}
	}

	return viper.GetViper(), nil
}
