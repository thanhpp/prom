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

WORKDIR /build/cmd/ccman

RUN go build -v -o ccman .

# RUN

FROM alpine

COPY --from=builder ["/build/cmd/ccman/ccman", "ccman"]
COPY --from=builder ["/build/cmd/ccman/dev.yml", "dev.yml"]

EXPOSE 8080

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.9.0/wait /wait
RUN chmod +x /wait

CMD /wait && /ccman