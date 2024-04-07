package common

import "os"

const (
	txtBalanceName = "balance.txt"
)

var (
	srcBalance      *os.File
	recordBalanceCh = make(chan string)
)

func initLogFile() {
	srcBalance = openFile(txtBalanceName)

	go recordBalanceRunner()
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

func RecordBalance(str string) {
	recordBalanceCh <- str
}
