package utils

import (
	"errors"
	"strings"

	"github.com/Unknwon/goconfig"
)

var config_dictionary = make(map[string]map[string]map[string]string)

func LoadDbConfig(section string) (string, error) {
	if config_keys, exist := config_dictionary["dbs"]; exist {
		if config_key, ok := config_keys[section]; ok {
			return config_key["connection"], nil
		} else {
			config_keys[section] = make(map[string]string)
		}
	} else {
		config_dictionary["dbs"] = make(map[string]map[string]string)
		config_dictionary["dbs"][section] = make(map[string]string)
	}

	DbConfPath := "conf/db.ini"
	Cfg, err := goconfig.LoadConfigFile(DbConfPath)
	if err != nil {
	}
	username, err := Cfg.GetValue("mysql."+section, "username")
	if err != nil {
	}
	password, err := Cfg.GetValue("mysql."+section, "password")
	if err != nil {
	}
	hostname, err := Cfg.GetValue("mysql."+section, "hostname")
	if err != nil {
	}
	port, err := Cfg.GetValue("mysql."+section, "port")
	if err != nil {
	}
	database, err := Cfg.GetValue("mysql."+section, "database")
	if err != nil {
	}
	charset, err := Cfg.GetValue("mysql."+section, "charset")
	if err != nil {
	}
	str := []string{username, ":", password, "@tcp(", hostname, ":", port, ")/", database, "?charset=", charset}
	config_dictionary["dbs"][section]["connection"] = strings.Join(str, "")
	return config_dictionary["dbs"][section]["connection"], err
}

func LoadConfig(section string) (map[string]string, error) {
	if config_cnfs, exist := config_dictionary["cnf"]; exist {
		if config_cnf, ok := config_cnfs[section]; ok {
			return config_cnf, nil
		}
	} else {
		config_dictionary["cnf"] = make(map[string]map[string]string)
	}
	var info = make(map[string]string)
	ConfigPath := "conf/config.ini"
	Cfg, err := goconfig.LoadConfigFile(ConfigPath)
	if err != nil {
		return info, errors.New("无法加载配置文件")
	}
	if section == "redis" { // redis配置
		info["hostname"], err = Cfg.GetValue(section, "hostname")
		info["password"], err = Cfg.GetValue(section, "password")
		info["port"], err = Cfg.GetValue(section, "port")
		info["num"], _ = Cfg.GetValue(section, "num")
	} else if section == "alidayu" { // alidayu阿里大鱼短信接口配置
		info["app_key"], err = Cfg.GetValue(section, "app_key")
		info["app_secret"], err = Cfg.GetValue(section, "app_secret")
		info["sms"], err = Cfg.GetValue(section, "sms")
	} else {
		return info, errors.New("无法获取键值")
	}
	config_dictionary["cnf"][section] = info
	return info, err
}
