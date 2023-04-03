FROM alpine:latest

# 复制在本地编译生成的二进制文件
WORKDIR /app
COPY bin/sgs-server .

# 将ssl证书复制到/app目录下
COPY nginx/cert .

# 运行
CMD ["./sgs-server"]