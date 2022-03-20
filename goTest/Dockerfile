FROM golang:alpine AS builder
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
RUN apk update 
RUN apk add --no-cache bash
WORKDIR /build
COPY . .

#COPY go.mod main.go ./
RUN go mod download
RUN go build -o main cmd/mutating-webhook/main.go

WORKDIR /dist
RUN cp /build/main .
RUN cp /build/cert/server.crt .
RUN cp /build/cert/server.key .

FROM scratch
COPY --from=builder /dist/main .
COPY --from=builder /dist/server.crt .
COPY --from=builder /dist/server.key .

ENTRYPOINT ["/main"]
