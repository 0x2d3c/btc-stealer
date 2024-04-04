package data

import (
	"btc-stealer/common"
	"log"

	"gorm.io/gorm"
)

type Eth struct {
	gorm.Model
	RootKey    string
	Address    string
	Mnemonic   string
	PrivateKey string
}

type Btc struct {
	gorm.Model
	Wif        string
	Address    string
	RootKey    string
	Mnemonic   string
	PrivateKey string
}

func SaveCoins(coins interface{}) {
	if err := common.DB.Create(coins).Error; err != nil {
		log.Println("save coins, err", err.Error(), "\ncoins:\n", coins)
	}
}
