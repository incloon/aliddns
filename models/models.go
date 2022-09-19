package models

type ConfigModel struct {
	AccessId            string // 阿里云的 Access Id
	AccessKey           string // 阿里云的 Access Key
	MainDomain          string // 需要更新的主域名，例如 iotserv.com
	SubDomainName       string // 需要更新的具体子域名，例如 www
	CheckUpdateInterval int    // 检查域名是否改变的时间间隔，单位秒，默认30秒
	Protocol            string // "ipv4"或"ipv6"或"all"，默认"all"
	NetworkAdapter      string // 网卡名，如果为空则获取公网IP
}
