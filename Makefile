.PHONY: help


help: ## Display this help screen
	@grep -E '^[a-z.A-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

check-dependencies: ## Check if all dependencies are installed
	@poetry --version
	@docker --version
	@poetry export

init: check-dependencies ## Initialize the project
	@poetry env activate
	@poetry install

add-dep: check-dependencies ## Install dependencies in the corresponding project. Ex: make add-dep model_training pandas = poetry add -G model_training pandas
	@poetry add -G $(word 2, $(MAKECMDGOALS)) $(word 3, $(MAKECMDGOALS))

requirements: check-dependencies ## Export the requirements.txt file for one of the projects
	@poetry export --without-hashes --format=requirements.txt --only $(word 2, $(MAKECMDGOALS)) -o requirements.txt
