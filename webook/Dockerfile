# 基础镜像
FROM docker.io/library/ubuntu:20.04
# 把编译后的打包进来这个镜像，放到工作目录 /app。你随便换
COPY webook /app/webook

# 设定工作目录
WORKDIR /app
# CMD 是执行命令


# 最佳 执行命令
ENTRYPOINT ["/app/webook"]


# 结合使用CMD和ENTRYPOINT的一个常见模式是，ENTRYPOINT指定一个脚本或二进制文件作为容器的入口点，而CMD则提供了这个入口点的参数。这样，你可以在运行容器时替换CMD指定的参数，而不需要更改ENTRYPOINT指定的入口点。