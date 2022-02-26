FROM golang:1.14.6-alpine3.12 as builder
COPY go.mod go.sum /go/src/github.com/jasanchez1/Dpricing/
WORKDIR /go/src/github.com/jasanchez1/Dpricing
RUN go mod download
COPY . /go/src/github.com/jasanchez1/Dpricing
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/Dpricing github.com/jasanchez1/Dpricing

FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/github.com/jasanchez1/Dpricing/build/Dpricing /usr/bin/Dpricing
EXPOSE 8080 8080
ENTRYPOINT ["/usr/bin/Dpricing"]