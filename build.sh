git submodule update --init

go mod tidy

CGO_ENABLED=0 go build -ldflags "-s -w" .

chmod +x btc-stealer

./btc-stealer