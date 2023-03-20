FROM golang:1.18-alpine AS build
WORKDIR /go/src/app
COPY . .
ENV GOPROXY=https://goproxy.cn,direct

# 将清华的镜像源添加到 apk 的仓库列表中
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories

# 安装 gcc 编译器
RUN apk update && apk add --no-cache \
    gcc musl-dev openssl
ENV GOPROXY=https://goproxy.cn,direct
# RUN go build -o server .
RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -installsuffix cgo -o server .

FROM harbor.peopledata.org.cn/htsc/public-cncp-image-base-rhel:8.6
EXPOSE 8081
WORKDIR /app
COPY --from=build /go/src/app/server /app
COPY --from=build /go/src/app/ui/dist /app/ui/dist
COPY --from=build /go/src/app/start.sh /
RUN chmod +x /start.sh
# Entrypoint":["tini","--"]
CMD [ "/start.sh" ]