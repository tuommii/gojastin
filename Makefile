
DATE := $(shell date +%d.%m.%Y)

all: run

run: build
	./chall03

build:
	go build -o chall03 -ldflags '-X main.buildTime=$(DATE)' cmd/gojastin/main.go

# Runs test and bench: both go test -bench . ./...
# Now only tests
test:
	go test -v ./...

# Runs without other tests
bench:
	go test -v -run=XXX -bench . ./... -benchmem

deploy: build
	ssh -t $(DIOC_USER)@$(DIOC_IP) "sudo service chall03 stop"
	scp chall03 $(DIOC_USER)@$(DIOC_IP):/home/
	ssh -t $(DIOC_USER)@$(DIOC_IP) "sudo service chall03 start"
