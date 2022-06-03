all: test swagger
	@:

tools:
	go install golang.org/x/tools/cmd/stringer@latest

swagger:
	cp api/v1/openapi.yaml docs/api_v1.yaml

test:
	@go test ./pkg/signature
