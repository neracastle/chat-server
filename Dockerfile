FROM golang:1.20-alpine3.19 as builder
LABEL authors="ivansemeniv"

COPY . /neracastle/chat/src
WORKDIR /neracastle/chat/src

RUN go mod download
RUN go build -o ./bin/chat_server cmd/grpc-server/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /neracastle/chat/src/bin/chat_server .

CMD ["./chat_server"]