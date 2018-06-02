DOCKER=docker
CURDIR=$(shell pwd)
CONTAINER=meteorologist
PORT=4000
IP=$(shell ipconfig getifaddr en0)

shell:
	$(DOCKER) run -it -p $(PORT):$(PORT) -v "$(CURDIR)":/go/src/github.com/lukaszsagol/meteorologist $(CONTAINER) /bin/bash
.PHONY: shell

build:
	$(DOCKER) build -t $(CONTAINER) .
.PHONY: build

clean:
	$(DOCKER) rmi -f $(CONTAINER)
.PHONY: clean
