#first stage - builder
FROM golang:1.12-stretch as builder
COPY . /App
WORKDIR /App
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -o App

#second stage
FROM alpine:latest
WORKDIR /root/
RUN apk add --no-cache tzdata
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /App .
CMD ["./App"]