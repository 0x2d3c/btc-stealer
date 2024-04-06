package common

type Proxy struct {
	Enable  bool   `json:"enable"`
	Address string `json:"address"`
}

type Config struct {
	Proxy           Proxy  `json:"proxy"`
	WordsList       string `json:"words_list"`
	EtherscanApiKey string `json:"etherscan_api_key"`
}

var (
	config *Config
)

func InitConfig(cfg *Config) {
	config = cfg

	initLogFile()
}

func GetETHScanAPIAddress() string {
	return "https://api.etherscan.io/api?module=account&action=balancemulti&tag=latest&apikey=" + config.EtherscanApiKey + "&address="
}
