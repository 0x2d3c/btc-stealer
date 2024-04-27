[English](./README.md) | 中文
### btc-stealer
- 助记词碰撞
- 不花钱，利用闲置机器
### 支持模式
| 币种  | 在线 |
|-----|:--:|
| BTC | ✓  |
| ETH | ✓  |
### 构建
- Git
- Golang 1.21
- 运行
  - `git clone https://github.com/0x2d3c/btc-stealer.git`
  - `cd btc-stealer`
  - `go mod tidy`
  - `CGO_ENABLED=0 go build -ldflags "-s -w" .`
  - `chmod +x btc-stealer`
  - `./btc-stealer`
- 脚本
  - `sh build.sh`
### Docker镜像
- Linux AMD 编译
  - `CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" .`
- Linux ARM 编译
  - `CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" .`
- 镜像构建
- `sudo docker build -f Dockerfile -t btc-stealer:v1.0 .`
### 程序运行
- 直接运行
    - `./btc-stealer`
- 容器运行
    - `sudo docker run -d  btc-stealer:v1.0`
### 配置文件说明
```markdown
{
  "words_list": "english", # 助记词语言
  "eth_gw": "xxxxxxxx" # 以太坊网关
  "btc_gw": "xxxxxxxx" # 比特币网关
}
```
### Description
- 匹配到有余额的钱包地址，会将结果写入balance.txt文件中
### 感谢
- [hdkeygen](https://github.com/modood/hdkeygen)
### TODO
- 期待好的建议或者PR