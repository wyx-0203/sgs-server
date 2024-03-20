FROM alpine:latest

RUN mkdir /app
COPY bin/sgs-server /app

CMD ["/app/sgs-server"]

# FROM golang:alpine

# WORKDIR /app

# COPY go.mod ./
# COPY go.sum ./
# RUN go mod download

# COPY * ./

# RUN go build -o /sgs-server

# CMD [ "/sgs-server" ]