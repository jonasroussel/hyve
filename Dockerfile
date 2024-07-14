FROM golang:1.22-alpine AS builder

WORKDIR /build

COPY . .

RUN go build -o /build/proxbee .

FROM alpine:latest

COPY --from=builder /build/proxbee /bin/proxbee

EXPOSE 80
EXPOSE 443

ENTRYPOINT ["/bin/proxbee"]
