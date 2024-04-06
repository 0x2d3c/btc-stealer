FROM alpine:3.17.3 AS builder

COPY btc-stealer /root

WORKDIR /root

CMD ["./btc-stealer"]