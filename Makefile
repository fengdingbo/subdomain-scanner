TARGET=./bin
ARCHS=amd64 386
LDFLAGS="-s -w"
BIN="sub-domain-scanner"
current:
	go build -o ${BIN}

windows:
	@for GOARCH in ${ARCHS}; do \
		echo "Building for windows $${GOARCH} ..." ; \
		mkdir -p ${TARGET} ; \
		GOOS=windows GOARCH=$${GOARCH} go build -ldflags=${LDFLAGS} \
		 -o ${TARGET}/${BIN}-for-$${GOARCH}.exe ; \
	done; \
	echo "Done."

clean:
	rm -Rf ${TARGET}
