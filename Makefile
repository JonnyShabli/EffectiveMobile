swagger-gen:
	swag init --parseDependency --parseDependencyLevel 0  --generalInfo ./internal/controller/subsHandlers.go
.PHONY: swagger-gen
