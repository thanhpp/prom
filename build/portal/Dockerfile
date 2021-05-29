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

WORKDIR /build/cmd/portal

RUN go build -v -o portal .

# RUN

FROM alpine

COPY --from=builder ["/build/cmd/portal/portal", "portal"]
COPY --from=builder ["/build/cmd/portal/dev.yml", "dev.yml"]

EXPOSE 12345

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.9.0/wait /wait
RUN chmod +x /wait

CMD /wait && /portal