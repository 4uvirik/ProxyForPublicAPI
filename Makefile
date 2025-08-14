APP_NAME := Proxy for public API
ENTRYPOINT := cmd/main.go

.PHONY: tests cover

tests:
	go test ./... -coverprofile=covarage.out
	go tool cover -func=covarage.out | grep total

cover:
	go tool cover -html=covarage.out -o coverage.html
	firefox coverage.html