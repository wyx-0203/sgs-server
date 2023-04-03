FROM nginx:alpine

WORKDIR /etc/nginx
COPY default.conf conf.d

# ssl证书 本地调试时需注释这一句
COPY cert cert

# webgl主页
COPY index.html html/

# webgl静态文件
WORKDIR /webgl
COPY webgl .