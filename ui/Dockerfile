# 构建 React 应用程序的镜像
FROM node:16-alpine AS build

# 设置工作目录并拷贝应用程序代码到镜像中
WORKDIR /app
COPY . .

# 安装依赖并编译应用程序
RUN npm install --registry=https://registry.npm.taobao.org && \
    npm run build


# 构建 Nginx 镜像
FROM harbor.peopledata.org.cn/htsc/public-cncp-image-base-nginx:1.0

# 将 React 应用程序的静态资源复制到 Nginx 的 HTML 根目录下
COPY --from=build /app/dist /usr/share/nginx/html

# 复制 Nginx 配置文件到镜像中
COPY nginx.conf /etc/nginx/nginx.conf

# 暴露 Nginx 的 HTTP 端口
EXPOSE 80

# 启动 Nginx 服务器
CMD ["nginx", "-g", "daemon off;"]
