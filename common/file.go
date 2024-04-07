package common

import (
	"bufio"
	"os"
	"strings"
)

// readWalletAddressFile the format of the file content is, one wallet address per line
func readWalletAddressFile(filenames []string) map[string]struct{} {
	kv := make(map[string]struct{})
	for _, filename := range filenames {
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
	return kv
}
