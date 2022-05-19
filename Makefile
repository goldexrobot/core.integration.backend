all: test swagger
	@:

tools:
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

swagger:
	cp api/v1/openapi.yaml docs/goldex_api_v1.yaml
	cp callback/openapi.yaml docs/business_callbacks.yaml

test:
	@go test ./signature
