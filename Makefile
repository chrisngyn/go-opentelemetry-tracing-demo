.PHONY: openapi
openapi: openapi-http

.PHONY: openapi-http
openapi-http:
	@./scripts/openapi-http.sh service1 internal/service1 main

.PHONY: proto
proto:
	buf generate api/protobuf