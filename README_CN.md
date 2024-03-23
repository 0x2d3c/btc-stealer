[English](./README.md) | 中文
### btc-stealer
- 一个简单的BTC/ETH助记词碰撞示例 & 不花钱的极低概率抽奖程序
- 若匹配密码库效果更佳
### 程序构建
- MySQL 8.0
- Golang 1.21
- 依赖拉取
    - `go mod tidy`
    - `git submodule update --init`
- 构建命令
    - `CGO_ENABLED=0  GOOS=linux  GOARCH=amd64  go build -ldflags "-s -w" .`
### 镜像构建
- `sudo docker build -f Dockerfile -t btc-stealer:v1.0 `
### 程序运行
- 直接运行
    - `btc-stealer`
- 容器运行
    - `sudo docker run -d  btc-stealer:v1.0`
### 感谢
- [hdkeygen](https://github.com/modood/hdkeygen)
### TODO
- 期待好的建议或者PR