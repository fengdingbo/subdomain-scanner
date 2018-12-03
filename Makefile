TARGET=./bin
ARCHS=amd64 386
LDFLAGS="-s -w"
#LDFLAGS="-s -w -X main.VERSION=0.3 -X 'main.GIT_HASH=`git log --stat |head -1|awk '{print $$2}'`' -X 'main.GO_VERSION=`go version`'"
BIN="sub-domain-scanner"
PARKAGE=`find ./bin|grep -v "tar.gz"`
current:
	go build -ldflags=${LDFLAGS}

windows:
	@for GOARCH in ${ARCHS}; do \
		echo "Building for windows $${GOARCH} ..." ; \
		mkdir -p ${TARGET} ; \
		GOOS=windows GOARCH=$${GOARCH} go build -ldflags=${LDFLAGS} \
		 -o ${TARGET}/${BIN}-for-windows-$${GOARCH}.exe ; \
	done; \
	echo "Done."

linux:
	@for GOARCH in ${ARCHS}; do \
		echo "Building for linux $${GOARCH} ..." ; \
		mkdir -p ${TARGET} ; \
		GOOS=linux GOARCH=$${GOARCH} go build -ldflags=${LDFLAGS} \
		 -o ${TARGET}/${BIN}-for-linux-$${GOARCH} ; \
	done; \
	echo "Done."

darwin:
	@for GOARCH in ${ARCHS}; do \
		echo "Building for darwin $${GOARCH} ..." ; \
		mkdir -p ${TARGET} ; \
		GOOS=darwin GOARCH=$${GOARCH} go build -ldflags=${LDFLAGS} \
		 -o ${TARGET}/${BIN}-for-darwin-$${GOARCH} ; \
	done; \
	echo "Done."
tag:
	@for S in ${PARKAGE}; do \
	tar  -zcvf $${S}.tar.gz $${S} ; \
	done; \

all: darwin linux windows

clean:
	rm -Rf ${TARGET}
