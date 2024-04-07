package btc

import (
	"log"

	"btc-stealer/common"
	"btc-stealer/data"
)

func RunBTCOfflineCheck() {
	cs := []bool{true, false}
	for {
		for _, bit := range common.Bits {
			mnemonic := common.Mnemonic(bit)
			for _, c := range cs {
				addressGenOnceBTC(bit, mnemonic, c)
			}
		}
	}
}

func addressGenOnceBTC(bitSize int, mnemonic string, compress bool) {
	coins, address := addressGenBTC(bitSize, mnemonic, compress)

	has, _ := common.OfflineBTCCheck(address)
	for _, wallet := range has {
		coin, ok := coins[wallet]
		if !ok {
			continue
		}
		common.RecordBalance(coin.String())
	}
}

func addressGenBTC(bitSize int, mnemonic string, compress bool) (map[string]*data.Btc, []string) {
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

		wif, addres, _, _, _, err := key.Encode(compress)
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

		wif, _, _, segwitNested, _, err := key.Encode(compress)
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

		wif, _, segwitBech32, _, _, err := key.Encode(compress)
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

		wif, _, _, _, taproot, err := key.Encode(compress)
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
