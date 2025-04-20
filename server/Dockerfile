FROM golang:1.22.12-alpine

WORKDIR /server
COPY . .

# Install curl so we can probe connectivity
RUN apk update --no-cache \
    && apk add --no-cache curl

# 1) Test Google reachability: if it fails in 0.5s, switch Alpine mirror to Aliyun
RUN if curl -s --connect-timeout 0.5 https://www.google.com/ >/dev/null; then \
        echo "✔ Google reachable → keeping default APK mirror"; \
    else \
        echo "⚠ Google unreachable → switching APK mirror to Aliyun"; \
        sed -i 's|dl-cdn.alpinelinux.org|mirrors.aliyun.com|g' /etc/apk/repositories; \
    fi \
    && apk update --no-cache \
    && apk add --no-cache ffmpeg

# 2) Again test Google reachability: if it fails, set GOPROXY before building
RUN if curl -s --connect-timeout 0.5 https://www.google.com/ >/dev/null; then \
        echo "✔ Google reachable → using direct Go module proxy"; \
    else \
        echo "⚠ Google unreachable → setting GOPROXY to goproxy.cn"; \
        go env -w GOPROXY=https://goproxy.cn,direct; \
    fi \
    && go mod tidy \
    && go build -o app cmd/main.go

EXPOSE 9000

CMD ["./app"]
