NAME = demo

all: assets
	go build -o $(NAME)

assets: clean
	go-bindata \
		-pkg=main \
		-o=bindata.go \
		-prefix="files" \
		files/...

clean:
	rm -fr $(NAME) bindata.go

.PHONY: all assets clean
