English | [中文](./README.zh.md)
### btc-stealer
- Mnemonic collision
- Utilizes idle machines without spending money
### Support
| Coins | Online |
|-------|:------:|
| BTC   |   ✓    |
| ETH   |   ✓    |
### Build
- Golang 1.21
- RUN
  - `git clone https://github.com/0x2d3c/btc-stealer.git`
  - `cd btc-stealer`
  - `go mod tidy`
  - `CGO_ENABLED=0 go build -ldflags "-s -w" .`
  - `chmod +x btc-stealer`
  - `./btc-stealer`
- Shell
  - `sh build.sh`
### Docker
- Linux arch AMD Build
  - `CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" .`
- Linux arch ARM Build
  - `CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" .`
- Build Image 
- `sudo docker build -f Dockerfile -t btc-stealer:v1.0 .`
### Execution
- Binary
    - `./btc-stealer`
- Docker
    - `sudo docker run -d btc-stealer:v1.0`
### Config
```markdown
{
  "words_list": "english", # mnemonic language
  "eth_gw": "xxxxxxxx" # eth gateway
  "btc_gw": "xxxxxxxx" # btc gateway
}
```
### Description
- If a wallet address with a balance is matched, the result will be written to the balance.txt file.
### Thanks
- [hdkeygen](https://github.com/modood/hdkeygen)
### TODO
- Looking forward to good suggestions or PRs