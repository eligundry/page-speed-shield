FROM golang:1.16-alpine AS builder

RUN apk update \
    && apk add --no-cache git ca-certificates

# Install dependencies
WORKDIR /src
ADD ./api/go.* /src/
RUN go mod download

# Build the app
ADD ./api/main.go /src/main.go
RUN CGO_ENABLED=0 go build -installsuffix 'static' -o /bin/page-speed-shield main.go

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /bin/page-speed-shield /bin/page-speed-shield 
ENTRYPOINT ["/bin/page-speed-shield"]
