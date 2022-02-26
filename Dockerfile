FROM golang:1.14.6-alpine3.12 as builder
COPY go.mod go.sum /go/src/gitlab.com/idoko/dpricing/
WORKDIR /go/src/gitlab.com/idoko/dpricing
RUN go mod download
COPY . /go/src/gitlab.com/idoko/dpricing
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/dpricing gitlab.com/idoko/dpricing

FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/gitlab.com/idoko/dpricing/build/dpricing /usr/bin/dpricing
EXPOSE 8080 8080
ENTRYPOINT ["/usr/bin/dpricing"]