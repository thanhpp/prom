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

WORKDIR /build/cmd/usrman

RUN go build -v -o usrman .

# RUN

FROM alpine

COPY --from=builder ["/build/cmd/usrman/usrman", "usrman"]
COPY --from=builder ["/build/cmd/usrman/dev.yml", "dev.yml"]

EXPOSE 8090

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.9.0/wait /wait
RUN chmod +x /wait

CMD /wait && /usrman