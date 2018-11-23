TARGET=./bin
BIN="sub-domain-scanner"
current:
	go build -o ${BIN}
clean:
	rm -Rf ${BIN}
