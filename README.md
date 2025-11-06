# Docker Quick Docs

[![Release](https://github.com/soulteary/docker-quick-docs/actions/workflows/release.yaml/badge.svg)](https://github.com/soulteary/docker-quick-docs/actions/workflows/release.yaml) [![Codecov](https://github.com/soulteary/docker-quick-docs/actions/workflows/codecov.yml/badge.svg)](https://github.com/soulteary/docker-quick-docs/actions/workflows/codecov.yml) [![CodeQL](https://github.com/soulteary/docker-quick-docs/actions/workflows/codeql.yml/badge.svg)](https://github.com/soulteary/docker-quick-docs/actions/workflows/codeql.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/soulteary/docker-quick-docs)](https://goreportcard.com/report/github.com/soulteary/docker-quick-docs)

本地部署、能够快速访问的文档工具，用来改善 GitHub Pages 文档访问体验。The project offers a lightweight way to host documentation locally and browse it at high speed without depending on GitHub Pages.

## Overview | 项目简介

- **EN**: Docker Quick Docs serves static documentation bundles from your local machine, removing latency introduced by public hosting and giving you full control of the browsing experience.
- **ZH**: Docker Quick Docs 可以在本地机器中快速部署静态文档，绕过公网访问的延迟，让你完全掌控浏览体验。

## Download | 获取方式

- **EN**: Download the executable that matches your operating system from the GitHub [Releases page](https://github.com/soulteary/docker-quick-docs/releases).
- **ZH**: 前往 GitHub [发布页面](https://github.com/soulteary/docker-quick-docs/releases) 下载与你操作系统匹配的可执行文件。

![](.github/assets.jpg)

- **EN**: Prefer Docker? Pull the ready-to-run image from Docker Hub.
- **ZH**: 如果习惯使用 Docker，可以直接从 Docker Hub 拉取镜像。

![](.github/dockerhub.jpg)

```bash
docker pull soulteary/docker-quick-docs:v0.1.7
# EN: or pull the latest tag automatically
# ZH: 或者直接拉取最新的镜像标签
docker pull soulteary/docker-quick-docs
```

## Usage | 使用示例

- **EN**: The following walkthrough localizes the documentation of [baidu/san](http://github.com/baidu/san) by cloning its GitHub Pages branch.
- **ZH**: 下面演示如何克隆 [baidu/san](http://github.com/baidu/san) 的 GitHub Pages 分支，将文档本地化部署。

```bash
git clone http://github.com/baidu/san --depth 1 --branch=gh-pages
Cloning into 'san'...
warning: redirecting to https://github.com/baidu/san/
remote: Enumerating objects: 405, done.
remote: Counting objects: 100% (405/405), done.
remote: Compressing objects: 100% (197/197), done.
remote: Total 405 (delta 154), reused 303 (delta 65), pack-reused 0
Receiving objects: 100% (405/405), 2.17 MiB | 5.18 MiB/s, done.
Resolving deltas: 100% (154/154), done.
```

- **EN**: Move the cloned site into the `docs` directory so the server can pick it up.
- **ZH**: 将克隆好的站点移动到 `docs` 目录中，便于程序读取。

```bash
mv san docs/
```

- **EN**: Start the server with either the native binary or the Docker container.
- **ZH**: 使用原生二进制或 Docker 容器启动服务。

```bash
# EN: native binary
# ZH: 原生二进制
./quick-docs

# EN: Docker container
# ZH: Docker 容器
docker run --rm -it -v `pwd`/docs:/app/docs -p 8080:8080 soulteary/docker-quick-docs:v0.1.7
```

- **EN**: When the application starts you should see logs similar to the following.
- **ZH**: 程序启动后，会出现类似如下的日志输出。

```bash
2024/01/04 10:38:31 Quick Docs
2024/01/04 10:38:31 未设置环境变量 `PORT`，使用默认端口：8080
```

- **EN**: Open your browser at `http://localhost:8080` to browse the imported documentation immediately.
- **ZH**: 打开浏览器访问 `http://localhost:8080`，即可快速浏览刚导入的文档。

## Configuration | 配置

- **EN**: Customize the listening port by defining the `PORT` environment variable before starting the program.
- **ZH**: 如果需要修改监听端口，在启动程序前设置 `PORT` 环境变量即可。

```bash
PORT=8080 ./quick-docs
# EN: or with Docker
# ZH: 或者使用 Docker
docker run --rm -it -e PORT=8080 -v `pwd`/docs:/app/docs -p 8080:8080 soulteary/docker-quick-docs:v0.1.7
```
