package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

var (
	count int64
	kv    = map[string]struct{}{}
)

func ReadFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		addr := strings.TrimSpace(scanner.Text())
		if len(addr) > 34 {
			continue
		}
		kv[addr] = struct{}{}
	}
}

func ReadFiles(filenames []string) {
	for _, name := range filenames {
		ReadFile(name)
	}
}

func addrLoad() {
	files := []string{
		"addr/Bitcoin/2023/04/p2pkh_Rich_Max_1.txt",
		"addr/Bitcoin/2023/04/p2pkh_Rich_Max_10.txt",
		"addr/Bitcoin/2023/04/p2pkh_Rich_Max_100.txt",
		"addr/Bitcoin/2023/04/p2pkh_Rich_Max_1000.txt",
		"addr/Bitcoin/2023/04/p2pkh_Rich_Max_10000.txt",
		"addr/Bitcoin/2023/04/p2pkh_Rich_Max_100000.txt",
	}

	ReadFiles(files)
}

func CountPrint() {
	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		fmt.Printf("碰撞: %9d, 时间: %s\n", count, time.Now().Format("2006-01-02 15:04:05"))
	}
}

func main() {
	addrLoad()

	go CountPrint()

	ctrl := make(chan struct{}, 100)
	for {
		ctrl <- struct{}{}

		go GenAddr(ctrl)
	}
}

func GenAddr(ch <-chan struct{}) {
	entropy, _ := bip39.NewEntropy(256)
	mnemonic, _ := bip39.NewMnemonic(entropy)

	// 这里替换为你自己的密码
	password := ""

	mtk, _ := bip32.NewMasterKey(bip39.NewSeed(mnemonic, password))
	sub, _ := mtk.NewChildKey(hdkeychain.HardenedKeyStart + 44)
	sub, _ = sub.NewChildKey(hdkeychain.HardenedKeyStart + 0)
	sub, _ = sub.NewChildKey(hdkeychain.HardenedKeyStart + 0)
	sub, _ = sub.NewChildKey(0)
	sub, _ = sub.NewChildKey(0)

	address, err := btcutil.NewAddressPubKey(mtk.PublicKey().Key, &chaincfg.MainNetParams)
	if err != nil {
		panic(err)
	}

	if _, ok := kv[address.EncodeAddress()]; ok {
		fmt.Printf("助记词：%s, 地址：%s, 私钥：%s\n", mnemonic, address.EncodeAddress(), mtk.B58Serialize())
	}

	atomic.AddInt64(&count, 1)

	<-ch
}
