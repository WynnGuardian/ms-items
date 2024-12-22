package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/wynnguardian/common/utils"
)

type PrivateConfig struct {
	Tokens struct {
		Self      string   `json:"self"`
		Whitelist []string `json:"whitelist"`
	} `json:"tokens"`
	DB struct {
		Hostname string `json:"hostname"`
		Port     int    `json:"port"`
		Password string `json:"password"`
		Username string `json:"username"`
		Database string `json:"database"`
	} `json:"database"`
}

type ServerConfig struct {
	Port int `json:"port"`
}

type HostsConfig struct {
	Discord string `json:"discord"`
}

type Config struct {
	Private PrivateConfig
	Hosts   HostsConfig
	Server  ServerConfig
}

var MainConfig *Config = &Config{}

func Load() {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	mustReadAndAssign(pwd, "config/private.json", &MainConfig.Private)
	mustReadAndAssign(pwd, "config/hosts.json", &MainConfig.Hosts)
	mustReadAndAssign(pwd, "config/server.json", &MainConfig.Server)
}

func mustReadAndAssign(pwd, relativeDir string, target interface{}) {
	f := utils.MustVal(os.ReadFile(fmt.Sprintf("%s/%s", pwd, relativeDir)))
	utils.Must(json.Unmarshal(f, &target))
}
