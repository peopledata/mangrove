# syntax=docker/dockerfile:1
FROM golang:1.18-alpine AS build
WORKDIR /src
COPY . .
ENV GOPROXY=https://goproxy.cn,direct

# 将清华的镜像源添加到 apk 的仓库列表中
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories

# 安装 gcc 编译器
RUN apk update && apk add --no-cache \
    gcc musl-dev openssl
#    && rm -rf /var/cache/apk/* \

ENV GOPROXY=https://goproxy.cn,direct
RUN go build -o server .

# FROM build AS development
# RUN ls -la
# ENV GOPROXY=https://goproxy.cn,direct
# CMD ["go", "run", "main.go"]

FROM harbor.peopledata.org.cn/htsc/public-cncp-image-base-rhel:8.6
EXPOSE 8081
WORKDIR /src
COPY --from=build /src/server /src/server
COPY --from=build /src/config.yaml /src/config.yaml
COPY --from=build /src/ui /src/ui
CMD ["/src/server"]
