FROM golang:1.22-alpine AS builder

WORKDIR /build

COPY . .

RUN go build -o /build/hyve .

FROM scratch

COPY --from=builder /build/hyve /usr/bin/hyve

ENV DATA_DIR=/var/lib/hyve
ENV USER_DIR=${DATA_DIR}/user
ENV STORE=file
ENV STORE_DIR=${DATA_DIR}/certificates

EXPOSE 80
EXPOSE 443

ENTRYPOINT ["/usr/bin/hyve"]
