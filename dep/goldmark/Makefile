.PHONY: test fuzz

test:
	go test -coverprofile=profile.out -coverpkg=utilware/dep/goldmark,utilware/dep/goldmark/ast,utilware/dep/goldmark/extension,utilware/dep/goldmark/extension/ast,utilware/dep/goldmark/parser,utilware/dep/goldmark/renderer,utilware/dep/goldmark/renderer/html,utilware/dep/goldmark/text,utilware/dep/goldmark/util ./...

cov: test
	go tool cover -html=profile.out

fuzz:
	which go-fuzz > /dev/null 2>&1 || (GO111MODULE=off go get -u github.com/dvyukov/go-fuzz/go-fuzz github.com/dvyukov/go-fuzz/go-fuzz-build; GO111MODULE=off go get -d github.com/dvyukov/go-fuzz-corpus; true)
	rm -rf ./fuzz/corpus
	rm -rf ./fuzz/crashers
	rm -rf ./fuzz/suppressions
	rm -f ./fuzz/fuzz-fuzz.zip
	cd ./fuzz && go-fuzz-build
	cd ./fuzz && go-fuzz
