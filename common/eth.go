package common

var ethWallets = make(map[string]struct{})

func initLoadETHWallet(files []string) {
	ethWallets = readWalletAddressFile(files)
}

func OfflineETHCheck(wallets []string) ([]string, []string) {
	var (
		has    []string
		hasNot []string
	)
	for _, wallet := range wallets {
		if _, ok := ethWallets[wallet]; ok {
			has = append(has, wallet)

			continue
		}
		hasNot = append(hasNot, wallet)
	}
	return has, hasNot
}
