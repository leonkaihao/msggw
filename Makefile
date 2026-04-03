$(eval RELEASE_TAG := $(shell jq -r '.release_tag' component.json))
$(eval VERSION := $(shell jq -r '.version' component.json))
$(eval REGISTRY := $(shell jq -r '.registry' component.json))
$(eval TARGET_NAME := $(shell jq -r '.name' component.json))

TEST_SUITES= config funcs model operator parser service symbol
outdir ?= $(CURDIR)/bin
testsoutdir ?= $(CURDIR)/tests

.PHONY: all build release test clean

build:
	go mod tidy;
	CGO_ENABLED=0 go build -o $(outdir)/$(TARGET_NAME) $(CURDIR)/app/msggw
release:
	git tag -a $(RELEASE_TAG) -m "Release $(RELEASE_TAG)"
	git push origin $(RELEASE_TAG)
rm-release:
	git tag -d 

container:
	docker build -t msggw:$(RELEASE_TAG) -f ./build/docker/Dockerfile_msggw .
	docker tag msggw:$(RELEASE_TAG) $(REGISTRY)/msggw:$(RELEASE_TAG)
	docker tag msggw:$(RELEASE_TAG) msggw:$(VERSION)
publish:
	docker push $(REGISTRY)/msggw:$(RELEASE_TAG)
	docker tag $(REGISTRY)/msggw:$(RELEASE_TAG) $(REGISTRY)/msggw:$(VERSION)-${BUILD_NUMBER}
	docker push $(REGISTRY)/msggw:$(VERSION)-${BUILD_NUMBER}
test:
	@for pkg in $(TEST_SUITES); do \
	go test -v -cover $(CURDIR)/pkg/$$pkg; \
	done;
test/bench:
	go mod tidy && go test -bench=. -benchmem ./pkg/...
test/pack:
	@for pkg in $(TEST_SUITES); do \
	go test -c -o=$(testsoutdir)/$(TARGET_NAME).$$pkg.ctrl.test $(CURDIR)/pkg/$$pkg --tags=ctrl || exit 1; \
	go test -c -o=$(testsoutdir)/$(TARGET_NAME).$$pkg.ctrl.integrationtest --tags=integration_ctrl $(CURDIR)/pkg/$$pkg || exit 1; \
	done;
	mkdir -p $(testsoutdir)
clean:
	rm -rf $(outdir)
	rm -rf $(testsoutdir)
	go clean -modcache
	docker rmi $(REGISTRY)/msggw:$(RELEASE_TAG)
	docker rmi $(REGISTRY)/msggw:$(VERSION)-${BUILD_NUMBER}
