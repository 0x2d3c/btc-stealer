package btc

import (
	"fmt"
	"log"

	"btc-stealer/common"
	"btc-stealer/data"
)

func AddressGenBTC(bitSize int, passphrase, mnemonic string, compress bool) {
	km, err := common.NewKeyManager(bitSize, passphrase, mnemonic)
	if err != nil {
		log.Fatal(err)
	}

	masterKey, err := km.GetMasterKey()
	if err != nil {
		log.Fatal(err)
	}

	rk := masterKey.B58Serialize()

	coins := make([]*data.Btc, 0)
	for i := 0; i < 10; i++ {
		key, err := km.GetKey(common.PurposeBIP44, common.CoinTypeBTC, 0, 0, uint32(i))
		if err != nil {
			log.Fatal(err)
		}

		wif, address, _, _, _, err := key.Encode(compress)
		if err != nil {
			log.Fatal(err)
		}

		coins = append(coins, &data.Btc{
			RootKey:  rk,
			Wif:      wif,
			Address:  address,
			Mnemonic: mnemonic,
		})
	}
	for i := 0; i < 10; i++ {
		key, err := km.GetKey(common.PurposeBIP49, common.CoinTypeBTC, 0, 0, uint32(i))
		if err != nil {
			log.Fatal(err)
		}

		wif, _, _, segwitNested, _, err := key.Encode(compress)
		if err != nil {
			log.Fatal(err)
		}

		coins = append(coins, &data.Btc{
			RootKey:  rk,
			Wif:      wif,
			Address:  segwitNested,
			Mnemonic: mnemonic,
		})
	}
	for i := 0; i < 10; i++ {
		key, err := km.GetKey(common.PurposeBIP84, common.CoinTypeBTC, 0, 0, uint32(i))
		if err != nil {
			log.Fatal(err)
		}

		wif, _, segwitBech32, _, _, err := key.Encode(compress)
		if err != nil {
			log.Fatal(err)
		}

		coins = append(coins, &data.Btc{
			RootKey:  rk,
			Wif:      wif,
			Address:  segwitBech32,
			Mnemonic: mnemonic,
		})
	}
	for i := 0; i < 10; i++ {
		key, err := km.GetKey(common.PurposeBIP86, common.CoinTypeBTC, 0, 0, uint32(i))
		if err != nil {
			log.Fatal(err)
		}

		wif, _, _, _, taproot, err := key.Encode(compress)
		if err != nil {
			log.Fatal(err)
		}

		coins = append(coins, &data.Btc{
			RootKey:  rk,
			Wif:      wif,
			Address:  taproot,
			Mnemonic: mnemonic,
		})
	}

	// TODO:
	// scan & check & save

	fmt.Printf("%+v \n", coins)
}
