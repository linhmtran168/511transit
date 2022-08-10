VERSION=$(shell cat ./VERSION)
IMG_TAG=linhmtran168/511transit:$(VERSION)

.PHONY: tidy back-dev front-dev front-install front-build

tidy:
	go mod tidy

back-dev:
	air

front-install:
	cd ./web && npm install

front-dev:
	cd ./web && npm run dev

front-build:
	cd ./web && npm run build

build-prod-img:
	docker build -t $(IMG_TAG) .

run-prod-img:
	docker run -p 8080:80 -e API_KEY=$(API_KEY) $(IMG_TAG)

# create_mock:
# 	@docker run -v "$(CURDIR)":/src \
# 		-w /src vektra/mockery --all  --output "./internal/todo/mocks" --dir "./internal/todo/"