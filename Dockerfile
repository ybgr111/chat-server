FROM golang:1.21-alpine AS builder

COPY . /github.com/ybgr111/chat-server/source/
WORKDIR /github.com/ybgr111/chat-server/source/

RUN go mod download
RUN go build -o ./bin/chat-server cmd/grpc_server/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/ybgr111/chat-server/source/bin/chat-server .

CMD ["./chat-server"]