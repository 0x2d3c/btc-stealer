[English](./README.md) | 中文
### btc-stealer
- BTC/ETH助记词碰撞
- 支持ETH链上余额查询
- 不花钱，利用闲置机器
### 依赖准备
- MySQL 8.0
- Golang 1.21
- 依赖拉取
    - `git clone https://github.com/0x2d3c/btc-stealer.git`
    - `cd btc-stealer && go mod tidy`
- 构建命令
    - `CGO_ENABLED=0  GOOS=linux  GOARCH=amd64  go build -ldflags "-s -w" .`
### 镜像构建
- `sudo docker build -f Dockerfile -t btc-stealer:v1.0 `
### 程序运行
- 直接运行
    - `btc-stealer`
- 容器运行
    - `sudo docker run -d  btc-stealer:v1.0`
### 配置文件说明
```markdown
{
  "proxy": {  # 代理配置
    "enable": true, # 是否开启代理
    "address": "http://0.0.0.0:2334" # 代理地址
  },
  "words_list": "english", # 助记词语言
  "etherscan_api_key": "xxxxxxxx" # api key配置
}
```
### 感谢
- [hdkeygen](https://github.com/modood/hdkeygen)
### TODO
- 期待好的建议或者PR