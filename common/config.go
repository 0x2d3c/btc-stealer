package common

import "fmt"

type DBAuth struct {
	IPPort   string `yaml:"ip_port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Eth struct {
	EtherscanApiKey           string `json:"etherscan_api_key"`
	ScanRequestFrequency      int64  `yaml:"scan_request_frequency"`
	AddressGenerationQuantity int64  `yaml:"address_generation_quantity"`
}

type Proxy struct {
	Enable  bool   `yaml:"enable"`
	Address string `yaml:"address"`
}

type Config struct {
	DB        DBAuth `yaml:"db"`
	Eth       Eth    `yaml:"eth"`
	Proxy     Proxy  `yaml:"proxy"`
	WordsList string `yaml:"words_list"`
}

var config *Config

func InitConfig(cfg *Config) {
	initGORM(cfg.DB)

	config = cfg
}

func GetScanRequestFrequency() int64 {
	return config.Eth.ScanRequestFrequency
}

func GetETHScanAPIAddress() string {
	return fmt.Sprintf("https://api.etherscan.io/api?module=account&action=balancemulti&tag=latest&apikey=%s&address=", config.Eth.EtherscanApiKey)
}
