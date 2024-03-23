package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/sha3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Purpose BIP43 - Purpose Field for Deterministic Wallets
// https://github.com/bitcoin/bips/blob/master/bip-0043.mediawiki
//
// Purpose is a constant set to 44' (or 0x8000002C) following the BIP43 recommendation.
// It indicates that the subtree of this node is used according to this specification.
//
// What does 44' mean in BIP44?
// https://bitcoin.stackexchange.com/questions/74368/what-does-44-mean-in-bip44
//
// 44' means that hardened keys should be used. The distinguisher for whether
// a key a given index is hardened is that the index is greater than 2^31,
// which is 2147483648. In hex, that is 0x80000000. That is what the apostrophe (') means.
// The 44 comes from adding it to 2^31 to get the final hardened key index.
// In hex, 44 is 2C, so 0x80000000 + 0x2C = 0x8000002C.
type Purpose = uint32

const (
	PurposeBIP44 Purpose = 0x8000002C // 44' BIP44
	PurposeBIP49 Purpose = 0x80000031 // 49' BIP49
	PurposeBIP84 Purpose = 0x80000054 // 84' BIP84
	PurposeBIP86 Purpose = 0x80000056 // 86' BIP86 //taprrot
)

// CoinType SLIP-0044 : Registered coin types for BIP-0044
// https://github.com/satoshilabs/slips/blob/master/slip-0044.md
type CoinType = uint32

const (
	CoinTypeBTC CoinType = 0x80000000
	CoinTypeLTC CoinType = 0x80000002
	CoinTypeETH CoinType = 0x8000003c
	CoinTypeEOS CoinType = 0x800000c2
)

const (
	Apostrophe uint32 = 0x80000000 // 0'
)

type Key struct {
	path     string
	bip32Key *bip32.Key
}

func (k *Key) Encode(compress bool) (wif, address, segwitBech32, segwitNested, taproot string, err error) {
	prvKey, _ := btcec.PrivKeyFromBytes(k.bip32Key.Key)
	return GenerateFromBytes(prvKey, compress)
}

// https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki
// bip44 define the following 5 levels in BIP32 path:
// m / purpose' / coin_type' / account' / change / address_index

func (k *Key) GetPath() string {
	return k.path
}

type KeyManager struct {
	mnemonic   string
	passphrase string
	keys       map[string]*bip32.Key
	mux        sync.Mutex
}

// NewKeyManager return new key manager
// bitSize has to be a multiple 32 and be within the inclusive range of {128, 256}
// 128: 12 phrases
// 256: 24 phrases
func NewKeyManager(bitSize int, passphrase, mnemonic string) (*KeyManager, error) {
	if mnemonic == "" {
		entropy, err := bip39.NewEntropy(bitSize)
		if err != nil {
			return nil, err
		}
		mnemonic, err = bip39.NewMnemonic(entropy)
		if err != nil {
			return nil, err
		}
	}

	km := &KeyManager{
		mnemonic:   mnemonic,
		passphrase: passphrase,
		keys:       make(map[string]*bip32.Key, 0),
	}
	return km, nil
}

func (km *KeyManager) GetMnemonic() string {
	return km.mnemonic
}

func (km *KeyManager) GetPassphrase() string {
	return km.passphrase
}

func (km *KeyManager) GetSeed() []byte {
	return bip39.NewSeed(km.GetMnemonic(), km.GetPassphrase())
}

func (km *KeyManager) getKey(path string) (*bip32.Key, bool) {
	km.mux.Lock()
	defer km.mux.Unlock()

	key, ok := km.keys[path]
	return key, ok
}

func (km *KeyManager) setKey(path string, key *bip32.Key) {
	km.mux.Lock()
	defer km.mux.Unlock()

	km.keys[path] = key
}

func (km *KeyManager) GetMasterKey() (*bip32.Key, error) {
	path := "m"

	key, ok := km.getKey(path)
	if ok {
		return key, nil
	}

	key, err := bip32.NewMasterKey(km.GetSeed())
	if err != nil {
		return nil, err
	}

	km.setKey(path, key)

	return key, nil
}

