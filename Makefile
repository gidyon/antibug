PROJECT_NAME := antibug
PKG := github.com/gidyon/antibug
SERVICE_PKG_BUILD := ${PKG}/cmd/gateway

API_IN_PATH := api/proto

API_OUT_PATH := pkg/api
SWAGGER_DOC_OUT_PATH := api/swagger


proto_compile_pathogen:
	protoc -I=$(API_IN_PATH) -I=third_party --go_out=plugins=grpc:$(API_OUT_PATH)/pathogen pathogen.proto &&\
	protoc -I=$(API_IN_PATH) -I=third_party --grpc-gateway_out=logtostderr=true:$(API_OUT_PATH)/pathogen pathogen.proto &&\
	protoc -I=$(API_IN_PATH) -I=third_party --swagger_out=logtostderr=true:$(SWAGGER_DOC_OUT_PATH) pathogen.proto

proto_compile_antimicrobial:
	protoc -I=$(API_IN_PATH) -I=third_party --go_out=plugins=grpc:$(API_OUT_PATH)/antimicrobial antimicrobial.proto &&\
	protoc -I=$(API_IN_PATH) -I=third_party --grpc-gateway_out=logtostderr=true:$(API_OUT_PATH)/antimicrobial antimicrobial.proto &&\
	protoc -I=$(API_IN_PATH) -I=third_party --swagger_out=logtostderr=true:$(SWAGGER_DOC_OUT_PATH) antimicrobial.proto

proto_compile_facility:
	protoc -I=$(API_IN_PATH) -I=third_party --go_out=plugins=grpc:$(API_OUT_PATH)/facility facility.proto &&\
	protoc -I=$(API_IN_PATH) -I=third_party --grpc-gateway_out=logtostderr=true:$(API_OUT_PATH)/facility facility.proto &&\
	protoc -I=$(API_IN_PATH) -I=third_party --swagger_out=logtostderr=true:$(SWAGGER_DOC_OUT_PATH) facility.proto

proto_compile_account:
	protoc -I=$(API_IN_PATH) -I=third_party --go_out=plugins=grpc:$(API_OUT_PATH)/account account.proto &&\
	protoc -I=$(API_IN_PATH) -I=third_party --grpc-gateway_out=logtostderr=true:$(API_OUT_PATH)/account account.proto &&\
	protoc -I=$(API_IN_PATH) -I=third_party --swagger_out=logtostderr=true:$(SWAGGER_DOC_OUT_PATH) account.proto

proto_compile_culture:
	protoc -I=$(API_IN_PATH) -I=third_party --go_out=plugins=grpc:$(API_OUT_PATH)/culture culture.proto &&\
	protoc -I=$(API_IN_PATH) -I=third_party --grpc-gateway_out=logtostderr=true:$(API_OUT_PATH)/culture culture.proto &&\
	protoc -I=$(API_IN_PATH) -I=third_party --swagger_out=logtostderr=true:$(SWAGGER_DOC_OUT_PATH) culture.proto

proto_compile_antibiogram:
	protoc -I=$(API_IN_PATH) -I=third_party --go_out=plugins=grpc:$(API_OUT_PATH)/antibiogram antibiogram.proto &&\
	protoc -I=$(API_IN_PATH) -I=third_party --grpc-gateway_out=logtostderr=true:$(API_OUT_PATH)/antibiogram antibiogram.proto &&\
	protoc -I=$(API_IN_PATH) -I=third_party --swagger_out=logtostderr=true:$(SWAGGER_DOC_OUT_PATH) antibiogram.proto

proto_compile_all: proto_compile_pathogen proto_compile_antimicrobial proto_compile_facility proto_compile_account proto_compile_culture proto_compile_antibiogram

run_app:
	go run cmd/gateway/*.go

run_account:
	cd cmd/modules/account && go build -o service && JWT_SIGNING_KEY=albahtep ./service -config-file=/home/gideon/go/src/github.com/gidyon/antibug/configs/account.dev.yml

run_antibiogram:
	cd cmd/modules/antibiogram && go build -o service && JWT_SIGNING_KEY=albahtep ./service -config-file=/home/gideon/go/src/github.com/gidyon/antibug/configs/antibiogram.dev.yml

run_antimicrobial:
	cd cmd/modules/antimicrobial && go build -o service && JWT_SIGNING_KEY=albahtep ./service -config-file=/home/gideon/go/src/github.com/gidyon/antibug/configs/antimicrobial.dev.yml

run_culture:
	cd cmd/modules/culture && go build -o service && JWT_SIGNING_KEY=albahtep ./service -config-file=/home/gideon/go/src/github.com/gidyon/antibug/configs/culture.dev.yml

run_facility:
	cd cmd/modules/facility && go build -o service && JWT_SIGNING_KEY=albahtep ./service -config-file=/home/gideon/go/src/github.com/gidyon/antibug/configs/facility.dev.yml

run_pathogen:
	cd cmd/modules/pathogen && go build -o service && JWT_SIGNING_KEY=albahtep ./service -config-file=/home/gideon/go/src/github.com/gidyon/antibug/configs/pathogen.dev.yml

setup_dev: ## Sets up a development environment for the digimed project
	@cd deployments/compose/dev &&\
	docker-compose up -d

setup_mysql: ## Sets up a development environment for the digimed project
	@cd deployments/compose/dev &&\
	docker-compose up -d mysql

setup_redis: ## Sets up a development environment for the digimed project
	@cd deployments/compose/dev &&\
	docker-compose up -d redis

teardown_dev: ## Tear down development environment for the digimed project
	@cd deployments/compose/dev &&\
	docker-compose down

compile_gateway:
	go build -i -v -o gateway $(SERVICE_PKG_BUILD)
	
docker_build_gateway: ## Create a docker image for the service
ifdef tag
	@docker build -t gidyon/$(PROJECT_NAME)-gateway:$(tag) .
else
	@docker build -t gidyon/$(PROJECT_NAME)-gateway:latest .
endif

docker_tag_gateway:
ifdef tag
	@docker tag gidyon/$(PROJECT_NAME)-gateway:$(tag) gidyon/$(PROJECT_NAME)-gateway:$(tag)
else
	@docker tag gidyon/$(PROJECT_NAME)-gateway:latest gidyon/$(PROJECT_NAME)-gateway:latest
endif

docker_push_gateway:
ifdef tag
	@docker push gidyon/$(PROJECT_NAME)-gateway:$(tag)
else
	@docker push gidyon/$(PROJECT_NAME)-gateway:latest
endif

build_and_push_gateway: compile_gateway docker_build_gateway docker_tag_gateway docker_push_gateway

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'