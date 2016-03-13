package client

import (
	"fmt"
	"strconv"

	ini "github.com/vaughan0/go-ini"
)

// common config
var (
	ServerAddr        string = "0.0.0.0"
	ServerPort        int64  = 7000
	LogFile           string = "console"
	LogWay            string = "console"
	LogLevel          string = "info"
	HeartBeatInterval int64  = 20
	HeartBeatTimeout  int64  = 60
)

var ProxyClients map[string]*ProxyClient = make(map[string]*ProxyClient)

func LoadConf(confFile string) (err error) {
	var tmpStr string
	var ok bool

	conf, err := ini.LoadFile(confFile)
	if err != nil {
		return err
	}

	// common
	tmpStr, ok = conf.Get("common", "server_addr")
	if ok {
		ServerAddr = tmpStr
	}

	tmpStr, ok = conf.Get("common", "server_port")
	if ok {
		ServerPort, _ = strconv.ParseInt(tmpStr, 10, 64)
	}

	tmpStr, ok = conf.Get("common", "log_file")
	if ok {
		LogFile = tmpStr
		if LogFile == "console" {
			LogWay = "console"
		} else {
			LogWay = "file"
		}
	}

	tmpStr, ok = conf.Get("common", "log_level")
	if ok {
		LogLevel = tmpStr
	}

	// proxies
	for name, section := range conf {
		if name != "common" {
			proxyClient := &ProxyClient{}
			proxyClient.Name = name

			proxyClient.Passwd, ok = section["passwd"]
			if !ok {
				return fmt.Errorf("Parse ini file error: proxy [%s] no passwd found", proxyClient.Name)
			}

			proxyClient.LocalIp, ok = section["local_ip"]
			if !ok {
				// use 127.0.0.1 as default
				proxyClient.LocalIp = "127.0.0.1"
			}

			portStr, ok := section["local_port"]
			if ok {
				proxyClient.LocalPort, err = strconv.ParseInt(portStr, 10, 64)
				if err != nil {
					return fmt.Errorf("Parse ini file error: proxy [%s] local_port error", proxyClient.Name)
				}
			} else {
				return fmt.Errorf("Parse ini file error: proxy [%s] local_port not found", proxyClient.Name)
			}

			ProxyClients[proxyClient.Name] = proxyClient
		}
	}

	if len(ProxyClients) == 0 {
		return fmt.Errorf("Parse ini file error: no proxy config found")
	}

	return nil
}
