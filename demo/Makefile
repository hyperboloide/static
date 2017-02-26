NAME = demo

all: assets
	go build -o $(NAME)

assets: clean
	go-bindata \
		-pkg=files \
		-o=files/bindata.go \
		-prefix="files" \
		files/...

clean:
	rm -fr $(NAME) files/bindata.go

.PHONY: all assets clean
