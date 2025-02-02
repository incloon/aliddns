package config

import (
	"fmt"
	"github.com/incloon/aliddns/models"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var ConfigFileName = "aliddns.yaml"
var ConfigFilePath = fmt.Sprintf("./%s", ConfigFileName)
var ConfigModel = &models.ConfigModel{
	AccessId:            "*AccessId",
	AccessKey:           "*AccessKey",
	MainDomain:          "*example.com",
	SubDomainName:       "*www",
	CheckUpdateInterval: 30,
	Protocol:            "all",
	NetworkAdapter:      "",
}

//将配置写入指定的路径的文件
func WriteConfigFile(ConfigMode *models.ConfigModel, path string) (err error) {
	configByte, err := yaml.Marshal(ConfigMode)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	if ioutil.WriteFile(path, configByte, 0644) == nil {
		return
	}
	return
}

func InitConfigFile() {
	//	生成配置文件模板
	err := os.MkdirAll(filepath.Dir(ConfigFilePath), 0644)
	if err != nil {
		return
	}
	err = WriteConfigFile(ConfigModel, ConfigFilePath)
	if err == nil {
		fmt.Println("config created")
		return
	}
	log.Fatalln("Write config file error, Please check if the program has write permission! Or create a configuration file manually.")
}

var SupportedProtocols = [3]string{"ipv4", "ipv6", "all"}

func UseConfigFile() {
	//配置文件存在
	log.Println("config file path: ", ConfigFilePath)
	content, err := ioutil.ReadFile(ConfigFilePath)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	err = yaml.Unmarshal(content, ConfigModel)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	// 判断 protocol 是否正确
	protocol := strings.ToLower(ConfigModel.Protocol)
	for _, val := range SupportedProtocols {
		if val == protocol {
			return
		}
	}
	log.Fatalf("Not support protocol type '%s'. Available values : [ipv4, ipv6, all]\n", ConfigModel.Protocol)
	return
}

func LoadSnapcraftConfigPath() {
	//是否是snapcraft应用，如果是则从snapcraft指定的工作目录保存配置文件
	appDataPath, havaAppDataPath := os.LookupEnv("SNAP_USER_DATA")
	if havaAppDataPath {
		ConfigFilePath = filepath.Join(appDataPath, ConfigFileName)
	}
}
