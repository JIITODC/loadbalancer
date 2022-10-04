all: docker-build run run-compose

.PHONY: \
	build \
	build-debug \
	clean \
	docker-build \
	run \
	run-compose

build:
	@$(MAKE) -C docker-lb build-bare

build-debug:
	@$(MAKE) -C docker-lb build-debug-bare

clean:
	@echo "Cleaning artefacts"
	@rm lb.out 2> /dev/null

docker-build:
	@$(MAKE) -C Backend docker-build-bare

run: build run-compose
	@echo "Running the application"
	@./docker-lb/lb.out

run-compose:
	@$(MAKE) -C Backend run-compose-bare
