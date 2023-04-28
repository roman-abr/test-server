FROM golang:alpine AS builder

WORKDIR /build

COPY . .

RUN go build

FROM alpine

WORKDIR /build

COPY --from=builder /build/test-server /build/test-server

CMD ["./test-server"]