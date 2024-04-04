package main

import (
	"btc-stealer/common"
	"btc-stealer/data"
	"btc-stealer/eth"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

func init() {

	file, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}

	var cfg common.Config
	if err = yaml.Unmarshal(file, &cfg); err != nil {
		panic(err)
	}

	common.InitConfig(&cfg)
	common.SetWords(cfg.WordsList)
	common.DB.AutoMigrate(&data.Eth{}, &data.Btc{})
}

func RunCheck() {

	eth.AddressGenETH(128)
	eth.AddressGenETH(256)
	//m128, m256 := common.Mnemonic(128), common.Mnemonic(256)

	//btc.AddressGenBTC(128, "", m128, false)
	//btc.AddressGenBTC(128, "", m128, true)
	//btc.AddressGenBTC(256, "", m256, false)
	//btc.AddressGenBTC(256, "", m256, true)
}

func main() {
	freq := time.Millisecond * time.Duration(common.GetScanRequestFrequency())
	ticker := time.NewTicker(freq)
	for range ticker.C {
		go RunCheck()
	}
}
