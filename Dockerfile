# syntax=docker/dockerfile:1
FROM golang:1.18-alpine AS build
WORKDIR /go/src/github.com/org/repo
COPY . .
ENV GOPROXY=https://goproxy.cn,direct
RUN go build -o server .

FROM build AS development
RUN ls -la
ENV GOPROXY=https://goproxy.cn,direct
CMD ["go", "run", "main.go"]

FROM alpine:3.12
EXPOSE 8000
COPY --from=build /go/src/github.com/org/repo/server /server
CMD ["/server"]
