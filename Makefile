GOMOD=link_checker
build: release
upx:
	upx -9 $(GOMOD)*
debug: compdbg
release: comprel upx
windows: export GOOS=windows
windows: export GOARCH=amd64
windows: release
comprel:
	go build -ldflags="-s -w" .
compdbg:
	go build -race -gcflags=all=-d=checkptr=0 .
clean:
	rm -fr $(GOMOD)*
