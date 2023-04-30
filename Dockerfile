FROM alpine:3.17.3 AS builder

COPY addr/Bitcoin/2023/04/p2pkh_Rich_Max_*.txt /root/addr/Bitcoin/2023/04/

COPY btc-stealer /root

WORKDIR /root

CMD ["./btc-stealer"]