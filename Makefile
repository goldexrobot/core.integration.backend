all: test swagger
	@:

tools:
	@:

swagger:
	@cp api/v1/openapi.yaml docs/api_v1.yaml

test:
	@go test ./pkg/signature
