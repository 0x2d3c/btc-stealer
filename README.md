English | [中文](./README.zh.md)
### btc-stealer
- Description:
   - A simple example demonstrating BTC/ETH mnemonic collision & a no-cost lottery program with extremely low probability.
   - It will be better if it matches the password library
- Program Construction:
   - MySQL Version: 8.0
   - Golang Version: 1.21
   - Dependency Installation:
      - `go mod tidy`
      - `git submodule update --init`
   - Build:
      - `CGO_ENABLED=0  GOOS=linux  GOARCH=amd64  go build -ldflags "-s -w" .`
- Docker:
   - `sudo docker build -f Dockerfile -t btc-stealer:v1.0`
- Run:
   - Direct Execution:
      - `btc-stealer`
   - Container Execution:
      - `sudo docker run -d  btc-stealer:v1.0`
### Thanks:
- [hdkeygen](https://github.com/modood/hdkeygen)
### TODO:
- Looking forward to good suggestions or PR