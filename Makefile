TARGET=./bin
ARCHS=amd64 386
LDFLAGS="-s -w"
#LDFLAGS="-s -w -X main.VERSION=0.3 -X 'main.GIT_HASH=`git log --stat |head -1|awk '{print $$2}'`' -X 'main.GO_VERSION=`go version`'"
BIN="sub-domain-scanner"
PARKAGE=`find ./bin -type d`
current:
	go build -ldflags=${LDFLAGS}

windows:
	@for GOARCH in ${ARCHS}; do \
		echo "Building for windows $${GOARCH} ..." ; \
		mkdir -p ${TARGET}/${BIN}-windows-$${GOARCH} ; \
		GOOS=windows GOARCH=$${GOARCH} go build -ldflags=${LDFLAGS} \
		 -o ${TARGET}/${BIN}-windows-$${GOARCH}/${BIN}.exe ; \
	done; \
	echo "Done."

linux:
	@for GOARCH in ${ARCHS}; do \
		echo "Building for linux $${GOARCH} ..." ; \
		mkdir -p ${TARGET}/${BIN}-linux-$${GOARCH} ; \
		GOOS=linux GOARCH=$${GOARCH} go build -ldflags=${LDFLAGS} \
		 -o ${TARGET}/${BIN}-linux-$${GOARCH}/${BIN} ; \
	done; \
	echo "Done."

darwin:
	@for GOARCH in ${ARCHS}; do \
		echo "Building for darwin $${GOARCH} ..." ; \
		mkdir -p ${TARGET}/${BIN}-darwin-$${GOARCH} ; \
		GOOS=darwin GOARCH=$${GOARCH} go build -ldflags=${LDFLAGS} \
		 -o ${TARGET}/${BIN}-darwin-$${GOARCH}/${BIN} ; \
	done; \
	echo "Done."
tag:
	@for S in ${PARKAGE}; do \
	cp -Rfv dict $${S}/ ; \
	tar  -zcvf $${S}.tar.gz $${S} ; \
	done; \

all: darwin linux windows

clean:
	rm -Rf ${TARGET}
