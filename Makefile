.PHONY: dev fmt golint test cover html clean
PACKAGES=`go list ./...`

app:
	@go run ./app/main.go

dev:
	@dev_appserver.py app-local.yaml --port=8080 --admin_port=8010

fmt:
	@for pkg in ${PACKAGES}; do \
	go fmt $$pkg; \
    done;

test:
	@RICHGO_FORCE_COLOR=1 ENV=test richgo test -v -mod=vendor -cover ./...
