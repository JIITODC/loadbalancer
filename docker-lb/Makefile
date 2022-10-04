.PHONY: build-bare

build-bare:
	@echo "Building the application"
	@go build -o lb.out

build-debug-bare:
	@echo "Building the application with debug symbols"
	@go build -gcflags "all=-N -l" -o lb.out
