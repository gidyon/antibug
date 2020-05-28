PROJECT_NAME := antibug
PKG := gtuhub.com/gidyon/$(PROJECT_NAME)

compile:
	go build -i -v -o service .

docker_build:
ifdef tag
	@docker build -t gidyon/$(PROJECT_NAME)-antimicrobial:$(tag) .
else
	@docker build -t gidyon/$(PROJECT_NAME)-antimicrobial:latest .
endif

docker_tag:
ifdef tag
	@docker tag gidyon/$(PROJECT_NAME)-antimicrobial:$(tag) gidyon/$(PROJECT_NAME)-antimicrobial:$(tag)
else
	@docker tag gidyon/$(PROJECT_NAME)-antimicrobial:latest gidyon/$(PROJECT_NAME)-antimicrobial:latest
endif

docker_push:
ifdef tag
	@docker push gidyon/$(PROJECT_NAME)-antimicrobial:$(tag)
else
	@docker push gidyon/$(PROJECT_NAME)-antimicrobial:latest
endif

build_image: docker_build docker_tag docker_push

build: compile docker_build docker_tag docker_push