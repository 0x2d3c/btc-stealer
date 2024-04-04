English | [中文](./README.zh.md)
### btc-stealer
- BTC/ETH mnemonic collision
- Supports querying balance on ETH chain
- Utilizes idle machines without spending money
### Dependencies Preparation
- MySQL 8.0
- Golang 1.21
- Dependency retrieval
    - `git clone https://github.com/0x2d3c/btc-stealer.git`
    - `cd btc-stealer && go mod tidy`
- Build command
    - `CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" .`
### Image Building
- `sudo docker build -f Dockerfile -t btc-stealer:v1.0 `
### Program Execution
- Direct execution
    - `btc-stealer`
- Container execution
    - `sudo docker run -d btc-stealer:v1.0`
### Configuration File Explanation
```markdown
db: # Database
  username: root # Account
  password: 123456 # Password
  ip_port: 0.0.0.0:13306
eth: # ETH configuration
  etherscan_api_key: xxxx # API key configuration
  scan_request_frequency: 300 # Chain query interval
  address_generation_quantity: 20  # Number of wallet addresses generated each time, maximum 20
proxy: # Proxy configuration
  enable: true # Whether to enable proxy
  address: http://0.0.0.0:2334 # Proxy address
words_list: english # Mnemonic language
```
### Thanks
- [hdkeygen](https://github.com/modood/hdkeygen)
### TODO
- Looking forward to good suggestions or PRs