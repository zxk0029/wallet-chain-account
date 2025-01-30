FROM golang:1.21.1-alpine3.18 as builder

RUN apk add --no-cache make ca-certificates gcc musl-dev linux-headers git jq bash

COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum

WORKDIR /app

RUN go mod download

ARG CONFIG=config.yml

# build wallet-chain-account with the shared go.mod & go.sum files
COPY . /app/wallet-chain-account

WORKDIR /app/wallet-chain-account

RUN make

FROM alpine:3.18

COPY --from=builder /app/wallet-chain-account/wallet-chain-account /usr/local/bin
COPY --from=builder /app/wallet-chain-account/${CONFIG} /etc/wallet-chain-account/

WORKDIR /app

ENTRYPOINT ["wallet-chain-account"]
CMD ["-c", "/etc/wallet-chain-account/config.yml"]
