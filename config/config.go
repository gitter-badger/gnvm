package config

import (
	// lib
	. "github.com/Kenshin/cprint"
	"github.com/tsuru/config"

	// go
	"os"
	"regexp"
	"strings"

	// local
	"gnvm/util"
)

var configPath, globalversion, latsetversion string

const (
	VERSION  = "0.1.4 beta"
	CONFIG   = ".gnvmrc"
	NEWLINE  = "\n"
	UNKNOWN  = "unknown"
	LATEST   = "latest"
	NODELIST = "index.json"

	REGISTRY     = "registry"
	REGISTRY_KEY = "registry: "
	REGISTRY_VAL = "http://nodejs.org/dist/"
	TAOBAO       = "http://npm.taobao.org/mirrors/node"

	NODEROOT     = "noderoot"
	NODEROOT_KEY = "noderoot: "
	NODEROOT_VAL = "root"

	GLOBAL_VERSION     = "globalversion"
	GLOBAL_VERSION_KEY = "globalversion: "
	GLOBAL_VERSION_VAL = UNKNOWN

	LATEST_VERSION     = "latestversion"
	LATEST_VERSION_KEY = "latestversion: "
	LATEST_VERSION_VAL = UNKNOWN

	//CURRENT_VERSION     = "currentversion"
	//CURRENT_VERSION_KEY = "currentversion: "
	//CURRENT_VERSION_VAL = UNKNOWN
)

func init() {

	// try catch
	defer func() {
		if err := recover(); err != nil {
			Error(ERROR, "gnvm.exe an error has occurred. please check. \nError: ", err)
			os.Exit(0)
		}
	}()

	// set config path
	configPath = util.GlobalNodePath + "\\" + CONFIG

	// config file is exist
	file, err := os.Open(configPath)
	defer file.Close()
	if err != nil && os.IsNotExist(err) {
		P(WARING, "config file %v is not exist.\n", configPath)
		createConfig()
	}

	// read config
	readConfig()

}

func createConfig() {

	// create file
	file, err := os.Create(configPath)
	defer file.Close()
	if err != nil {
		P(ERROR, "config file create Error: %v\n", err.Error())
		return
	}

	// get <root>/node.exe version
	version, err := util.GetNodeVersion(util.GlobalNodePath + "\\")
	if err != nil {
		P(WARING, "not found global node.exe version, please use '%v'. See '%v'.\n", "gnvm install x.xx.xx -g", "gnvm help install")
		globalversion = GLOBAL_VERSION_VAL
	} else {
		globalversion = version
	}

	//write init config
	_, fileErr := file.WriteString(REGISTRY_KEY + REGISTRY_VAL + NEWLINE + NODEROOT_KEY + util.GlobalNodePath + NEWLINE + GLOBAL_VERSION_KEY + globalversion + NEWLINE + LATEST_VERSION_KEY + LATEST_VERSION_VAL)
	if fileErr != nil {
		P(ERROR, "write config file Error: %v\n", fileErr.Error())
		return
	}

	P(DEFAULT, "Config file %v create success.\n", configPath)
	//P(NOTICE, "if you first run gnvm.exe, please use %v or %v.", "gnvm config INIT", "gnvm config registry TAOBAO", "\n")

}

func readConfig() {
	if err := config.ReadConfigFile(configPath); err != nil {
		P(ERROR, "read config file fail, please use '%v'. \nError: %v\n", "gnvm config INIT", err.Error())
		return
	}
}

func SetConfig(key string, value interface{}) string {

	if key == "registry" {

		if !strings.HasPrefix(value.(string), "http://") {
			P(WARING, "%v need %v", value.(string), "http://", "\n")
			value = "http://" + value.(string)
		}

		reg, _ := regexp.Compile(`^https?:\/\/(w{3}\.)?(\w+\.)+([a-zA-Z]{2,})(:\d{1,4})?\/?($)?`)

		switch {
		case !reg.MatchString(value.(string)):
			P(ERROR, "registry value %v must url valid.\n", value.(string))
			return ""
		case !strings.HasSuffix(value.(string), "/"):
			value = value.(string) + "/"
		}
	}

	// set new value
	config.Set(key, value)

	// delete old config
	if err := os.Remove(configPath); err != nil {
		P(ERROR, "remove config file Error: %v\n", err.Error())
	}

	// write new config
	if err := config.WriteConfigFile(configPath, 0777); err != nil {
		P(ERROR, "write config file Error: %v\n", err.Error())
	}

	return value.(string)

}

func GetConfig(key string) string {
	value, err := config.GetString(key)
	if err != nil {
		P(ERROR, "get config Error: %v\n", err.Error())
		value = UNKNOWN
	}
	return value
}

func ReSetConfig() {
	if newValue := SetConfig(REGISTRY, REGISTRY_VAL); newValue != "" {
		P(NOTICE, "%v      init success, new value is %v\n", REGISTRY, newValue)
	}
	if newValue := SetConfig(NODEROOT, util.GlobalNodePath); newValue != "" {
		P(NOTICE, "%v      init success, new value is %v\n", NODEROOT, newValue)
	}
	version, err := util.GetNodeVersion(util.GlobalNodePath + "\\")
	if err != nil {
		P(WARING, "not found global node.exe version, please use '%v'. See '%v'.\n", "gnvm install x.xx.xx -g", "gnvm help install")
		globalversion = GLOBAL_VERSION_VAL
	} else {
		globalversion = version
	}
	if newValue := SetConfig(GLOBAL_VERSION, globalversion); newValue != "" {
		P(NOTICE, "%v init success, new value is %v\n", GLOBAL_VERSION, newValue)
	}
	/*
		url := REGISTRY_VAL + "latest/" + util.SHASUMS
		P(NOTICE, "get node.exe latest version from %v, please wait.", url, "\n")
		if latest := util.GetLatestVersion(url); latest != "" {
			latsetversion = latest
		} else {
			latsetversion = LATEST_VERSION_VAL
		}
		if newValue := SetConfig(LATEST_VERSION, latsetversion); newValue != "" {
			P(NOTICE, "%v init success, new value is %v\n", LATEST_VERSION, newValue)
		}
	*/
}