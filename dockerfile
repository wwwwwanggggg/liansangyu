# 使用官方 Go 镜像作为基础镜像
FROM golang:1.24 AS builder

# 设置 Go 模块代理（注意等号后不能有空格！）
ENV GOPROXY=https://goproxy.io,direct \
    GO111MODULE=on

# 设置工作目录
WORKDIR /app

# 将当前目录的所有文件复制到容器中
COPY . .

# 下载依赖并构建应用程序
RUN go mod tidy && go build -o main .

# 使用更小的基础镜像运行应用程序
FROM debian:bullseye-slim

# 设置工作目录
WORKDIR /app

# 安装必要的依赖项（如 MySQL 客户端）
RUN apt-get update && apt-get install -y mysql-client && rm -rf /var/lib/apt/lists/*

# 复制构建的二进制文件和环境变量文件
COPY --from=builder /app/main .
COPY .env .

# 暴露应用程序的端口（根据你的应用程序监听的端口修改）
EXPOSE 8080

# 设置环境变量（可选）
ENV ENV_FILE_PATH=.env

# 启动应用程序
CMD ["./main"]