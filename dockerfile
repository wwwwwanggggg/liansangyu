# 第一阶段：构建环境（基于 Ubuntu 22.04）
FROM ubuntu:22.04 AS builder

# 安装 Go 工具链和必要依赖
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    golang-1.24 \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# 设置 Go 环境变量
ENV GOPROXY=https://goproxy.io,direct \
    GO111MODULE=on \
    CGO_ENABLED=0

# 复制代码并构建
WORKDIR /app
COPY . .
RUN go mod tidy && \
    go build -tags netgo,osusergo -ldflags="-s -w" -o main .

# 第二阶段：运行环境（基于 Ubuntu 22.04）
FROM ubuntu:22.04

# 安装运行时依赖（如 MySQL 客户端）
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    default-mysql-client \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# 复制二进制文件和配置文件
WORKDIR /app
COPY --from=builder /app/main .
COPY .env .

# 暴露端口并设置启动命令
EXPOSE 8080
ENV ENV_FILE_PATH=.env
CMD ["./main