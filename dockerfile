# Build go
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
ENV CGO_ENABLED=0
RUN go mod download
RUN go build -v -o AikoR -trimpath -ldflags "-s -w -buildid=" ./AikoR

# Release
FROM  alpine
RUN  apk --update --no-cache add tzdata ca-certificates \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN mkdir /etc/AikoR/
COPY --from=builder /app/AikoR /usr/local/bin

ENTRYPOINT [ "AikoR", "--config", "/etc/AikoR/aiko.yml"]