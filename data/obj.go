package data

type Eth struct {
	RootKey    string
	Address    string
	Mnemonic   string
	PrivateKey string
}

func (e *Eth) String() string {
	return "Address:" + e.Address + "\nPrivateKey:" + e.PrivateKey + "\nRootKey:" + e.RootKey + "" + "\nMnemonic:" + e.Mnemonic
}

func (e *Eth) RecordString() string {
	return "Address:" + e.Address + "\nPrivateKey:" + e.PrivateKey
}

type Btc struct {
	Wif        string
	Address    string
	RootKey    string
	Mnemonic   string
	PrivateKey string
}

func (b *Btc) String() string {
	return "Address:" + b.Address + "\nPrivateKey:" + b.PrivateKey + "\nRootKey:" + b.RootKey + "" + "\nWif:" + b.Wif + "\nMnemonic:" + b.Mnemonic
}