func (km *KeyManager) GetPurposeKey(purpose uint32) (*bip32.Key, error) {
	path := fmt.Sprintf(`m/%d'`, purpose-Apostrophe)

	key, ok := km.getKey(path)
	if ok {
		return key, nil
	}

	parent, err := km.GetMasterKey()
	if err != nil {
		return nil, err
	}

	key, err = parent.NewChildKey(purpose)
	if err != nil {
		return nil, err
	}

	km.setKey(path, key)

	return key, nil
}

func (km *KeyManager) GetCoinTypeKey(purpose, coinType uint32) (*bip32.Key, error) {
	path := fmt.Sprintf(`m/%d'/%d'`, purpose-Apostrophe, coinType-Apostrophe)

	key, ok := km.getKey(path)
	if ok {
		return key, nil
	}

	parent, err := km.GetPurposeKey(purpose)
	if err != nil {
		return nil, err
	}

	key, err = parent.NewChildKey(coinType)
	if err != nil {
		return nil, err
	}

	km.setKey(path, key)

	return key, nil
}

func (km *KeyManager) GetAccountKey(purpose, coinType, account uint32) (*bip32.Key, error) {
	path := fmt.Sprintf(`m/%d'/%d'/%d'`, purpose-Apostrophe, coinType-Apostrophe, account)

	key, ok := km.getKey(path)
	if ok {
		return key, nil
	}

	parent, err := km.GetCoinTypeKey(purpose, coinType)
	if err != nil {
		return nil, err
	}

	key, err = parent.NewChildKey(account + Apostrophe)
	if err != nil {
		return nil, err
	}

	km.setKey(path, key)

	return key, nil
}

// GetChangeKey ...
// https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki#change
// change constant 0 is used for external chain
// change constant 1 is used for internal chain (also known as change addresses)
func (km *KeyManager) GetChangeKey(purpose, coinType, account, change uint32) (*bip32.Key, error) {
	path := fmt.Sprintf(`m/%d'/%d'/%d'/%d`, purpose-Apostrophe, coinType-Apostrophe, account, change)

	key, ok := km.getKey(path)
	if ok {
		return key, nil
	}

	parent, err := km.GetAccountKey(purpose, coinType, account)
	if err != nil {
		return nil, err
	}

	key, err = parent.NewChildKey(change)
	if err != nil {
		return nil, err
	}

	km.setKey(path, key)

	return key, nil
}

func (km *KeyManager) GetKey(purpose, coinType, account, change, index uint32) (*Key, error) {
	path := fmt.Sprintf(`m/%d'/%d'/%d'/%d/%d`, purpose-Apostrophe, coinType-Apostrophe, account, change, index)

	key, ok := km.getKey(path)
	if ok {
		return &Key{path: path, bip32Key: key}, nil
	}

	parent, err := km.GetChangeKey(purpose, coinType, account, change)
	if err != nil {
		return nil, err
	}

	key, err = parent.NewChildKey(index)
	if err != nil {
		return nil, err
	}

	km.setKey(path, key)

	return &Key{path: path, bip32Key: key}, nil
}

func GenerateFromBytes(prvKey *btcec.PrivateKey, compress bool) (wif, address, segwitBech32, segwitNested, taproot string, err error) {
	// generate the wif(wallet import format) string
	btcwif, err := btcutil.NewWIF(prvKey, &chaincfg.MainNetParams, compress)
	if err != nil {
		return "", "", "", "", "", err
	}
	wif = btcwif.String()

	// generate a normal p2pkh address
	serializedPubKey := btcwif.SerializePubKey()
	addressPubKey, err := btcutil.NewAddressPubKey(serializedPubKey, &chaincfg.MainNetParams)
	if err != nil {
		return "", "", "", "", "", err
	}
	address = addressPubKey.EncodeAddress()

	// generate a normal p2wkh address from the pubkey hash
	witnessProg := btcutil.Hash160(serializedPubKey)
	addressWitnessPubKeyHash, err := btcutil.NewAddressWitnessPubKeyHash(witnessProg, &chaincfg.MainNetParams)
	if err != nil {
		return "", "", "", "", "", err
	}
	segwitBech32 = addressWitnessPubKeyHash.EncodeAddress()

	// generate an address which is
	// backwards compatible to Bitcoin nodes running 0.6.0 onwards, but
	// allows us to take advantage of segwit's scripting improvments,
	// and malleability fixes.
	serializedScript, err := txscript.PayToAddrScript(addressWitnessPubKeyHash)
	if err != nil {
		return "", "", "", "", "", err
	}
	addressScriptHash, err := btcutil.NewAddressScriptHash(serializedScript, &chaincfg.MainNetParams)
	if err != nil {
		return "", "", "", "", "", err
	}
	segwitNested = addressScriptHash.EncodeAddress()

	//taproot
	tapKey := txscript.ComputeTaprootKeyNoScript(prvKey.PubKey())
	addressTaproot, err := btcutil.NewAddressTaproot(schnorr.SerializePubKey(tapKey), &chaincfg.MainNetParams)
	if err != nil {
		return "", "", "", "", "", err
	}
	taproot = addressTaproot.EncodeAddress()

	return wif, address, segwitBech32, segwitNested, taproot, nil
}

