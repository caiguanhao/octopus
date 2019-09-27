octopus:
	go build -ldflags '-extldflags "-static"'

clean:
	rm -f octopus

.PHONY: clean
