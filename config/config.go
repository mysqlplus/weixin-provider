package config

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/toolkits/file"
)

type HttpConfig struct {
	Listen string `json:"listen"`
	Token  string `json:"token"`
}

type WeixinConfig struct {
	Tokenurl   string `json:"tokenurl"`
	Sendurl    string `json:"sendurl"`
	Agentid    int    `json:"agentid"`
	Corpid     string `json:"corpid"`
	Corpsecret string `json:"corpsecret"`
}

type GlobalConfig struct {
	Debug  bool          `json:"debug"`
	Http   *HttpConfig   `json:"http"`
	Weixin *WeixinConfig `json:"weixin"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	configLock = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

func Parse(cfg string) error {
	if cfg == "" {
		return fmt.Errorf("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		return fmt.Errorf("configuration file %s is nonexistent", cfg)
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		return fmt.Errorf("read configuration file %s fail %s", cfg, err.Error())
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		return fmt.Errorf("parse configuration file %s fail %s", cfg, err.Error())
	}

	configLock.Lock()
	defer configLock.Unlock()
	config = &c

	log.Println("load configuration file", cfg, "successfully")
	return nil

}