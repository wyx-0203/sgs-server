FROM nginx:alpine

WORKDIR /etc/nginx
COPY default.conf conf.d
COPY index.html html/

# WORKDIR /
# COPY . webgl
WORKDIR /webgl
COPY webgl .
# COPY index.html .