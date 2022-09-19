package utils

import (
	"github.com/incloon/aliddns/config"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
)

//u4="http://ipv4.ident.me http://ipv4.icanhazip.com http://nsupdate.info/myip http://whatismyip.akamai.com http://ipv4.myip.dk/api/info/IPv4Address http://checkip4.spdyn.de http://v4.ipv6-test.com/api/myip.php http://checkip.amazonaws.com http://ipinfo.io/ip http://bot.whatismyipaddress.com http://ipv4.ident.me http://ipv4.icanhazip.com http://nsupdate.info/myip http://whatismyip.akamai.com http://ipv4.myip.dk/api/info/IPv4Address http://checkip4.spdyn.de http://v4.ipv6-test.com/api/myip.php http://checkip.amazonaws.com http://ipinfo.io/ip http://bot.whatismyipaddress.com"
//u6="http://ipv6.ident.me http://ipv6.icanhazip.com http://ipv6.ident.me http://ipv6.icanhazip.com http://ipv6.yunohost.org http://v6.ipv6-test.com/api/myip.php http://ipv6.ident.me http://ipv6.icanhazip.com http://ipv6.yunohost.org http://v6.ipv6-test.com/api/myip.php http://ipv6.ident.me http://ipv6.icanhazip.com http://ipv6.ident.me http://ipv6.icanhazip.com http://ipv6.yunohost.org http://v6.ipv6-test.com/api/myip.php http://ipv6.ident.me http://ipv6.icanhazip.com http://ipv6.yunohost.org http://v6.ipv6-test.com/api/myip.php"

var Ipv4APIUrls = []string{
	"http://whatismyip.akamai.com",
	"http://v4.ipv6-test.com/api/myip.php",
	"http://checkip.amazonaws.com",
	"api.ipify.org",
	"canhazip.com",
	"ident.me",
	"whatismyip.akamai.com",
	"myip.dnsomatic.com",
	"http://members.3322.org/dyndns/getip",
	"http://ifconfig.me/ip",
	"http://ip.3322.net",
	"https://myexternalip.com/raw",
	"http://ipv4.ident.me",
	"http://ipv4.icanhazip.com",
	"http://nsupdate.info/myip",
	"http://ipv4.myip.dk/api/info/IPv4Address",
	"http://checkip4.spdyn.de",
	"http://ipinfo.io/ip",
}
var Ipv6APIUrls = []string{
	"http://v6.ipv6-test.com/api/myip.php",
	"http://bbs6.ustc.edu.cn/cgi-bin/myip",
	"http://ipv6.ident.me",
	"http://ipv6.icanhazip.com",
	"http://ipv6.yunohost.org",
}

func GetMyPublicIpv4() string {
	for _, url := range Ipv4APIUrls {
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("get public ipv4 err: %s", err)
			continue
		}
		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("get public ipv4 err: %s", err)
			_ = resp.Body.Close()
			continue
		}
		ipv4 := strings.Replace(string(bytes), "\n", "", -1)
		ip := net.ParseIP(ipv4)
		if ip != nil {
			log.Println("got public ipv4 addr: ", ip.String())
			_ = resp.Body.Close()
			return ip.String()
		}
	}
	log.Println("fail to fetch ipv4!")
	return ""
}

func GetMyPublicIpv6() string {
	for _, url := range Ipv6APIUrls {
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("get public ipv6 err: %s", err)
			continue
		}
		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("get public ipv6 err: %s", err)
			_ = resp.Body.Close()
			continue
		}
		tmp := strings.Replace(string(bytes), "document.write('", "", -1)
		tmp = strings.Replace(tmp, "');", "", -1)
		ipv6 := strings.Replace(tmp, "\n", "", -1)
		ip := net.ParseIP(ipv6)
		if ip != nil {
			log.Println("got public ipv6 addr: ", ip.String())
			_ = resp.Body.Close()
			return ip.String()
		}
	}
	log.Println("fail to fetch ipv6!")
	return ""
}

func GetNetworkAdapterAddr(name string) []net.Addr {
	inter, err := net.InterfaceByName(name)
	if err != nil {
		panic(err)
	}
	// 判断网卡是否开启，过滤本地环回接口
	if inter.Flags&net.FlagUp == 0{
		log.Printf("network adapter %s is disable", name)
		return []net.Addr{}
	}
	// 获取网卡下所有的地址
	addrs, err := inter.Addrs()
	if err != nil {
		log.Printf("network adapter %s doesn't have avaliable address", name)
		return []net.Addr{}
	}
	return addrs
}

func GetLocalIpv4(name string) string {
	addrs := GetNetworkAdapterAddr(name)
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				log.Printf("got %s ipv4 addr: %s", name, ipnet.IP.String())
				return ipnet.IP.String()
			}
		}
	}
	log.Printf("network adapter %s doesn't have ipv4", name)
	return ""
}

func GetLocalIpv6(name string) string {
	addrs := GetNetworkAdapterAddr(name)
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() == nil {
				log.Printf("got %s ipv6 addr: %s", name, ipnet.IP.String())
				return ipnet.IP.String()
			}
		}
	}
	log.Printf("network adapter %s doesn't have ipv6", name)
	return ""
}

func GetSubDomains(mainDomian string) (*alidns.DescribeDomainRecordsResponse, error) {
	client, err := GetAliYunClient()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	request := alidns.CreateDescribeDomainRecordsRequest()
	request.Scheme = "https"
	request.DomainName = mainDomian
	return client.DescribeDomainRecords(request)
}

func UpdateSubDomain(subDomain *alidns.Record) error {
	client, err := GetAliYunClient()
	if err != nil {
		log.Println(err)
		return err
	}
	request := alidns.CreateUpdateDomainRecordRequest()
	request.Scheme = "https"
	request.RecordId = subDomain.RecordId
	request.RR = subDomain.RR
	request.Type = subDomain.Type
	request.Value = subDomain.Value
	request.TTL = requests.NewInteger64(subDomain.TTL)

	_, err = client.UpdateDomainRecord(request)
	if err != nil {
		log.Print("UpdateDomainRecord:", err)
		return err
	}
	return nil
}

func AddSubDomainRecord(subDomain *alidns.Record) error {
	client, err := GetAliYunClient()
	if err != nil {
		log.Println(err)
		return err
	}

	request := alidns.CreateAddDomainRecordRequest()
	request.Scheme = "https"
	request.DomainName = subDomain.DomainName
	request.RR = subDomain.RR
	request.Type = subDomain.Type
	request.Value = subDomain.Value
	request.TTL = requests.NewInteger64(subDomain.TTL)

	_, err = client.AddDomainRecord(request)
	if err != nil {
		log.Print("AddSubDomainRecord:", err)
		return err
	}
	return nil
}

func GetAliYunClient() (*alidns.Client, error) {
	return alidns.NewClientWithAccessKey("cn-hangzhou", config.ConfigModel.AccessId, config.ConfigModel.AccessKey)
}
