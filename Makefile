.PHONY: test cover

cover:
	go test -coverprofile=cov.out ./...
	go tool cover -html=cov.out -o cover.html