
DATE := $(shell date +%d.%m.%Y)

all: run

run: build
	./chall03

build:
	go build -o chall03 -ldflags '-X main.buildTime=$(DATE)' cmd/gojastin/main.go

test:
	go test ./...

# Runs without other tests, remove -run=XXX
bench:
	go test -run=XXX -bench . ./...

deploy: build
	ssh -t $(DIOC_USER)@$(DIOC_IP) "sudo service chall03 stop"
	scp chall03 $(DIOC_USER)@$(DIOC_IP):/home/
	ssh -t $(DIOC_USER)@$(DIOC_IP) "sudo service chall03 start"
