FROM golang:alpine AS builder
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
RUN apk update 
RUN apk add --no-cache bash
WORKDIR /build
COPY . .
RUN ls
#COPY go.mod main.go ./
RUN go mod download
RUN go build -o main .

WORKDIR /dist

RUN cp /build/main .
RUN ls
RUN cp /build/certs/webhook-server-tls.crt .
RUN cp /build/certs/webhook-server-tls.key .

FROM scratch
COPY --from=builder /dist/main .
COPY --from=builder /dist/webhook-server-tls.crt .
COPY --from=builder /dist/webhook-server-tls.key .

ENTRYPOINT ["/main"]
