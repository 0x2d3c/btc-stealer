package common

var btcWallets = make(map[string]struct{})

func initLoadBTCWallet(files []string) {
	btcWallets = readWalletAddressFile(files)
}

func OfflineBTCCheck(wallets []string) ([]string, []string) {
	var (
		has    []string
		hasNot []string
	)
	for _, wallet := range wallets {
		if _, ok := btcWallets[wallet]; ok {
			has = append(has, wallet)

			continue
		}
		hasNot = append(hasNot, wallet)
	}
	return has, hasNot
}
