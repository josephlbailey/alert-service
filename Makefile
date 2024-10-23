prep-dev-env:
	docker compose -f ./deploy/local/compose.yaml down -v && docker compose -f ./deploy/local/compose.yaml up -d --wait

dev-run:
	ENVIRONMENT=dev go run main.go

sqlc-gen:
	pushd ./internal/db && sqlc generate && popd

mock-gen:
	mockgen -source=./internal/db/store.go -destination=./internal/mock/store.go --package=mock