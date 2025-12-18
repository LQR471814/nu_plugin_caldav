codegen:
	go run ./internal/codegen \
		-pkg nuconv \
		-out ./internal/nuconv/nuconv.go
	sqlc generate

