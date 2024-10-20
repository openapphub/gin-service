# 使用官方 Go 镜像作为基础镜像
FROM golang:1.20-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制 go mod 和 sum 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o openapphub ./cmd/api

# 使用轻量级的 alpine 镜像
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 从 builder 阶段复制编译好的应用
COPY --from=builder /app/openapphub .
COPY --from=builder /app/.env .

# 暴露端口
EXPOSE 3000

# 运行应用
CMD ["./openapphub"]