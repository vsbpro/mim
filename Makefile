all: docker_image build_mim

docker_image: build_mim

build_mim:
	go build -o ./mim.exe cmd/main.go

docker_image: build_mim
	sudo docker build -t mim:latest .