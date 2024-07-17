FROM golang:1.22-alpine AS builder

WORKDIR /build

COPY . .

RUN go build -o /build/proxbee .

FROM alpine:latest

COPY --from=builder /build/proxbee /usr/bin/proxbee

ENV DATA_DIR=/var/lib/proxbee
ENV USER_DIR=${DATA_DIR}/user
ENV STORE=file
ENV STORE_DIR=${DATA_DIR}/certificates

EXPOSE 80
EXPOSE 443

ENTRYPOINT ["/usr/bin/proxbee"]
