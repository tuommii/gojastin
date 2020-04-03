
DATE := $(shell date +%d.%m.%Y)

all: run

run: build
	./chall03

build:
	go build -o chall03 -ldflags '-X main.buildTime=$(DATE)' cmd/gojastin/main.go

deploy: build
	mv main chall03
	ssh -t $(DIOC_USER)@$(DIOC_IP) "sudo service chall03 stop"
	scp chall03 $(DIOC_USER)@$(DIOC_IP):/home/
	ssh -t $(DIOC_USER)@$(DIOC_IP) "sudo service chall03 start"
