.PHONY: docker-build-bare

docker-build-bare:
	@echo "Building the docker image"
	@docker build -t load-balancer:latest .

run-compose-bare:
	@echo "Spawning fake load instances"
	@docker-compose up
