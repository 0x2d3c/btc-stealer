package common

type Config struct {
	WordsList string `json:"words_list"`
	ETHGW     string `json:"eth_gw"`
	BTCGW     string `json:"btc_gw"`
}

var config *Config

func InitConfig(cfg *Config) {
	config = cfg

	initLogFile()
}
