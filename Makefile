.PHONY: build  # make build, here build is not a real file
build:
	./build/build.sh go-nuva ./bin/go-nuva

linux:
	GOOS=linux GOARCH=amd64 ./build/build.sh go-nuva ./bin/go-nuva

darwin:
	GOOS=darwin GOARCH=arm64 ./build/build.sh go-nuva ./bin/go-nuva

clean:
	rm -rf ./bin