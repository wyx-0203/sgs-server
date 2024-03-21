up:
	docker-compose up -d --build

down:
	docker-compose down

api_:
	env GOOS=linux CGO_ENABLED=0 go build -o api/bin/sgs-server api/main.go
	minikube image build -t api api/

room:
	minikube image build -t room service-room/

webgl:
	echo y | rm -r nginx/webgl/*
	cp -r ../sgs-unity/Build/WebGL/* nginx/webgl