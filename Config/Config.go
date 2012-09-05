package Config

import (
	"Utils"
	"encoding/xml"
	"os"
	"flag"
	"time"
)


type redisConfig struct {
	Addr, Port string
	PoolSize   int
}
type VsConfig struct {
	XMLName    xml.Name `xml:"Config"`
	Redis      redisConfig
	TickTime time.Duration
}

var vsConfig VsConfig
var ConfigFile string
func init() {
	flag.StringVar(&ConfigFile,"c","","config file path")
	if ConfigFile==""{
		ConfigFile="/home/dingbo/workspace/message/src/test.xml"
	}
	ParseXml(ConfigFile)
}
/**
解析xml文件
**/
func ParseXml(configFile string){
	file, err := os.Open(configFile)
	if err != nil {
		Utils.LogPanicErr(err)
		return
	}
	xmlObj := xml.NewDecoder(file)
	err = xmlObj.Decode(&vsConfig)
	if err != nil {
		Utils.LogPanicErr(err)
		return
	}
	Utils.LogInfo("parse xml=%v\n",vsConfig)
}
/**
得到redis的配置
**/
func GetRedisConfig() redisConfig {
	return vsConfig.Redis
}

func GetTickTime()time.Duration{
	return vsConfig.TickTime
}