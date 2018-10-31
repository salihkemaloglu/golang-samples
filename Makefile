# import deploy config
# You can change the default deploy config with `make cnf="deploy_special.env" release`
dpl ?= deploy.env
include $(dpl)
# HELP
.PHONY: help
help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

# DOCKER TASKS
# Build the container
build: ## Build the container
	docker build -t $(APP_NAME) .

build-nc: ## Build the container without caching
	docker build --no-cache -t $(APP_NAME) .
## Run container on docker-compose environment
up:run-dev run-pro 
up-dev:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up
up-pro:
	docker-compose -f docker-compose.yml -f docker-compose.pro.yml up
## stop container on docker-compose environment
stop:stop-dev stop-pro 
stop-dev:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml stop
stop-pro:
	docker-compose -f docker-compose.yml -f docker-compose.pro.yml stop
## down container on docker-compose environment
down:down-dev down-pro 
down-dev:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml down
down-pro:
	docker-compose -f docker-compose.yml -f docker-compose.pro.yml down	

# Docker publish
push: push-local push-remote 
# Publish the local repository
push-local: 
	docker push localhost:5000/$(APP_NAME)
## Publish the remote repository
push-remote: 
	docker push $(DOCKER_REPO)/$(APP_NAME)
# Docker pull
pull:pull-local pull-remote
# Pull the local repository
pull-local:
	docker pull localhost:5000/$(APP_NAME)
# Pull the remote repository	
pull-remote:
	docker pull $(DOCKER_REPO)/$(APP_NAME)	
# Docker tagging image
tag-image:  
	docker tag $(APP_NAME) localhost:5000/$(APP_NAME)
#Local registry for docker images 
run-registry:
	docker run -d -p 5000:5000 --restart=always --name registry registry:2
stop-registry:
	docker stop registry
rm-registry:
	docker stop registry && docker rm -v registry
rmi-registry:
	docker rmi registry:2				
