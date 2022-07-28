# Build go
FROM golang:latest AS builder
WORKDIR /app
COPY . .
ENV CGO_ENABLED=0
RUN go mod download && \
    go env -w GOFLAGS=-buildvcs=false && \
    go build -v -o AikoR -trimpath -ldflags "-s -w -buildid=" ./AikoR

# Release
FROM alpine:latest 
RUN apk --update --no-cache add tzdata ca-certificates && \
    cp /usr/share/zoneinfo/Asia/Ho_Chi_Minh /etc/localtime && \
    mkdir /etc/AikoR/
COPY --from=builder /app/AikoR /usr/local/bin

ENTRYPOINT [ "AikoR", "--config", "/etc/AikoR/aiko.yml"]