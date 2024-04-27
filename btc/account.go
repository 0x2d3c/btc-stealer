package btc

import (
	"log"

	"btc-stealer/common"
	"btc-stealer/data"
)

func AddressBTCCheck() {
	for _, bit := range common.Bits {
		mnemonic := common.Mnemonic(bit)
		addressGenOnceBTC(bit, mnemonic)
	}
}

func addressGenOnceBTC(bitSize int, mnemonic string) {
	coins, address := addressGenBTC(bitSize, mnemonic)

	for _, addr := range address {
		for {
			has, err := common.HttpBTC(addr)
			if err != nil {
				continue
			}

			if has {
				coin, ok := coins[addr]
				if !ok {
					break
				}
				common.RecordBalance(coin.String())
			}

			break
		}
	}
}

func addressGenBTC(bitSize int, mnemonic string) (map[string]*data.Btc, []string) {
	km, err := common.NewKeyManager(bitSize, "", mnemonic)
	if err != nil {
		log.Fatal(err)
	}

	masterKey, err := km.GetMasterKey()
	if err != nil {
		log.Fatal(err)
	}

	rk := masterKey.B58Serialize()

	coins := make(map[string]*data.Btc)

	var address []string
	{
		key, err := km.GetKey(common.PurposeBIP44, common.CoinTypeBTC, 0, 0, 0)
		if err != nil {
			log.Fatal(err)
		}

		wif, addres, _, _, _, err := key.Encode(false)
		if err != nil {
			log.Fatal(err)
		}

		coins[addres] = &data.Btc{
			RootKey:  rk,
			Wif:      wif,
			Address:  addres,
			Mnemonic: mnemonic,
		}

		address = append(address, addres)
	}

	{
		key, err := km.GetKey(common.PurposeBIP49, common.CoinTypeBTC, 0, 0, 0)
		if err != nil {
			log.Fatal(err)
		}

		wif, _, _, segwitNested, _, err := key.Encode(false)
		if err != nil {
			log.Fatal(err)
		}

		coins[segwitNested] = &data.Btc{
			RootKey:  rk,
			Wif:      wif,
			Address:  segwitNested,
			Mnemonic: mnemonic,
		}

		address = append(address, segwitNested)
	}

	{
		key, err := km.GetKey(common.PurposeBIP84, common.CoinTypeBTC, 0, 0, 0)
		if err != nil {
			log.Fatal(err)
		}

		wif, _, segwitBech32, _, _, err := key.Encode(false)
		if err != nil {
			log.Fatal(err)
		}

		coins[segwitBech32] = &data.Btc{
			RootKey:  rk,
			Wif:      wif,
			Address:  segwitBech32,
			Mnemonic: mnemonic,
		}

		address = append(address, segwitBech32)
	}

	{
		key, err := km.GetKey(common.PurposeBIP86, common.CoinTypeBTC, 0, 0, 0)
		if err != nil {
			log.Fatal(err)
		}

		wif, _, _, _, taproot, err := key.Encode(false)
		if err != nil {
			log.Fatal(err)
		}

		coins[taproot] = &data.Btc{
			RootKey:  rk,
			Wif:      wif,
			Address:  taproot,
			Mnemonic: mnemonic,
		}

		address = append(address, taproot)
	}

	return coins, address
}
