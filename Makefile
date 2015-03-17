NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m
DEPS = $(go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)

init:
	@echo "$(OK_COLOR)==> This project uses Godep, downloading...$(NO_COLOR)"
	go get github.com/tools/godep
	go get github.com/google/go-github/github
	go get github.com/franela/goblin
	go get github.com/onsi/gomega
	go get golang.org/x/net/context
	go get golang.org/x/oauth2
	go get github.com/gorilla/schema
	go get github.com/kelseyhightower/envconfig

format:
	@echo "$(OK_COLOR)==> Formatting$(NO_COLOR)"
	godep go fmt ./...

test:
	@echo "$(OK_COLOR)==> Testing$(NO_COLOR)"
	godep go test -short $(filter-out $@,$(MAKECMDGOALS)) ./...

test-int:
	@echo "$(OK_COLOR)==> Integration Testing$(NO_COLOR)"
	godep go test $(filter-out $@,$(MAKECMDGOALS)) ./...

install: format
	@echo "$(OK_COLOR)==> Building and Installing$(NO_COLOR)"
	godep go install

clean:
	rm -rf $(GOPATH)/pkg/* $(GOPATH)/bin/marvin

run: install
	@echo "$(OK_COLOR)==> Running$(NO_COLOR)"
	$(GOPATH)/bin/marvin
