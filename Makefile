NOW := $(shell date +%Y.%m.%d.%H.%M.%S)

octopus:
	go build -v -ldflags '-X main.Version=$(NOW) -extldflags "-static"'

package: octopus
	MD5=$(shell openssl md5 octopus | awk '{print $$2}') && \
	cp octopus octopus-$$MD5 && \
	tar cfvz octopus-$$MD5.tar.gz octopus-$$MD5

package-docker:
	docker build --platform linux/386 -t octopus .
	docker run --rm -v="$$PWD:/host" octopus sh -c "cp octopus-*.tar.gz /host"

clean:
	rm -f octopus octopus-*

.PHONY: clean
