package cfg

import (
	"time"

	"github.com/spf13/viper"
)

type Cfg struct {
	Mysql   CfgMysql
	Fserver CfgServer
}

type CfgMysql struct {
	Driver   string        `json:"driver" validate:"required"`
	Host     string        `json:"host" validate:"required"`
	Port     string        `json:"port" validate:"required"`
	User     string        `json:"user" validate:"required"`
	Password string        `json:"password"`
	Dbname   string        `json:"dbname" validate:"required"`
	Maxopen  int           `json:"maxopen"`
	Maxlife  time.Duration `json:"maxlife"`
	Maxidle  int           `json:"maxidle"`
	Maxtime  time.Duration `json:"maxtime"`
}

type CfgServer struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

func LoadConfig(cfg *Cfg) error {
	viper.AddConfigPath("cfg")
	viper.AddConfigPath("../cfg")
	viper.SetConfigName("cfg")
	viper.SetConfigType("yml")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.Unmarshal(cfg)
	if err != nil {
		return err
	}

	return nil
}
