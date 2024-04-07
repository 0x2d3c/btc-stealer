package common

type Proxy struct {
	Enable  bool   `json:"enable"`
	Address string `json:"address"`
}

type Wallet struct {
	ETH []string `json:"eth"`
	BTC []string `json:"btc"`
}

type Config struct {
	Mode            int    `json:"mode"`
	Proxy           Proxy  `json:"proxy"`
	Wallet          Wallet `json:"wallet"`
	WordsList       string `json:"words_list"`
	EtherscanApiKey string `json:"etherscan_api_key"`
}

var (
	config *Config
)

const (
	ModeMix = iota
	ModeOnline
	ModeOffline
)

func InitConfig(cfg *Config) {
	config = cfg

	initLogFile()

	switch cfg.Mode {
	case ModeMix, ModeOffline:
		// both wallet file need config
		if len(cfg.Wallet.ETH) == 0 || len(cfg.Wallet.ETH) == 0 {
			panic("invalid mode config")
		}

		initLoadETHWallet(cfg.Wallet.ETH)
		initLoadBTCWallet(cfg.Wallet.BTC)
	case ModeOnline:
	default:
		panic("invalid mode")
	}
}
func GetMode() int {
	return config.Mode
}

func GetETHScanAPIAddress() string {
	return "https://api.etherscan.io/api?module=account&action=balancemulti&tag=latest&apikey=" + config.EtherscanApiKey + "&address="
}
