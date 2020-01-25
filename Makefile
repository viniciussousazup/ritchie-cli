# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=rit
CMD_PATH=./cmd/main.go
DIST=dist
DIST_MAC=$(DIST)/mac
DIST_LINUX=$(DIST)/linux
DIST_WIN=$(DIST)/windows
VERSION=$(RELEASE_VERSION)
GIT_REMOTE=https://$(GIT_USERNAME):$(GIT_PASSWORD)@github.com/ZupIT/ritchie-cli
MODULE=$(shell go list -m)
DATE=$(shell date +%D_%H:%M)

build:
	mkdir -p $(DIST_MAC) $(DIST_LINUX) $(DIST_WIN)
	#LINUX
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags '-X $(MODULE)/pkg/cmd.Version=$(VERSION) -X $(MODULE)/pkg/cmd.BuildDate=$(DATE)' -o ./$(DIST_LINUX)/$(BINARY_NAME) -v $(CMD_PATH)
	#MAC
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -ldflags '-X $(MODULE)/pkg/cmd.Version=$(VERSION) -X $(MODULE)/pkg/cmd.BuildDate=$(DATE)' -o ./$(DIST_MAC)/$(BINARY_NAME) -v $(CMD_PATH)
	#WINDOWS 64
	GOOS=windows GOARCH=amd64 $(GOBUILD) -ldflags '-X $(MODULE)/pkg/cmd.Version=$(VERSION) -X $(MODULE)/pkg/cmd.BuildDate=$(DATE)' -o ./$(DIST_WIN)/$(BINARY_NAME).exe -v $(CMD_PATH)

build-dev:
	mkdir -p $(DIST_MAC) $(DIST_LINUX)
	#LINUX
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags '-X $(MODULE)/pkg/cmd.Version=dev -X $(MODULE)/pkg/cmd.BuildDate=$(DATE) -X $(MODULE)/pkg/env.ServerUrl=https://ritchie-server-dev.itiaws.dev -X $(MODULE)/pkg/env.Environment=dev' -o ./$(DIST_LINUX)/$(BINARY_NAME) -v $(CMD_PATH)
	#MAC
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -ldflags '-X $(MODULE)/pkg/cmd.Version=dev -X $(MODULE)/pkg/cmd.BuildDate=$(DATE) -X $(MODULE)/pkg/env.ServerUrl=https://ritchie-server-dev.itiaws.dev -X $(MODULE)/pkg/env.Environment=dev' -o ./$(DIST_MAC)/$(BINARY_NAME) -v $(CMD_PATH)

build-qa:
	mkdir -p $(DIST_MAC) $(DIST_LINUX)
	#LINUX
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags '-X $(MODULE)/pkg/cmd.Version=dev -X $(MODULE)/pkg/cmd.BuildDate=$(DATE) -X $(MODULE)/pkg/env.ServerUrl=https://ritchie-server.itiaws.dev -X $(MODULE)/pkg/env.Environment=qa' -o ./$(DIST_LINUX)/$(BINARY_NAME) -v $(CMD_PATH)
	#MAC
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -ldflags '-X $(MODULE)/pkg/cmd.Version=dev -X $(MODULE)/pkg/cmd.BuildDate=$(DATE) -X $(MODULE)/pkg/env.ServerUrl=https://ritchie-server.itiaws.dev -X $(MODULE)/pkg/env.Environment=qa' -o ./$(DIST_MAC)/$(BINARY_NAME) -v $(CMD_PATH)

test:
	$(GOTEST) -short ./...

test-v:
	$(GOTEST) -v ./...

release:
	envsubst < "Release.md.template" > "Release.md"
	git config --global user.email "$(GIT_EMAIL)"
	git config --global user.name "$(GIT_USER)"
	git add .
	git commit -m "release"
	git push $(GIT_REMOTE) HEAD:release-$(RELEASE_VERSION)
	git tag -a $(RELEASE_VERSION) -m "release"
	git push $(GIT_REMOTE) $(RELEASE_VERSION)
	curl --user $(GIT_USERNAME):$(GIT_PASSWORD) -X POST https://api.github.com/repos/ZupIT/ritchie-cli/pulls -H 'Content-Type: application/json' -d '{ "title": "Release $(RELEASE_VERSION) merge", "body": "Release $(RELEASE_VERSION) merge with master", "head": "release-$(RELEASE_VERSION)", "base": "master" }'
	aws s3 sync dist s3://ritchie-cli-bucket234376412767550/$(RELEASE_VERSION) --include "*"
	echo "$(RELEASE_VERSION)" > stable.txt
	aws s3 sync . s3://ritchie-cli-bucket234376412767550/ --exclude "*" --include "stable.txt"

publish:
	echo "Do nothing"

clean:
	rm -rf $(DIST_MAC) $(DIST_LINUX) $(DIST_WIN)
