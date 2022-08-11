VERSION=$(shell cat ./VERSION)
PROTO_FILES=$(shell find ./api/protos/ -iname '*.proto')
IMG_TAG=linhmtran168/511transit:$(VERSION)

.PHONY: install-deps install-go-deps tidy back-dev back-test back-gen-mock back-gen-proto front-dev front-install front-build

install-deps: install-go-deps front-install

install-go-deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install github.com/cosmtrek/air@latest
	go install github.com/golang/mock/mockgen@v1.6.0	

tidy:
	go mod tidy

back-dev:
	air

back-test:
	go test ./...

back-gen-mock:
	go generate ./...

back-gen-proto:
	protoc --proto_path=api/protos/ --go_out=. --go_opt=module=github.com/linhmtran168/511transit $(PROTO_FILES)

lint:
	golangci-lint run

front-install:
	cd ./web && npm install

front-dev:
	cd ./web && npm run dev

front-build:
	cd ./web && npm run build

build-prod-img:
	docker build -t $(IMG_TAG) .

run-prod-img:
	docker run -p 8080:8080 -e API_KEY=$(API_KEY) $(IMG_TAG)

to-fly:
	flyctl deploy
