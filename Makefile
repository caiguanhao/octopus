NOW := $(shell date +%Y.%m.%d.%H.%M.%S)

octopus:
	go build -ldflags '-X main.Version=$(NOW) -extldflags "-static"'

package: octopus
	MD5=$(shell openssl md5 octopus | awk '{print $$2}') && \
	cp octopus octopus-$$MD5 && \
	tar cfvz octopus-$$MD5.tar.gz octopus-$$MD5

clean:
	rm -f octopus octopus-*

.PHONY: clean
