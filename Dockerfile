#使用golang 1.18作为蓝本
FROM golang:1.18 as build
#设置工作目录
WORKDIR /app/
#拷贝所有文件到工作目录
COPY . .
#下载库
ENV GO111MODULE=on
ENV GOPROXY="https://mirrors.aliyun.com/goproxy/"
#编译
RUN go build -o pod-controller

RUN mkdir -p /app/logs
#如果有配置
#COPY --from=build /go/config.yml /
#运行
CMD ["/app/pod-controller"]


