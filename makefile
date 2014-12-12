all: presets.go

presets.go: 
	go get github.com/jteeuwen/go-bindata/...
	go-bindata -pkg=punkt -o presets.go -ignore="README|\\.py" data
	
clean: 
	rm presets.go

.PHONY: all
.PHONY: $(BINARIES)
