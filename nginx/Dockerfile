FROM nginx:alpine

WORKDIR /etc/nginx
COPY default.conf conf.d

# webgl主页
COPY index.html html/

# webgl静态文件
WORKDIR /webgl
COPY webgl .