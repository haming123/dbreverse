package app

import (
	"github.com/haming123/wego/wini"
)

type DbParam struct {
	DbHost string `ini:"db_host"`
	DbPort string `ini:"db_port"`
	DbName string `ini:"db_name"`
	DbUser string `ini:"db_user"`
	DbPwd  string `ini:"db_pwd"`
}

type AppConfig struct {
	DbDriver   string  `ini:"db_driver"`
	PkgName    string  `ini:"pkg_name"`
	CreateTime string  `ini:"create_time"`
	UseTag     bool    `ini:"use_field_tag"`
	UsePool    bool    `ini:"use_pool"`
	DbCfg      DbParam `ini:"db"`
}

var AppCfg AppConfig

func ReadAppConfig(conf_file ...string) (*AppConfig, error) {
	data, err := wini.InitConfigData(conf_file...)
	if err != nil {
		return &AppCfg, err
	}

	err = data.GetStruct(&AppCfg)
	if err != nil {
		return &AppCfg, err
	}

	if AppCfg.PkgName == "" {
		AppCfg.PkgName = "model"
	}

	AppCfg.UseTag = true
	return &AppCfg, nil
}
