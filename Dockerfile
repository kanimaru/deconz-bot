# Dockerfile for Rasperry PI
FROM golang:1.19-alpine as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN apk update \
        && apk upgrade \
        && apk add --no-cache ca-certificates git \
        && update-ca-certificates 2>/dev/null || true
RUN CGO_ENABLED=0 GOARCH=arm64 GOARM=5 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .

FROM scratch
COPY --from=builder /build/main /app/
COPY --from=builder /build/view /app/view
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /app
EXPOSE 9800
CMD ["/app/main"]
