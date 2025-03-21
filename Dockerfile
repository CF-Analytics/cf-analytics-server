# 1️⃣ 第一阶段：构建 Go 应用
FROM golang:1.23.2 AS builder

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制代码并编译（可更改 GOOS/GOARCH 实现跨平台）
COPY . .

# 使用 CGO 禁用 + 静态编译，提高兼容性
# docker build --build-arg TARGETOS=linux --build-arg TARGETARCH=amd64 -t myapp:latest .
# docker build --build-arg TARGETOS=windows --build-arg TARGETARCH=amd64 -t myapp-windows .
# docker build --build-arg TARGETOS=darwin --build-arg TARGETARCH=arm64 -t myapp-mac .
ARG TARGETOS
ARG TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go build -o app .

# 2️⃣ 第二阶段：创建轻量级运行环境
FROM alpine:latest

# 设置工作目录
WORKDIR /root/

# 复制编译后的二进制文件
COPY --from=builder /app/app .

# 运行应用
CMD ["./app"]
