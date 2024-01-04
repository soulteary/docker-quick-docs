FROM golang:1.21.5-alpine3.19 as builder
RUN echo '' > /etc/apk/repositories && \
    echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.19/main"         >> /etc/apk/repositories && \
    echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.19/community"    >> /etc/apk/repositories && \
    echo "Asia/Shanghai" > /etc/timezone
RUN apk add upx
WORKDIR /build
ENV CGO_ENABLED=0
COPY . .
RUN go build -ldflags "-w -s"
RUN upx -9 -o docker-quick-docs.minify docker-quick-docs && mv docker-quick-docs.minify docker-quick-docs

FROM alpine:3.19
RUN echo '' > /etc/apk/repositories && \
    echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.19/main"         >> /etc/apk/repositories && \
    echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.19/community"    >> /etc/apk/repositories && \
    echo "Asia/Shanghai" > /etc/timezone
RUN apk add openssl && rm -rf /var/cache/apk/*
WORKDIR /app
COPY --from=Builder /build/docker-quick-docs /bin/docker-quick-docs
ENTRYPOINT ["quick-docs"]