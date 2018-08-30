.PHONY: container publish serve serve-container clean

app        := tailor
static-app := build/linux-amd64/$(app)

$(app): *.go
	go build -o $@
