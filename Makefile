all: test callback-swagger golang-apiv1
	@:

tools:
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

backend:
	$(info Backend:)
	@$(MAKE) --no-print-directory -C backend

terminal:
	$(info Terminal:)
	@$(MAKE) --no-print-directory -C terminal

test:
	$(info Go test)
	@go test ./signature

callback-swagger:
	$(info Swagger for callbacks)
	@MSYS_NO_PATHCONV=1 docker run --rm -v $(shell pwd):/goldex:rw -it quay.io/goswagger/swagger generate spec -m -w /goldex/callback -o /goldex/docs/swagger/backend/callback.swagger.json
	
golang-apiv1:
	$(info Golang API v1 from proto)
	@$(MAKE) --no-print-directory -C api/v1/golang

