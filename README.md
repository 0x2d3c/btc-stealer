English | [中文](./README.zh.md)
### btc-stealer
- Mnemonic collision
- Utilizes idle machines without spending money
### Support
| Coins |    Mix    | Online | Offline  |
|-------|:---------:|:------:|:--------:|
| BTC   |     x     |   x    |    ✓     |
| ETH   |     ✓     |   ✓    |    ✓     |
### Build
- Golang 1.21
- RUN
  - `git clone https://github.com/0x2d3c/btc-stealer.git`
  - `cd btc-stealer`
  - `git submodule update --init`
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
  "mode":2, # running mode, 0:mix 1:online 2:offline
  "wallet":{
    "eth":["wallet_path/xxx.txt"],
    "btc":["wallet_path/xxx.txt"]
  },
  "proxy": {  # proxy configuration
    "enable": true,
    "address": "http://0.0.0.0:2334"
  },
  "words_list": "english", # mnemonic language
  "etherscan_api_key": "xxxxxxxx" # API key configuration, online mode need it
}
```
### Description
- There will be no console output when running this program. If a wallet address with a balance is matched, the result will be written to the balance.txt file.
### Thanks
- [hdkeygen](https://github.com/modood/hdkeygen)
### TODO
- Looking forward to good suggestions or PRs