APP_NAME=monitor-agent
ENTRY=./cmd/agent

build:
	go build -o bin/$(APP_NAME) $(ENTRY)

run:
	go run $(ENTRY)

clean:
	rm -rf bin/$(APP_NAME)

