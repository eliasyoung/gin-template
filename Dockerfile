# --- Stage 1: Build Stage ---
FROM golang:1.25-alpine AS builder

# 1. 安装构建时必要的工具
# - git: 用于拉取 Go modules 中可能存在的 git 依赖
# - ca-certificates: 用于所有 HTTPS 连接 (go mod download, git clone)
RUN apk add --no-cache git ca-certificates

# 2. 设置工作目录
WORKDIR /app

# 3. 设置 Go 环境变量，优化依赖下载
# - GOPROXY: 使用国内代理，大幅提升下载速度和稳定性
# - CGO_ENABLED=0: 创建静态链接的二进制文件，这是制作最小化生产镜像的关键
ENV GOPROXY=https://goproxy.cn,direct
ENV CGO_ENABLED=0

# 4. 复制并下载依赖
# 将 go.mod 和 go.sum 分开复制，以最大化利用 Docker 的层缓存机制。
# 只要这两个文件不变，下面的 RUN 指令就不会重新执行。
COPY go.mod go.sum ./
RUN go mod download

# 5. 复制所有源代码
COPY . .

# 6. 编译应用
# - ARG: 允许在构建时从外部传入目标操作系统和架构，增强了交叉编译的灵活性
# - -ldflags="-w -s": 裁剪掉调试信息，显著减小最终二进制文件的大小
ARG TARGETOS=linux
ARG TARGETARCH=amd64
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags="-w -s" -o /app/main ./cmd/main/main.go


# --- Stage 2: Final/Production Stage ---
# 使用一个极简的、经过安全扫描的官方镜像作为基础
# debian:stable-slim 是一个比 alpine 兼容性更好、比 ubuntu 小得多的优秀选择
FROM debian:12-slim
# 或者，如果您对体积有极致要求，且确认没有glibc依赖，可以用 FROM alpine:latest

# 1. 安装运行时的必要依赖
# - ca-certificates: 保证您的应用在运行时可以调用外部 HTTPS API (如 X-Node API)
# - tzdata: 提供时区信息，以便应用能正确处理和显示时间
# - gosu: 用于安全执行 entrypoint.sh中的指令
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    tzdata \
    gosu \
    && rm -rf /var/lib/apt/lists/*

# 2. 设置容器时区为上海 (CST)
# 这主要影响日志记录和 time.Now() 的默认行为，是一个很好的实践
ENV TZ=Asia/Shanghai

# 3. 创建一个非 root 用户和组，以遵循最小权限原则，增强安全性
# 该用户的 HOME 为工作目录 /app
RUN addgroup --system appgroup && adduser --system --ingroup appgroup --home /app appuser

# 4. 设置工作目录
WORKDIR /app

# 5. 从 builder 阶段复制已编译好的二进制文件及启动脚本 entrypoint.sh
COPY --from=builder /app/main .
COPY entrypoint.sh /entrypoint.sh

# 6. 用 root 修复权限
RUN chmod +x /entrypoint.sh && \
    chown -R appuser:appgroup /app /entrypoint.sh

# 7. 暴露应用监听的端口 (请根据您的配置修改)
EXPOSE 8080

# 8. 设置容器启动时执行的命令
ENTRYPOINT ["/entrypoint.sh"]
