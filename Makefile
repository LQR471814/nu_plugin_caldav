codegen:
	go run ./internal/codegen \
		-pkg nuconv \
		-out ./internal/nuconv/nuconv.gen.go
	sqlc generate

