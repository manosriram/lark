run:
	go build -o lark && ./lark source.lark
build:
	go build -o lark
test:
	go test -v
