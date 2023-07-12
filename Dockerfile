# 使用官方的Go基础镜像作为构建环境
FROM golang:1.17 AS builder

# 设置工作目录
WORKDIR /app

# 将Go模块文件复制到工作目录
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 将源代码复制到工作目录
COPY . .

# 构建可执行文件
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 使用官方的Alpine基础镜像作为运行环境
FROM alpine:latest

# 在Alpine中安装ca-certificates以支持HTTPS
RUN apk --no-cache add ca-certificates

# 将可执行文件从构建环境复制到运行环境
COPY --from=builder /app/main /app/main

# 设置工作目录
WORKDIR /app

# 暴露端口
EXPOSE 58180

# 运行程序
CMD ["/app/main"]
