VAR1 := $(shell cat .env)

dev:
	@echo "Starting development server"
	@echo $(VAR1)
	@go run main.go
