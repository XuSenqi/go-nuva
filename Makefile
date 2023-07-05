# run from repository root

TARGET=./bin/go-nuva
MAINFILE=main.go

.PHONY: build  # make build, here build is not a real file
build:
	rm -f ${TARGET}
	MAINFILE=${MAINFILE} bash ./build.sh go-nuva ${TARGET}

clean:
	rm -rf ./bin