// encodeEthereum encodes the private key and address for Ethereum.
func encodeEthereum(privateKeyBytes []byte) (privateKey, address string) {
	_, pubKey := btcec.PrivKeyFromBytes(privateKeyBytes)

	publicKey := pubKey.ToECDSA()
	publicKeyBytes := append(publicKey.X.Bytes(), publicKey.Y.Bytes()...)

	// Ethereum uses the last 20 bytes of the keccak256 hash of the public key
	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKeyBytes)
	addr := hash.Sum(nil)
	addr = addr[len(addr)-20:]

	return hex.EncodeToString(privateKeyBytes), eip55checksum(fmt.Sprintf("0x%x", addr))
}

// eip55checksum implements the EIP55 checksum address encoding.
// this function is copied from the go-ethereum library: go-ethereum/common/types.go checksumHex method
func eip55checksum(address string) string {
	buf := []byte(address)
	sha := sha3.NewLegacyKeccak256()
	sha.Write(buf[2:])
	hash := sha.Sum(nil)
	for i := 2; i < len(buf); i++ {
		hashByte := hash[(i-2)/2]
		if i%2 == 0 {
			hashByte = hashByte >> 4
		} else {
			hashByte &= 0xf
		}
		if buf[i] > '9' && hashByte > 7 {
			buf[i] -= 32
		}
	}
	return string(buf[:])
}

