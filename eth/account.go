package eth

import (
	"btc-stealer/common"
	"btc-stealer/data"
	"encoding/hex"
	"fmt"
	"log"
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

func addressGenETH(bitSize, size int) (map[string]*data.Eth, []string) {
	address := make([]string, size)
	coins := make(map[string]*data.Eth, size)
	for i := range address {
		coin := addressGenETHMaster(bitSize, common.Mnemonic(bitSize), "")

		address[i] = coin.Address
		coins[coin.Address] = coin
	}
	return coins, address
}

func AddressETHCheck() {
	for _, bit := range common.Bits {
		addressGenBitETH(bit)
	}
}

func addressGenBitETH(bitSize int) {
	size := 20

	coins, address := addressGenETH(bitSize, size)
	for _, addr := range address {
		coin, ok := coins[addr]
		if !ok {
			return
		}
		for {
			res := common.HttpETH(addr)

			var done bool
			switch {
			case res == "":
				done = true
			case strings.HasPrefix(res, "0x"):
				if res != "0x0" {
					common.RecordBalance(coin.String())
				}
				done = true
			}

			if done {
				break
			}
		}
	}
}

func addressGenETHMaster(bitSize int, mnemonic, passphrase string) *data.Eth {
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

func addressGenETHMasterAndSub(bitSize int, mnemonic, passphrase string) map[string]*data.Eth {
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
