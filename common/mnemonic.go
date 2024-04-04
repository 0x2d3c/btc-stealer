package common

import (
	"github.com/tyler-smith/go-bip39"
	"github.com/tyler-smith/go-bip39/wordlists"
)

const (
	ChineseSimplified  = "chinese_simplified"
	ChineseTraditional = "chinese_traditional"
	Czech              = "czech"
	English            = "english"
	French             = "french"
	Italian            = "italian"
	Japanese           = "japanese"
	Korean             = "korean"
	Spanish            = "spanish"
)

func SetWords(lang string) {
	switch lang {
	case Czech:
		bip39.SetWordList(wordlists.Czech)
	case French:
		bip39.SetWordList(wordlists.French)
	case Korean:
		bip39.SetWordList(wordlists.Korean)
	case English:
		bip39.SetWordList(wordlists.English)
	case Italian:
		bip39.SetWordList(wordlists.Italian)
	case Spanish:
		bip39.SetWordList(wordlists.Spanish)
	case Japanese:
		bip39.SetWordList(wordlists.Japanese)
	case ChineseSimplified:
		bip39.SetWordList(wordlists.ChineseSimplified)
	case ChineseTraditional:
		bip39.SetWordList(wordlists.ChineseTraditional)
	}
}

func Mnemonic(bitSize int) string {
	entropy, err := bip39.NewEntropy(bitSize)
	if err != nil {
		panic(err)
	}

	words, err := bip39.NewMnemonic(entropy)
	if err != nil {
		panic(err)
	}

	return words
}