func mnemonic(bitSize int) string {
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

const (
	BTC = iota
	ETH
)

type Coins struct {
	gorm.Model
	Address    string
	Mnemonic   string
	Passphrase string
	RootKey    string
	PrivateKey string
	Wif        string
	Typ        int
}

func AddressGenBTC(bitSize int, passphrase, mnemonic string, compress bool) {
	km, err := NewKeyManager(bitSize, passphrase, mnemonic)
	if err != nil {
		log.Fatal(err)
	}

	masterKey, err := km.GetMasterKey()
	if err != nil {
		log.Fatal(err)
	}

	rk := masterKey.B58Serialize()

	coins := make([]*Coins, 0)
	for i := 0; i < 10; i++ {
		key, err := km.GetKey(PurposeBIP44, CoinTypeBTC, 0, 0, uint32(i))
		if err != nil {
			log.Fatal(err)
		}

		wif, address, _, _, _, err := key.Encode(compress)
		if err != nil {
			log.Fatal(err)
		}

		if _, ok := kv[address]; !ok {
			continue
		}

		coins = append(coins, &Coins{
			RootKey:    rk,
			Wif:        wif,
			Typ:        BTC,
			Address:    address,
			Mnemonic:   mnemonic,
			Passphrase: passphrase,
		})
	}
	for i := 0; i < 10; i++ {
		key, err := km.GetKey(PurposeBIP49, CoinTypeBTC, 0, 0, uint32(i))
		if err != nil {
			log.Fatal(err)
		}

		wif, _, _, segwitNested, _, err := key.Encode(compress)
		if err != nil {
			log.Fatal(err)
		}

		if _, ok := kv[segwitNested]; !ok {
			continue
		}

		coins = append(coins, &Coins{
			RootKey:    rk,
			Wif:        wif,
			Typ:        BTC,
			Address:    segwitNested,
			Mnemonic:   mnemonic,
			Passphrase: passphrase,
		})
	}
	for i := 0; i < 10; i++ {
		key, err := km.GetKey(PurposeBIP84, CoinTypeBTC, 0, 0, uint32(i))
		if err != nil {
			log.Fatal(err)
		}

		wif, _, segwitBech32, _, _, err := key.Encode(compress)
		if err != nil {
			log.Fatal(err)
		}

		if _, ok := kv[segwitBech32]; !ok {
			continue
		}

		coins = append(coins, &Coins{
			RootKey:    rk,
			Wif:        wif,
			Typ:        BTC,
			Address:    segwitBech32,
			Mnemonic:   mnemonic,
			Passphrase: passphrase,
		})
	}
	for i := 0; i < 10; i++ {
		key, err := km.GetKey(PurposeBIP86, CoinTypeBTC, 0, 0, uint32(i))
		if err != nil {
			log.Fatal(err)
		}

		wif, _, _, _, taproot, err := key.Encode(compress)
		if err != nil {
			log.Fatal(err)
		}

		if _, ok := kv[taproot]; !ok {
			continue
		}

		coins = append(coins, &Coins{
			RootKey:    rk,
			Wif:        wif,
			Typ:        BTC,
			Address:    taproot,
			Mnemonic:   mnemonic,
			Passphrase: passphrase,
		})
	}

	if len(coins) != 0 {
		if err = db.Create(coins).Error; err != nil {

		}
	}
}

func AddressGenETH(bitSize int, mnemonic, passphrase string) {
	km, err := NewKeyManager(bitSize, passphrase, mnemonic)
	if err != nil {
		log.Fatal(err)
	}

	masterKey, err := km.GetMasterKey()
	if err != nil {
		log.Fatal(err)
	}

	rk := masterKey.B58Serialize()

	coins := make([]*Coins, 0)
	for i := 0; i < 10; i++ {
		key, err := km.GetKey(PurposeBIP44, CoinTypeETH, 0, 0, uint32(i))
		if err != nil {
			log.Fatal(err)
		}

		privateKey, address := encodeEthereum(key.bip32Key.Key)
		if _, ok := kv[address]; !ok {
			continue
		}

		coins = append(coins, &Coins{
			Typ:        ETH,
			RootKey:    rk,
			Address:    address,
			Mnemonic:   mnemonic,
			PrivateKey: privateKey,
			Passphrase: passphrase,
		})
	}

	if len(coins) != 0 {
		if err = db.Create(coins).Error; err != nil {

		}
	}
}

var (
	db *gorm.DB
	kv = map[string]struct{}{}
)

func init() {
	dsn := "user:password@tcp(0.0.0.0:3306)/task?charset=utf8mb4&parseTime=True&loc=Local"

	src, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db = src

	db.AutoMigrate(&Coins{})

	ReadFiles()
}

func ReadFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		kv[strings.TrimSpace(scanner.Text())] = struct{}{}
	}
}

func ReadFiles() {
	files := []string{
		"addr/Bitcoin/2023/04/p2pkh_Rich_Max_1.txt",
		"addr/Bitcoin/2023/04/p2pkh_Rich_Max_10.txt",
		"addr/Bitcoin/2023/04/p2pkh_Rich_Max_100.txt",
		"addr/Bitcoin/2023/04/p2pkh_Rich_Max_1000.txt",
		"addr/Bitcoin/2023/04/p2pkh_Rich_Max_10000.txt",
		"addr/Bitcoin/2023/04/p2pkh_Rich_Max_100000.txt",
		"addr/Bitcoin/2023/04/RichAddr_Max-1.txt",
		"addr/Bitcoin/2023/04/RichAddr_Max-10.txt",
		"addr/Bitcoin/2023/04/RichAddr_Max-100.txt",
		"addr/Bitcoin/2023/04/RichAddr_Max-1000.txt",
		"addr/Bitcoin/2023/04/RichAddr_Max-10000.txt",
		"addr/Bitcoin/2023/04/RichAddr_Max_100000.txt",
		"addr/Bitcoin/rich_100.txt",
		"addr/ETHEREUM/EthRich.txt",
		"addr/ETHEREUM/rich_01.txt",
	}
	for _, name := range files {
		ReadFile(name)
	}
}

func task() {
	ticker := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-ticker.C:
			m128, m256 := mnemonic(128), mnemonic(256)

			AddressGenBTC(128, "<none>", m128, false)
			AddressGenBTC(128, "<none>", m128, true)
			AddressGenBTC(256, "<none>", m256, false)
			AddressGenBTC(256, "<none>", m256, true)
			AddressGenETH(128, "<none>", m128)
			AddressGenETH(256, "<none>", m256)
		}
	}
}

func main() {
	task()
}
