run: 
	@go run -mod=vendor .

test: 
	@go test -mod=vendor -v