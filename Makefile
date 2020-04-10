NOW := $(shell date +%Y.%m.%d.%H.%M.%S)

octopus:
	go build -ldflags '-X main.Version=$(NOW) -extldflags "-static"'

clean:
	rm -f octopus

.PHONY: clean
