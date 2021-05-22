CMD := registry
IMAGE := raptor/registry

build:
	go build -o $(CMD) ./cmd/main.go
	
clean:
	rm $(CMD)
	
deploy:
	kubectl create namespace raptor
	kubectl apply -f ./deployments

docker:
	docker build --tag $(IMAGE) -f ./build/Dockerfile .

lint: 
	golangci-lint run ./...
	
run:
	make build
	./$(CMD)

ut:
	go test -v -count=1 -race -gcflags=-l -timeout=30s ./...

.PHONY: build clean deploy docker lint run ut
