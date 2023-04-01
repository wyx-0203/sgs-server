up:
	docker-compose up -d --build

down:
	docker-compose down

build:
	env GOOS=linux CGO_ENABLED=0 go build -o bin/sgs-server

webgl:
	echo y | rm -r nginx/webgl/*
	cp -r ../sgs-unity/Build/WebGL/* nginx/webgl