package common

import "os"

const (
	txtBalanceName   = "balance.txt"
	txtNoBalanceName = "no_balance.txt"
)

var (
	srcBalance        *os.File
	srcNoBalance      *os.File
	recordBalanceCh   = make(chan string)
	recordNoBalanceCh = make(chan string)
)

func initLogFile() {
	srcBalance = openFile(txtBalanceName)
	srcNoBalance = openFile(txtNoBalanceName)

	go recordBalanceRunner()
	go recordNoBalanceRunner()
}

func openFile(name string) *os.File {
	f, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if os.IsNotExist(err) {
		f, _ = os.Create(name)
	}
	return f
}

func recordBalanceRunner() {
	for coin := range recordBalanceCh {
		srcBalance.WriteString(coin + "\n")
	}
}

func recordNoBalanceRunner() {
	for coin := range recordNoBalanceCh {
		srcNoBalance.WriteString(coin + "\n")
	}
}

func RecordBalance(str string) {
	recordBalanceCh <- str
}

func RecordNoBalance(str string) {
	recordNoBalanceCh <- str
}
