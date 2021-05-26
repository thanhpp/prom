# BUILD
FROM golang:1.16-alpine AS builder

# set ENV variables
ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64

# Build directory
WORKDIR /build

COPY ["./go.mod", "./go.sum", "./"]

RUN go mod download

COPY [".", "."]

WORKDIR /build/cmd/noti

RUN go build -v -o noti .

# RUN

FROM alpine

WORKDIR /run

COPY --from=builder ["/build/cmd/noti/noti", "noti"]
COPY --from=builder ["/build/cmd/noti/dev.yml", "dev.yml"]

EXPOSE 23456

ENV WAIT_VERSION 2.9.0
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/$WAIT_VERSION/wait /wait
RUN chmod +x /wait

CMD ./wait && ./noti