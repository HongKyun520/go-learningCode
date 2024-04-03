# 打包命令
GOOS=linux GOARCH=arm go build -o webook .

# 制作镜像
docker build -t flycash/webook:v0.0.1 .