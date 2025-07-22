swagger-gen:
	swag init --outputTypes json --generalInfo ./internal/controller/subsHandlers.go
.PHONY: swagger-gen
