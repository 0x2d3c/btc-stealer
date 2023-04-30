CGO_ENABLED=0  GOOS=linux  GOARCH=amd64  go build -ldflags "-s -w" .

upx --best --lzma btc-stealer