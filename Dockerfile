# syntax=docker/dockerfile:1
FROM golang:1.18-alpine AS build
WORKDIR /go/src/github.com/org/repo
COPY . .
ENV GOPROXY=https://goproxy.cn,direct

# 将清华的镜像源添加到 apk 的仓库列表中
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories

# 安装 gcc 编译器
RUN apk update && apk add --no-cache \
    gcc musl-dev openssl
#    && rm -rf /var/cache/apk/* \

RUN go build -o server .

FROM build AS development
RUN ls -la
ENV GOPROXY=https://goproxy.cn,direct
CMD ["go", "run", "main.go"]

FROM repo-dev.htsc/public-cncp-image-base-local/rhel:8.6
EXPOSE 8081
COPY --from=build /go/src/github.com/org/repo/server /server
COPY --from=build /go/src/github.com/org/repo/config.yaml /config.yaml
CMD ["/server"]
